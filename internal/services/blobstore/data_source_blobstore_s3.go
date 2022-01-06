package blobstore

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlobstoreS3() *schema.Resource {
	return &schema.Resource{
		Description: "Use this to create a Nexus blobstore.",

		Read: dataSourceBlobstoreS3Read,
		Schema: map[string]*schema.Schema{
			"id": common.DataSourceID,
			"name": {
				Description: "Blobstore name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"blob_count": {
				Computed:    true,
				Description: "Count of blobs",
				Type:        schema.TypeInt,
			},
			"soft_quota": getDataSourceSoftQuotaSchema(),
			"total_size_in_bytes": {
				Computed:    true,
				Description: "The total size of the blobstore in Bytes",
				Type:        schema.TypeInt,
			},
			"bucket_configuration": {
				Description: "The S3 bucket configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_bucket_connection": {
							Description: "Additional connection configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Description: "A custom endpoint URL for third party object stores using the S3 API.",
										Computed:    true,
										Type:        schema.TypeString,
									},
									"force_path_style": {
										Description: "Setting this flag will result in path-style access being used for all requests.",
										Computed:    true,
										Type:        schema.TypeBool,
									},
									"signer_type": {
										Description: "An API signature version which may be required for third party object stores using the S3 API.",
										Computed:    true,
										Type:        schema.TypeString,
									},
								},
							},
							Computed: true,
							Type:     schema.TypeList,
						},
						"bucket": {
							Description: "The S3 bucket configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Description: "The AWS region to create a new S3 bucket in or an existing S3 bucket's region",
										Computed:    true,
										Type:        schema.TypeString,
									},
									"name": {
										Description: "The name of the S3 bucket",
										Computed:    true,
										Type:        schema.TypeString,
									},
									"prefix": {
										Description: "The S3 blob store (i.e S3 object) key prefix",
										Computed:    true,
										Type:        schema.TypeString,
									},
									"expiration": {
										Description: "How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)",
										Computed:    true,
										Type:        schema.TypeInt,
									},
								},
							},
							Computed: true,
							Type:     schema.TypeList,
						},
						"bucket_security": {
							Description: "Additional security configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_id": {
										Description: "An IAM access key ID for granting access to the S3 bucket",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"secret_access_key": {
										Description: "The secret access key associated with the specified IAM access key ID",
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
									"role": {
										Description: "An IAM role to assume in order to access the S3 bucket",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"session_token": {
										Description: "An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket",
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
								},
							},
							Computed: true,
							Type:     schema.TypeList,
						},
						"encryption": {
							Description: "Additional bucket encryption configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"encryption_key": {
										Description: "The encryption key.",
										Computed:    true,
										Type:        schema.TypeString,
									},
									"encryption_type": {
										Description: "The type of S3 server side encryption to use.",
										Computed:    true,
										Type:        schema.TypeString,
									},
								},
							},
							Computed: true,
							Type:     schema.TypeList,
						},
					},
				},
				Computed: true,
				Type:     schema.TypeList,
			},
		},
	}
}

func dataSourceBlobstoreS3Read(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceBlobstoreS3Read(resourceData, m)
}
