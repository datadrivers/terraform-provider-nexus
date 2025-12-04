package repository_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryRHosted() repository.RHostedRepository {
	writePolicy := repository.StorageWritePolicyAllow

	return repository.RHostedRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		},
		Component: &repository.Component{
			ProprietaryComponents: true,
		},
	}
}

func testAccResourceRepositoryRHostedConfig(repo repository.RHostedRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryRHostedTemplate := template.Must(template.New("RHostedRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryRHosted))
	if err := resourceRepositoryRHostedTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryRHosted(t *testing.T) {
	repo := testAccResourceRepositoryRHosted()
	resourceName := "nexus_repository_r_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryRHostedConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(resourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(resourceName, "online", strconv.FormatBool(repo.Online)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(resourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
						resource.TestCheckResourceAttr(resourceName, "storage.0.write_policy", string(*repo.Storage.WritePolicy)),
						resource.TestCheckResourceAttr(resourceName, "component.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "component.0.proprietary_components", strconv.FormatBool(repo.Component.ProprietaryComponents)),
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
