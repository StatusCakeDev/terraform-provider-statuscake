package statuscake

import (
	"fmt"
	"strconv"

	"log"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/schema"
)

func castInterfaceToSliceStrings(in interface{}) []string {
	input := in.([]interface{})
	res := make([]string, len(input))

	for i, element := range input {
		res[i] = element.(string)
	}
	return res
}

func resourceStatusCakeTest() *schema.Resource {
	return &schema.Resource{
		Create: CreateTest,
		Update: UpdateTest,
		Delete: DeleteTest,
		Read:   ReadTest,

		Schema: map[string]*schema.Schema{
			"test_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"website_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"website_url": {
				Type:     schema.TypeString,
				Required: true,
			},

			"contact_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"check_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},

			"test_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40,
			},

			"confirmations": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"trigger_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"custom_header": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"user_agent": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"uptime": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"node_locations": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"ping_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"basic_user": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"basic_pass": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"public": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"logo_image": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"branding": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"website_host": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"virus": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"find_string": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"do_not_find": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"real_browser": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"test_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status_codes": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "204, 205, 206, 303, 400, 401, 403, 404, 405, 406, " +
					"408, 410, 413, 444, 429, 494, 495, 496, 499, 500, 501, 502, 503, " +
					"504, 505, 506, 507, 508, 509, 510, 511, 521, 522, 523, 524, 520, " +
					"598, 599",
			},

			"use_jar": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"post_raw": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"final_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"follow_redirect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func CreateTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	newTest := &statuscake.Test{
		WebsiteName:    d.Get("website_name").(string),
		WebsiteURL:     d.Get("website_url").(string),
		CheckRate:      d.Get("check_rate").(int),
		TestType:       d.Get("test_type").(string),
		Paused:         d.Get("paused").(bool),
		Timeout:        d.Get("timeout").(int),
		ContactID:      d.Get("contact_id").(int),
		Confirmation:   d.Get("confirmations").(int),
		Port:           d.Get("port").(int),
		TriggerRate:    d.Get("trigger_rate").(int),
		CustomHeader:   d.Get("custom_header").(string),
		UserAgent:      d.Get("user_agent").(string),
		Status:         d.Get("status").(string),
		Uptime:         d.Get("uptime").(float64),
		NodeLocations:  castInterfaceToSliceStrings(d.Get("node_locations")),
		PingURL:        d.Get("ping_url").(string),
		BasicUser:      d.Get("basic_user").(string),
		BasicPass:      d.Get("basic_pass").(string),
		Public:         d.Get("public").(int),
		LogoImage:      d.Get("logo_image").(string),
		Branding:       d.Get("branding").(int),
		WebsiteHost:    d.Get("website_host").(string),
		Virus:          d.Get("virus").(int),
		FindString:     d.Get("find_string").(string),
		DoNotFind:      d.Get("do_not_find").(bool),
		RealBrowser:    d.Get("real_browser").(int),
		TestTags:       d.Get("test_tags").(string),
		StatusCodes:    d.Get("status_codes").(string),
		UseJar:         d.Get("use_jar").(int),
		PostRaw:        d.Get("post_raw").(string),
		FinalEndpoint:  d.Get("final_endpoint").(string),
		FollowRedirect: d.Get("follow_redirect").(bool),
	}

	log.Printf("[DEBUG] Creating new StatusCake Test: %s", d.Get("website_name").(string))

	response, err := client.Tests().Update(newTest)
	if err != nil {
		return fmt.Errorf("Error creating StatusCake Test: %s", err.Error())
	}

	d.Set("test_id", fmt.Sprintf("%d", response.TestID))
	d.Set("status", response.Status)
	d.Set("uptime", fmt.Sprintf("%.3f", response.Uptime))
	d.SetId(fmt.Sprintf("%d", response.TestID))

	return ReadTest(d, meta)
}

func UpdateTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	params := getStatusCakeTestInput(d)

	log.Printf("[DEBUG] StatusCake Test Update for %s", d.Id())
	_, err := client.Tests().Update(params)
	if err != nil {
		return fmt.Errorf("Error Updating StatusCake Test: %s", err.Error())
	}
	return nil
}

func DeleteTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	testId, parseErr := strconv.Atoi(d.Id())
	if parseErr != nil {
		return parseErr
	}
	log.Printf("[DEBUG] Deleting StatusCake Test: %s", d.Id())
	err := client.Tests().Delete(testId)
	if err != nil {
		return err
	}

	return nil
}

func ReadTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	testId, parseErr := strconv.Atoi(d.Id())
	if parseErr != nil {
		return parseErr
	}
	testResp, err := client.Tests().Detail(testId)
	if err != nil {
		return fmt.Errorf("Error Getting StatusCake Test Details for %s: Error: %s", d.Id(), err)
	}
	d.Set("website_name", testResp.WebsiteName)
	d.Set("website_url", testResp.WebsiteURL)
	d.Set("check_rate", testResp.CheckRate)
	d.Set("test_type", testResp.TestType)
	d.Set("paused", testResp.Paused)
	d.Set("timeout", testResp.Timeout)
	d.Set("contact_id", testResp.ContactID)
	d.Set("confirmations", testResp.Confirmation)
	d.Set("port", testResp.Port)
	d.Set("trigger_rate", testResp.TriggerRate)
	d.Set("custom_header", testResp.CustomHeader)
	d.Set("user_agent", testResp.UserAgent)
	d.Set("status", testResp.Status)
	d.Set("uptime", testResp.Uptime)
	if err := d.Set("node_locations", testResp.NodeLocations); err != nil {
		return fmt.Errorf("[WARN] Error setting node locations: %s", err)
	}
	d.Set("ping_url", testResp.PingURL)
	d.Set("basic_user", testResp.BasicUser)
	d.Set("basic_pass", testResp.BasicPass)
	d.Set("public", testResp.Public)
	d.Set("logo_image", testResp.LogoImage)
	d.Set("branding", testResp.Branding)
	d.Set("website_host", testResp.WebsiteHost)
	d.Set("virus", testResp.Virus)
	d.Set("find_string", testResp.FindString)
	d.Set("do_not_find", testResp.DoNotFind)
	d.Set("real_browser", testResp.RealBrowser)
	d.Set("test_tags", testResp.TestTags)
	d.Set("status_codes", testResp.StatusCodes)
	d.Set("use_jar", testResp.UseJar)
	d.Set("post_raw", testResp.PostRaw)
	d.Set("final_endpoint", testResp.FinalEndpoint)
	d.Set("follow_redirect", testResp.FollowRedirect)

	return nil
}

func getStatusCakeTestInput(d *schema.ResourceData) *statuscake.Test {
	testId, parseErr := strconv.Atoi(d.Id())
	if parseErr != nil {
		log.Printf("[DEBUG] Error Parsing StatusCake TestID: %s", d.Id())
	}
	test := &statuscake.Test{
		TestID: testId,
	}
	if v, ok := d.GetOk("website_name"); ok {
		test.WebsiteName = v.(string)
	}
	if v, ok := d.GetOk("website_url"); ok {
		test.WebsiteURL = v.(string)
	}
	if v, ok := d.GetOk("check_rate"); ok {
		test.CheckRate = v.(int)
	}
	if v, ok := d.GetOk("contact_id"); ok {
		test.ContactID = v.(int)
	}
	if v, ok := d.GetOk("test_type"); ok {
		test.TestType = v.(string)
	}
	if v, ok := d.GetOk("paused"); ok {
		test.Paused = v.(bool)
	}
	if v, ok := d.GetOk("timeout"); ok {
		test.Timeout = v.(int)
	}
	if v, ok := d.GetOk("confirmations"); ok {
		test.Confirmation = v.(int)
	}
	if v, ok := d.GetOk("port"); ok {
		test.Port = v.(int)
	}
	if v, ok := d.GetOk("trigger_rate"); ok {
		test.TriggerRate = v.(int)
	}
	if v, ok := d.GetOk("custom_header"); ok {
		test.CustomHeader = v.(string)
	}
	if v, ok := d.GetOk("user_agent"); ok {
		test.UserAgent = v.(string)
	}
	if v, ok := d.GetOk("node_locations"); ok {
		test.NodeLocations = castInterfaceToSliceStrings(v)
	}
	if v, ok := d.GetOk("ping_url"); ok {
		test.PingURL = v.(string)
	}
	if v, ok := d.GetOk("basic_user"); ok {
		test.BasicUser = v.(string)
	}
	if v, ok := d.GetOk("basic_pass"); ok {
		test.BasicPass = v.(string)
	}
	if v, ok := d.GetOk("public"); ok {
		test.Public = v.(int)
	}
	if v, ok := d.GetOk("logo_image"); ok {
		test.LogoImage = v.(string)
	}
	if v, ok := d.GetOk("branding"); ok {
		test.Branding = v.(int)
	}
	if v, ok := d.GetOk("website_host"); ok {
		test.WebsiteHost = v.(string)
	}
	if v, ok := d.GetOk("virus"); ok {
		test.Virus = v.(int)
	}
	if v, ok := d.GetOk("find_string"); ok {
		test.FindString = v.(string)
	}
	if v, ok := d.GetOk("do_not_find"); ok {
		test.DoNotFind = v.(bool)
	}
	if v, ok := d.GetOk("real_browser"); ok {
		test.RealBrowser = v.(int)
	}
	if v, ok := d.GetOk("test_tags"); ok {
		test.TestTags = v.(string)
	}
	if v, ok := d.GetOk("status_codes"); ok {
		test.StatusCodes = v.(string)
	}
	if v, ok := d.GetOk("use_jar"); ok {
		test.UseJar = v.(int)
	}
	if v, ok := d.GetOk("post_raw"); ok {
		test.PostRaw = v.(string)
	}
	if v, ok := d.GetOk("final_endpoint"); ok {
		test.FinalEndpoint = v.(string)
	}
	if v, ok := d.GetOk("follow_redirect"); ok {
		test.FollowRedirect = v.(bool)
	}

	return test
}
