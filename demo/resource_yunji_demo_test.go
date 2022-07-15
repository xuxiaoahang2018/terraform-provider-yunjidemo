package demo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// 所有的测试都至少要两个Steps ，一个测试创建，一个测试更新。
// 而且所有的值必须要覆盖到，不论创建还是更新。
// 如果单测都无法100% 覆盖，那么何谈高质量交付。这到了现场测试100% 必出问题。所以测试覆盖率应该至少达到90% 以上。
// provider 代码会直接应用于生产，所以单元测试是重中之中。
// 当程序代码写完了时候，这个项目只算完成了20%， 剩下的80%应该是去编写测试去增加程序健壮性。
func TestAccYunjiDemo_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPackResourceDestroy,
		Steps: []resource.TestStep{
			{
				// 这里config 里面配置的就是我们实际使用provider时写入的hcl 代码
				Config: testYunjiDemoPackConfig,
				// 这个测试就是断言 当我们的hcl 配置代码执行 terraform apply 后，写入tfstate 的值和预期的是否相符合
				// 所以我们所有在 schema.ResourceData.Set 到tfstate 的参数都要全部去测试，不能遗漏一个！
				Check: resource.ComposeTestCheckFunc(
					testEndpointExist("yunjidemo_demo.test"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "instance_name", "yunji"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "disk_size", "100"),
					// list, set, 类型的值测试断言可以参考官方文档
					// https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/teststep
					// .# 断言list 里面有多少value
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.#", "2"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.0.port", "81"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.0.protocol", "http"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.1.port", "82"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.1.protocol", "https"),
					// .% 断言map 中有多少key
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "memory.%", "2"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "memory.memory_size", "1024"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "memory.memory_unit", "test"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "config_json", "  {\"yunji\":\"test\"}\n"),
					// .# 断言set 里面有多少value
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "set_demo.#", "2"),
				),
			},
			{
				// 这里config 里面配置的就是我们对上一步测试的hcl 代码进行更新，去测试我们Resource Update 更新的接口
				Config: testYunjiDemoUpdate,
				Check: resource.ComposeTestCheckFunc(
					testEndpointExist("yunjidemo_demo.test"),
					// 这里的测试断言，是执行apply 进行更新后获得的是不是预期更新后的值。
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "instance_name", "yunji"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "disk_size", "200"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.#", "2"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.0.port", "8888"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.0.protocol", "https"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.1.port", "9999"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "networks.1.protocol", "http"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "memory.%", "2"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "memory.memory_size", "2048"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "memory.memory_unit", "test222"),
					resource.TestCheckResourceAttr("yunjidemo_demo.test", "config_json", "  {\"yunji\":\"test\"}\n"),
				),
			},
			// 测试资源是否可以正常import 导入
			{
				ResourceName:      "yunjidemo_demo.test",
				ImportState:       true,
				ImportStateVerify: true,
			},


		},
	})
}

var testYunjiDemoPackConfig = `
resource "yunjidemo_demo" "test" {
  instance_name  = "yunji"
  disk_size = 100
  # List 或者Set 类型我们可以通过下述方式去写
  networks {
	port = 81
    protocol = "http"
  }
  networks {
	port = 82
    protocol = "https"
  }
  memory = {
    memory_size = "1024"
	memory_unit = "test"
  }
  set_demo {
	set_test = "yunji"
  }
  set_demo {
	set_test = "yunji22"
  }

  # 复杂不确定的json 我们可以以EOF 这种文本格式输入
  # 参考 https://www.terraform.io/language/expressions/strings
  config_json =<<EOF
  {"yunji":"test"}
EOF
}
`


var testYunjiDemoUpdate = `
resource "yunjidemo_demo" "test" {
  instance_name  = "yunji"
  disk_size = 200
  networks {
	port = 8888
    protocol = "https"
  }
  networks {
	port = 9999
    protocol = "http"
  }
  memory = {
    memory_size = "2048"
	memory_unit = "test222"
  }
  set_demo {
	set_test = "yunji"
  }
  set_demo {
	set_test = "yunji22"
  }
  # 复杂不确定的json 我们可以以EOF 这种文本格式输入
  # 参考 https://www.terraform.io/language/expressions/strings
  config_json =<<EOF
  {"yunji":"test"}
EOF
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
