package statuscake

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

//ContactGroup represent the data received by the API with GET
type ContactGroup struct {
	GroupName    string   `json:"GroupName"    url:"GroupName,omitempty"`
	Emails       []string `json:"Emails"`
	EmailsPut    string   `url:"Email,omitempty"`
	Mobiles      string   `json:"Mobiles"      url:"Mobile,omitempty"`
	Boxcar       string   `json:"Boxcar"       url:"Boxcar,omitempty"`
	Pushover     string   `json:"Pushover"     url:"Pushover,omitempty"`
	ContactID    int      `json:"ContactID"    url:"ContactID,omitempty"`
	DesktopAlert string   `json:"DesktopAlert" url:"DesktopAlert,omitempty"`
	PingURL      string   `json:"PingURL"      url:"PingURL,omitempty"`
}

//Response represent the data received from the API
type Response struct {
	Success  bool   `json:"Success"`
	Message  string `json:"Message"`
	InsertID int    `json:"InsertID"`
}

//ContactGroups represent the actions done with the API
type ContactGroups interface {
	All() ([]*ContactGroup, error)
	Detail(int) (*ContactGroup, error)
	Update(*ContactGroup) (*ContactGroup, error)
	Delete(int) error
	Create(*ContactGroup) (*ContactGroup, error)
}

func findContactGroup(responses []*ContactGroup, id int) (*ContactGroup, error) {
	var response *ContactGroup
	for _, elem := range responses {
		if (*elem).ContactID == id {
			return elem, nil
		}
	}
	return response, fmt.Errorf("%d Not found", id)
}

type contactGroups struct {
	client apiClient
}

//NewContactGroups return a new ssls
func NewContactGroups(c apiClient) ContactGroups {
	return &contactGroups{
		client: c,
	}
}

//All return a list of all the ContactGroup from the API
func (tt *contactGroups) All() ([]*ContactGroup, error) {
	rawResponse, err := tt.client.get("/ContactGroups", nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting StatusCake contactGroups: %s", err.Error())
	}
	var getResponse []*ContactGroup
	err = json.NewDecoder(rawResponse.Body).Decode(&getResponse)
	if err != nil {
		return nil, err
	}
	return getResponse, err
}

//Detail return the ContactGroup corresponding to the id
func (tt *contactGroups) Detail(id int) (*ContactGroup, error) {
	responses, err := tt.All()
	if err != nil {
		return nil, err
	}
	myContactGroup, errF := findContactGroup(responses, id)
	if errF != nil {
		return nil, errF
	}
	return myContactGroup, nil
}

//Update update the API with cg and create one if cg.ContactID=0 then return the corresponding ContactGroup
func (tt *contactGroups) Update(cg *ContactGroup) (*ContactGroup, error) {

	if cg.ContactID == 0 {
		return tt.Create(cg)
	}
	cg.EmailsPut = strings.Join(cg.Emails, ",")
	var v url.Values

	v, _ = query.Values(*cg)

	rawResponse, err := tt.client.put("/ContactGroups/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake ContactGroup: %s", err.Error())
	}

	var response Response
	err = json.NewDecoder(rawResponse.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("%s", response.Message)
	}

	return cg, nil
}

//Delete delete the ContactGroup which ID is id
func (tt *contactGroups) Delete(id int) error {
	_, err := tt.client.delete("/ContactGroups/Update", url.Values{"ContactID": {fmt.Sprint(id)}})
	return err
}

//CreatePartial create the ContactGroup whith the data in cg and return the ContactGroup created
func (tt *contactGroups) Create(cg *ContactGroup) (*ContactGroup, error) {
	cg.ContactID = 0
	cg.EmailsPut = strings.Join(cg.Emails, ",")
	var v url.Values
	v, _ = query.Values(*cg)

	rawResponse, err := tt.client.put("/ContactGroups/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake ContactGroup: %s", err.Error())
	}

	var response Response
	err = json.NewDecoder(rawResponse.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("%s", response.Message)
	}

	cg.ContactID = response.InsertID

	return cg, nil
}
