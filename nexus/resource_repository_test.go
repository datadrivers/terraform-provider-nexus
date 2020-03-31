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
        write_policy = "ALLOW"
	}
}
`, name, aptDistribution, aptSigningKEypair, aptSigningPassphrase)
}

func TestAccRepositoryBowerHosted(t *testing.T) {
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
		write_policy = "ALLOW"
	}
}`, name, rewritePackageURLs)
}

func TestAccRepositoryDockerHostedWithPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-hosted-%s", acctest.RandString(10))
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
		write_policy = "ALLOW"
	}
}`, name, strconv.FormatBool(online), httpPort, httpsPort)
}
func TestAccRepositoryDockerHostedWithoutPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-hosted-%s", acctest.RandString(10))
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
		write_policy = "ALLOW"
	}
}`, name, strconv.FormatBool(online))
}

func TestAccRepositoryDockerProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-proxy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryDockerProxy(repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.docker_proxy", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.docker_proxy", "format", "docker"),
					resource.TestCheckResourceAttr("nexus_repository.docker_proxy", "type", "proxy"),
				),
			},
		},
	})
}

func testAccRepositoryDockerProxy(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_proxy" {
	name   = "%s"
	type   = "proxy"
	format = "docker"

	docker {
		force_basic_auth = true
		v1enabled        = false
	}

    docker_proxy {
		index_type = "HUB"
		index_url  = "http://www.example.com"
	}

    http_client {
		authentication {
			type = "username"
		}
	}

	negative_cache {
		enabled = true
		ttl     = 1440
	}

	proxy {
		remote_url  = "https://index.docker.io"
	}

	storage {
		blob_store_name = "default"
		write_policy    = "ALLOW"
	}
}`, name)
}

func TestAccRepositoryDockerGroup(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-group-%s", acctest.RandString(10))
	memberRepoName := fmt.Sprintf("test-repo-docker-group-member-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryDockerProxy(memberRepoName) + testAccRepositoryDockerGroup(repoName, memberRepoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.docker_group", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.docker_group", "format", "docker"),
					resource.TestCheckResourceAttr("nexus_repository.docker_group", "type", "group"),
				),
			},
		},
	})
}

func testAccRepositoryDockerGroup(name string, memberRepoName string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_group" {
	name   = "%s"
	format = "docker"
	type   = "group"
	online = true
	
	group {
		member_names = [nexus_repository.docker_proxy.name]
	}
	
	docker {
		force_basic_auth = true
		http_port        = 8082
		https_port       = 0
		v1enabled        = false
	}
	
	storage {
		blob_store_name                = "default"
		strict_content_type_validation = true
	}
}`, name)
}
