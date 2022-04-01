package demo

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "硬盘大小",
			},
		},
	}
}

func resourceDemoCreate(data *schema.ResourceData, meta interface{}) error {
	// 获取服务地址
	conf := meta.(*Configuration)
	endpoint := conf.endpoint

	// 构造HTTP 请求，调用API接口
	client := &http.Client{}
	postData := url.Values{}
	postData.Add("instance_name", data.Get("instance_name").(string))
	postData.Add("disk_size", data.Get("disk_size").(string))
	request, err := http.NewRequest("POST", endpoint, strings.NewReader(postData.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	_, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return readErr
	}
	// code here to error handle
	defer response.Body.Close()
	// TODO
	data.SetId("weiyi_demo_id")
	return resourceDemoRead(data, meta)
}

func resourceDemoRead(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	endpoint := conf.endpoint
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("%sget?id=%s",endpoint, data.Id()), nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return readErr
	}
	defer response.Body.Close()
	var tempMap map[string]interface{}
	json.Unmarshal(body, &tempMap)
	data.Set("instance_name", tempMap["instance_name"])
	data.Set("disk_size", tempMap["disk_size"])
	return nil
}

func resourceDemoUpdate(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	endpoint := conf.endpoint
	postData := url.Values{}
	if data.HasChange("instance_name") {
		postData.Add("instance_name", data.Get("instance_name").(string))
	}
	if data.HasChange("disk_size") {
		postData.Add("disk_size", data.Get("disk_size").(string))
	}
	client := &http.Client{}
	request, err := http.NewRequest("PUT", fmt.Sprintf("%supdate?id=%s", endpoint, data.Id()),
		strings.NewReader(postData.Encode()))
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return resourceDemoRead(data, meta)
}

func resourceDemoDelete(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	endpoint := conf.endpoint
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%sdelete?id=%s",endpoint, data.Id()), nil)
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}