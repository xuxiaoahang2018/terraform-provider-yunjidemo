package demo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccYunjiDemo_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPackResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testYunjiDemoPackConfig,
				Check: resource.ComposeTestCheckFunc(
					testEndpointExist("yunjidemo_demo.test"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "instance_name", "yunji"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "disk_size", "100"),
				),
			},
		},
	})
}

var testYunjiDemoPackConfig = `
resource "yunjidemo_demo" "test" {
  instance_name  = "yunji"
  disk_size = 100
}
`

func testEndpointExist(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}
		return nil
	}
}

func testAccCheckPackResourceDestroy(s *terraform.State) error {
	return nil
}
