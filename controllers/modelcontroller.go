package controllers

import (
	"context"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

// UpdateModel Update a model by its ID
func UpdateModel(ctx context.Context, id string, model ModelRequest) error {
	res := &ardoq.Model{}
	errResponse := new(ardoq.Error)
	model.ID = id
	resp, err := client().Patch("model/"+id).
		BodyProvider(ardoqBodyProvider{request: model}).
		Receive(res, errResponse)
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
