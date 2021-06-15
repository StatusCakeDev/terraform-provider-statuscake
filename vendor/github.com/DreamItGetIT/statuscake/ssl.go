package statuscake

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

//Ssl represent the data received by the API with GET
type Ssl struct {
	ID             string              `json:"id"                 url:"id,omitempty"`
	Domain         string              `json:"domain"             url:"domain,omitempty"`
	Checkrate      int                 `json:"checkrate"          url:"checkrate,omitempty"`
	ContactGroupsC string              `                          url:"contact_groups,omitempty"`
	AlertAt        string              `json:"alert_at"           url:"alert_at,omitempty"`
	AlertReminder  bool                `json:"alert_reminder"     url:"alert_expiry,omitempty"`
	AlertExpiry    bool                `json:"alert_expiry"       url:"alert_reminder,omitempty"`
	AlertBroken    bool                `json:"alert_broken"       url:"alert_broken,omitempty"`
	AlertMixed     bool                `json:"alert_mixed"        url:"alert_mixed,omitempty"`
	Paused         bool                `json:"paused"`
	IssuerCn       string              `json:"issuer_cn"`
	CertScore      string              `json:"cert_score"`
	CipherScore    string              `json:"cipher_score"`
	CertStatus     string              `json:"cert_status"`
	Cipher         string              `json:"cipher"`
	ValidFromUtc   string              `json:"valid_from_utc"`
	ValidUntilUtc  string              `json:"valid_until_utc"`
	MixedContent   []map[string]string `json:"mixed_content"`
	Flags          map[string]bool     `json:"flags"`
	ContactGroups  []string            `json:"contact_groups"`
	LastReminder   int                 `json:"last_reminder"`
	LastUpdatedUtc string              `json:"last_updated_utc"`
}

//PartialSsl represent  a ssl test creation or modification
type PartialSsl struct {
	ID             int
	Domain         string
	Checkrate      string
	ContactGroupsC string
	AlertAt        string
	AlertExpiry    bool
	AlertReminder  bool
	AlertBroken    bool
	AlertMixed     bool
}

type createSsl struct {
	ID             int              `url:"id,omitempty"`
	Domain         string           `url:"domain"         json:"domain"`
	Checkrate      jsonNumberString `url:"checkrate"      json:"checkrate"`
	ContactGroupsC string           `url:"contact_groups" json:"contact_groups"`
	AlertAt        string           `url:"alert_at"       json:"alert_at"`
	AlertExpiry    bool             `url:"alert_expiry"   json:"alert_expiry"`
	AlertReminder  bool             `url:"alert_reminder" json:"alert_reminder"`
	AlertBroken    bool             `url:"alert_broken"   json:"alert_broken"`
	AlertMixed     bool             `url:"alert_mixed"    json:"alert_mixed"`
}

func (cs *createSsl) fromPartial(p *PartialSsl) {
	cs.ID = p.ID
	cs.Domain = p.Domain
	cs.Checkrate = jsonNumberString(p.Checkrate)
	cs.ContactGroupsC = p.ContactGroupsC
	cs.AlertAt = p.AlertAt
	cs.AlertExpiry = p.AlertExpiry
	cs.AlertReminder = p.AlertReminder
	cs.AlertBroken = p.AlertBroken
	cs.AlertMixed = p.AlertMixed
}

func (cs *createSsl) toPartial(p *PartialSsl) {
	p.ID = cs.ID
	p.Domain = cs.Domain
	p.Checkrate = string(cs.Checkrate)
	p.ContactGroupsC = cs.ContactGroupsC
	p.AlertAt = cs.AlertAt
	p.AlertExpiry = cs.AlertExpiry
	p.AlertReminder = cs.AlertReminder
	p.AlertBroken = cs.AlertBroken
	p.AlertMixed = cs.AlertMixed
}

type updateSsl struct {
	ID             int              `url:"id"`
	Domain         string           `url:"domain"         json:"domain"`
	Checkrate      jsonNumberString `url:"checkrate"      json:"checkrate"`
	ContactGroupsC string           `url:"contact_groups" json:"contact_groups"`
	AlertAt        string           `url:"alert_at"       json:"alert_at"`
	AlertExpiry    bool             `url:"alert_expiry"   json:"alert_expiry"`
	AlertReminder  bool             `url:"alert_reminder" json:"alert_reminder"`
	AlertBroken    bool             `url:"alert_broken"   json:"alert_broken"`
	AlertMixed     bool             `url:"alert_mixed"    json:"alert_mixed"`
}

func (us *updateSsl) fromPartial(p *PartialSsl) {
	us.ID = p.ID
	us.Domain = p.Domain
	us.Checkrate = jsonNumberString(p.Checkrate)
	us.ContactGroupsC = p.ContactGroupsC
	us.AlertAt = p.AlertAt
	us.AlertExpiry = p.AlertExpiry
	us.AlertReminder = p.AlertReminder
	us.AlertBroken = p.AlertBroken
	us.AlertMixed = p.AlertMixed
}

func (us *updateSsl) toPartial(p *PartialSsl) {
	p.ID = us.ID
	p.Domain = us.Domain
	p.Checkrate = string(us.Checkrate)
	p.ContactGroupsC = us.ContactGroupsC
	p.AlertAt = us.AlertAt
	p.AlertExpiry = us.AlertExpiry
	p.AlertReminder = us.AlertReminder
	p.AlertBroken = us.AlertBroken
	p.AlertMixed = us.AlertMixed
}

type sslUpdateResponse struct {
	Success bool        `json:"Success"`
	Message interface{} `json:"Message"`
}

type sslCreateResponse struct {
	Success bool        `json:"Success"`
	Message interface{} `json:"Message"`
	Input   createSsl   `json:"Input"`
}

//Ssls represent the actions done wit the API
type Ssls interface {
	All() ([]*Ssl, error)
	completeSsl(*PartialSsl) (*Ssl, error)
	Detail(string) (*Ssl, error)
	Update(*PartialSsl) (*Ssl, error)
	UpdatePartial(*PartialSsl) (*PartialSsl, error)
	Delete(ID string) error
	CreatePartial(*PartialSsl) (*PartialSsl, error)
	Create(*PartialSsl) (*Ssl, error)
}

func consolidateSsl(s *Ssl) {
	(*s).ContactGroupsC = strings.Trim(strings.Join(strings.Fields(fmt.Sprint((*s).ContactGroups)), ","), "[]")
}

func findSsl(responses []*Ssl, id string) (*Ssl, error) {
	var response *Ssl
	for _, elem := range responses {
		if (*elem).ID == id {
			return elem, nil
		}
	}
	return response, fmt.Errorf("%s Not found", id)
}

func (tt *ssls) completeSsl(s *PartialSsl) (*Ssl, error) {
	full, err := tt.Detail(strconv.Itoa((*s).ID))
	if err != nil {
		return nil, err
	}
	(*full).ContactGroups = strings.Split((*s).ContactGroupsC, ",")
	return full, nil
}

//Partial return a PartialSsl corresponding to the Ssl
func Partial(s *Ssl) (*PartialSsl, error) {
	if s == nil {
		return nil, fmt.Errorf("s is nil")
	}
	id, err := strconv.Atoi(s.ID)
	if err != nil {
		return nil, err
	}
	return &PartialSsl{
		ID:             id,
		Domain:         s.Domain,
		Checkrate:      strconv.Itoa(s.Checkrate),
		ContactGroupsC: s.ContactGroupsC,
		AlertReminder:  s.AlertReminder,
		AlertExpiry:    s.AlertExpiry,
		AlertBroken:    s.AlertBroken,
		AlertMixed:     s.AlertMixed,
		AlertAt:        s.AlertAt,
	}, nil

}

type ssls struct {
	client apiClient
}

//NewSsls return a new ssls
func NewSsls(c apiClient) Ssls {
	return &ssls{
		client: c,
	}
}

//All return a list of all the ssl from the API
func (tt *ssls) All() ([]*Ssl, error) {
	rawResponse, err := tt.client.get("/SSL", nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting StatusCake Ssl: %s", err.Error())
	}
	var getResponse []*Ssl
	err = json.NewDecoder(rawResponse.Body).Decode(&getResponse)
	if err != nil {
		return nil, err
	}

	for ssl := range getResponse {
		consolidateSsl(getResponse[ssl])
	}

	return getResponse, err
}

//Detail return the ssl corresponding to the id
func (tt *ssls) Detail(id string) (*Ssl, error) {
	responses, err := tt.All()
	if err != nil {
		return nil, err
	}
	mySsl, errF := findSsl(responses, id)
	if errF != nil {
		return nil, errF
	}
	return mySsl, nil
}

//Update update the API with s and create one if s.ID=0 then return the corresponding Ssl
func (tt *ssls) Update(s *PartialSsl) (*Ssl, error) {
	var err error
	s, err = tt.UpdatePartial(s)
	if err != nil {
		return nil, err
	}
	return tt.completeSsl(s)
}

//UpdatePartial update the API with s and create one if s.ID=0 then return the corresponding PartialSsl
func (tt *ssls) UpdatePartial(s *PartialSsl) (*PartialSsl, error) {
	if (*s).ID == 0 {
		return tt.CreatePartial(s)
	}

	var v url.Values
	{
		us := updateSsl{}
		us.fromPartial(s)
		v, _ = query.Values(us)
	}

	rawResponse, err := tt.client.put("/SSL/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake Ssl: %s", err.Error())
	}

	var updateResponse sslUpdateResponse
	err = json.NewDecoder(rawResponse.Body).Decode(&updateResponse)
	if err != nil {
		return nil, err
	}

	if !updateResponse.Success {
		return nil, fmt.Errorf("%s", updateResponse.Message.(string))
	}

	return s, nil
}

//Delete delete the ssl which ID is id
func (tt *ssls) Delete(id string) error {
	_, err := tt.client.delete("/SSL/Update", url.Values{"id": {fmt.Sprint(id)}})
	if err != nil {
		return err
	}

	return nil
}

//Create create the ssl whith the data in s and return the Ssl created
func (tt *ssls) Create(s *PartialSsl) (*Ssl, error) {
	var err error
	s, err = tt.CreatePartial(s)
	if err != nil {
		return nil, err
	}
	return tt.completeSsl(s)
}

//CreatePartial create the ssl whith the data in s and return the PartialSsl created
func (tt *ssls) CreatePartial(s *PartialSsl) (*PartialSsl, error) {
	(*s).ID = 0
	var v url.Values
	{
		cs := createSsl{}
		cs.fromPartial(s)
		v, _ = query.Values(cs)
	}

	rawResponse, err := tt.client.put("/SSL/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake Ssl: %s", err.Error())
	}

	var createResponse sslCreateResponse
	err = json.NewDecoder(rawResponse.Body).Decode(&createResponse)
	if err != nil {
		return nil, err
	}

	if !createResponse.Success {
		return nil, fmt.Errorf("%s", createResponse.Message.(string))
	}
	createResponse.Input.toPartial(s)
	(*s).ID = int(createResponse.Message.(float64))

	return s, nil
}
