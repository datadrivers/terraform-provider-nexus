package repository_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryRubygemsGroup() repository.RubyGemsGroupRepository {
	return repository.RubyGemsGroupRepository{
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

func testAccResourceRepositoryRubygemsGroupConfig(repo repository.RubyGemsGroupRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryRubygemsGroupTemplate := template.Must(template.New("RubygemsGroupRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryRubygemsGroup))
	if err := resourceRepositoryRubygemsGroupTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryRubygemsGroup(t *testing.T) {
	repoHosted := testAccResourceRepositoryRubygemsHosted()
	repo := testAccResourceRepositoryRubygemsGroup()
	repo.Group.MemberNames = append(repo.Group.MemberNames, repoHosted.Name)
	resourceName := "nexus_repository_rubygems_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryRubygemsHostedConfig(repoHosted) + testAccResourceRepositoryRubygemsGroupConfig(repo),
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
