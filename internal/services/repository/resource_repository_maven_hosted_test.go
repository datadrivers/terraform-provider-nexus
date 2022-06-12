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

func testAccResourceRepositoryMavenHosted() repository.MavenHostedRepository {
	writePolicy := repository.StorageWritePolicyAllow
	versionPolicy := repository.MavenVersionPolicyRelease
	layoutPolicy := repository.MavenLayoutPolicyStrict
	contentDisposition := repository.MavenContentDispositionAttachment

	return repository.MavenHostedRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Maven: repository.Maven{
			VersionPolicy:      &versionPolicy,
			LayoutPolicy:       &layoutPolicy,
			ContentDisposition: &contentDisposition,
		},
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
			WritePolicy:                 &writePolicy,
		},
		Cleanup: &repository.Cleanup{
			PolicyNames: []string{"cleanup-weekly"},
		},
		Component: &repository.Component{
			ProprietaryComponents: true,
		},
	}
}

func testAccResourceRepositoryMavenHostedConfig(repo repository.MavenHostedRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryMavenHostedTemplate := template.Must(template.New("MavenHostedRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryMavenHosted))
	if err := resourceRepositoryMavenHostedTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryMavenHosted(t *testing.T) {
	repo := testAccResourceRepositoryMavenHosted()
	resourceName := "nexus_repository_maven_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryMavenHostedConfig(repo),
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
						resource.TestCheckResourceAttr(resourceName, "cleanup.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cleanup.0.policy_names.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "cleanup.0.policy_names.0", repo.Cleanup.PolicyNames[0]),
						resource.TestCheckResourceAttr(resourceName, "component.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "component.0.proprietary_components", strconv.FormatBool(repo.Component.ProprietaryComponents)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "maven.0.version_policy", string(*repo.Maven.VersionPolicy)),
						resource.TestCheckResourceAttr(resourceName, "maven.0.layout_policy", string(*repo.Maven.LayoutPolicy)),
						resource.TestCheckResourceAttr(resourceName, "maven.0.content_disposition", string(*repo.Maven.ContentDisposition)),
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
