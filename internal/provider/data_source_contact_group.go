package provider

import (
	"context"
	"net/http"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	intvalidation "github.com/StatusCakeDev/terraform-provider-statuscake/internal/provider/validation"
)

func dataSourceStatusCakeContactGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStatusCakeContactGroupRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Contact group ID",
				ValidateFunc: intvalidation.StringIsNumerical,
			},
			"email_addresses": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of email addresses",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"integrations": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of integration IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mobile_numbers": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of international format mobile phone numbers",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the contact group",
			},
			"ping_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL or IP address of an endpoint to push uptime events. Currently this only supports HTTP GET endpoints",
			},
		},
	}
}

func dataSourceStatusCakeContactGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Get("id").(string)

	res, err := client.GetContactGroup(ctx, id).Execute()

	// If the resource is not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get contact group with ID: %s", err)
	}

	if err := d.Set("email_addresses", flattenContactGroupEmailAddresses(res.Data.EmailAddresses, d)); err != nil {
		return diag.Errorf("failed to read email addresses: %s", err)
	}

	if err := d.Set("integrations", flattenContactGroupIntegrations(res.Data.Integrations, d)); err != nil {
		return diag.Errorf("failed to read integrations: %s", err)
	}

	if err := d.Set("mobile_numbers", flattenContactGroupMobileNumbers(res.Data.MobileNumbers, d)); err != nil {
		return diag.Errorf("failed to read mobile numbers: %s", err)
	}

	if err := d.Set("name", flattenContactGroupName(res.Data.Name, d)); err != nil {
		return diag.Errorf("failed to read name: %s", err)
	}

	if err := d.Set("ping_url", flattenContactGroupPingURL(res.Data.PingURL, d)); err != nil {
		return diag.Errorf("failed to ping url: %s", err)
	}

	d.SetId(res.Data.ID)
	return nil
}
