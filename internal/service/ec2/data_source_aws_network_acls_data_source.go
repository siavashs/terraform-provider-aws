package ec2

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
	"github.com/terraform-providers/terraform-provider-aws/internal/keyvaluetags"
	"github.com/terraform-providers/terraform-provider-aws/internal/tags"
)

func DataSourceNetworkACLs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsNetworkAclsRead,
		Schema: map[string]*schema.Schema{
			"filter": ec2CustomFiltersSchema(),

			"tags": tags.TagsSchemaComputed(),

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceAwsNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).EC2Conn

	req := &ec2.DescribeNetworkAclsInput{}

	if v, ok := d.GetOk("vpc_id"); ok {
		req.Filters = buildEC2AttributeFilterList(
			map[string]string{
				"vpc-id": v.(string),
			},
		)
	}

	filters, filtersOk := d.GetOk("filter")
	tags, tagsOk := d.GetOk("tags")

	if tagsOk {
		req.Filters = append(req.Filters, buildEC2TagFilterList(
			keyvaluetags.New(tags.(map[string]interface{})).Ec2Tags(),
		)...)
	}

	if filtersOk {
		req.Filters = append(req.Filters, buildEC2CustomFilterList(
			filters.(*schema.Set),
		)...)
	}

	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		req.Filters = nil
	}

	log.Printf("[DEBUG] DescribeNetworkAcls %s\n", req)
	resp, err := conn.DescribeNetworkAcls(req)
	if err != nil {
		return err
	}

	if resp == nil || len(resp.NetworkAcls) == 0 {
		return errors.New("no matching network ACLs found")
	}

	networkAcls := make([]string, 0)

	for _, networkAcl := range resp.NetworkAcls {
		networkAcls = append(networkAcls, aws.StringValue(networkAcl.NetworkAclId))
	}

	d.SetId(meta.(*client.AWSClient).Region)

	if err := d.Set("ids", networkAcls); err != nil {
		return fmt.Errorf("Error setting network ACL ids: %w", err)
	}

	return nil
}