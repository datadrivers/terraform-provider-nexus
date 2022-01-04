package blobstore

import (
	"fmt"
	"log"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/tools"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBlobstoreS3() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus S3 blobstore.",

		Create: resourceBlobstoreS3Create,
		Read:   resourceBlobstoreS3Read,
		Update: resourceBlobstoreS3Update,
		Delete: resourceBlobstoreS3Delete,
		Exists: resourceBlobstoreS3Exists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Blobstore name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"available_space_in_bytes": {
				Description: "Available space in Bytes",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"blob_count": {
				Description: "Count of blobs",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"soft_quota": getResourceSoftQuotaSchema(),
			"total_size_in_bytes": {
				Description: "The total size of the blobstore in Bytes",
				Type:        schema.TypeInt,
				Computed:    true,
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
									"max_connection_pool_size": {
										Description: "Setting this value will override the default connection pool size of Nexus of the s3 client for this blobstore.",
										Optional:    true,
										Type:        schema.TypeInt,
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
										Description: "How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)",
										Required:    true,
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
										Optional:    true,
									},
									"secret_access_key": {
										Description: "The secret access key associated with the specified IAM access key ID",
										Type:        schema.TypeString,
										Optional:    true,
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
							Description: "Additional bucket encryption configurations",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"encryption_key": {
										Description: "The encryption key.",
										Optional:    true,
										Type:        schema.TypeString,
									},
									"encryption_type": {
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
				Required: true,
				Type:     schema.TypeList,
			},
		},
	}
}

func getBlobstoreS3FromResourceData(d *schema.ResourceData) blobstore.S3 {
	bucketConfigurationList := d.Get("bucket_configuration").([]interface{})
	bucketConfiguration := bucketConfigurationList[0].(map[string]interface{})
	bucketList := bucketConfiguration["bucket"].([]interface{})
	bucket := bucketList[0].(map[string]interface{})

	bs := blobstore.S3{
		Name: d.Get("name").(string),
		BucketConfiguration: blobstore.S3BucketConfiguration{
			Bucket: blobstore.S3Bucket{
				Expiration: int32(bucket["expiration"].(int)),
				Name:       bucket["name"].(string),
				Prefix:     bucket["prefix"].(string),
				Region:     bucket["region"].(string),
			},
		},
	}

	if _, ok := bucketConfiguration["advanced_bucket_connection"]; ok {
		advancedBucketConfigurationList := bucketConfiguration["advanced_bucket_connection"].([]interface{})
		if len(advancedBucketConfigurationList) > 0 {
			advancedBucketConfiguration := advancedBucketConfigurationList[0].(map[string]interface{})

			bs.BucketConfiguration.AdvancedBucketConnection = &blobstore.S3AdvancedBucketConnection{
				Endpoint:       advancedBucketConfiguration["endpoint"].(string),
				SignerType:     advancedBucketConfiguration["signer_type"].(string),
				ForcePathStyle: tools.GetBoolPointer(advancedBucketConfiguration["force_path_style"].(bool)),
			}
		}
	}

	if _, ok := bucketConfiguration["bucket_security"]; ok {
		bucketSecurityList := bucketConfiguration["bucket_security"].([]interface{})
		if len(bucketSecurityList) > 0 && bucketSecurityList[0] != nil {
			bucketSecurity := bucketSecurityList[0].(map[string]interface{})

			bs.BucketConfiguration.BucketSecurity = &blobstore.S3BucketSecurity{
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

			bs.BucketConfiguration.Encryption = &blobstore.S3Encryption{
				Key:  encryption["encryption_key"].(string),
				Type: encryption["encryption_type"].(string),
			}
		}
	}

	if _, ok := d.GetOk("soft_quota"); ok {
		softQuotaList := d.Get("soft_quota").([]interface{})
		softQuotaConfig := softQuotaList[0].(map[string]interface{})

		bs.SoftQuota = &blobstore.SoftQuota{
			Limit: int64(softQuotaConfig["limit"].(int)),
			Type:  softQuotaConfig["type"].(string),
		}
	}

	return bs
}

func resourceBlobstoreS3Create(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreS3FromResourceData(d)

	if err := nexusClient.BlobStore.S3.Create(&bs); err != nil {
		return err
	}

	d.SetId(bs.Name)
	d.Set("name", bs.Name)

	return resourceBlobstoreS3Read(d, m)
}

func resourceBlobstoreS3Read(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.S3.Get(resourceData.Id())
	log.Print(bs)
	if err != nil {
		return err
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := nexusClient.BlobStore.List()
	if err != nil {
		return err
	}
	for _, generic := range genericBlobstores {
		if generic.Name == bs.Name {
			genericBlobstoreInformation = generic
		}
	}

	if bs == nil {
		resourceData.SetId("")
		return nil
	}

	if err := resourceData.Set("name", bs.Name); err != nil {
		return err
	}
	if err := resourceData.Set("available_space_in_bytes", genericBlobstoreInformation.AvailableSpaceInBytes); err != nil {
		return err
	}
	if err := resourceData.Set("blob_count", genericBlobstoreInformation.BlobCount); err != nil {
		return err
	}
	if err := resourceData.Set("total_size_in_bytes", genericBlobstoreInformation.TotalSizeInBytes); err != nil {
		return err
	}
	if err := resourceData.Set("bucket_configuration", flattenBucketConfiguration(&bs.BucketConfiguration, resourceData)); err != nil {
		return fmt.Errorf("error reading bucket configuration: %s", err)
	}

	if bs.SoftQuota != nil {
		if err := resourceData.Set("soft_quota", flattenSoftQuota(bs.SoftQuota)); err != nil {
			return fmt.Errorf("error reading soft quota: %s", err)
		}
	}

	return nil
}

func resourceBlobstoreS3Update(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreS3FromResourceData(resourceData)
	if err := nexusClient.BlobStore.S3.Update(resourceData.Id(), &bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreS3Delete(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	if err := nexusClient.BlobStore.S3.Delete(resourceData.Id()); err != nil {
		return err
	}

	resourceData.SetId("")

	return nil
}

func resourceBlobstoreS3Exists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.S3.Get(resourceData.Id())
	return bs != nil, err
}

func flattenBucketConfiguration(bucketConfig *blobstore.S3BucketConfiguration, resourceData *schema.ResourceData) []map[string]interface{} {
	if bucketConfig == nil {
		return nil
	}
	data := map[string]interface{}{
		"advanced_bucket_connection": flattenAdvancedBucketConnection(bucketConfig.AdvancedBucketConnection),
		"bucket":                     flattenBucket(bucketConfig.Bucket),
		"bucket_security":            flattenBucketSecurity(bucketConfig.BucketSecurity, resourceData),
		"encryption":                 flattenEncryption(bucketConfig.Encryption),
	}
	return []map[string]interface{}{data}
}

func flattenAdvancedBucketConnection(bucketConnection *blobstore.S3AdvancedBucketConnection) []map[string]interface{} {
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

func flattenBucket(bucket blobstore.S3Bucket) []map[string]interface{} {
	data := map[string]interface{}{
		"expiration": bucket.Expiration,
		"name":       bucket.Name,
		"prefix":     bucket.Prefix,
		"region":     bucket.Region,
	}
	return []map[string]interface{}{data}
}

func flattenBucketSecurity(bucketSecurity *blobstore.S3BucketSecurity, resourceData *schema.ResourceData) []map[string]interface{} {
	if bucketSecurity == nil {
		return nil
	}
	data := map[string]interface{}{
		"access_key_id":     bucketSecurity.AccessKeyID,
		"role":              bucketSecurity.Role,
		"secret_access_key": resourceData.Get("bucket_configuration.0.bucket_security.0.secret_access_key"), // secret_access_key",
		"session_token":     bucketSecurity.SessionToken,
	}
	return []map[string]interface{}{data}
}

func flattenEncryption(encryption *blobstore.S3Encryption) []map[string]interface{} {
	if encryption == nil {
		return nil
	}
	data := map[string]interface{}{
		"encryption_key":  encryption.Key,
		"encryption_type": encryption.Type,
	}
	return []map[string]interface{}{data}
}
