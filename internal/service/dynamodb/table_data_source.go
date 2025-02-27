package dynamodb

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_dynamodb_table")
func DataSourceTable() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceTableRead,
		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attribute": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrType: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m[names.AttrName].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			"billing_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"global_secondary_index": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"non_key_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"projection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"read_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"write_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"hash_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_secondary_index": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"non_key_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"projection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m[names.AttrName].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Required: true,
			},
			"point_in_time_recovery": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrEnabled: {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"range_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrKMSKeyARN: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"server_side_encryption": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrEnabled: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						names.AttrKMSKeyARN: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"stream_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"stream_label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_view_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags: tftags.TagsSchemaComputed(),
			"ttl": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						names.AttrEnabled: {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"write_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DynamoDBConn()
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Get(names.AttrName).(string)

	result, err := conn.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Dynamodb Table (%s): %s", name, err)
	}

	if result == nil || result.Table == nil {
		return sdkdiag.AppendErrorf(diags, "reading Dynamodb Table (%s): not found", name)
	}

	table := result.Table

	d.SetId(aws.StringValue(table.TableName))

	d.Set(names.AttrARN, table.TableArn)
	d.Set(names.AttrName, table.TableName)
	d.Set("deletion_protection_enabled", table.DeletionProtectionEnabled)

	if table.BillingModeSummary != nil {
		d.Set("billing_mode", table.BillingModeSummary.BillingMode)
	} else {
		d.Set("billing_mode", dynamodb.BillingModeProvisioned)
	}

	if table.ProvisionedThroughput != nil {
		d.Set("write_capacity", table.ProvisionedThroughput.WriteCapacityUnits)
		d.Set("read_capacity", table.ProvisionedThroughput.ReadCapacityUnits)
	}

	if err := d.Set("attribute", flattenTableAttributeDefinitions(table.AttributeDefinitions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting attribute: %s", err)
	}

	for _, attribute := range table.KeySchema {
		if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeHash {
			d.Set("hash_key", attribute.AttributeName)
		}

		if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeRange {
			d.Set("range_key", attribute.AttributeName)
		}
	}

	if err := d.Set("local_secondary_index", flattenTableLocalSecondaryIndex(table.LocalSecondaryIndexes)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting local_secondary_index: %s", err)
	}

	if err := d.Set("global_secondary_index", flattenTableGlobalSecondaryIndex(table.GlobalSecondaryIndexes)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting global_secondary_index: %s", err)
	}

	if table.StreamSpecification != nil {
		d.Set("stream_view_type", table.StreamSpecification.StreamViewType)
		d.Set("stream_enabled", table.StreamSpecification.StreamEnabled)
	} else {
		d.Set("stream_view_type", "")
		d.Set("stream_enabled", false)
	}

	d.Set("stream_arn", table.LatestStreamArn)
	d.Set("stream_label", table.LatestStreamLabel)

	if err := d.Set("server_side_encryption", flattenTableServerSideEncryption(table.SSEDescription)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting server_side_encryption: %s", err)
	}

	if err := d.Set("replica", flattenReplicaDescriptions(table.Replicas)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting replica: %s", err)
	}

	if table.TableClassSummary != nil {
		d.Set("table_class", table.TableClassSummary.TableClass)
	} else {
		d.Set("table_class", dynamodb.TableClassStandard)
	}

	pitrOut, err := conn.DescribeContinuousBackupsWithContext(ctx, &dynamodb.DescribeContinuousBackupsInput{
		TableName: aws.String(d.Id()),
	})
	// When a Table is `ARCHIVED`, DescribeContinuousBackups returns `TableNotFoundException`
	if err != nil && !tfawserr.ErrCodeEquals(err, "UnknownOperationException", dynamodb.ErrCodeTableNotFoundException) {
		return sdkdiag.AppendErrorf(diags, "describing DynamoDB Table (%s) Continuous Backups: %s", d.Id(), err)
	}

	if err := d.Set("point_in_time_recovery", flattenPITR(pitrOut)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting point_in_time_recovery: %s", err)
	}

	ttlOut, err := conn.DescribeTimeToLiveWithContext(ctx, &dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "describing DynamoDB Table (%s) Time to Live: %s", d.Id(), err)
	}

	if err := d.Set("ttl", flattenTTL(ttlOut)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting ttl: %s", err)
	}

	tags, err := ListTags(ctx, conn, d.Get(names.AttrARN).(string))
	// When a Table is `ARCHIVED`, ListTags returns `ResourceNotFoundException`
	if err != nil && !(tfawserr.ErrMessageContains(err, "UnknownOperationException", "Tagging is not currently supported in DynamoDB Local.") || tfresource.NotFound(err)) {
		return sdkdiag.AppendErrorf(diags, "listing tags for DynamoDB Table (%s): %s", d.Get(names.AttrARN).(string), err)
	}

	if err := d.Set(names.AttrTags, tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}
