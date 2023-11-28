package blobstore

import (
	"github.com/dre2004/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenSoftQuota(softQuota *blobstore.SoftQuota) []map[string]interface{} {
	if softQuota == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"limit": softQuota.Limit,
			"type":  softQuota.Type,
		},
	}
}

func flattenAzureBucketConfiguration(bucketConfig *blobstore.AzureBucketConfiguration, resourceData *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"account_name":   bucketConfig.AccountName,
			"authentication": flattenBucketConfigurationAuthentication(&bucketConfig.Authentication, resourceData),
			"container_name": bucketConfig.ContainerName,
		},
	}
}

func flattenBucketConfigurationAuthentication(authenticationConfig *blobstore.AzureBucketConfigurationAuthentication, resourceData *schema.ResourceData) []map[string]interface{} {
	data := map[string]interface{}{
		"authentication_method": string(authenticationConfig.AuthenticationMethod),
	}
	if accountKey, ok := resourceData.GetOk("bucket_configuration.0.authentication.0.account_key"); ok {
		data["account_key"] = accountKey
	}

	return []map[string]interface{}{data}
}

func flattenS3BucketConfiguration(bucketConfig *blobstore.S3BucketConfiguration, resourceData *schema.ResourceData) []map[string]interface{} {
	if bucketConfig == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"advanced_bucket_connection": flattenAdvancedBucketConnection(bucketConfig.AdvancedBucketConnection),
			"bucket":                     flattenBucket(bucketConfig.Bucket),
			"bucket_security":            flattenBucketSecurity(bucketConfig.BucketSecurity, resourceData),
			"encryption":                 flattenEncryption(bucketConfig.Encryption),
		},
	}
}

func flattenAdvancedBucketConnection(bucketConnection *blobstore.S3AdvancedBucketConnection) []map[string]interface{} {
	if bucketConnection == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"endpoint":         bucketConnection.Endpoint,
			"force_path_style": bucketConnection.ForcePathStyle,
			"signer_type":      bucketConnection.SignerType,
		},
	}
}

func flattenBucket(bucket blobstore.S3Bucket) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"expiration": bucket.Expiration,
			"name":       bucket.Name,
			"prefix":     bucket.Prefix,
			"region":     bucket.Region,
		},
	}
}

func flattenBucketSecurity(bucketSecurity *blobstore.S3BucketSecurity, resourceData *schema.ResourceData) []map[string]interface{} {
	if bucketSecurity == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"access_key_id":     bucketSecurity.AccessKeyID,
			"role":              bucketSecurity.Role,
			"secret_access_key": resourceData.Get("bucket_configuration.0.bucket_security.0.secret_access_key"), // secret_access_key",
			"session_token":     bucketSecurity.SessionToken,
		},
	}
}

func flattenEncryption(encryption *blobstore.S3Encryption) []map[string]interface{} {
	if encryption == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"encryption_key":  encryption.Key,
			"encryption_type": encryption.Type,
		},
	}
}
