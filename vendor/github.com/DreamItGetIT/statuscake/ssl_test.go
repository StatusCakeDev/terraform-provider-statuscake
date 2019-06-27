package statuscake

import (
	"testing"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
)

func TestSsl_All(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "sslListAllOk.json",
	}
	tt := NewSsls(c)
	ssls, err := tt.All()
	require.Nil(err)

	assert.Equal("/SSL", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
	assert.Len(ssls, 3)
	mixed := make(map[string]string)
	flags := make(map[string]bool)
	flags["is_extended"] = false
	flags["has_pfs"] = true
	flags["is_broken"] = false
	flags["is_expired"] = false
	flags["is_missing"] = false
	flags["is_revoked"] = false
	flags["has_mixed"] = false
	expectedTest := &Ssl{
		ID: "143615",
		Checkrate: 2073600,
		Paused: false,
		Domain: "https://www.exemple.com",
		IssuerCn: "Let's Encrypt Authority X3",
		CertScore: "95",
		CipherScore: "100",
		CertStatus: "CERT_OK",
		Cipher: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		ValidFromUtc: "2019-05-28 01:22:00",
		ValidUntilUtc: "2019-08-26 01:22:00",
		MixedContent: []map[string]string{},
		Flags: flags,
		ContactGroups: []string{},
		ContactGroupsC: "",
		AlertAt: "7,18,2019",
		LastReminder: 2019,
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		LastUpdatedUtc: "2019-06-20 10:11:03",
	}
	assert.Equal(expectedTest, ssls[0])

	expectedTest.ID="143617"
	expectedTest.LastUpdatedUtc="2019-06-20 10:23:20"
	assert.Equal(expectedTest, ssls[2])

	expectedTest.ID="143616"
	expectedTest.LastUpdatedUtc="2019-06-20 10:23:14"
	mixed["type"]="img"
	mixed["src"]="http://example.com/image.gif"
	expectedTest.MixedContent=[]map[string]string{mixed}
	expectedTest.ContactGroupsC="12,13,34"
	expectedTest.ContactGroups=[]string{"12","13","34"}
	assert.Equal(expectedTest, ssls[1])
}

func TestSsls_Detail_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "sslListAllOk.json",
	}
	tt := NewSsls(c)

	ssl, err := tt.Detail("143616")
	require.Nil(err)
	assert.Equal("/SSL", c.sentRequestPath)
	assert.Equal("GET", c.sentRequestMethod)
	assert.Nil(c.sentRequestValues)
	
	mixed := make(map[string]string)
	flags := make(map[string]bool)

	mixed["type"]="img"
	mixed["src"]="http://example.com/image.gif"
	
	flags["is_extended"] = false
	flags["has_pfs"] = true
	flags["is_broken"] = false
	flags["is_expired"] = false
	flags["is_missing"] = false
	flags["is_revoked"] = false
	flags["has_mixed"] = false
	expectedTest := &Ssl{
		ID: "143616",
		Checkrate: 2073600,
		Paused: false,
		Domain: "https://www.exemple.com",
		IssuerCn: "Let's Encrypt Authority X3",
		CertScore: "95",
		CipherScore: "100",
		CertStatus: "CERT_OK",
		Cipher: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		ValidFromUtc: "2019-05-28 01:22:00",
		ValidUntilUtc: "2019-08-26 01:22:00",
		MixedContent: []map[string]string{mixed},
		Flags: flags,
		ContactGroups: []string{"12","13","34"},
		ContactGroupsC: "12,13,34",
		AlertAt: "7,18,2019",
		LastReminder: 2019,
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		LastUpdatedUtc: "2019-06-20 10:23:14",
	}
	
	assert.Equal(expectedTest, ssl)
}

func TestSsls_CreatePartial_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "sslCreateOk.json",
	}
	tt := NewSsls(c)
	partial := &PartialSsl{
		Domain: "https://www.exemple.com",
		Checkrate: "2073600",
		ContactGroupsC: "",
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		AlertAt: "7,18,2019",
	}
	expectedRes := &PartialSsl {
		ID: 143616,
		Domain: "https://www.exemple.com",
		Checkrate: "2073600",
		ContactGroupsC: "",
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		AlertAt: "7,18,2019",
	}
	res, err := tt.CreatePartial(partial)
	require.Nil(err)
	assert.Equal("/SSL/Update", c.sentRequestPath)
	assert.Equal("PUT", c.sentRequestMethod)
	assert.Equal(c.sentRequestValues,url.Values(url.Values{"domain":[]string{"https://www.exemple.com"}, "checkrate":[]string{"2073600"}, "contact_groups":[]string{""}, "alert_at":[]string{"7,18,2019"}, "alert_expiry":[]string{"true"}, "alert_reminder":[]string{"true"}, "alert_broken":[]string{"true"}, "alert_mixed":[]string{"true"}}))

	assert.Equal(expectedRes, res)
}

func TestSsls_UpdatePartial_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "sslUpdateOk.json",
	}
	tt := NewSsls(c)
	partial := &PartialSsl{
		ID: 143616,
		Domain: "https://www.exemple.com",
		Checkrate: "2073600",
		ContactGroupsC: "",
		AlertReminder: false,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		AlertAt: "7,18,2019",
	}
	expectedRes := &PartialSsl {
		ID: 143616,
		Domain: "https://www.exemple.com",
		Checkrate: "2073600",
		ContactGroupsC: "",
		AlertReminder: false,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		AlertAt: "7,18,2019",
	}
	res, err := tt.UpdatePartial(partial)
	require.Nil(err)
	assert.Equal(expectedRes, res)
	assert.Equal("/SSL/Update", c.sentRequestPath)
	assert.Equal("PUT", c.sentRequestMethod)
	assert.Equal(c.sentRequestValues,url.Values(url.Values{"id":[]string{"143616"},"domain":[]string{"https://www.exemple.com"}, "checkrate":[]string{"2073600"}, "contact_groups":[]string{""}, "alert_at":[]string{"7,18,2019"}, "alert_expiry":[]string{"true"}, "alert_reminder":[]string{"false"}, "alert_broken":[]string{"true"}, "alert_mixed":[]string{"true"}}))
}

func TestSsl_complete_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "sslListAllOk.json",
	}
	tt := NewSsls(c)

	partial := &PartialSsl {
		ID: 143616,
		Domain: "https://www.exemple.com",
		Checkrate: "2073600",
		ContactGroupsC: "12,13,34",
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		AlertAt: "7,18,2019",
	}
	full, err := tt.completeSsl(partial)
	require.Nil(err)
	mixed := make(map[string]string)
	flags := make(map[string]bool)

	mixed["type"]="img"
	mixed["src"]="http://example.com/image.gif"
	
	flags["is_extended"] = false
	flags["has_pfs"] = true
	flags["is_broken"] = false
	flags["is_expired"] = false
	flags["is_missing"] = false
	flags["is_revoked"] = false
	flags["has_mixed"] = false
	expectedTest := &Ssl{
		ID: "143616",
		Checkrate: 2073600,
		Paused: false,
		Domain: "https://www.exemple.com",
		IssuerCn: "Let's Encrypt Authority X3",
		CertScore: "95",
		CipherScore: "100",
		CertStatus: "CERT_OK",
		Cipher: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		ValidFromUtc: "2019-05-28 01:22:00",
		ValidUntilUtc: "2019-08-26 01:22:00",
		MixedContent: []map[string]string{mixed},
		Flags: flags,
		ContactGroups: []string{"12","13","34"},
		ContactGroupsC: "12,13,34",
		AlertAt: "7,18,2019",
		LastReminder: 2019,
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		LastUpdatedUtc: "2019-06-20 10:23:14",
	}
	
	assert.Equal(expectedTest, full)
	
}

func TestSsl_partial_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	
	mixed := make(map[string]string)
	flags := make(map[string]bool)

	mixed["type"]="img"
	mixed["src"]="http://example.com/image.gif"
	
	flags["is_extended"] = false
	flags["has_pfs"] = true
	flags["is_broken"] = false
	flags["is_expired"] = false
	flags["is_missing"] = false
	flags["is_revoked"] = false
	flags["has_mixed"] = false
	full := &Ssl{
		ID: "143616",
		Checkrate: 2073600,
		Paused: false,
		Domain: "https://www.exemple.com",
		IssuerCn: "Let's Encrypt Authority X3",
		CertScore: "95",
		CipherScore: "100",
		CertStatus: "CERT_OK",
		Cipher: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		ValidFromUtc: "2019-05-28 01:22:00",
		ValidUntilUtc: "2019-08-26 01:22:00",
		MixedContent: []map[string]string{mixed},
		Flags: flags,
		ContactGroups: []string{"12","13","34"},
		ContactGroupsC: "12,13,34",
		AlertAt: "7,18,2019",
		LastReminder: 2019,
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		LastUpdatedUtc: "2019-06-20 10:23:14",
	}
	expectedTest:=&PartialSsl {
		ID: 143616,
		Domain: "https://www.exemple.com",
		Checkrate: "2073600",
		ContactGroupsC: "12,13,34",
		AlertReminder: true,
		AlertExpiry: true,
		AlertBroken: true,
		AlertMixed: true,
		AlertAt: "7,18,2019",
	}
	partial,err:=Partial(full)
	require.Nil(err)
	assert.Equal(expectedTest, partial)
	
}

func TestSsls_Delete_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := &fakeAPIClient{
		fixture: "sslDeleteOk.json",
	}
	tt := NewSsls(c)

	err := tt.Delete("143616")
	require.Nil(err)
	assert.Equal("/SSL/Update", c.sentRequestPath)
	assert.Equal("DELETE", c.sentRequestMethod)
	assert.Equal(c.sentRequestValues,url.Values(url.Values{"id":[]string{"143616"},},))
}
