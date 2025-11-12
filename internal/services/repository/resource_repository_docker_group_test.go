package repository_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"text/template"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryDockerGroup(name string) repository.DockerGroupRepository {
	return repository.DockerGroupRepository{
		Name:   name,
		Online: true,
		Docker: repository.Docker{
			ForceBasicAuth: true,
			HTTPPort:       tools.GetIntPointer(rand.Intn(999) + 32000),
			HTTPSPort:      tools.GetIntPointer(rand.Intn(999) + 33000),
			V1Enabled:      false,
		},
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
		Group: repository.GroupDeploy{
			MemberNames: []string{},
		},
	}
}

func testAccResourceRepositoryDockerGroupConfig(repo repository.DockerGroupRepository) string {
	buf := &bytes.Buffer{}
	resourceRepositoryDockerGroupTemplate := template.Must(template.New("DockerGroupRepository").Funcs(acceptance.TemplateFuncMap).Parse(acceptance.TemplateStringRepositoryDockerGroup))
	if err := resourceRepositoryDockerGroupTemplate.Execute(buf, repo); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestAccResourceRepositoryDockerGroup(t *testing.T) {
	nameHosted := fmt.Sprintf("acceptance-%s", acctest.RandString(10))
	nameGroup := fmt.Sprintf("acceptance-%s", acctest.RandString(10))

	repoHosted := testAccResourceRepositoryDockerHosted(nameHosted)
	repoGroup := testAccResourceRepositoryDockerGroup(nameGroup)

	subdomain := ""
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "false" {
		subdomain = nameGroup
	}
	repoGroup.Docker.Subdomain = &subdomain
	repoGroup.Group.MemberNames = append(repoGroup.Group.MemberNames, repoHosted.Name)

	resourceName := "nexus_repository_docker_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryDockerHostedConfig(repoHosted) + testAccResourceRepositoryDockerGroupConfig(repoGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repoGroup.Name),
						resource.TestCheckResourceAttr(resourceName, "name", repoGroup.Name),
						resource.TestCheckResourceAttr(resourceName, "online", strconv.FormatBool(repoGroup.Online)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.blob_store_name", repoGroup.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(resourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repoGroup.Storage.StrictContentTypeValidation)),
						resource.TestCheckResourceAttr(resourceName, "group.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "group.0.member_names.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "group.0.member_names.0", repoGroup.Group.MemberNames[0]),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "docker.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "docker.0.force_basic_auth", strconv.FormatBool(repoGroup.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.http_port", strconv.Itoa(*repoGroup.Docker.HTTPPort)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.https_port", strconv.Itoa(*repoGroup.Docker.HTTPSPort)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.v1_enabled", strconv.FormatBool(repoGroup.Docker.V1Enabled)),
						resource.TestCheckResourceAttr(resourceName, "docker.0.subdomain", subdomain),
					),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           repoGroup.Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
