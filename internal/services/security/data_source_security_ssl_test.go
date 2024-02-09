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

func TestAccDatasourceSecuritySSL(t *testing.T) {
	dataSourceName := "data.nexus_security_ssl.acceptance"

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
				Config: testAccDatasourceSecuritySSLConfig(certReq),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", cert.Id),
					resource.TestCheckResourceAttr(dataSourceName, "fingerprint", cert.Fingerprint),
					resource.TestCheckResourceAttr(dataSourceName, "serial_number", cert.SerialNumber),
					resource.TestCheckResourceAttr(dataSourceName, "issuer_common_name", cert.IssuerCommonName),
					resource.TestCheckResourceAttr(dataSourceName, "issuer_organization", cert.IssuerOrganization),
					resource.TestCheckResourceAttr(dataSourceName, "issuer_organization_unit", cert.IssuerOrganizationUnit),
					resource.TestCheckResourceAttr(dataSourceName, "subject_common_name", cert.SubjectCommonName),
					resource.TestCheckResourceAttr(dataSourceName, "subject_organization", cert.SubjectOrganization),
					resource.TestCheckResourceAttr(dataSourceName, "subject_organization_unit", cert.SubjectOrganizationUnit),
					resource.TestCheckResourceAttr(dataSourceName, "pem", cert.Pem),
					resource.TestCheckResourceAttr(dataSourceName, "issued_on", strconv.FormatInt(cert.IssuedOn, 10)),
					resource.TestCheckResourceAttr(dataSourceName, "expires_on", strconv.FormatInt(cert.ExpiresOn, 10)),
				),
			},
		},
	})
}

func testAccDatasourceSecuritySSLConfig(certReq security.CertificateRequest) string {
	return fmt.Sprintf(`
data "nexus_security_ssl" "acceptance" {
  host = "%s"
  port = %s
}
`, certReq.Host, strconv.Itoa(certReq.Port))
}
