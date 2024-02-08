package security

import (
	"fmt"
	"strings"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecuritySSLTruststore() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to add an SSL certificate to the nexus Truststore",

		Create: resourceSecuritySSLTruststoreCertCreate,
		Read:   resourceSecuritySSLTruststoreCertRead,
		Update: nil, // changing the certificate should re-create it
		Delete: resourceSecuritySSLTruststoreCertDelete,
		Exists: resourceSecuritySSLTruststoreCertExists,

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"pem": {
				Description: "The cert in PEM format",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			"fingerprint": {
				Description: "The fingerprint of the cert",
				Required:    false,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func getSecuritySSLTruststoreCertFromResourceData(d *schema.ResourceData) security.SSLCertificate {
	return security.SSLCertificate{
		Pem:         d.Get("pem").(string),
		Fingerprint: d.Get("fingerprint").(string),
	}
}

func resourceSecuritySSLTruststoreCertCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	cert := getSecuritySSLTruststoreCertFromResourceData(d)

	if err := client.Security.SSL.AddCertificate(&cert); err != nil {
		return err
	}

	err := resourceSecuritySSLTruststoreCertRead(d, m)
	if err != nil {
		return fmt.Errorf("failed retrieving freshly created certificate '%s'", cert.Pem)
	}
	return err
}

func resourceSecuritySSLTruststoreCertRead(d *schema.ResourceData, m interface{}) error {
	cert, err := securitySSLTruststoreFindCert(d.Get("pem").(string), m)
	if err != nil {
		return err
	}

	if cert == nil {
		return nil
	}

	d.Set("fingerprint", cert.Fingerprint)
	d.SetId(cert.Id)

	return nil
}

func securitySSLTruststoreFindCert(pem string, m interface{}) (*security.SSLCertificate, error) {
	client := m.(*nexus.NexusClient)

	certs, err := client.Security.SSL.ListCertificates()
	if err != nil {
		return nil, err
	}

	if certs == nil || len(*certs) < 1 {
		return nil, nil
	}

	for _, cert := range *certs {
		if strings.ReplaceAll(cert.Pem, "\n", "") == strings.ReplaceAll(pem, "\n", "") {
			return &cert, nil
		}
	}

	return nil, nil
}

func resourceSecuritySSLTruststoreCertDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.SSL.RemoveCertificate(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecuritySSLTruststoreCertExists(d *schema.ResourceData, m interface{}) (bool, error) {
	cert, err := securitySSLTruststoreFindCert(d.Get("pem").(string), m)
	return cert != nil, err
}
