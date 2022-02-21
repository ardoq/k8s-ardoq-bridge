package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"time"
)

// CreateFields a model by its ID
func CreateFields(id string, fields []FieldRequest) error {
	res := &[]ardoq.Field{}
	errResponse := new(ardoq.Error)
	fields = completeFields(id, fields)
	for _, v := range fields {
		requestStarted := time.Now()
		resp, err := client().Post("field").
			BodyProvider(ardoqBodyProvider{request: v}).
			Receive(res, errResponse)
		metrics.RequestLatency.WithLabelValues("update").Observe(time.Since(requestStarted).Seconds())
		if errResponse.Code == 409 {
			continue
		} else if err != nil {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			klog.Error(err)
			return errors.Wrap(err, "could not create field")
		}
		if errResponse.NotOk() {
			metrics.RequestStatusCode.WithLabelValues("error").Inc()
			errResponse.Code = resp.StatusCode
			return errResponse
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
