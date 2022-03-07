package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Error(err)
		return errors.Wrap(err, "could not get model")
	}
	if errResponse.NotOk() {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		errResponse.Code = resp.StatusCode
		return errResponse
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	return nil
}
