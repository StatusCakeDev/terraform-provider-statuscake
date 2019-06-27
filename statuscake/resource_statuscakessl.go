package statuscake

import (
	"fmt"
	"strconv"

	"log"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStatusCakeSsl() *schema.Resource {
	return &schema.Resource{
		Create: CreateSsl,
		Update: UpdateSsl,
		Delete: DeleteSsl,
		Read:   ReadSsl,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ssl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"contact_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"contact_groups_c": {
				Type:     schema.TypeString,
				Required: true,
			},

			"checkrate": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"alert_at": {
				Type:     schema.TypeString,
				Required: true,
			},

			"alert_reminder": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alert_expiry": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alert_broken": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alert_mixed": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"paused": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"issuer_cn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cert_score": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cipher_score": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cert_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cipher": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"valid_from_utc": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"valid_until_utc": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mixed_content": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
				Computed: true,
			},

			"flags": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeBool},
				Computed: true,
			},

			"last_reminder": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"last_updated_utc": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateSsl(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	newSsl := &statuscake.PartialSsl{
		Domain:         d.Get("domain").(string),
		Checkrate:      strconv.Itoa(d.Get("checkrate").(int)),
		ContactGroupsC: d.Get("contact_groups_c").(string),
		AlertReminder:  d.Get("alert_reminder").(bool),
		AlertExpiry:    d.Get("alert_expiry").(bool),
		AlertBroken:    d.Get("alert_broken").(bool),
		AlertMixed:     d.Get("alert_mixed").(bool),
		AlertAt:        d.Get("alert_at").(string),
	}

	log.Printf("[DEBUG] Creating new StatusCake Ssl: %s", d.Get("domain").(string))

	response, err := statuscake.NewSsls(client).Create(newSsl)
	if err != nil {
		fmt.Println(newSsl)
		fmt.Println(client)
		return fmt.Errorf("Error creating StatusCake Ssl: %s", err.Error())
	}

	d.Set("ssl_id", response.ID)
	d.Set("contact_groups", response.ContactGroups)
	d.Set("paused", response.Paused)
	d.Set("issuer_cn", response.IssuerCn)
	d.Set("cert_score", response.CertScore)
	d.Set("cipher_score", response.CipherScore)
	d.Set("cert_status", response.CertStatus)
	d.Set("cipher", response.Cipher)
	d.Set("valid_from_utc", response.ValidFromUtc)
	d.Set("valid_until_utc", response.ValidUntilUtc)
	d.Set("mixed_content", response.MixedContent)
	d.Set("flags", response.Flags)
	d.Set("last_reminder", response.LastReminder)
	d.Set("last_updated_utc", response.LastUpdatedUtc)
	d.SetId(response.ID)

	return ReadSsl(d, meta)
}

func UpdateSsl(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	params := getStatusCakeSslInput(d)

	log.Printf("[DEBUG] StatusCake Ssl Update for %s", d.Id())
	_, err := statuscake.NewSsls(client).Update(params)
	if err != nil {
		return fmt.Errorf("Error Updating StatusCake Ssl: %s", err.Error())
	}
	return nil
}

func DeleteSsl(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	log.Printf("[DEBUG] Deleting StatusCake Ssl: %s", d.Id())
	err := statuscake.NewSsls(client).Delete(d.Id())

	return err
}

func ReadSsl(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	response, err := statuscake.NewSsls(client).Detail(d.Id())
	if err != nil {
		return fmt.Errorf("Error Getting StatusCake Ssl Details for %s: Error: %s", d.Id(), err)
	}
	d.Set("domain", response.Domain)
	d.Set("checkrate", response.Checkrate)
	d.Set("contact_groups_c", response.ContactGroupsC)
	d.Set("alert_reminder", response.AlertReminder)
	d.Set("alert_expiry", response.AlertExpiry)
	d.Set("alert_broken", response.AlertBroken)
	d.Set("alert_mixed", response.AlertMixed)
	d.Set("alert_at", response.AlertAt)
	d.Set("ssl_id", response.ID)
	d.Set("contact_groups", response.ContactGroups)
	d.Set("paused", response.Paused)
	d.Set("issuer_cn", response.IssuerCn)
	d.Set("cert_score", response.CertScore)
	d.Set("cipher_score", response.CipherScore)
	d.Set("cert_status", response.CertStatus)
	d.Set("cipher", response.Cipher)
	d.Set("valid_from_utc", response.ValidFromUtc)
	d.Set("valid_until_utc", response.ValidUntilUtc)
	d.Set("mixed_content", response.MixedContent)
	d.Set("flags", response.Flags)
	d.Set("last_reminder", response.LastReminder)
	d.Set("last_updated_utc", response.LastUpdatedUtc)
	d.SetId(response.ID)

	return nil
}

func getStatusCakeSslInput(d *schema.ResourceData) *statuscake.PartialSsl {
	sslId, parseErr := strconv.Atoi(d.Id())
	if parseErr != nil {
		log.Printf("[DEBUG] Error Parsing StatusCake Id: %s", d.Id())
	}
	ssl := &statuscake.PartialSsl{
		ID: sslId,
	}

	if v, ok := d.GetOk("domain"); ok {
		ssl.Domain = v.(string)
	}

	if v, ok := d.GetOk("checkrate"); ok {
		ssl.Checkrate = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("contact_groups_c"); ok {
		ssl.ContactGroupsC = v.(string)
	}

	if v, ok := d.GetOk("alert_reminder"); ok {
		ssl.AlertReminder = v.(bool)
	}

	if v, ok := d.GetOk("alert_expiry"); ok {
		ssl.AlertExpiry = v.(bool)
	}

	if v, ok := d.GetOk("alert_broken"); ok {
		ssl.AlertBroken = v.(bool)
	}

	if v, ok := d.GetOk("alert_mixed"); ok {
		ssl.AlertMixed = v.(bool)
	}

	if v, ok := d.GetOk("alert_at"); ok {
		ssl.AlertAt = v.(string)
	}

	return ssl
}
