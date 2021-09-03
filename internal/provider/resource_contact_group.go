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

func resourceStatusCakeContactGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakeContactGroupCreate,
		ReadContext:   resourceStatusCakeContactGroupRead,
		UpdateContext: resourceStatusCakeContactGroupUpdate,
		DeleteContext: resourceStatusCakeContactGroupDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email_addresses": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of email addresses",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: intvalidation.IsEmailAddress,
				},
			},
			"integrations": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of integration IDs",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: intvalidation.StringIsNumerical,
				},
			},
			"mobile_numbers": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of international format mobile phone numbers",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the contact group",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"ping_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "URL or IP address of an endpoint to push uptime events. Currently this only supports HTTP GET endpoints",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
	}
}

func resourceStatusCakeContactGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})

	emailAddresses, err := expandContactGroupEmailAddresses(d.Get("email_addresses"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("email_addresses") {
		body["email_addresses"] = emailAddresses
	}

	integrations, err := expandContactGroupIntegrations(d.Get("integrations"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("integrations") {
		body["integrations"] = integrations
	}

	mobileNumbers, err := expandContactGroupMobileNumbers(d.Get("mobile_numbers"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("mobile_numbers") {
		body["mobile_numbers"] = mobileNumbers
	}

	name, err := expandContactGroupName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	url, err := expandContactGroupPingURL(d.Get("ping_url"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("ping_url") {
		body["ping_url"] = url
	}

	log.Print("[DEBUG] Creating StatusCake contact group")
	log.Printf("[DEBUG] Request body: %+v", body)

	res, err := client.CreateContactGroupWithData(ctx, body).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create contact group: %w", err))
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeContactGroupRead(ctx, d, meta)
}

func resourceStatusCakeContactGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	res, err := client.GetContactGroup(ctx, id).Execute()

	// If the resource it not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get contact group with ID: %w", err))
	}

	if err := d.Set("email_addresses", flattenContactGroupEmailAddresses(res.Data.EmailAddresses, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read email addresses: %+v", err))
	}

	if err := d.Set("integrations", flattenContactGroupIntegrations(res.Data.Integrations, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read integrations: %+v", err))
	}

	if err := d.Set("mobile_numbers", flattenContactGroupMobileNumbers(res.Data.MobileNumbers, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read mobile numbers: %+v", err))
	}

	if err := d.Set("name", flattenContactGroupName(res.Data.Name, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to read name: %+v", err))
	}

	if err := d.Set("ping_url", flattenContactGroupPingURL(res.Data.PingURL, d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to ping url: %+v", err))
	}

	return nil
}

func resourceStatusCakeContactGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})
	id := d.Id()

	emailAddresses, err := expandContactGroupEmailAddresses(d.Get("email_addresses"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("email_addresses") {
		body["email_addresses"] = emailAddresses
	}

	integrations, err := expandContactGroupIntegrations(d.Get("integrations"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("integrations") {
		body["integrations"] = integrations
	}

	mobileNumbers, err := expandContactGroupMobileNumbers(d.Get("mobile_numbers"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("mobile_numbers") {
		body["mobile_numbers"] = mobileNumbers
	}

	name, err := expandContactGroupName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	url, err := expandContactGroupPingURL(d.Get("ping_url"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("ping_url") {
		body["ping_url"] = url
	}

	log.Printf("[DEBUG] Updating StatusCake contact group with ID: %s", id)
	log.Printf("[DEBUG] Request body: %+v", body)

	if err := client.UpdateContactGroupWithData(ctx, id, body).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update contact group: %w", err))
	}

	return resourceStatusCakeContactGroupRead(ctx, d, meta)
}

func resourceStatusCakeContactGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake contact group with ID: %s", id)

	if err := client.DeleteContactGroup(ctx, id).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete contact group with id %s: %w", id, err))
	}

	return nil
}

func expandContactGroupEmailAddresses(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenContactGroupEmailAddresses(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandContactGroupIntegrations(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenContactGroupIntegrations(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandContactGroupMobileNumbers(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenContactGroupMobileNumbers(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandContactGroupPingURL(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenContactGroupPingURL(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandContactGroupName(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenContactGroupName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
