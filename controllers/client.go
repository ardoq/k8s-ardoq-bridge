package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	UserAgentPrefix = "k8s-ardoq-bridge"
	Version         = "0.0.0"
)

type ardoqDecoder struct {
}

func (a ardoqDecoder) Decode(resp *http.Response, v interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	var data interface{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	// check if StatusCode is OK, if not add StatusCode to the Error
	// so that ErrorReponse.Errors() return true, and the actual decoded response get shown in terraform
	if code := resp.StatusCode; 200 <= code && code <= 299 {
	} else {
		// apparently there's an error,
		var errResponse *ardoq.Error = v.(*ardoq.Error)
		errResponse.Code = resp.StatusCode
	}

	return mapstructure.WeakDecode(data, v)
}

type ardoqBodyProvider struct {
	request interface{}
	//fields  interface{}
}

func (a ardoqBodyProvider) Body() (io.Reader, error) {
	// request := a.request.(ComponentRequest)

	// marshal component
	requestJSON, _ := json.Marshal(a.request)

	// create new map as destination for both Unmarshal methods to combine the data
	flatRequest := make(map[string]interface{})
	err := json.Unmarshal(requestJSON, &flatRequest)
	if err != nil {
		return nil, err
	}

	//if len(a.fields.(map[string]interface{})) > 0 {
	//	// marshal component.Fields
	//	fieldsJSON, _ := json.Marshal(a.fields)
	//	err = json.Unmarshal(fieldsJSON, &flatRequest)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(flatRequest)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (a ardoqBodyProvider) ContentType() string {
	return "application/json"
}
func client() *sling.Sling {
	type OrgSearchQuery struct {
		Org string `url:"org,omitempty"`
	}

	return sling.New().Base(baseUri).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, Version)).
		Set("Authorization", fmt.Sprintf("Token token=%s", apiKey)).ResponseDecoder(ardoqDecoder{}).
		QueryStruct(&OrgSearchQuery{Org: org})
}
func RestyClient() *resty.Request {
	return resty.New().SetBaseURL(baseUri).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}).R().
		SetAuthToken(apiKey).
		SetQueryParam("org", org)
}

type BodyProvider struct {
	request interface{}
	fields  interface{}
}

func (a BodyProvider) Body() io.Reader {
	requestJSON, _ := json.Marshal(a.request)
	flatRequest := make(map[string]interface{})
	err := json.Unmarshal(requestJSON, &flatRequest)
	if err != nil {
		log.Error(err)
		return nil
	}

	if len(a.fields.(map[string]interface{})) > 0 {
		// marshal component.Fields
		fieldsJSON, _ := json.Marshal(a.fields)
		err = json.Unmarshal(fieldsJSON, &flatRequest)
		if err != nil {
			log.Error(err)
			return nil
		}
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(flatRequest)
	if err != nil {
		log.Error(err)
		return nil
	}
	return buf
}
