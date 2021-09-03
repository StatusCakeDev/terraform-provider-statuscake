package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
	intvalidation "github.com/StatusCakeDev/terraform-provider-statuscake/internal/provider/validation"
)

func resourceStatusCakeMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakeMaintenanceWindowCreate,
		ReadContext:   resourceStatusCakeMaintenanceWindowRead,
		UpdateContext: resourceStatusCakeMaintenanceWindowUpdate,
		DeleteContext: resourceStatusCakeMaintenanceWindowDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"end": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "End of the maintenance window (RFC3339 format)",
				ValidateFunc: validation.IsRFC3339Time,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the maintenance window",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"repeat_interval": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "never",
				Description:  "How often the maintenance window should occur",
				ValidateFunc: validation.StringInSlice(statuscake.MaintenanceWindowRepeatIntervalValues(), false),
			},
			"start": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Start of the maintenance window (RFC3339 format)",
				ValidateFunc: validation.IsRFC3339Time,
			},
			"tags": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of tags used to include matching uptime checks in this maintenance window",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				AtLeastOneOf: []string{"tags", "tests"},
			},
			"tests": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of uptime check IDs explicitly included in this maintenance window",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: intvalidation.StringIsNumerical,
				},
				AtLeastOneOf: []string{"tags", "tests"},
			},
			"timezone": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Standard timezone associated with this maintenance window",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceStatusCakeMaintenanceWindowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})

	end, err := expandMaintenanceWindowEnd(d.Get("end"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("end") {
		body["end_at"] = end
	}

	name, err := expandMaintenanceWindowName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	interval, err := expandMaintenanceWindowRepeatInterval(d.Get("repeat_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("repeat_interval") {
		body["repeat_interval"] = interval
	}

	start, err := expandMaintenanceWindowStart(d.Get("start"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("start") {
		body["start_at"] = start
	}

	tags, err := expandMaintenanceWindowTags(d.Get("tags"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tags") {
		body["tags"] = tags
	}

	tests, err := expandMaintenanceWindowTests(d.Get("tests"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tests") {
		body["tests"] = tests
	}

	timezone, err := expandMaintenanceWindowTimezone(d.Get("timezone"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("timezone") {
		body["timezone"] = timezone
	}

	log.Print("[DEBUG] Creating StatusCake maintenance window")
	log.Printf("[DEBUG] Request body: %+v", body)

	res, err := client.CreateMaintenanceWindowWithData(ctx, body).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create maintenance window: %w", err))
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeMaintenanceWindowRead(ctx, d, meta)
}

func resourceStatusCakeMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	res, err := client.GetMaintenanceWindow(ctx, id).Execute()

	// If the resource it not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get maintenance window test with ID: %w", err))
	}

	if err := d.Set("end", flattenMaintenanceWindowEnd(res.Data.End, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read end: %+v", err))
	}

	if err := d.Set("name", flattenMaintenanceWindowName(res.Data.Name, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read name: %+v", err))
	}

	if err := d.Set("repeat_interval", flattenMaintenanceWindowRepeatInterval(res.Data.RepeatInterval, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read repeat interval: %+v", err))
	}

	if err := d.Set("start", flattenMaintenanceWindowStart(res.Data.Start, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read start: %+v", err))
	}

	if err := d.Set("tags", flattenMaintenanceWindowTags(res.Data.Tags, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read tags: %+v", err))
	}

	if err := d.Set("tests", flattenMaintenanceWindowTests(res.Data.Tests, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read tests: %+v", err))
	}

	if err := d.Set("timezone", flattenMaintenanceWindowTimezone(res.Data.Timezone, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read timezone: %+v", err))
	}

	return nil
}

func resourceStatusCakeMaintenanceWindowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})
	id := d.Id()

	end, err := expandMaintenanceWindowEnd(d.Get("end"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("end") {
		body["end_at"] = end
	}

	name, err := expandMaintenanceWindowName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	interval, err := expandMaintenanceWindowRepeatInterval(d.Get("repeat_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("repeat_interval") {
		body["repeat_interval"] = interval
	}

	start, err := expandMaintenanceWindowStart(d.Get("start"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("start") {
		body["start_at"] = start
	}

	tags, err := expandMaintenanceWindowTags(d.Get("tags"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tags") {
		body["tags"] = tags
	}

	tests, err := expandMaintenanceWindowTests(d.Get("tests"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tests") {
		body["tests"] = tests
	}

	timezone, err := expandMaintenanceWindowTimezone(d.Get("timezone"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("timezone") {
		body["timezone"] = timezone
	}

	log.Printf("[DEBUG] Updating StatusCake maintenance window with ID: %s", id)
	log.Printf("[DEBUG] Request body: %+v", body)

	if err := client.UpdateMaintenanceWindowWithData(ctx, id, body).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update maintenance window: %w", err))
	}

	return resourceStatusCakeMaintenanceWindowRead(ctx, d, meta)
}

func resourceStatusCakeMaintenanceWindowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake maintenance window with ID: %s", id)

	if err := client.DeleteMaintenanceWindow(ctx, id).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete maintenance window with id %s: %w", id, err))
	}

	return nil
}

func expandMaintenanceWindowEnd(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return time.Parse(time.RFC3339, v.(string))
}

func flattenMaintenanceWindowEnd(v interface{}, d *schema.ResourceData) interface{} {
	t := v.(time.Time)
	return t.Format(time.RFC3339)
}

func expandMaintenanceWindowName(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenMaintenanceWindowName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandMaintenanceWindowRepeatInterval(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.MaintenanceWindowRepeatInterval(v.(string)), nil
}

func flattenMaintenanceWindowRepeatInterval(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandMaintenanceWindowStart(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return time.Parse(time.RFC3339, v.(string))
}

func flattenMaintenanceWindowStart(v interface{}, d *schema.ResourceData) interface{} {
	t := v.(time.Time)
	return t.Format(time.RFC3339)
}

func expandMaintenanceWindowTags(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenMaintenanceWindowTags(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandMaintenanceWindowTests(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenMaintenanceWindowTests(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandMaintenanceWindowTimezone(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenMaintenanceWindowTimezone(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
