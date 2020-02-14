package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryAptHosted(t *testing.T) {
	t.Parallel()

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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "format", "apt"),
					resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "type", "hosted"),
				),
			},
			{
				ResourceName:      "nexus_repository.apt_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// apt_signing not returned by API
				ImportStateVerifyIgnore: []string{"apt_signing"},
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
	bowerRewritePackageURLs := true

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryBowerHosted(repoName, bowerRewritePackageURLs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "format", "bower"),
					resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "type", "hosted"),
				),
			},
			{
				ResourceName:      "nexus_repository.bower_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// bower attribute not returned by API
				ImportStateVerifyIgnore: []string{"bower"},
			},
		},
	})
}

func testAccRepositoryBowerHosted(name string, rewritePackageURLs bool) string {
	return fmt.Sprintf(`
resource "nexus_repository" "bower_hosted" {
	name   = "%s"
	format = "bower"
	type   = "hosted"

	bower {
		rewrite_package_urls = %v
	}

	storage {

	}
}`, name, rewritePackageURLs)
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.docker_hosted", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.docker_hosted", "format", "docker"),
					resource.TestCheckResourceAttr("nexus_repository.docker_hosted", "type", "hosted"),
				),
			},
			{
				ResourceName:      "nexus_repository.docker_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
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
		http_port        = %d
		https_port       = %d
		force_basic_auth = true
		v1enabled        = true
	}

    docker_proxy {
		index_url  = "https://index.docker.io"
		index_type = "HUB"
	}

    http_client {
		authentication {
			type = "username"
		}
	}

	negative_cache {

	}

	proxy {
		remote_url = "https://registry.npmjs.org"
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name, httpPort, httpsPort)
}
