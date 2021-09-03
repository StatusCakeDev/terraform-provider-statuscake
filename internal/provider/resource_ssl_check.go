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

func resourceStatusCakeSSLCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakeSSLCheckCreate,
		ReadContext:   resourceStatusCakeSSLCheckRead,
		UpdateContext: resourceStatusCakeSSLCheckUpdate,
		DeleteContext: resourceStatusCakeSSLCheckDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"alert_config": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Alert configuration block",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_at": &schema.Schema{
							Type:        schema.TypeSet,
							Required:    true,
							MinItems:    3,
							MaxItems:    3,
							Description: "List representing when alerts should be sent (days). Must be exactly 3 numerical values",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntAtLeast(1),
							},
						},
						"on_broken": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alerts when SSL certificate issues are found",
						},
						"on_expiry": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alerts when the SSL certificate is to expire",
						},
						"on_mixed": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alerts when mixed content is found",
						},
						"on_reminder": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alert reminders",
						},
					},
				},
			},
			"check_interval": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of seconds between checks",
				ValidateFunc: intvalidation.Int32InSlice(statuscake.SSLTestCheckRateValues()),
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
			"follow_redirects": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to follow redirects when testing. Disabled by default",
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
							Description:  "URL of the server under test",
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
						"hostname": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Hostname of the server under test",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"paused": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the check should be run",
			},
			"user_agent": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Custom user agent string set when testing",
			},
		},
	}
}

func resourceStatusCakeSSLCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})

	config, err := expandSSLCheckAlertConfig(d.Get("alert_config"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("alert_config") {
		body = merge(body, config.(map[string]interface{}))
	}

	checkInterval, err := expandSSLCheckInterval(d.Get("check_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("check_interval") {
		body["check_rate"] = checkInterval
	}

	contactGroups, err := expandSSLCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	followRedirects, err := expandSSLCheckFollowRedirects(d.Get("follow_redirects"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("follow_redirects") {
		body["follow_redirects"] = followRedirects
	}

	monitoredResource, err := expandSSLCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	paused, err := expandSSLCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	userAgent, err := expandSSLCheckHostname(d.Get("user_agent"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("user_agent") {
		body["user_agent"] = userAgent
	}

	log.Print("[DEBUG] Creating StatusCake SSL check")
	log.Printf("[DEBUG] Request body: %+v", body)

	res, err := client.CreateSslTestWithData(ctx, body).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create SSL check: %w", err))
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeSSLCheckRead(ctx, d, meta)
}

func resourceStatusCakeSSLCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	res, err := client.GetSslTest(ctx, id).Execute()

	// If the resource it not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get SSL check with ID: %w", err))
	}

	if err := d.Set("alert_config", flattenSSLCheckAlertConfig(res.Data, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read alert config: %+v", err))
	}

	if err := d.Set("check_interval", flattenSSLCheckInterval(res.Data.CheckRate, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read check interval: %+v", err))
	}

	if err := d.Set("contact_groups", flattenSSLCheckContactGroups(res.Data.ContactGroups, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read contact groups: %+v", err))
	}

	if err := d.Set("follow_redirects", flattenSSLCheckFollowRedirects(res.Data.FollowRedirects, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read follow redirects: %+v", err))
	}

	if err := d.Set("monitored_resource", flattenSSLCheckMonitoredResource(res.Data, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read monitored resource: %+v", err))
	}

	if err := d.Set("paused", flattenSSLCheckPaused(res.Data.Paused, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read paused: %+v", err))
	}

	if err := d.Set("user_agent", flattenSSLCheckUserAgent(res.Data.UserAgent, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read user agent: %+v", err))
	}

	return nil
}

func resourceStatusCakeSSLCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})
	id := d.Id()

	config, err := expandSSLCheckAlertConfig(d.Get("alert_config"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("alert_config") {
		body = merge(body, config.(map[string]interface{}))
	}

	checkInterval, err := expandSSLCheckInterval(d.Get("check_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("check_interval") {
		body["check_rate"] = checkInterval
	}

	contactGroups, err := expandSSLCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	followRedirects, err := expandSSLCheckFollowRedirects(d.Get("follow_redirects"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("follow_redirects") {
		body["follow_redirects"] = followRedirects
	}

	monitoredResource, err := expandSSLCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	paused, err := expandSSLCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	userAgent, err := expandSSLCheckHostname(d.Get("user_agent"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("user_agent") {
		body["user_agent"] = userAgent
	}

	log.Printf("[DEBUG] Updating StatusCake SSL check with ID: %s", id)
	log.Printf("[DEBUG] Request body: %+v", body)

	if err := client.UpdateSslTestWithData(ctx, id, body).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update SSL check: %w", err))
	}

	return resourceStatusCakeSSLCheckRead(ctx, d, meta)
}

func resourceStatusCakeSSLCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake SSL check with ID: %s", id)

	if err := client.DeleteSslTest(ctx, id).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete SSL check with id %s: %w", id, err))
	}

	return nil
}

func expandSSLCheckAddress(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenSSLCheckAddress(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckAlertAt(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertInt32Set(v.(*schema.Set)), nil
}

func flattenSSLCheckAlertAt(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckAlertConfig(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	alertAt, err := expandSSLCheckAlertAt(original["alert_at"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.alert_at") {
		transformed["alert_at"] = alertAt
	}

	broken, err := expandSSLCheckOnBroken(original["on_broken"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.broken") {
		transformed["alert_broken"] = broken
	}

	expiry, err := expandSSLCheckOnExpiry(original["on_expiry"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.on_expiry") {
		transformed["alert_expiry"] = expiry
	}

	mixed, err := expandSSLCheckOnMixed(original["on_mixed"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.on_mixed") {
		transformed["alert_mixed"] = mixed
	}

	reminder, err := expandSSLCheckOnReminder(original["on_reminder"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("alert_config.0.on_reminder") {
		transformed["alert_reminder"] = reminder
	}

	return transformed, nil
}

func flattenSSLCheckAlertConfig(v interface{}, d *schema.ResourceData) interface{} {
	original := v.(statuscake.SSLTest)
	return []map[string]interface{}{
		map[string]interface{}{
			"alert_at":    flattenSSLCheckAlertAt(original.AlertAt, d),
			"on_broken":   flattenSSLCheckOnBroken(original.AlertBroken, d),
			"on_expiry":   flattenSSLCheckOnExpiry(original.AlertExpiry, d),
			"on_mixed":    flattenSSLCheckOnMixed(original.AlertMixed, d),
			"on_reminder": flattenSSLCheckOnReminder(original.AlertReminder, d),
		},
	}
}

func expandSSLCheckInterval(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.SSLTestCheckRate(v.(int)), nil
}

func flattenSSLCheckInterval(v interface{}, d *schema.ResourceData) interface{} {
	return int(v.(statuscake.SSLTestCheckRate))
}

func expandSSLCheckContactGroups(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenSSLCheckContactGroups(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckFollowRedirects(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenSSLCheckFollowRedirects(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckHostname(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenSSLCheckHostname(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckMonitoredResource(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	address, err := expandSSLCheckAddress(original["address"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("monitored_resource.0.address") {
		transformed["website_url"] = address
	}

	hostname, err := expandSSLCheckHostname(original["hostname"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("monitored_resource.0.hostname") {
		transformed["hostname"] = hostname
	}

	return transformed, nil
}

func flattenSSLCheckMonitoredResource(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.SSLTest)
	return []map[string]interface{}{
		map[string]interface{}{
			"address":  flattenSSLCheckAddress(data.WebsiteURL, d),
			"hostname": flattenSSLCheckHostname(data.Hostname, d),
		},
	}
}

func expandSSLCheckOnBroken(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenSSLCheckOnBroken(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckOnExpiry(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenSSLCheckOnExpiry(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckOnMixed(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenSSLCheckOnMixed(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckOnReminder(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenSSLCheckOnReminder(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckPaused(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenSSLCheckPaused(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandSSLCheckUserAgent(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenSSLCheckUserAgent(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
