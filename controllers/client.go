package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"io"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func RestyClient() *resty.Request {

	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests || r.StatusCode() == http.StatusBadGateway || r.StatusCode() == http.StatusGatewayTimeout
			},
		).
		SetBaseURL(getBaseUri()).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"x-org":        getOrg(),
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// Log all requests at trace level for debugging

			log.Tracef("%s %s -> HTTP %d (%s)",
				resp.Request.Method,
				resp.Request.URL,
				resp.StatusCode(),
				resp.Status())

			// Log response body, trimming if too large
			if len(resp.Body()) > 0 {
				body := string(resp.Body())
				if len(body) > 1000 {
					body = body[:1000] + "... (truncated)"
				}
				log.Tracef("Response body: %s", body)
			}

			// Check if the HTTP status code indicates an error
			// Note: 404 and 409 are not treated as errors (404 checks existence, 409 is conflict/already exists)
			if resp.StatusCode() >= 400 && resp.StatusCode() != http.StatusNotFound && resp.StatusCode() != http.StatusConflict {
				// Always log errors regardless of log level
				log.Errorf("HTTP Error: %s %s -> %d (%s)",
					resp.Request.Method,
					resp.Request.URL,
					resp.StatusCode(),
					resp.Status())

				// Log response body for errors, trimming if too large
				if len(resp.Body()) > 0 {
					body := string(resp.Body())
					if len(body) > 1000 {
						body = body[:1000] + "... (truncated)"
					}
					log.Errorf("Error response body: %s", body)
				}

				// Extract error message if available
				errorMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), resp.Status())
				if resp.Error() != nil {
					if httpErr, ok := resp.Error().(*HttpError); ok && httpErr.Message != "" {
						errorMsg = fmt.Sprintf("%s - %s", errorMsg, httpErr.Message)
					}
				}
				return fmt.Errorf("%s", errorMsg)
			}
			return nil
		})

	requestClient := client.R().SetAuthToken(getApiKey()).SetError(new(HttpError))

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
