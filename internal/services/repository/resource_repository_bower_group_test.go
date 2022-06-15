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

func testAccResourceRepositoryBowerGroup() repository.BowerGroupRepository {
	return repository.BowerGroupRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
		Group: repository.Group{
			MemberNames: []string{},
		},
	}
}

func testAccResourceRepositoryBowerGroupConfig(repo repository.BowerGroupRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryBowerGroupTemplate := template.Must(template.New("BowerGroupRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryBowerGroup))
	if err := resourceRepositoryBowerGroupTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryBowerGroup(t *testing.T) {
	repoHosted := testAccResourceRepositoryBowerHosted()
	repo := testAccResourceRepositoryBowerGroup()
	repo.Group.MemberNames = append(repo.Group.MemberNames, repoHosted.Name)
	resourceName := "nexus_repository_bower_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryBowerHostedConfig(repoHosted) + testAccResourceRepositoryBowerGroupConfig(repo),
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
						resource.TestCheckResourceAttr(resourceName, "group.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "group.0.member_names.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "group.0.member_names.0", repo.Group.MemberNames[0]),
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
