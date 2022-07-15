package demo

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
	"time"
)

func dataSourceYunjiAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYunjiAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceYunjiAccountRead(data *schema.ResourceData, meta interface{}) error {
	conf := meta.(*Configuration)
	endpoint := conf.endpoint
	client := &http.Client{}
	name := data.Get("name").(string)
	request, err := http.NewRequest("GET", fmt.Sprintf("%sdata_source?name=%s",endpoint, name), nil)
	if err != nil {
		return WrapError(err)
	}
	// 所有的API 都要进行错误重试
	wait := IncrementalWait(3*time.Second, 5*time.Second)
	response := &http.Response{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError{
		response, err := client.Do(request)
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

	response, err = client.Do(request)
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
	data.Set("name", tempMap["name"])
	// 资源import 的时候也需要资源id去导入，所以datasource 资源也必须设置资源唯一id
	data.SetId(tempMap["id"].(string))

	return nil
}

