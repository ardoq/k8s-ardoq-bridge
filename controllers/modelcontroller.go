package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"time"
)

// UpdateModel Update a model by its ID
func UpdateModel(id string, model ModelRequest) error {
	res := &ardoq.Model{}
	errResponse := new(ardoq.Error)
	model.ID = id
	requestStarted := time.Now()
	resp, err := client().Patch("model/"+id).
		BodyProvider(ardoqBodyProvider{request: model}).
		Receive(res, errResponse)
	metrics.RequestLatency.WithLabelValues("update").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		klog.Error(err)
		return errors.Wrap(err, "could not get model")
	}
	if errResponse.NotOk() {
		errResponse.Code = resp.StatusCode
		return errResponse
	}
	return nil
}
