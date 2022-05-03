package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"io"
)

func RestyClient() *resty.Request {
	client := resty.New()
	request := client.SetBaseURL(baseUri).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}).R().
		SetAuthToken(apiKey).
		SetQueryParam("org", org).
		SetError(new(HttpError))

	return request
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
