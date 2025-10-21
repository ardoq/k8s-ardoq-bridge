package controllers

import (
	"K8SArdoqBridge/app/lib/metrics"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// UpdateModel Update a model by its ID
func UpdateModel(id string, model ModelRequest) error {
	model.ID = id
	requestStarted := time.Now()
	_, err := RestyClient().SetBody(BodyProvider{
		request: model,
	}.Body()).Patch("model/" + id)
	metrics.RequestLatency.WithLabelValues("update").Observe(time.Since(requestStarted).Seconds())
	if err != nil {
		metrics.RequestStatusCode.WithLabelValues("error").Inc()
		log.Error(err)
		return errors.Wrap(err, "could not get model")
	}
	metrics.RequestStatusCode.WithLabelValues("success").Inc()
	return nil
}
