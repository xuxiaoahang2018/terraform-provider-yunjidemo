package demo

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
)

func resourceDemo() *schema.Resource {
	return &schema.Resource{
		// functions for the various actions
		Create: resourceDemoCreate,
		Read:   resourceDemoRead,
		Update: resourceDemoUpdate,
		Delete: resourceDemoDelete,

		Schema: map[string]*schema.Schema{
			// resource arguments and their specifications go here
			"instance_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "服务器名称",
			},
			"disk_size": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "硬盘大小",
			},
		},
	}
}

func resourceDemoCreate(data *schema.ResourceData, meta interface{}) error {
	body := &map[string]interface{}{
		"instance_name":         data.Get("instance_name"),
		"disk_size":           data.Get("head"),
	}
	bodyEncode, err := json.Marshal(*body)
	if err != nil {
		log.Fatal("fail to marshal message body")
		return err
	}
	// 获取服务地址
	conf := meta.(*Configuration)
	endpoint := conf.endpoint

	// 构造HTTP 请求，调用API接口
	client := &http.Client{}
	var bodyBuffer *bytes.Buffer
	bodyBuffer = bytes.NewBuffer(bodyEncode)
	request, err := http.NewRequest("POST", endpoint, bodyBuffer)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	_, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	// code here to error handle
	defer response.Body.Close()
	data.SetId("weiyi_demo_id")
	return nil
}

func resourceDemoRead(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDemoUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDemoDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}