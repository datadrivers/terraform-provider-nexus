package blobstore

import (
	"fmt"
	"log"

	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/blobstore"
	blobstoreSchema "github.com/dre2004/terraform-provider-nexus/internal/schema/blobstore"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBlobstoreAzure() *schema.Resource {
	return &schema.Resource{
		Description: `~> PRO Feature

Use this resource to create a Nexus Azure blobstore.`,

		Create: resourceBlobstoreAzureCreate,
		Read:   resourceBlobstoreAzureRead,
		Update: resourceBlobstoreAzureUpdate,
		Delete: resourceBlobstoreAzureDelete,
		Exists: resourceBlobstoreAzureExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id":                  common.ResourceID,
			"name":                blobstoreSchema.ResourceName,
			"blob_count":          blobstoreSchema.ResourceBlobCount,
			"soft_quota":          blobstoreSchema.ResourceSoftQuota,
			"total_size_in_bytes": blobstoreSchema.ResourceTotalSizeInBytes,
			"bucket_configuration": {
				Description: "The Azure specific configuration details for the Azure object that'll contain the blob store",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Description: "Account name found under Access keys for the storage account",
							Required:    true,
							Type:        schema.TypeString,
						},
						"authentication": {
							Description: "The Azure specific authentication details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_method": {
										Description:  "The type of Azure authentication to use. Possible values: `ACCOUNTKEY` and `MANAGEDIDENTITY`",
										Required:     true,
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{string(blobstore.AzureAuthenticationMethodAccountKey), string(blobstore.AzureAuthenticationMethodManagedIdentity)}, false),
									},
									"account_key": {
										Description: "The account key. Required if `authentication_method` is `ACCOUNTKEY`",
										Optional:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MaxItems: 1,
							Required: true,
							Type:     schema.TypeList,
						},
						"container_name": {
							Description: "The name of an existing container to be used for storage",
							Required:    true,
							Type:        schema.TypeString,
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

func getBlobstoreAzureFromResourceData(d *schema.ResourceData) blobstore.Azure {
	bucketConfiguration := d.Get("bucket_configuration").([]interface{})[0].(map[string]interface{})
	authenticationConfig := bucketConfiguration["authentication"].([]interface{})[0].(map[string]interface{})

	bs := blobstore.Azure{
		Name: d.Get("name").(string),
		BucketConfiguration: blobstore.AzureBucketConfiguration{
			AccountName: bucketConfiguration["account_name"].(string),
			Authentication: blobstore.AzureBucketConfigurationAuthentication{
				AuthenticationMethod: blobstore.AzureAuthenticationMethod(authenticationConfig["authentication_method"].(string)),
			},
			ContainerName: bucketConfiguration["container_name"].(string),
		},
	}

	if accountKey, ok := authenticationConfig["account_key"]; ok {
		bs.BucketConfiguration.Authentication.AccountKey = accountKey.(string)
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

func resourceBlobstoreAzureCreate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreAzureFromResourceData(resourceData)

	if err := nexusClient.BlobStore.Azure.Create(&bs); err != nil {
		return err
	}

	resourceData.SetId(bs.Name)
	resourceData.Set("name", bs.Name)

	return resourceBlobstoreAzureRead(resourceData, m)
}

func resourceBlobstoreAzureRead(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.Azure.Get(resourceData.Id())
	log.Printf("[DEBUG] BlobStore:\n%+v\n", bs)
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
	if err := resourceData.Set("blob_count", genericBlobstoreInformation.BlobCount); err != nil {
		return err
	}
	if err := resourceData.Set("total_size_in_bytes", genericBlobstoreInformation.TotalSizeInBytes); err != nil {
		return err
	}
	if err := resourceData.Set("bucket_configuration", flattenAzureBucketConfiguration(&bs.BucketConfiguration, resourceData)); err != nil {
		return fmt.Errorf("error reading bucket configuration: %s", err)
	}

	if bs.SoftQuota != nil {
		if err := resourceData.Set("soft_quota", flattenSoftQuota(bs.SoftQuota)); err != nil {
			return fmt.Errorf("error reading soft quota: %s", err)
		}
	}

	return nil
}

func resourceBlobstoreAzureUpdate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreAzureFromResourceData(resourceData)
	if err := nexusClient.BlobStore.Azure.Update(resourceData.Id(), &bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreAzureDelete(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	if err := nexusClient.BlobStore.Azure.Delete(resourceData.Id()); err != nil {
		return err
	}

	resourceData.SetId("")

	return nil
}

func resourceBlobstoreAzureExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.Azure.Get(resourceData.Id())
	return bs != nil, err
}
