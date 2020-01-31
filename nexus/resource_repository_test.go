package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryAptHosted(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoAptDistribution := "bionic"
	repoAptSigningKeypair := acctest.RandString(10)
	repoAptSigningPassphrase := acctest.RandString(10)
	repoCleanupPolicyNames := []string{"weekly-cleanup"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryAptHosted(repoName, repoAptDistribution, repoAptSigningKeypair, repoAptSigningPassphrase, repoCleanupPolicyNames),
			},
		},
	})
}

func testAccRepositoryAptHosted(name string, aptDistribution string, aptSigningKEypair string, aptSigningPassphrase string, cleanupPolicyNames []string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "apt_hosted" {
	name   = "%s"
	format = "apt"
	type   = "hosted"

    apt {
		distribution = "%s"
	}

    apt_signing {
		keypair    = "%s"
		passphrase = "%s"
	}

	storage {

	}
}
`, name, aptDistribution, aptSigningKEypair, aptSigningPassphrase)
}

func TestAccRpositoryBowerHosted(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryBowerHosted(repoName),
			},
		},
	})
}

func testAccRepositoryBowerHosted(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "bower_hosted" {
	name   = "%s"
	format = "%s"
	type   = "%s"

	bower {
		
	}

	storage {

	}
}`, name, "bower", "hosted")
}

func TestAccRepositoryDockerHostedWithPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoOnline := true
	repoHTTPPort := acctest.RandIntRange(32767, 49152)
	repoHTTPSPort := acctest.RandIntRange(49153, 65535)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryDockerHostedWithPorts(repoName, repoOnline, repoHTTPPort, repoHTTPSPort),
			},
		},
	})
}

func testAccRepositoryDockerHostedWithPorts(name string, online bool, httpPort int, httpsPort int) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_hosted" {
	name   = "%s"
	format = "docker"
	type   = "hosted"
	online = %s

	docker {
		http_port        = %d
		https_port       = %d
		force_basic_auth = true
		v1enabled        = true
	}

	storage {

	}
}`, name, strconv.FormatBool(online), httpPort, httpsPort)
}
func TestAccRepositoryDockerHostedWithoutPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoOnline := true

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryDockerHostedWithoutPorts(repoName, repoOnline),
			},
		},
	})
}

func testAccRepositoryDockerHostedWithoutPorts(name string, online bool) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_hosted" {
	name   = "%s"
	format = "docker"
	type   = "hosted"
	online = %s

	docker {
		force_basic_auth = true
		v1enabled        = true
	}

	storage {

	}
}`, name, strconv.FormatBool(online))
}

func TestAccRepositoryDockerProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoHTTPPort := acctest.RandIntRange(32767, 49152)
	repoHTTPSPort := acctest.RandIntRange(49153, 65535)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryDockerProxy(repoName, repoHTTPPort, repoHTTPSPort),
			},
		},
	})
}

func testAccRepositoryDockerProxy(name string, httpPort int, httpsPort int) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_proxy" {
	name   = "%s"
	type   = "proxy"
	format = "docker"

	docker {
		http_port = %d
		https_port = %d
		force_basic_auth = true
		v1enabled = true
	}

    docker_proxy {
		index_url  = "https://index.docker.io"
		index_type = "HUB"
	}

    http_client {

	}

	negative_cache {

	}

	proxy {
		remote_url = "https://registry.npmjs.org"
	}

	storage {

	}
}`, name, httpPort, httpsPort)
}
