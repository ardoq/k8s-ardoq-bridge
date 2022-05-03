package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

func ardRestClient() *ardoq.APIClient {
	a, err := ardoq.NewRestClient(baseUri, apiKey, org, "v0.0.0")
	if err != nil {
		fmt.Printf("cannot create new restclient %s", err)
		os.Exit(1)
	}
	return a
}

func RestyClient() *resty.Request {
	requestClient := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests || r.StatusCode() == http.StatusBadGateway || r.StatusCode() == http.StatusGatewayTimeout
			},
		).
		SetBaseURL(baseUri).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}).R().
		SetAuthToken(apiKey).
		SetQueryParam("org", org).
		SetError(new(HttpError))

	return requestClient
}

func Decode(resp []byte, v interface{}) error {
	var data interface{}
	err := json.Unmarshal(resp, &data)
	if err != nil {
		return err
	}
	return mapstructure.WeakDecode(data, v)
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
	if a.fields != nil {
		if len(a.fields.(map[string]interface{})) > 0 {
			// marshal component.Fields
			fieldsJSON, _ := json.Marshal(a.fields)
			err = json.Unmarshal(fieldsJSON, &flatRequest)
			if err != nil {
				log.Error(err)
				return nil
			}
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
