package demo

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
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
	data.Set("name", tempMap["name"])
	data.SetId(tempMap["id"].(string))

	return nil
}

