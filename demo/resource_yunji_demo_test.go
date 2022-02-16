package demo

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider

func TestAccDemo_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPackResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDemoConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("yunji_demo.test", "instance_name", "yunji"),
					resource.TestCheckResourceAttr("gateway_pack.test", "disk_size", "100"),
				),
			},
		},
	})
}

var testDemoConfig = `
resource "yunji_demo" "test" {
	instance_name = "yunji"
	disk_size = 100
}

`
func testAccCheckPackResourceDestroy(s *terraform.State) error {
	return nil
}
