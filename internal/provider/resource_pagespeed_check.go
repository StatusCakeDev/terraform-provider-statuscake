package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
	intvalidation "github.com/StatusCakeDev/terraform-provider-statuscake/internal/provider/validation"
)

func resourceStatusCakePagespeedCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakePagespeedCheckCreate,
		ReadContext:   resourceStatusCakePagespeedCheckRead,
		UpdateContext: resourceStatusCakePagespeedCheckUpdate,
		DeleteContext: resourceStatusCakePagespeedCheckDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"alert_config": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Alert configuration block. Omitting this block disabled all alerts",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_bigger": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "An alert will be sent if the size of the page is larger than this value (kb).",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"alert_slower": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "An alert will be sent if the load time of the page exceeds this value (ms).",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"alert_smaller": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "An alert will be sent if the size of the page is smaller than this value (kb).",
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
			"check_interval": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of seconds between checks",
				ValidateFunc: intvalidation.Int32InSlice(statuscake.PagespeedTestCheckRateValues()),
			},
			"contact_groups": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of contact group IDs",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: intvalidation.StringIsNumerical,
				},
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assigned monitoring location on which checks will be run",
			},
			"monitored_resource": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Monitored resource configuration block. The describes server under test",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "URL or IP address of the website under test",
							ValidateFunc: validation.Any(validation.IsURLWithHTTPorHTTPS, validation.IsIPAddress),
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the check",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"paused": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the check should be run",
			},
			"region": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Region on which to run checks",
				ValidateFunc: validation.StringInSlice(statuscake.PagespeedTestRegionValues(), false),
			},
		},
	}
}

func resourceStatusCakePagespeedCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})

	config, err := expandPagespeedCheckAlertConfig(d.Get("alert_config"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("alert_config") {
		body = merge(body, config.(map[string]interface{}))
	}

	checkInterval, err := expandPagespeedCheckInterval(d.Get("check_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("check_interval") {
		body["check_rate"] = checkInterval
	}

	contactGroups, err := expandPagespeedCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	name, err := expandPagespeedCheckName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	monitoredResource, err := expandPagespeedCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	paused, err := expandPagespeedCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	region, err := expandPagespeedCheckRegion(d.Get("region"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("region") {
		body["region"] = region
	}

	log.Print("[DEBUG] Creating StatusCake pagespeed check")
	log.Printf("[DEBUG] Request body: %+v", body)

	res, err := client.CreatePagespeedTestWithData(ctx, body).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create pagespeed check: %w", err))
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakePagespeedCheckRead(ctx, d, meta)
}

func resourceStatusCakePagespeedCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	res, err := client.GetPagespeedTest(ctx, id).Execute()

	// If the resource it not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get pagespeed check with ID: %w", err))
	}

	if err := d.Set("alert_config", flattenPagespeedCheckAlertConfig(res.Data, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read alert config: %+v", err))
	}

	if err := d.Set("check_interval", flattenPagespeedCheckInterval(res.Data.CheckRate, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read check interval: %+v", err))
	}

	if err := d.Set("contact_groups", flattenPagespeedCheckContactGroups(res.Data.ContactGroups, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read contact groups: %+v", err))
	}

	if err := d.Set("location", flattenPagespeedCheckLocation(res.Data.Location, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read location: %+v", err))
	}

	if err := d.Set("monitored_resource", flattenPagespeedCheckMonitoredResource(res.Data, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read monitored resource: %+v", err))
	}

	if err := d.Set("name", flattenPagespeedCheckName(res.Data.Name, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read name: %+v", err))
	}

	if err := d.Set("paused", flattenPagespeedCheckPaused(res.Data.Paused, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read paused: %+v", err))
	}

	return nil
}

func resourceStatusCakePagespeedCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})
	id := d.Id()

	config, err := expandPagespeedCheckAlertConfig(d.Get("alert_config"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("alert_config") {
		body = merge(body, config.(map[string]interface{}))
	}

	checkInterval, err := expandPagespeedCheckInterval(d.Get("check_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("check_interval") {
		body["check_rate"] = checkInterval
	}

	contactGroups, err := expandPagespeedCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	name, err := expandPagespeedCheckName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	monitoredResource, err := expandPagespeedCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	paused, err := expandPagespeedCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	region, err := expandPagespeedCheckRegion(d.Get("region"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("region") {
		body["region"] = region
	}

	log.Printf("[DEBUG] Updating StatusCake pagespeed check with ID: %s", id)
	log.Printf("[DEBUG] Request body: %+v", body)

	if err := client.UpdatePagespeedTestWithData(ctx, id, body).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update pagespeed check: %w", err))
	}

	return resourceStatusCakePagespeedCheckRead(ctx, d, meta)
}

func resourceStatusCakePagespeedCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake pagespeed check with ID: %s", id)

	if err := client.DeletePagespeedTest(ctx, id).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete pagespeed check with id %s: %w", id, err))
	}

	return nil
}

func expandPagespeedCheckAddress(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenPagespeedCheckAddress(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckAlertBigger(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenPagespeedCheckAlertBigger(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckAlertConfig(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	bigger, err := expandPagespeedCheckAlertBigger(original["alert_bigger"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.alert_bigger") {
		transformed["alert_bigger"] = bigger
	}

	slower, err := expandPagespeedCheckAlertSlower(original["alert_slower"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.alert_slower") {
		transformed["alert_slower"] = slower
	}

	smaller, err := expandPagespeedCheckAlertSmaller(original["alert_smaller"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.alert_smaller") {
		transformed["alert_smaller"] = smaller
	}

	return transformed, nil
}

func flattenPagespeedCheckAlertConfig(v interface{}, d *schema.ResourceData) interface{} {
	original := v.(statuscake.PagespeedTest)
	return []map[string]interface{}{
		map[string]interface{}{
			"alert_bigger":  flattenPagespeedCheckAlertBigger(original.AlertBigger, d),
			"alert_slower":  flattenPagespeedCheckAlertSlower(original.AlertSlower, d),
			"alert_smaller": flattenPagespeedCheckAlertSmaller(original.AlertSmaller, d),
		},
	}
}

func expandPagespeedCheckAlertSlower(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int64(v.(int)), nil
}

func flattenPagespeedCheckAlertSlower(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckAlertSmaller(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenPagespeedCheckAlertSmaller(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckInterval(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.PagespeedTestCheckRate(v.(int)), nil
}

func flattenPagespeedCheckInterval(v interface{}, d *schema.ResourceData) interface{} {
	return int(v.(statuscake.PagespeedTestCheckRate))
}

func expandPagespeedCheckContactGroups(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenPagespeedCheckContactGroups(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenPagespeedCheckLocation(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckMonitoredResource(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	address, err := expandPagespeedCheckAddress(original["address"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("monitored_resource.0.address") {
		transformed["website_url"] = address
	}

	return transformed, nil
}

func flattenPagespeedCheckMonitoredResource(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.PagespeedTest)
	return []map[string]interface{}{
		map[string]interface{}{
			"address": flattenPagespeedCheckAddress(data.WebsiteURL, d),
		},
	}
}

func expandPagespeedCheckName(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenPagespeedCheckName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckPaused(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenPagespeedCheckPaused(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandPagespeedCheckRegion(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.PagespeedTestRegion(v.(string)), nil
}
