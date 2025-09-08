package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	intdiag "github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider/diag"
	intvalidation "github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider/validation"
)

func resourceStatusCakeHeartbeatCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakeHeartbeatCheckCreate,
		ReadContext:   resourceStatusCakeHeartbeatCheckRead,
		UpdateContext: resourceStatusCakeHeartbeatCheckUpdate,
		DeleteContext: resourceStatusCakeHeartbeatCheckDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"check_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the heartbeat check",
			},
			"contact_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of contact group IDs",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: intvalidation.StringIsNumerical,
				},
			},
			"monitored_resource": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Monitored resource configuration block. This describes the server under test",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Name of the hosting provider",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the check",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the check should be run",
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of seconds since the last ping before the check is considered down.",
				ValidateFunc: validation.IntBetween(30, 172800),
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of tags",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceStatusCakeHeartbeatCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})

	contactGroups, err := expandHeartbeatCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	monitoredResource, err := expandHeartbeatCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	name, err := expandHeartbeatCheckName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	paused, err := expandHeartbeatCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	period, err := expandHeartbeatCheckPeriod(d.Get("period"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("period") {
		body["period"] = period
	}

	tags, err := expandHeartbeatCheckTags(d.Get("tags"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tags") {
		body["tags"] = tags
	}

	log.Printf("[DEBUG] Request body: %+v", body)

	res, err := client.CreateHeartbeatTestWithData(ctx, body).Execute()
	if err != nil {
		return intdiag.FromErr("failed to create heartbeat check", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeHeartbeatCheckRead(ctx, d, meta)
}

func resourceStatusCakeHeartbeatCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	res, err := client.GetHeartbeatTest(ctx, id).Execute()

	// If the resource is not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get heartbeat check with ID: %s, error: %s", id, err)
	}

	if err := d.Set("contact_groups", flattenHeartbeatCheckContactGroups(res.Data.ContactGroups, d)); err != nil {
		return diag.Errorf("failed to read contact groups: %s", err)
	}

	if err := d.Set("monitored_resource", flattenHeartbeatCheckMonitoredResource(res.Data, d)); err != nil {
		return diag.Errorf("failed to read monitored resource: %s", err)
	}

	if err := d.Set("name", flattenHeartbeatCheckName(res.Data.Name, d)); err != nil {
		return diag.Errorf("failed to read name: %s", err)
	}

	if err := d.Set("paused", flattenHeartbeatCheckPaused(res.Data.Paused, d)); err != nil {
		return diag.Errorf("failed to read paused: %s", err)
	}

	if err := d.Set("period", flattenHeartbeatCheckPeriod(res.Data.Period, d)); err != nil {
		return diag.Errorf("failed to read period: %s", err)
	}

	if err := d.Set("check_url", flattenHeartbeatCheckURL(res.Data.WebsiteURL, d)); err != nil {
		return diag.Errorf("failed to read check URL: %s", err)
	}

	if err := d.Set("tags", flattenHeartbeatCheckTags(res.Data.Tags, d)); err != nil {
		return diag.Errorf("failed to read tags: %s", err)
	}

	return nil
}

func resourceStatusCakeHeartbeatCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})
	id := d.Id()

	contactGroups, err := expandHeartbeatCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	monitoredResource, err := expandHeartbeatCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	name, err := expandHeartbeatCheckName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	paused, err := expandHeartbeatCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	period, err := expandHeartbeatCheckPeriod(d.Get("period"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("period") {
		body["period"] = period
	}

	tags, err := expandHeartbeatCheckTags(d.Get("tags"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tags") {
		body["tags"] = tags
	}

	log.Printf("[DEBUG] Updating StatusCake heartbeat check with ID: %s", id)
	log.Printf("[DEBUG] Request body: %+v", body)

	if err := client.UpdateHeartbeatTestWithData(ctx, id, body).Execute(); err != nil {
		return intdiag.FromErr(fmt.Sprintf("failed to update heartbeat check with id %s", id), err)
	}

	return resourceStatusCakeHeartbeatCheckRead(ctx, d, meta)
}

func resourceStatusCakeHeartbeatCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake heartbeat check with ID: %s", id)

	if err := client.DeleteHeartbeatTest(ctx, id).Execute(); err != nil {
		return intdiag.FromErr(fmt.Sprintf("failed to delete heartbeat check with id %s", id), err)
	}

	return nil
}

func expandHeartbeatCheckContactGroups(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenHeartbeatCheckContactGroups(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandHeartbeatCheckHost(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenHeartbeatCheckHost(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandHeartbeatCheckMonitoredResource(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		// If the monitored resource is not set then return an empty map. This is
		// necessary for the Heartbeat API only because the monitored_resource block
		// is optional. Therefore when the entire block is removed then the "host"
		// field is not set.
		//
		// At present this causes the API to return an error. This is a bug in the
		// API and does have a fix ready to go.
		return map[string]interface{}{
			"host": "",
		}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	host, err := expandHeartbeatCheckHost(original["host"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("monitored_resource.0.host") {
		transformed["host"] = host
	}

	return transformed, nil
}

func flattenHeartbeatCheckMonitoredResource(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.HeartbeatTest)

	host := flattenHeartbeatCheckHost(data.Host, d)
	if !isValid(host) {
		return nil
	}

	return []map[string]interface{}{
		{
			"host": host,
		},
	}
}

func expandHeartbeatCheckName(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenHeartbeatCheckName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandHeartbeatCheckPaused(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenHeartbeatCheckPaused(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandHeartbeatCheckPeriod(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenHeartbeatCheckPeriod(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenHeartbeatCheckURL(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandHeartbeatCheckTags(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenHeartbeatCheckTags(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
