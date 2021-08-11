package directconnect_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/internal/acctest"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
)

func TestAccAWSDxConnectionAssociation_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, directconnect.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsDxConnectionAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxConnectionAssociationConfig(sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsDxConnectionAssociationExists("aws_dx_connection_association.test"),
				),
			},
		},
	})
}

func TestAccAWSDxConnectionAssociation_multiConns(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, directconnect.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAwsDxConnectionAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxConnectionAssociationConfig_multiConns(sdkacctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsDxConnectionAssociationExists("aws_dx_connection_association.test1"),
					testAccCheckAwsDxConnectionAssociationExists("aws_dx_connection_association.test2"),
				),
			},
		},
	})
}

func testAccCheckAwsDxConnectionAssociationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*client.AWSClient).DirectConnectConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_dx_connection_association" {
			continue
		}

		input := &directconnect.DescribeConnectionsInput{
			ConnectionId: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeConnections(input)
		if err != nil {
			return err
		}
		for _, v := range resp.Connections {
			if *v.ConnectionId == rs.Primary.ID && v.LagId != nil {
				return fmt.Errorf("Dx Connection (%s) is not dissociated with Lag", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckAwsDxConnectionAssociationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccDxConnectionAssociationConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_dx_connection" "test" {
  name      = "tf-dx-%s"
  bandwidth = "1Gbps"
  location  = "EqSe2-EQ"
}

resource "aws_dx_lag" "test" {
  name                  = "tf-dx-%s"
  connections_bandwidth = "1Gbps"
  location              = "EqSe2-EQ"
  force_destroy         = true
}

resource "aws_dx_connection_association" "test" {
  connection_id = aws_dx_connection.test.id
  lag_id        = aws_dx_lag.test.id
}
`, rName, rName)
}

func testAccDxConnectionAssociationConfig_multiConns(rName string) string {
	return fmt.Sprintf(`
resource "aws_dx_connection" "test1" {
  name      = "tf-dxconn1-%s"
  bandwidth = "1Gbps"
  location  = "EqSe2-EQ"
}

resource "aws_dx_connection" "test2" {
  name      = "tf-dxconn2-%s"
  bandwidth = "1Gbps"
  location  = "EqSe2-EQ"
}

resource "aws_dx_lag" "test" {
  name                  = "tf-dx-%s"
  connections_bandwidth = "1Gbps"
  location              = "EqSe2-EQ"
  force_destroy         = true
}

resource "aws_dx_connection_association" "test1" {
  connection_id = aws_dx_connection.test1.id
  lag_id        = aws_dx_lag.test.id
}

resource "aws_dx_connection_association" "test2" {
  connection_id = aws_dx_connection.test2.id
  lag_id        = aws_dx_lag.test.id
}
`, rName, rName, rName)
}