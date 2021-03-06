package datadog

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/s33dunda/go-datadog-api"
	"strings"
)

func resourceDatadogLogsPipelineOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatadogLogsPipelineOrderCreate,
		Update: resourceDatadogLogsPipelineOrderUpdate,
		Read:   resourceDatadogLogsPipelineOrderRead,
		Delete: resourceDatadogLogsPipelineOrderDelete,
		Exists: resourceDatadogLogsPipelineOrderExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
			"pipelines": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDatadogLogsPipelineOrderCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceDatadogLogsPipelineOrderUpdate(d, meta)
}

func resourceDatadogLogsPipelineOrderRead(d *schema.ResourceData, meta interface{}) error {
	ddList, err := meta.(*datadog.Client).GetLogsPipelineList()
	if err != nil {
		return err
	}

	if err = d.Set("pipelines", ddList.PipelineIds); err != nil {
		return err
	}

	return nil
}

func resourceDatadogLogsPipelineOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	var ddPipelineList datadog.LogsPipelineList
	tfList := d.Get("pipelines").([]interface{})
	ddList := make([]string, len(tfList))
	for i, id := range tfList {
		ddList[i] = id.(string)
	}
	ddPipelineList.PipelineIds = ddList
	var tfId string
	if name, exists := d.GetOk("name"); exists {
		tfId = name.(string)
	}
	if _, err := meta.(*datadog.Client).UpdateLogsPipelineList(&ddPipelineList); err != nil {
		// Cannot map pipelines to existing ones
		if strings.Contains(err.Error(), "422 Unprocessable Entity") {
			ddPipelineOrder, getErr := meta.(*datadog.Client).GetLogsPipelineList()
			if getErr != nil {
				return fmt.Errorf("error updating logs pipeline list: (%s)", err.Error())
			}
			return fmt.Errorf("cannot map pipelines to existing ones\n existing pipelines: %s\n pipeline to be updated: %s",
				ddPipelineOrder.PipelineIds,
				ddList)
		}
		return fmt.Errorf("error updating logs pipeline list: (%s)", err.Error())
	}
	d.SetId(tfId)
	return resourceDatadogLogsPipelineOrderRead(d, meta)
}

// The deletion of pipeline order is not supported from config API.
// This function simply delete the pipeline order resource from terraform state.
func resourceDatadogLogsPipelineOrderDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceDatadogLogsPipelineOrderExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}
