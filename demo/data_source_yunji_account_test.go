package demo


import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)


func TestAccYunjiAccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckYunjiAccountDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.yunjidemo_account.current", "id"),
					resource.TestCheckResourceAttrSet("data.yunjidemo_account.current", "name"),
				),
			},
		},
	})
}

const testAccCheckYunjiAccountDataSourceBasic = `
data "yunjidemo_account" "current" {
	name = "ecs"
}
`
