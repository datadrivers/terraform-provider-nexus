package deprecated

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceBlobstore() *schema.Resource {
	return &schema.Resource{
		Description: "Use this to create a Nexus blobstore.",

		Read: dataSourceBlobstoreRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Description: "The type of the blobstore",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "Blobstore name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"path": {
				Description: "The path to the blobstore contents",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"available_space_in_bytes": {
				Computed:    true,
				Description: "Available space in Bytes",
				Type:        schema.TypeInt,
			},
			"blob_count": {
				Computed:    true,
				Description: "Count of blobs",
				Type:        schema.TypeInt,
			},
			"bucket_configuration": {
				Description: "The S3 bucket configuration. Needed for blobstore type 'S3'",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_bucket_connection": {
							Description: "Additional connection configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Description: "A custom endpoint URL for third party object stores using the S3 API.",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"force_path_style": {
										Default:     false,
										Description: "Setting this flag will result in path-style access being used for all requests.",
										Optional:    true,
										Type:        schema.TypeBool,
									},
									"signer_type": {
										Description: "An API signature version which may be required for third party object stores using the S3 API.",
										Optional:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
						"bucket": {
							Description: "The S3 bucket configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Description: "The AWS region to create a new S3 bucket in or an existing S3 bucket's region",
										Required:    true,
										Type:        schema.TypeString,
									},
									"name": {
										Description: "The name of the S3 bucket",
										Required:    true,
										Type:        schema.TypeString,
									},
									"prefix": {
										Description: "The S3 blob store (i.e S3 object) key prefix",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"expiration": {
										Default:     0,
										Description: "How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)",
										Optional:    true,
										Type:        schema.TypeInt,
									},
								},
							},
							MaxItems: 1,
							Required: true,
							Type:     schema.TypeList,
						},
						"bucket_security": {
							Description: "Additional security configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key_id": {
										Description: "An IAM access key ID for granting access to the S3 bucket",
										Type:        schema.TypeString,
										Required:    true,
									},
									"secret_access_key": {
										Description: "The secret access key associated with the specified IAM access key ID",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"role": {
										Description: "An IAM role to assume in order to access the S3 bucket",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"session_token": {
										Description: "An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket",
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
									},
								},
							},
							MaxItems: 1,
							Required: true,
							Type:     schema.TypeList,
						},
						"encryption": {
							Description: "Additional bucket encryption configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"encryption_key": {
										Description: "The encryption key.",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"encryption_type": {
										Default:      "s3ManagedEncryption",
										Description:  "The type of S3 server side encryption to use.",
										Optional:     true,
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"s3ManagedEncryption", "kmsManagedEncryption"}, false),
									},
								},
							},
							MaxItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"soft_quota": {
				Description: "Soft quota of the blobstore",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Description: "The limit in Bytes. Minimum value is 1000000",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"type": {
							Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"total_size_in_bytes": {
				Computed:    true,
				Description: "The total size of the blobstore in Bytes",
				Type:        schema.TypeInt,
			},
		},
	}
}

func dataSourceBlobstoreRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))

	return resourceBlobstoreRead(d, m)
}
