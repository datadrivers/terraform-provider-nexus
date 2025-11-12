package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecuritySSLTrustStore() *schema.Resource {
	return &schema.Resource{
		Description: `Use this data source to retrieve ALL certificates in the Nexus truststore.`,

		Read: dataSourceSecuritySslTrustStoreRead,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Unique identifier of the certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"fingerprint": {
							Description: "Fingerprint of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"serial_number": {
							Description: "Serial number of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issuer_common_name": {
							Description: "Common name of the issuer of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issuer_organization": {
							Description: "Organization of the issuer of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issuer_organization_unit": {
							Description: "Organization unit of the issuer of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subject_common_name": {
							Description: "Common name of the subject of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subject_organization": {
							Description: "Organization of the subject of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subject_organization_unit": {
							Description: "Organization unit of the subject of the retrieved certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"pem": {
							Description: "PEM encoded certificate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issued_on": {
							Description: "Timestamp for when the certificate was issued",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"expires_on": {
							Description: "Timestamp for when the certificate expires",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSecuritySslTrustStoreRead(d *schema.ResourceData, m interface{}) error {
	d.SetId("global")
	client := m.(*nexus.NexusClient)
	certificates, err := client.Security.SSL.ListCertificates()
	if err != nil {
		return err
	}
	var certList []interface{}

	for _, cert := range *certificates {
		certMap := make(map[string]interface{})
		certMap["id"] = cert.Id
		certMap["fingerprint"] = cert.Fingerprint
		certMap["serial_number"] = cert.SerialNumber
		certMap["issuer_common_name"] = cert.IssuerCommonName
		certMap["issuer_organization"] = cert.IssuerOrganization
		certMap["issuer_organization_unit"] = cert.IssuerOrganizationUnit
		certMap["subject_common_name"] = cert.SubjectCommonName
		certMap["subject_organization"] = cert.SubjectOrganization
		certMap["subject_organization_unit"] = cert.SubjectOrganizationUnit
		certMap["pem"] = cert.Pem
		certMap["issued_on"] = cert.IssuedOn
		certMap["expires_on"] = cert.ExpiresOn

		certList = append(certList, certMap)
	}

	if err := d.Set("certificates", certList); err != nil {
		return err
	}

	return nil
}
