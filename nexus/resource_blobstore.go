package nexus

import (
	"fmt"
	"log"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBlobstore() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlobstoreCreate,
		Read:   resourceBlobstoreRead,
		Update: resourceBlobstoreUpdate,
		Delete: resourceBlobstoreDelete,
		Exists: resourceBlobstoreExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"S3", "File"}, false),
			},

			"name": {
				Description: "Blobstore name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				ConflictsWith: []string{"bucket_configuration"},
				Description:   "The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory",
				Type:          schema.TypeString,
				Optional:      true,
			},
			"available_space_in_bytes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blob_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bucket_configuration": {
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_bucket_connection": {
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
							Optional: true,
							Type:     schema.TypeList,
						},
						"encryption": {
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Default:     0,
							Description: "The limit in MB.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						"type": {
							Description:  "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"spaceRemainingQuota", "spaceUsedQuota"}, false),
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"total_size_in_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getBlobstoreFromResourceData(d *schema.ResourceData) nexus.Blobstore {
	bs := nexus.Blobstore{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
	}

	switch bs.Type {
	case nexus.BlobstoreTypeFile:
		bs.Path = d.Get("path").(string)
	case nexus.BlobstoreTypeS3:
		if _, ok := d.GetOk("bucket_configuration"); ok {
			bucketConfigurationList := d.Get("bucket_configuration").([]interface{})
			bucketConfiguration := bucketConfigurationList[0].(map[string]interface{})

			bs.BlobstoreS3BucketConfiguration = &nexus.BlobstoreS3BucketConfiguration{}

			if _, ok := bucketConfiguration["advanced_bucket_connection"]; ok {
				advancedBucketConfigurationList := bucketConfiguration["advanced_bucket_connection"].([]interface{})
				if len(advancedBucketConfigurationList) > 0 {
					advancedBucketConfiguration := advancedBucketConfigurationList[0].(map[string]interface{})

					bs.BlobstoreS3BucketConfiguration.BlobstoreS3AdvancedBucketConnection = &nexus.BlobstoreS3AdvancedBucketConnection{
						Endpoint:       advancedBucketConfiguration["endpoint"].(string),
						SignerType:     advancedBucketConfiguration["signer_type"].(string),
						ForcePathStyle: advancedBucketConfiguration["force_path_style"].(bool),
					}
				}
			}

			if _, ok := bucketConfiguration["bucket"]; ok {
				bucketList := bucketConfiguration["bucket"].([]interface{})
				bucket := bucketList[0].(map[string]interface{})

				bs.BlobstoreS3BucketConfiguration.BlobstoreS3Bucket = &nexus.BlobstoreS3Bucket{
					Expiration: bucket["expiration"].(int),
					Name:       bucket["name"].(string),
					Prefix:     bucket["prefix"].(string),
					Region:     bucket["region"].(string),
				}
			}

			if _, ok := bucketConfiguration["bucket_security"]; ok {
				bucketSecurityList := bucketConfiguration["bucket_security"].([]interface{})
				if len(bucketSecurityList) > 0 && bucketSecurityList[0] != nil {
					bucketSecurity := bucketSecurityList[0].(map[string]interface{})

					bs.BlobstoreS3BucketConfiguration.BlobstoreS3BucketSecurity = &nexus.BlobstoreS3BucketSecurity{
						AccessKeyID:     bucketSecurity["access_key_id"].(string),
						Role:            bucketSecurity["role"].(string),
						SecretAccessKey: bucketSecurity["secret_access_key"].(string),
						SessionToken:    bucketSecurity["session_token"].(string),
					}
				}
			}

			if _, ok := bucketConfiguration["encryption"]; ok {
				encryptionList := bucketConfiguration["encryption"].([]interface{})
				if len(encryptionList) > 0 {
					encryption := encryptionList[0].(map[string]interface{})

					bs.BlobstoreS3BucketConfiguration.BlobstoreS3Encryption = &nexus.BlobstoreS3Encryption{
						Key:  encryption["encryption_key"].(string),
						Type: encryption["encryption_type"].(string),
					}
				}
			}
		}
	}

	if _, ok := d.GetOk("soft_quota"); ok {
		softQuotaList := d.Get("soft_quota").([]interface{})
		softQuotaConfig := softQuotaList[0].(map[string]interface{})

		bs.BlobstoreSoftQuota = &nexus.BlobstoreSoftQuota{
			Limit: softQuotaConfig["limit"].(int),
			Type:  softQuotaConfig["type"].(string),
		}
	}

	return bs
}

func resourceBlobstoreCreate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs := getBlobstoreFromResourceData(d)

	if err := nexusClient.BlobstoreCreate(bs); err != nil {
		return err
	}

	d.SetId(bs.Name)
	d.Set("name", bs.Name)

	return resourceBlobstoreRead(d, m)
}

func resourceBlobstoreRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs, err := nexusClient.BlobstoreRead(d.Id())
	log.Print(bs)
	if err != nil {
		return err
	}

	if bs == nil {
		d.SetId("")
		return nil
	}

	d.Set("available_space_in_bytes", bs.AvailableSpaceInBytes)
	d.Set("blob_count", bs.BlobCount)
	d.Set("name", bs.Name)
	d.Set("path", bs.Path)
	d.Set("total_size_in_bytes", bs.TotalSizeInBytes)
	d.Set("type", bs.Type)

	if bs.BlobstoreS3BucketConfiguration != nil {
		if err := d.Set("bucket_configuration", flattenBlobstoreBucketConfiguration(bs.BlobstoreS3BucketConfiguration, d)); err != nil {
			return fmt.Errorf("Error reading bucket configuration: %s", err)
		}
	}

	if bs.BlobstoreSoftQuota != nil {
		if err := d.Set("soft_quota", flattenBlobstoreSoftQuota(bs.BlobstoreSoftQuota)); err != nil {
			return fmt.Errorf("Error reading soft quota: %s", err)
		}
	}

	return nil
}

func resourceBlobstoreUpdate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	bs := getBlobstoreFromResourceData(d)
	if err := nexusClient.BlobstoreUpdate(d.Id(), bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if err := nexusClient.BlobstoreDelete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceBlobstoreExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	bs, err := nexusClient.BlobstoreRead(d.Id())
	return bs != nil, err
}

func flattenBlobstoreSoftQuota(softQuota *nexus.BlobstoreSoftQuota) []map[string]interface{} {
	if softQuota == nil {
		return nil
	}
	data := map[string]interface{}{
		"limit": softQuota.Limit,
		"type":  softQuota.Type,
	}
	return []map[string]interface{}{data}
}

func flattenBlobstoreBucketConfiguration(bucketConfig *nexus.BlobstoreS3BucketConfiguration, d *schema.ResourceData) []map[string]interface{} {
	if bucketConfig == nil {
		return nil
	}
	data := map[string]interface{}{
		"advanced_bucket_connection": flattenBlobstoreS3AdvancedBucketConnection(bucketConfig.BlobstoreS3AdvancedBucketConnection),
		"bucket":                     flattenBlobstoreS3Bucket(bucketConfig.BlobstoreS3Bucket),
		"bucket_security":            flattenBlobstoreS3BucketSecurity(bucketConfig.BlobstoreS3BucketSecurity, d),
		"encryption":                 flattenBlobstoreS3Encryption(bucketConfig.BlobstoreS3Encryption),
	}
	return []map[string]interface{}{data}
}

func flattenBlobstoreS3AdvancedBucketConnection(bucketConnection *nexus.BlobstoreS3AdvancedBucketConnection) []map[string]interface{} {
	if bucketConnection == nil {
		return nil
	}
	data := map[string]interface{}{
		"endpoint":         bucketConnection.Endpoint,
		"force_path_style": bucketConnection.ForcePathStyle,
		"signer_type":      bucketConnection.SignerType,
	}
	return []map[string]interface{}{data}
}

func flattenBlobstoreS3Bucket(bucekt *nexus.BlobstoreS3Bucket) []map[string]interface{} {
	if bucekt == nil {
		return nil
	}
	data := map[string]interface{}{
		"expiration": bucekt.Expiration,
		"name":       bucekt.Name,
		"prefix":     bucekt.Prefix,
		"region":     bucekt.Region,
	}
	return []map[string]interface{}{data}
}

func flattenBlobstoreS3BucketSecurity(bucketSecurity *nexus.BlobstoreS3BucketSecurity, d *schema.ResourceData) []map[string]interface{} {
	if bucketSecurity == nil {
		return nil
	}
	data := map[string]interface{}{
		"access_key_id":     bucketSecurity.AccessKeyID,
		"role":              bucketSecurity.Role,
		"secret_access_key": d.Get("bucket_configuration.0.bucket_security.0.secret_access_key"), // secret_access_key",
		"session_token":     bucketSecurity.SessionToken,
	}
	return []map[string]interface{}{data}
}

func flattenBlobstoreS3Encryption(encryption *nexus.BlobstoreS3Encryption) []map[string]interface{} {
	if encryption == nil {
		return nil
	}
	data := map[string]interface{}{
		"encryption_key":  encryption.Key,
		"encryption_type": encryption.Type,
	}
	return []map[string]interface{}{data}
}
