package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	intdiag "github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider/diag"
	intvalidation "github.com/StatusCakeDev/terraform-provider-statuscake/v2/internal/provider/validation"
)

const (
	matcherContains   = "CONTAINS_STRING"
	matcherNoContains = "NOT_CONTAINS_STRING"
)

func isHTTPCheckType(t statuscake.UptimeTestType) bool {
	return t == statuscake.UptimeTestTypeHEAD ||
		t == statuscake.UptimeTestTypeHTTP
}

func isTCPCheckType(t statuscake.UptimeTestType) bool {
	return t == statuscake.UptimeTestTypeSMTP ||
		t == statuscake.UptimeTestTypeSSH ||
		t == statuscake.UptimeTestTypeTCP
}

func resourceStatusCakeUptimeCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakeUptimeCheckCreate,
		ReadContext:   resourceStatusCakeUptimeCheckRead,
		UpdateContext: resourceStatusCakeUptimeCheckUpdate,
		DeleteContext: resourceStatusCakeUptimeCheckDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"check_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of seconds between checks",
				ValidateFunc: intvalidation.Int32InSlice(statuscake.UptimeTestCheckRateValues()),
			},
			"confirmation": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				Description:  "Number of confirmation servers to confirm downtime before an alert is triggered",
				ValidateFunc: validation.IntBetween(0, 3),
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
			"dns_check": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "DNS check configuration block. Only one of `dns_check`, `http_check`, `icmp_check`, and `tcp_check` may be specified",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_ips": {
							Type:        schema.TypeSet,
							Required:    true,
							MinItems:    1,
							Description: "List of IP addresses to compare against returned DNS records",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsIPAddress,
							},
						},
						"dns_server": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "FQDN or IP address of the nameserver to query",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"http_check": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "HTTP check configuration block. Only one of `dns_check`, `http_check`, `icmp_check`, and `tcp_check` may be specified",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_authentication": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Basic Authentication (RFC7235) configuration block",
							Elem: &schema.Resource{
								Schema: basicAuthSchema(),
							},
						},
						"content_matchers": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Content matcher configuration block. This is used to assert values within the response of the request",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "String to look for within the response. Considered down if not found",
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"include_headers": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Include header content in string match search",
									},
									"matcher": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      matcherContains,
										Description:  "Whether to consider the check as down if the content is present within the response",
										ValidateFunc: validation.StringInSlice([]string{matcherContains, matcherNoContains}, false),
									},
								},
							},
						},
						"enable_cookies": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable cookie storage",
						},
						"final_endpoint": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Specify where the redirect chain should end",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"follow_redirects": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to follow redirects when testing. Disabled by default",
						},
						"request_headers": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Represents headers to be sent when making requests",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"request_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "HTTP",
							Description:  "Type of HTTP check. Either HTTP, or HEAD",
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HEAD"}, false),
						},
						"request_payload": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Payload submitted with the request. Setting this updates the check to use the HTTP POST verb. Only one of `request_payload` or `request_payload_raw` may be specified",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ConflictsWith: []string{"http_check.0.request_payload_raw"},
						},
						"request_payload_raw": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "Raw payload submitted with the request. Setting this updates the check to use the HTTP POST verb. Only one of `request_payload` or `request_payload_raw` may be specified",
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"http_check.0.request_payload"},
						},
						"status_codes": {
							Type:        schema.TypeSet,
							Computed:    true,
							Optional:    true,
							MinItems:    1,
							Description: "List of status codes that trigger an alert. If not specified then the default status codes are used. Once set, the default status codes cannot be restored and ommitting this field does not clear the attribute",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: intvalidation.StringIsNumerical,
							},
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      15,
							Description:  "The number of seconds to wait to receive the first byte",
							ValidateFunc: validation.IntBetween(5, 75),
						},
						"user_agent": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Custom user agent string set when testing",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"validate_ssl": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to send an alert if the SSL certificate is soon to expire",
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"icmp_check": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "ICMP check configuration block. Only one of `dns_check`, `http_check`, `icmp_check`, and `tcp_check` may be specified",
				// There are no special fields for an ICMP check. All that is required
				// is the address which is supplied in the `monitored_resource` block.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Dummy attribute to allow for a nested block. This field should not be changed",
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"locations": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of assigned monitoring locations on which to run checks",
				Elem: &schema.Resource{
					Schema: locationSchema(),
				},
			},
			"monitored_resource": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Monitored resource configuration block. This describes the server under test",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "URL, FQDN, or IP address of the server under test",
							ValidateFunc: validation.StringIsNotEmpty,
						},
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
			"regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of regions on which to run checks. The values required for this parameter can be retrieved from the `GET /v1/uptime-locations` endpoint",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
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
			"tcp_check": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "TCP check configuration block. Only one of `dns_check`, `http_check`, `icmp_check`, and `tcp_check` may be specified",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Authentication configuration block",
							Elem: &schema.Resource{
								Schema: basicAuthSchema(),
							},
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Destination port for TCP checks",
							ValidateFunc: validation.IsPortNumber,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "TCP",
							Description:  "Type of TCP check. Either SMTP, SSH or TCP",
							ValidateFunc: validation.StringInSlice([]string{"SMTP", "SSH", "TCP"}, false),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      15,
							Description:  "The number of seconds to wait to receive the first byte",
							ValidateFunc: validation.IntBetween(5, 75),
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"trigger_rate": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "The number of minutes to wait before sending an alert",
				ValidateFunc: validation.IntBetween(0, 60),
			},
		},
	}
}

// basicAuthSchema returns the schema describing a basic authentication. Since
// basic auth can be found in multiple check types its structure has been
// encapsulated within a function.
func basicAuthSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
	}
}

func resourceStatusCakeUptimeCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})

	checkInterval, err := expandUptimeCheckInterval(d.Get("check_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("check_interval") {
		body["check_rate"] = checkInterval
	}

	confirmation, err := expandUptimeCheckConfirmation(d.Get("confirmation"), d)
	if err != nil {
		return diag.FromErr(err)
	}
	body["confirmation"] = confirmation

	contactGroups, err := expandUptimeCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	dnsCheck, err := expandUptimeCheckDNSCheck(d.Get("dns_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("dns_check") {
		body = merge(body, dnsCheck.(map[string]interface{}))
	}

	httpCheck, err := expandUptimeCheckHTTPCheck(d.Get("http_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("http_check") {
		body = merge(body, httpCheck.(map[string]interface{}))
	}

	icmpCheck, err := expandUptimeCheckICMPCheck(d.Get("icmp_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("icmp_check") {
		body = merge(body, icmpCheck.(map[string]interface{}))
	}

	monitoredResource, err := expandUptimeCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	name, err := expandUptimeCheckName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	paused, err := expandUptimeCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	regions, err := expandUptimeCheckRegions(d.Get("regions"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("regions") {
		body["regions"] = regions
	}

	tags, err := expandUptimeCheckTags(d.Get("tags"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tags") {
		body["tags"] = tags
	}

	tcpCheck, err := expandUptimeCheckTCPCheck(d.Get("tcp_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tcp_check") {
		body = merge(body, tcpCheck.(map[string]interface{}))
	}

	triggerRate, err := expandUptimeCheckTriggerRate(d.Get("trigger_rate"), d)
	if err != nil {
		return diag.FromErr(err)
	}
	body["trigger_rate"] = triggerRate

	log.Print("[DEBUG] Creating StatusCake uptime check")
	log.Printf("[DEBUG] Request body: %+v", body)

	res, err := client.CreateUptimeTestWithData(ctx, body).Execute()
	if err != nil {
		return intdiag.FromErr("failed to create uptime check", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeUptimeCheckRead(ctx, d, meta)
}

func resourceStatusCakeUptimeCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	res, err := client.GetUptimeTest(ctx, id).Execute()

	// If the resource is not found then remove it from the state.
	if err, ok := err.(statuscake.APIError); ok && err.Status == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get uptime check with ID: %s", err)
	}

	if err := d.Set("check_interval", flattenUptimeCheckInterval(res.Data.CheckRate, d)); err != nil {
		return diag.Errorf("failed to read check interval: %s", err)
	}

	if err := d.Set("confirmation", flattenUptimeCheckConfirmation(res.Data.Confirmation, d)); err != nil {
		return diag.Errorf("failed to read confirmation: %s", err)
	}

	if err := d.Set("contact_groups", flattenUptimeCheckContactGroups(res.Data.ContactGroups, d)); err != nil {
		return diag.Errorf("failed to read contact groups: %s", err)
	}

	if err := d.Set("dns_check", flattenUptimeCheckDNSCheck(res.Data, d)); err != nil {
		return diag.Errorf("failed to read DNS check: %s", err)
	}

	if err := d.Set("http_check", flattenUptimeCheckHTTPCheck(res.Data, d)); err != nil {
		return diag.Errorf("failed to read HTTP check: %s", err)
	}

	if err := d.Set("icmp_check", flattenUptimeCheckICMPCheck(res.Data, d)); err != nil {
		return diag.Errorf("failed to read ICMP check: %s", err)
	}

	if err := d.Set("monitored_resource", flattenUptimeCheckMonitoredResource(res.Data, d)); err != nil {
		return diag.Errorf("failed to read monitored resource: %s", err)
	}

	if err := d.Set("name", flattenUptimeCheckName(res.Data.Name, d)); err != nil {
		return diag.Errorf("failed to read name: %s", err)
	}

	if err := d.Set("paused", flattenUptimeCheckPaused(res.Data.Paused, d)); err != nil {
		return diag.Errorf("failed to read paused: %s", err)
	}

	if err := d.Set("locations", flattenMonitoringLocations(res.Data.Servers, d)); err != nil {
		return diag.Errorf("failed to read locations: %s", err)
	}

	if err := d.Set("tags", flattenUptimeCheckTags(res.Data.Tags, d)); err != nil {
		return diag.Errorf("failed to read tags: %s", err)
	}

	if err := d.Set("tcp_check", flattenUptimeCheckTCPCheck(res.Data, d)); err != nil {
		return diag.Errorf("failed to read TCP check: %s", err)
	}

	if err := d.Set("trigger_rate", flattenUptimeCheckTriggerRate(res.Data.TriggerRate, d)); err != nil {
		return diag.Errorf("failed to read trigger rate: %s", err)
	}

	return nil
}

func resourceStatusCakeUptimeCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	body := make(map[string]interface{})
	id := d.Id()

	checkInterval, err := expandUptimeCheckInterval(d.Get("check_interval"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("check_interval") {
		body["check_rate"] = checkInterval
	}

	confirmation, err := expandUptimeCheckConfirmation(d.Get("confirmation"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("confirmation") {
		body["confirmation"] = confirmation
	}

	contactGroups, err := expandUptimeCheckContactGroups(d.Get("contact_groups"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("contact_groups") {
		body["contact_groups"] = contactGroups
	}

	dnsCheck, err := expandUptimeCheckDNSCheck(d.Get("dns_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("dns_check") {
		body = merge(body, dnsCheck.(map[string]interface{}))
	}

	httpCheck, err := expandUptimeCheckHTTPCheck(d.Get("http_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("http_check") {
		body = merge(body, httpCheck.(map[string]interface{}))
	}

	icmpCheck, err := expandUptimeCheckICMPCheck(d.Get("icmp_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("icmp_check") {
		body = merge(body, icmpCheck.(map[string]interface{}))
	}

	monitoredResource, err := expandUptimeCheckMonitoredResource(d.Get("monitored_resource"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("monitored_resource") {
		body = merge(body, monitoredResource.(map[string]interface{}))
	}

	name, err := expandUptimeCheckName(d.Get("name"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("name") {
		body["name"] = name
	}

	paused, err := expandUptimeCheckPaused(d.Get("paused"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("paused") {
		body["paused"] = paused
	}

	regions, err := expandUptimeCheckRegions(d.Get("regions"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("regions") {
		body["regions"] = regions
	}

	tags, err := expandUptimeCheckTags(d.Get("tags"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tags") {
		body["tags"] = tags
	}

	tcpCheck, err := expandUptimeCheckTCPCheck(d.Get("tcp_check"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("tcp_check") {
		body = merge(body, tcpCheck.(map[string]interface{}))
	}

	triggerRate, err := expandUptimeCheckTriggerRate(d.Get("trigger_rate"), d)
	if err != nil {
		return diag.FromErr(err)
	} else if d.HasChange("trigger_rate") {
		body["trigger_rate"] = triggerRate
	}

	log.Printf("[DEBUG] Updating StatusCake uptime check with ID: %s", id)
	log.Printf("[DEBUG] Request body: %+v", body)

	if err := client.UpdateUptimeTestWithData(ctx, id, body).Execute(); err != nil {
		return intdiag.FromErr(fmt.Sprintf("failed to update uptime check with id %s", id), err)
	}

	return resourceStatusCakeUptimeCheckRead(ctx, d, meta)
}

func resourceStatusCakeUptimeCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.Client)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake uptime check with ID: %s", id)

	if err := client.DeleteUptimeTest(ctx, id).Execute(); err != nil {
		return intdiag.FromErr(fmt.Sprintf("failed to delete uptime check with id %s", id), err)
	}

	return nil
}

func expandUptimeCheckAddress(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckAddress(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckBasicAuthentication(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	password, err := expandUptimeCheckPassword(original["password"], d)
	if err != nil {
		return nil, err
	} else if d.HasChanges("http_check.0.basic_authentication.0.password", "tcp_check.0.authentication.0.password") {
		transformed["basic_password"] = password
	}

	username, err := expandUptimeCheckUsername(original["username"], d)
	if err != nil {
		return nil, err
	} else if d.HasChanges("http_check.0.basic_authentication.0.username", "tcp_check.0.authentication.0.username") {
		transformed["basic_username"] = username
	}

	return transformed, nil
}

func flattenUptimeCheckBasicAuthentication(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckConfirmation(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenUptimeCheckConfirmation(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckContactGroups(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenUptimeCheckContactGroups(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckContentMatchers(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	content, err := expandUptimeCheckContent(original["content"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.content_matchers.0.content") {
		transformed["find_string"] = content
	}

	includeHeaders, err := expandUptimeCheckIncludeHeaders(original["include_headers"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.content_matchers.0.include_headers") {
		transformed["include_header"] = includeHeaders
	}

	invert, err := expandUptimeCheckMatcher(original["matcher"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.content_matchers.0.matcher") {
		transformed["do_not_find"] = invert
	}

	return transformed, nil
}

func flattenUptimeCheckContentMatchers(v interface{}, d *schema.ResourceData) interface{} {
	original := v.(statuscake.UptimeTest)
	if original.FindString == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"content":         flattenUptimeCheckContent(original.FindString, d),
			"include_headers": flattenUptimeCheckIncludeHeaders(original.IncludeHeader, d),
			"matcher":         flattenUptimeCheckMatcher(original.DoNotFind, d),
		},
	}
}

func expandUptimeCheckContent(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckContent(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckDNSCheck(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	transformed["test_type"] = statuscake.UptimeTestTypeDNS

	ips, err := expandUptimeCheckDNSIPs(original["dns_ips"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("dns_check.0.dns_ips") {
		transformed["dns_ips"] = ips
	}

	server, err := expandUptimeCheckDNSServer(original["dns_server"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("dns_check.0.dns_server") {
		transformed["dns_server"] = server
	}

	return transformed, nil
}

func flattenUptimeCheckDNSCheck(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.UptimeTest)
	if data.TestType != statuscake.UptimeTestTypeDNS {
		return nil
	}

	return []map[string]interface{}{
		{
			"dns_ips":    flattenUptimeCheckDNSIPs(data.DNSIPs, d),
			"dns_server": flattenUptimeCheckDNSServer(data.DNSServer, d),
		},
	}
}

func expandUptimeCheckDNSIPs(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenUptimeCheckDNSIPs(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckDNSServer(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckDNSServer(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckEnableCookies(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenUptimeCheckEnableCookies(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckFinalEndpoint(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckFinalEndpoint(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckFollowRedirects(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenUptimeCheckFollowRedirects(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckHTTPCheck(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	auth, err := expandUptimeCheckBasicAuthentication(original["basic_authentication"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.basic_authentication") {
		transformed = merge(transformed, auth.(map[string]interface{}))
	}

	matchers, err := expandUptimeCheckContentMatchers(original["content_matchers"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.content_matchers") {
		transformed = merge(transformed, matchers.(map[string]interface{}))
	}

	enableCookies, err := expandUptimeCheckEnableCookies(original["enable_cookies"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.enable_cookies") {
		transformed["use_jar"] = enableCookies
	}

	finalEndpoint, err := expandUptimeCheckFinalEndpoint(original["final_endpoint"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.final_endpoint") {
		transformed["final_endpoint"] = finalEndpoint
	}

	followRedirects, err := expandUptimeCheckFollowRedirects(original["follow_redirects"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.follow_redirects") {
		transformed["follow_redirects"] = followRedirects
	}

	headers, err := expandUptimeCheckRequestHeaders(original["request_headers"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.request_headers") {
		transformed["custom_header"] = headers
	}

	method, err := expandUptimeCheckRequestMethod(original["request_method"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.request_method") {
		transformed["test_type"] = method
	}

	payload, err := expandUptimeCheckRequestPayload(original["request_payload"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.request_payload") {
		transformed["post_body"] = payload
	}

	raw, err := expandUptimeCheckRequestRaw(original["request_payload_raw"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.request_payload_raw") {
		transformed["post_raw"] = raw
	}

	codes, err := expandUptimeCheckStatusCodes(original["status_codes"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.status_codes") {
		transformed["status_codes_csv"] = codes
	}

	timeout, err := expandUptimeCheckTimeout(original["timeout"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.timeout") {
		transformed["timeout"] = timeout
	}

	userAgent, err := expandUptimeCheckUserAgent(original["user_agent"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.user_agent") {
		transformed["user_agent"] = userAgent
	}

	validateSSL, err := expandUptimeCheckValidateSSL(original["validate_ssl"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("http_check.0.validate_ssl") {
		transformed["enable_ssl_alert"] = validateSSL
	}

	return transformed, nil
}

func flattenUptimeCheckHTTPCheck(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.UptimeTest)
	if !isHTTPCheckType(data.TestType) {
		return nil
	}

	return []map[string]interface{}{
		{
			"basic_authentication": flattenUptimeCheckBasicAuthentication(d.Get("http_check.0.basic_authentication"), d),
			"content_matchers":     flattenUptimeCheckContentMatchers(data, d),
			"enable_cookies":       flattenUptimeCheckEnableCookies(data.UseJAR, d),
			"final_endpoint":       flattenUptimeCheckFinalEndpoint(data.FinalEndpoint, d),
			"follow_redirects":     flattenUptimeCheckFollowRedirects(data.FollowRedirects, d),
			"request_headers":      flattenUptimeCheckRequestHeaders(data.CustomHeader, d),
			"request_method":       flattenUptimeCheckRequestMethod(data.TestType, d),
			"request_payload":      flattenUptimeCheckRequestPayload(data.PostBody, d),
			"request_payload_raw":  flattenUptimeCheckRequestRaw(data.PostRaw, d),
			"status_codes":         flattenUptimeCheckStatusCodes(data.StatusCodes, d),
			"timeout":              flattenUptimeCheckTimeout(data.Timeout, d),
			"user_agent":           flattenUptimeCheckUserAgent(data.UserAgent, d),
			"validate_ssl":         flattenUptimeCheckValidateSSL(data.EnableSSLAlert, d),
		},
	}
}

func expandUptimeCheckHost(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckHost(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckICMPCheck(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return map[string]interface{}{
		"test_type": statuscake.UptimeTestTypePING,
	}, nil
}

func flattenUptimeCheckICMPCheck(v interface{}, _ *schema.ResourceData) interface{} {
	data := v.(statuscake.UptimeTest)
	if data.TestType != statuscake.UptimeTestTypePING {
		return nil
	}

	return []map[string]interface{}{
		{
			"enabled": true,
		},
	}
}

func expandUptimeCheckIncludeHeaders(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenUptimeCheckIncludeHeaders(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckInterval(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.UptimeTestCheckRate(v.(int)), nil
}

func flattenUptimeCheckInterval(v interface{}, d *schema.ResourceData) interface{} {
	return int(v.(statuscake.UptimeTestCheckRate))
}

func expandUptimeCheckMatcher(v interface{}, _ *schema.ResourceData) (interface{}, error) {
	if v.(string) == matcherContains {
		return false, nil
	}
	return true, nil
}

func flattenUptimeCheckMatcher(v interface{}, _ *schema.ResourceData) interface{} {
	if !v.(bool) {
		return matcherContains
	}
	return matcherNoContains
}

func expandUptimeCheckMonitoredResource(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	address, err := expandUptimeCheckAddress(original["address"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("monitored_resource.0.address") {
		transformed["website_url"] = address
	}

	host, err := expandUptimeCheckHost(original["host"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("monitored_resource.0.host") {
		transformed["host"] = host
	}

	return transformed, nil
}

func flattenUptimeCheckMonitoredResource(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.UptimeTest)
	return []map[string]interface{}{
		{
			"address": flattenUptimeCheckAddress(data.WebsiteURL, d),
			"host":    flattenUptimeCheckHost(data.Host, d),
		},
	}
}

func expandUptimeCheckName(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckPassword(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckPassword(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckPaused(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenUptimeCheckPaused(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckPort(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenUptimeCheckPort(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckProtocol(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.UptimeTestType(v.(string)), nil
}

func flattenUptimeCheckProtocol(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckRequestHeaders(v interface{}, _ *schema.ResourceData) (interface{}, error) {
	if !isValid(v) {
		return "", nil
	}

	b, err := json.Marshal(v.(map[string]interface{}))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func flattenUptimeCheckRequestHeaders(v interface{}, _ *schema.ResourceData) interface{} {
	var headers map[string]interface{}
	if err := json.Unmarshal([]byte(stringElem(v)), &headers); err != nil {
		return map[string]interface{}{}
	}
	return headers
}

func expandUptimeCheckRequestMethod(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return statuscake.UptimeTestType(v.(string)), nil
}

func flattenUptimeCheckRequestMethod(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckRequestPayload(v interface{}, _ *schema.ResourceData) (interface{}, error) {
	if !isValid(v) {
		return "", nil
	}

	b, err := json.Marshal(v.(map[string]interface{}))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func flattenUptimeCheckRequestPayload(v interface{}, _ *schema.ResourceData) interface{} {
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(stringElem(v)), &body); err != nil {
		return map[string]interface{}{}
	}
	return body
}

func expandUptimeCheckRequestRaw(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckRequestRaw(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckRegions(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringList(v.([]interface{})), nil
}

func expandUptimeCheckStatusCodes(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return strings.Join(convertStringSet(v.(*schema.Set)), ","), nil
}

func flattenUptimeCheckStatusCodes(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckTags(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return convertStringSet(v.(*schema.Set)), nil
}

func flattenUptimeCheckTags(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckTCPCheck(v interface{}, d *schema.ResourceData) (interface{}, error) {
	l := v.([]interface{})

	if len(l) == 0 || l[0] == nil {
		return map[string]interface{}{}, nil
	}

	original := l[0].(map[string]interface{})
	transformed := make(map[string]interface{})

	auth, err := expandUptimeCheckBasicAuthentication(original["authentication"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("tcp_check.0.authentication") {
		transformed = merge(transformed, auth.(map[string]interface{}))
	}

	port, err := expandUptimeCheckPort(original["port"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("tcp_check.0.port") {
		transformed["port"] = port
	}

	protocol, err := expandUptimeCheckProtocol(original["protocol"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("tcp_check.0.protocol") {
		transformed["test_type"] = protocol
	}

	timeout, err := expandUptimeCheckTimeout(original["timeout"], d)
	if err != nil {
		return nil, err
	} else if d.HasChange("tcp_check.0.timeout") {
		transformed["timeout"] = timeout
	}

	return transformed, nil
}

func flattenUptimeCheckTCPCheck(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.UptimeTest)
	if !isTCPCheckType(data.TestType) {
		return nil
	}

	return []map[string]interface{}{
		{
			"authentication": flattenUptimeCheckBasicAuthentication(d.Get("tcp_check.0.authentication"), d),
			"port":           flattenUptimeCheckPort(data.Port, d),
			"protocol":       flattenUptimeCheckProtocol(data.TestType, d),
			"timeout":        flattenUptimeCheckTimeout(data.Timeout, d),
		},
	}
}

func expandUptimeCheckTimeout(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenUptimeCheckTimeout(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckTriggerRate(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return int32(v.(int)), nil
}

func flattenUptimeCheckTriggerRate(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckUserAgent(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckUserAgent(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckUsername(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(string), nil
}

func flattenUptimeCheckUsername(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandUptimeCheckValidateSSL(v interface{}, d *schema.ResourceData) (interface{}, error) {
	return v.(bool), nil
}

func flattenUptimeCheckValidateSSL(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
