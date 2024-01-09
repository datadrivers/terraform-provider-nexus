package security_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/tools"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourcesecuritySSLTruststore(t *testing.T) {
	resName := "data.nexus_security_ssl_truststore.acceptance"

	certX509, _ := tools.TestRetrieveCert("https://google.com")
	certFingerPrint, _ := tools.TestGetCertificateFingerprint(certX509)
	certPEM, _ := tools.TestPemEncode(certX509.Raw, "CERTIFICATE")

	cert := security.SSLCertificate{
		Id:                      certFingerPrint,
		Fingerprint:             certFingerPrint,
		SerialNumber:            certX509.SerialNumber.String(),
		IssuerCommonName:        certX509.Issuer.CommonName,
		IssuerOrganization:      strings.Join(certX509.Issuer.Organization, "\n"),
		IssuerOrganizationUnit:  strings.Join(certX509.Issuer.OrganizationalUnit, "\n"),
		SubjectCommonName:       certX509.Subject.CommonName,
		SubjectOrganization:     strings.Join(certX509.Subject.Organization, "\n"),
		SubjectOrganizationUnit: strings.Join(certX509.Subject.OrganizationalUnit, "\n"),
		Pem:                     certPEM,
		IssuedOn:                certX509.NotBefore.UnixMilli(),
		ExpiresOn:               certX509.NotAfter.UnixMilli(),
	}

	certReq := security.CertificateRequest{
		Host: "google.com",
		Port: 443,
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
				Config: testAccResourceSecuritySSLTruststoreConfig(cert) + testAccDatasourceSecuritySSLTruststoreConfig(certReq),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", cert.Id),
					resource.TestCheckResourceAttr(resName, "fingerprint", cert.Fingerprint),
					resource.TestCheckResourceAttr(resName, "serial_number", cert.SerialNumber),
					resource.TestCheckResourceAttr(resName, "issuer_common_name", cert.IssuerCommonName),
					resource.TestCheckResourceAttr(resName, "issuer_organization", cert.IssuerOrganization),
					resource.TestCheckResourceAttr(resName, "issuer_organization_unit", cert.IssuerOrganizationUnit),
					resource.TestCheckResourceAttr(resName, "subject_common_name", cert.SubjectCommonName),
					resource.TestCheckResourceAttr(resName, "subject_organization", cert.SubjectOrganization),
					resource.TestCheckResourceAttr(resName, "subject_organization_unit", cert.SubjectOrganizationUnit),
					resource.TestCheckResourceAttr(resName, "pem", cert.Pem),
					resource.TestCheckResourceAttr(resName, "issued_on", strconv.FormatInt(cert.IssuedOn, 10)),
					resource.TestCheckResourceAttr(resName, "expires_on", strconv.FormatInt(cert.ExpiresOn, 10)),
				),
			},
		},
	})
}

func testAccDatasourceSecuritySSLTruststoreConfig(certReq security.CertificateRequest) string {
	return fmt.Sprintf(`
data "nexus_security_ssl_truststore" "acceptance" {
  host = "%s"
  port = %s
}
`, certReq.Host, strconv.Itoa(certReq.Port))
}
