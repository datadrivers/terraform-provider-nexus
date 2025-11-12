package security_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/tools"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcesecuritySSLTruststore(t *testing.T) {
	resName := "nexus_security_ssl_truststore.acceptance"

	certPEM, _, certFingerprint, _, _ := tools.TestGenerateRandomCertificate()
	cert := security.SSLCertificate{
		Id:          certFingerprint,
		Pem:         certPEM,
		Fingerprint: certFingerprint,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecuritySSLTruststoreConfig(cert),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resName, "id", cert.Id),
					resource.TestCheckResourceAttr(resName, "pem", cert.Pem+"\n"),
					resource.TestCheckResourceAttr(resName, "fingerprint", cert.Fingerprint),
				),
			},
		},
	})
}

func testAccResourceSecuritySSLTruststoreConfig(cert security.SSLCertificate) string {
	return fmt.Sprintf(`
resource "nexus_security_ssl_truststore" "acceptance" {
  pem = <<EOT
%s
EOT
}
`, cert.Pem)
}
