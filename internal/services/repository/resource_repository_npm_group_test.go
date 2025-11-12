package repository_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryNpmGroup() repository.NpmGroupRepository {
	return repository.NpmGroupRepository{
		Name:   fmt.Sprintf("test-repo-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
		Group: repository.GroupDeploy{
			MemberNames: []string{},
		},
	}
}

func testAccResourceRepositoryNpmGroupConfig(repo repository.NpmGroupRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryNpmGroupTemplate := template.Must(template.New("NpmGroupRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryNpmGroup))
	if err := resourceRepositoryNpmGroupTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryNpmGroup(t *testing.T) {
	repoHosted := testAccResourceRepositoryNpmHosted()
	repo := testAccResourceRepositoryNpmGroup()
	repo.Group.MemberNames = append(repo.Group.MemberNames, repoHosted.Name)

	writableMember := ""
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "false" {
		writableMember = repoHosted.Name
	}

	repo.Group.WritableMember = &writableMember
	resourceName := "nexus_repository_npm_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryNpmHostedConfig(repoHosted) + testAccResourceRepositoryNpmGroupConfig(repo),
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
						resource.TestCheckResourceAttr(resourceName, "group.0.writable_member", writableMember),
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
