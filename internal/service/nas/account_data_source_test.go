package nas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/internal/acctest"
)

func TestAccAWSBillingServiceAccount_basic(t *testing.T) {
	dataSourceName := "data.aws_billing_service_account.main"

	billingAccountID := "386209384616"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsBillingServiceAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", billingAccountID),
					acctest.CheckResourceAttrGlobalARNAccountID(dataSourceName, "arn", billingAccountID, "iam", "root"),
				),
			},
		},
	})
}

const testAccCheckAwsBillingServiceAccountConfig = `
data "aws_billing_service_account" "main" {}
`