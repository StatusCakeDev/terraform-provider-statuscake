package statuscake

import (
	"fmt"
	"strconv"

	"log"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/schema"
)

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

			"user_agent": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"use_jar": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"post_raw": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"find_string": {
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
				Default:  false,
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
		UserAgent:      d.Get("user_agent").(string),
		UseJar:         d.Get("use_jar").(bool),
		PostRaw:        d.Get("post_raw").(string),
		FindString:     d.Get("find_string").(string),
		FinalEndpoint:  d.Get("final_endpoint").(string),
		FollowRedirect: d.Get("follow_redirect").(bool),
	}

	log.Printf("[DEBUG] Creating new StatusCake Test: %s", d.Get("website_name").(string))

	response, err := client.Tests().Update(newTest)
	if err != nil {
		return fmt.Errorf("Error creating StatusCake Test: %s", err.Error())
	}

	d.Set("test_id", fmt.Sprintf("%d", response.TestID))
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
	d.Set("user_agent", testResp.UserAgent)
	d.Set("use_jar", testResp.UseJar)
	d.Set("post_raw", testResp.PostRaw)
	d.Set("find_string", testResp.FindString)
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
	if v, ok := d.GetOk("contact_id"); ok {
		test.ContactID = v.(int)
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
	if v, ok := d.GetOk("user_agent"); ok {
		test.UserAgent = v.(string)
	}
	if v, ok := d.GetOk("use_jar"); ok {
		test.UseJar = v.(bool)
	}
	if v, ok := d.GetOk("post_raw"); ok {
		test.PostRaw = v.(string)
	}
	if v, ok := d.GetOk("find_string"); ok {
		test.FindString = v.(string)
	}
	if v, ok := d.GetOk("final_endpoint"); ok {
		test.FinalEndpoint = v.(string)
	}
	if v, ok := d.GetOk("follow_redirect"); ok {
		test.FollowRedirect = v.(bool)
	}

	defaultStatusCodes := "204, 205, 206, 303, 400, 401, 403, 404, 405, 406, " +
		"408, 410, 413, 444, 429, 494, 495, 496, 499, 500, 501, 502, 503, " +
		"504, 505, 506, 507, 508, 509, 510, 511, 521, 522, 523, 524, 520, " +
		"598, 599"

	test.StatusCodes = defaultStatusCodes

	return test
}
