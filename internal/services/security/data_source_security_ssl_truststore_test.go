package security_test

import (
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/tools"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourcesecuritySSLTrustStore(t *testing.T) {
	resName := "data.nexus_security_ssl_truststore.acceptance"

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
				Check:  nil,
			},
			{
				Config: testAccDataSourceSecuritySSLTruststoreConfig(cert),
				Check: resource.ComposeTestCheckFunc(
					// Assuming the first certificate in the list is the one we're testing
					resource.TestCheckResourceAttr(resName, "certificates.0.id", cert.Id),
					resource.TestCheckResourceAttr(resName, "certificates.0.fingerprint", cert.Fingerprint),
					resource.TestCheckResourceAttr(resName, "certificates.0.pem", cert.Pem),
				),
			},
		},
	})
}

func testAccDataSourceSecuritySSLTruststoreConfig(cert security.SSLCertificate) string {
	return `data "nexus_security_ssl_truststore" "acceptance" {}`
}
