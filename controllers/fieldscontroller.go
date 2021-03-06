package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// CreateFields a model by its ID
func CreateFields(id string, fields []FieldRequest) error {
	fields = completeFields(id, fields)
	for _, v := range fields {
		requestStarted := time.Now()
		resp, err := RestyClient().SetBody(BodyProvider{
			request: v,
		}.Body()).Post("field")
		metrics.RequestLatency.WithLabelValues("update").Observe(time.Since(requestStarted).Seconds())
		if resp.StatusCode() == 409 {
			continue
		} else if err != nil {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			log.Error(err)
			return errors.Wrap(err, "could not create field")
		}
		metrics.RequestStatusCode.WithLabelValues("success").Inc()
	}
	return nil
}
func completeFields(modelId string, fields []FieldRequest) []FieldRequest {
	for k, v := range fields {
		if len(v.ComponentType) > 0 {
			for x, y := range v.ComponentType {
				fields[k].ComponentType[x] = lookUpTypeId(y)
			}
		}
		fields[k].Model = modelId
	}
	return fields
}
