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
					// 这种服务端设置回来的值如果不确定的话，我们可以用 TestCheckResourceAttrSet 简单断言目标key 是否存在即可。
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
