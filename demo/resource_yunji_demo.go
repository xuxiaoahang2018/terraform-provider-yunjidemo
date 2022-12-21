package demo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func resourceDemo() *schema.Resource {
	return &schema.Resource{
		Create: resourceDemoCreate,
		Read:   resourceDemoRead,
		Update: resourceDemoUpdate,
		Delete: resourceDemoDelete,
		// importer 该方法可以让用户通过 SetId 的资源Id 对资源进行导入，例如 terraform import example_thing.foo abc123
		// 参考 https://www.terraform.io/plugin/sdkv2/resources/import
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// 首先terraform 的变量名只支持小驼峰命名，并且只支持小写
			"instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew 参数，意味着该参数不可改变，当改变该参数，会把整个资源会删除重建，注意所有的ForceNew 参数
				// 一定要在文档中给予标注，否则不予以merge 代码及发布。请重视！
				ForceNew:    true,
				Description: "服务器名称",
			},
			"disk_size": &schema.Schema{
				Type: schema.TypeInt,
				// Required/Optional 该参数是否必填，在文档中也必须标注！
				Required:    true,
				Description: "硬盘大小",
			},
			"networks": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "网络",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// map 类型不要在嵌套list, map, set等复杂类型，里面key 只能使用string基础类型
			// 参考 https://www.terraform.io/plugin/sdkv2/schemas/schema-types#typemap
			"memory": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "内存",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory_size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"memory_init": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// 当我们传入的是一个json 类型，但是这个json的内容是可变的不确定的，我们就可以用string
			// 字段来接收这个可变的json.
			"config_json": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  "自定义配置",
			},
			"uuid": {
				Type: schema.TypeString,
				// computed 参数定义一般是需要服务端计算返回的不确定值，例如uuid, 时间戳等等，
				// 这种参数在使用main.tf 编写hcl 时候无需传入，但是我们要在查询方法里面，根据服务端查询获得的值
				// 将返回值Set 回tfstate 种。
				Computed:    true,
				Description: "自定义配置",
			},
			"set_demo": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"set_test": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDemoCreate(d *schema.ResourceData, meta interface{}) error {
	// 获取服务地址
	conf := meta.(*Configuration)
	// 构造HTTP 请求，调用API接口
	client := &http.Client{}
	yunji := YunjiDemo{}
	yunji.InstanceName = d.Get("instance_name").(string)
	yunji.DiskSize = d.Get("disk_size").(int)
	networks := d.Get("networks").([]interface{})
	for _, v := range networks {
		n := Network{}
		value := v.(map[string]interface{})
		n.Port = value["port"].(int)
		n.Protocol = value["protocol"].(string)
		yunji.Networks = append(yunji.Networks, n)
	}
	// 这种List, Map, Set 类型定义的参数，我们通过d.GetOk 去获得
	temp, ok := d.GetOk("memory")
	if ok {
		memory := temp.(map[string]interface{})
		yunji.Memory.MemorySize = memory["memory_size"].(string)
		yunji.Memory.MemoryUnit = memory["memory_unit"].(string)
	}

	// 像这种传递进来的不确定的json 类型的，我们使用string 字符串类型接收以后，可以通过反序列拿到实际的值
	config := d.Get("config_json").(string)
	if config != "" {
		var config_json map[string]interface{}
		if err := json.Unmarshal([]byte(config), &config_json); err != nil {
			return WrapError(err)
		}
		yunji.ConfigJson = config
	}
	if _, ok := d.GetOk("set_demo"); ok {
		// Set 参数解析
		setDemo := d.Get("set_demo").(*schema.Set).List()
		for _, v := range setDemo {
			result := v.(map[string]interface{})
			n := SetDemo{}
			n.SetTest = result["set_test"].(string)
			yunji.SetDemo = append(yunji.SetDemo, n)
		}
	}

	// 将请求结构体进行json 格式序列化, 并通过http.request库发送http 请求。
	req, err := json.Marshal(yunji)
	if err != nil {
		return WrapError(err)
	}
	request, err := http.NewRequest("POST", conf.endpoint+"create", bytes.NewBuffer(req))
	if err != nil {
		return WrapError(err)
	}
	request.Header.Set("Content-Type", "application/json")
	wait := IncrementalWait(3*time.Second, 5*time.Second)
	response := &http.Response{}
	// 所有的 API 请求都要使用 设置重试的时间和场景。 包括下面的READ, UPDATE, DELETE
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Println("[YUNJI] 开始创建了")
		fmt.Println("------------")
		fmt.Println(request.Body)
		fmt.Println("------------")
		response, err = client.Do(request)
		if err != nil {
			if NeedRetry(response.StatusCode, err) {
				wait()
				//RetryableError 是创建可从给定错误重试的 RetryError 的助手。
				return resource.RetryableError(err)
			}
			//NonRetryableError 是创建 RetryError 的助手，它_not_可以从给定的错误中重试。
			return resource.NonRetryableError(err)
		}
		return nil
	})

	_, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return WrapError(readErr)
	}
	defer response.Body.Close()
	// 这里要注意， SetId 的目的是为了给资源一个唯一id， 下面的read, update和delete 都需要
	// 通过这里Set的唯一id 去查询， 所以创建资源这里必须SetId,

	d.SetId("weiyi_demo_id")

	return resourceDemoRead(d, meta)
}

// TODO 强调所有在Read 方法中Set 的值在文档中也应要写到对应的地方。
func resourceDemoRead(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	client := &http.Client{}
	log.Println("[YUNJI] 创建成功开始等待")

	// 这里我们就通过创建时候Set 的Set的资源唯一id去对资源进行查询。
	request, err := http.NewRequest("GET", fmt.Sprintf("%sget?id=%s", conf.endpoint, data.Id()), nil)
	if err != nil {
		return WrapError(err)
	}
	fmt.Println("#1", request.Method)
	//fmt.Println("#2", request.GetBody)
	fmt.Println("#3", request.URL)
	fmt.Println("#4", request.Body)

	response, err := client.Do(request)
	if err != nil {
		return WrapError(err)
	}
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return WrapError(readErr)
	}
	defer response.Body.Close()
	var tempMap map[string]interface{}
	json.Unmarshal(body, &tempMap)

	// 这里我们根据服务端的返回去把服务端返回值重新Set回去。
	// 我们要求把所有我们在schema.Resource 里面定义的参数全部设置回去，举例我们定义了 A,B,C 参数，通过查询资源发现服务端返回了资源中参数A,B 的值，那么我们这里就要通过 schema.ResourceData.Set(A, value)
	// schema.ResourceData.Set(B, value) 等将参数重新Set 回去，所有Set 回去的值在文档中output也要体现出来。

	// 这里解释一下为什么我们要把这些值都Set 回去，假设我们定义一个资源名称参数 name 输入值为aa , 但是实际创建资源后, 资源名称为 bb ，这明显服务端所创资源我们预期不符合
	// 如果不Set 回来的化，tfstate 文件中默认参数值就是输入值，是感知不到这个服务端错误的，所以我们应该尽可能的把我们的定义的参数字段，重新Set 回去。
	data.Set("instance_name", tempMap["instance_name"])
	data.Set("disk_size", tempMap["disk_size"])
	data.Set("memory", tempMap["memory"])
	data.Set("config_json", tempMap["config_json"])

	// 这里要强调一点，所有schema.ResourceData.Set 的参数，必须是我们在入参schema.Resource.Schema 中定义的参数，否则是无法Set成功的。
	data.Set("networks", tempMap["networks"])

	// 像这种定义为compute 的值，就是根据服务端计算返回的值，用户无需传入。
	data.Set("uuid", tempMap["uuid"])
	data.Set("set_demo", tempMap["set_demo"])

	return nil
}

func resourceDemoUpdate(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	endpoint := conf.endpoint
	yunji := YunjiDemo{}
	update := false
	// 在更新这里，我们通过HasChange 函数判断用户的hcl 代码有无变更，如果有更新，捕获更新并去调用update api 接口
	if data.HasChange("config_json") {
		update = true
		config := data.Get("config_json").(string)
		if config != "" {
			var config_json map[string]interface{}
			if err := json.Unmarshal([]byte(config), &config_json); err != nil {
				return WrapError(err)
			}
			yunji.ConfigJson = config
		}
	}
	if data.HasChange("disk_size") {
		update = true
		yunji.DiskSize = data.Get("disk_size").(int)
	}
	if data.HasChange("networks") {
		update = true
		networks := data.Get("networks").([]interface{})
		for _, v := range networks {
			n := Network{}
			value := v.(map[string]interface{})
			n.Port = value["port"].(int)
			n.Protocol = value["protocol"].(string)
			yunji.Networks = append(yunji.Networks, n)
		}
	}
	if data.HasChange("set_demo") {
		if _, ok := data.GetOk("set_demo"); ok {
			// Set 参数解析
			setDemo := data.Get("set_demo").(*schema.Set).List()
			for _, v := range setDemo {
				result := v.(map[string]interface{})
				n := SetDemo{}
				n.SetTest = result["set_test"].(string)
				yunji.SetDemo = append(yunji.SetDemo, n)
			}
		}
	}
	if data.HasChange("memory") {
		update = true
		temp, ok := data.GetOk("memory")
		if ok {
			memory := temp.(map[string]interface{})
			yunji.Memory.MemorySize = memory["memory_size"].(string)
			yunji.Memory.MemoryUnit = memory["memory_unit"].(string)
		}
	}
	if update {
		client := &http.Client{}
		req, err := json.Marshal(yunji)
		if err != nil {
			return WrapError(err)
		}
		// 这里我们依旧根据 Create 方法里面设置的资源的唯一ID 去更新资源配置。
		request, err := http.NewRequest("PUT", fmt.Sprintf("%supdate?id=%s", endpoint, data.Id()), bytes.NewBuffer(req))
		if err != nil {
			return WrapError(err)
		}
		_, err = client.Do(request)
		if err != nil {
			return WrapError(err)
		}
	}
	// 在更新完成依旧调用Read 方法，把资源的最新资源重新写回tfstate
	return resourceDemoRead(data, meta)
}

func resourceDemoDelete(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	endpoint := conf.endpoint
	uuid := data.Get("uuid")
	fmt.Println("uuid: ", uuid)
	client := &http.Client{}
	// 这里我们依旧根据 Create 方法里面设置的资源的唯一ID 去删除资源。
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%sdelete?id=%s", endpoint, data.Id()), nil)
	if err != nil {
		return WrapError(err)
	}
	_, err = client.Do(request)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

type YunjiDemo struct {
	InstanceName string    `json:"instance_name"`
	DiskSize     int       `json:"disk_size"`
	Networks     []Network `json:"networks"`
	Memory       struct {
		MemorySize string `json:"memory_size"`
		MemoryUnit string `json:"memory_unit"`
	} `json:"memory"`
	ConfigJson string    `json:"config_json"`
	SetDemo    []SetDemo `json:"set_demo"`
}

type Network struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type SetDemo struct {
	SetTest string `json:"set_test"`
}
