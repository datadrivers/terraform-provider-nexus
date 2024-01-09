package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecuritySSLTruststore() *schema.Resource {
	return &schema.Resource{
		Description: `Use this data source to retrieve a SSL certificate via Nexus.`,

		Read: dataSourceSecuritySslTrustStoreRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"host": {
				Description: "Hostname for looking up certificate",
				Type:        schema.TypeString,
				Required:    true,
			},
			"port": {
				Description: "Port for looking up certificate",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     "443",
			},
			"fingerprint": {
				Description: "fingerprint field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"serial_number": {
				Description: "serialNumber field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"issuer_common_name": {
				Description: "issuerCommonName field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"issuer_organization": {
				Description: "issuerOrganization field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"issuer_organization_unit": {
				Description: "issuerOrganizationUnit field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"subject_common_name": {
				Description: "subjectCommonName field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"subject_organization": {
				Description: "subjectOrganization field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"subject_organization_unit": {
				Description: "subjectOrganizationUnit field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"pem": {
				Description: "pem field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"issued_on": {
				Description: "issuedOn field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeInt,
			},
			"expires_on": {
				Description: "expiresOn field of the retrieved cert",
				Computed:    true,
				Type:        schema.TypeInt,
			},
		},
	}
}

func dataSourceSecuritySslTrustStoreRead(d *schema.ResourceData, m interface{}) error {
	certReq := security.CertificateRequest{
		Host: d.Get("host").(string),
		Port: d.Get("port").(int),
	}

	client := m.(*nexus.NexusClient)
	cert, err := client.Security.SSL.GetCertificate(&certReq)
	if err != nil {
		return err
	}

	//log.Printf("[DEBUG] Found cert:\n%+v\n", cert)

	d.SetId(cert.Id)
	d.Set("fingerprint", cert.Fingerprint)
	d.Set("serial_number", cert.SerialNumber)
	d.Set("issuer_common_name", cert.IssuerCommonName)
	d.Set("issuer_organization", cert.IssuerOrganization)
	d.Set("issuer_organization_unit", cert.IssuerOrganizationUnit)
	d.Set("subject_common_name", cert.SubjectCommonName)
	d.Set("subject_organization", cert.SubjectOrganization)
	d.Set("subject_organization_unit", cert.SubjectOrganizationUnit)
	d.Set("pem", cert.Pem)
	d.Set("issued_on", cert.IssuedOn)
	d.Set("expires_on", cert.ExpiresOn)

	return nil
}
