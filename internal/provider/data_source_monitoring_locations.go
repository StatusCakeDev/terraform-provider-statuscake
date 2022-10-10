package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type monitoringLocationsFunc func(context.Context, *statuscake.Client, string) (statuscake.MonitoringLocations, error)

func dataSourceStatusCakeMonitoringLocations(fn monitoringLocationsFunc) *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStatusCakeMonitoringLocationsRead(fn),

		Schema: map[string]*schema.Schema{
			"region_code": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Location region code",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"locations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of monitoring locations",
				Elem: &schema.Resource{
					Schema: locationSchema(),
				},
			},
		},
	}
}

// locationsSchema returns the schema describing a single monitoring locations.
// Since locations features within multiple resources its structure has been
// encapsulated within a function.
func locationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location description",
		},
		"ipv4": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location IPv4 address",
		},
		"ipv6": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location IPv6 address",
		},
		"region": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location region",
		},
		"region_code": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location region code",
		},
		"status": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Location status",
		},
	}
}

func dataSourceStatusCakeMonitoringLocationsRead(fn monitoringLocationsFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*statuscake.Client)

		res, err := fn(ctx, client, d.Get("region_code").(string))
		if err != nil {
			return diag.Errorf("failed to list monitoring locations: %s", err)
		}

		if err := d.Set("locations", flattenMonitoringLocations(res.Data, d)); err != nil {
			return diag.Errorf("error setting monitoring locations: %s", err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

		return nil
	}
}

func listUptimeMonitoringLocations(ctx context.Context, client *statuscake.Client, location string) (statuscake.MonitoringLocations, error) {
	req := client.ListUptimeMonitoringLocations(ctx)

	if len(location) != 0 {
		req = req.Location(location)
	}

	return req.Execute()
}

func listPagespeedMonitoringLocations(ctx context.Context, client *statuscake.Client, location string) (statuscake.MonitoringLocations, error) {
	req := client.ListPagespeedMonitoringLocations(ctx)

	if len(location) != 0 {
		req = req.Location(location)
	}

	return req.Execute()
}

func flattenMonitoringLocationsDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringLocationsIPv4(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringLocationsIPv6(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringLocations(v interface{}, d *schema.ResourceData) interface{} {
	data := v.([]statuscake.MonitoringLocation)

	locations := make([]interface{}, len(data))
	for idx, location := range data {
		locations[idx] = flattenMonitoringLocation(location, d)
	}

	return locations
}

func flattenMonitoringLocation(v interface{}, d *schema.ResourceData) interface{} {
	data := v.(statuscake.MonitoringLocation)

	return map[string]interface{}{
		"description": flattenMonitoringLocationsDescription(data.Description, d),
		"ipv4":        flattenMonitoringLocationsIPv4(data.IPv4, d),
		"ipv6":        flattenMonitoringLocationsIPv6(data.IPv6, d),
		"region":      flattenMonitoringLocationsRegion(data.Region, d),
		"region_code": flattenMonitoringLocationsRegionCode(data.RegionCode, d),
		"status":      flattenMonitoringLocationsStatus(data.Status, d),
	}
}

func flattenMonitoringLocationsRegion(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringLocationsRegionCode(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringLocationsStatus(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
