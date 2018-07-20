package statuscake

import (
	"strings"
)

type autheticationErrorResponse struct {
	ErrNo int
	Error string
}

type updateResponse struct {
	Issues   interface{} `json:"Issues"`
	Success  bool        `json:"Success"`
	Message  string      `json:"Message"`
	InsertID int         `json:"InsertID"`
}

type deleteResponse struct {
	Success bool   `json:"Success"`
	Error   string `json:"Error"`
}

type detailResponse struct {
	Method          string   `json:"Method"`
	TestID          int      `json:"TestID"`
	TestType        string   `json:"TestType"`
	Paused          bool     `json:"Paused"`
	WebsiteName     string   `json:"WebsiteName"`
	URI             string   `json:"URI"`
	ContactID       int      `json:"ContactID"`
	Status          string   `json:"Status"`
	Uptime          float64  `json:"Uptime"`
	CustomHeader    string   `json:"CustomHeader"`
	UserAgent       string   `json:"UserAgent"`
	CheckRate       int      `json:"CheckRate"`
	Timeout         int      `json:"Timeout"`
	LogoImage       string   `json:"LogoImage"`
	Confirmation    int      `json:"Confirmation,string"`
	WebsiteHost     string   `json:"WebsiteHost"`
	NodeLocations   []string `json:"NodeLocations"`
	FindString      string   `json:"FindString"`
	DoNotFind       bool     `json:"DoNotFind"`
	LastTested      string   `json:"LastTested"`
	NextLocation    string   `json:"NextLocation"`
	Port            int      `json:"Port"`
	Processing      bool     `json:"Processing"`
	ProcessingState string   `json:"ProcessingState"`
	ProcessingOn    string   `json:"ProcessingOn"`
	DownTimes       int      `json:"DownTimes,string"`
	Sensitive       bool     `json:"Sensitive"`
	TriggerRate     int      `json:"TriggerRate,string"`
	UseJar          int      `json:"UseJar"`
	PostRaw         string   `json:"PostRaw"`
	FinalEndpoint   string   `json:"FinalEndpoint"`
	FollowRedirect  bool     `json:"FollowRedirect"`
	StatusCodes     []string `json:"StatusCodes"`
}

func (d *detailResponse) test() *Test {
	return &Test{
		TestID:         d.TestID,
		TestType:       d.TestType,
		Paused:         d.Paused,
		WebsiteName:    d.WebsiteName,
		WebsiteURL:     d.URI,
		CustomHeader:   d.CustomHeader,
		UserAgent:      d.UserAgent,
		ContactID:      d.ContactID,
		Status:         d.Status,
		Uptime:         d.Uptime,
		CheckRate:      d.CheckRate,
		Timeout:        d.Timeout,
		LogoImage:      d.LogoImage,
		Confirmation:   d.Confirmation,
		WebsiteHost:    d.WebsiteHost,
		NodeLocations:  d.NodeLocations,
		FindString:     d.FindString,
		DoNotFind:      d.DoNotFind,
		Port:           d.Port,
		TriggerRate:    d.TriggerRate,
		UseJar:         d.UseJar,
		PostRaw:        d.PostRaw,
		FinalEndpoint:  d.FinalEndpoint,
		FollowRedirect: d.FollowRedirect,
		StatusCodes:    strings.Join(d.StatusCodes[:], ","),
	}
}
