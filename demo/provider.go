package demo

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider init  block
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// provider arguments and their specifications go here
			"endpoint": {
				Type:        schema.TypeString,
				Optional:  true,
				Default: "http://127.0.0.1:8888/",
				DefaultFunc: schema.EnvDefaultFunc("GW_ENDPOINT", nil),
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "" {
						errors = append(errors, fmt.Errorf("Endpoint must not be an empty string"))
					}
					return
				},
			},
		},
		// map terraform dsl resources to functions
		ResourcesMap: map[string]*schema.Resource{
			// 这里命名格式是provider 名称+ 下划线 + 资源名称，其他命名方式provider将无法正确找到函数路径。
			"yunjidemo_demo": resourceDemo(),

		},
		DataSourcesMap: map[string]*schema.Resource{
			"yunjidemo_account": dataSourceYunjiAccount(),
		},
		// provider configuration function
		ConfigureFunc: configureProvider,
	}
}


// Configuration struct
type Configuration struct {
	endpoint string
}

// configure provider options
func configureProvider(data *schema.ResourceData) (interface{}, error) {
	// pass options from terraform DSL to the client
	endpoint := data.Get("endpoint").(string)
	// code to error handle
	return &Configuration{
		endpoint: endpoint,
	}, nil
}
