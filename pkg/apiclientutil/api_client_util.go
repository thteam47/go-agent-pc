package apiclientutil

import (
	"encoding/json"

	"github.com/thteam47/common-libs/errutil"
)

type GrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(data []byte) error {
	errMessage := string(data)
	grpcError := GrpcError{}
	if err := json.Unmarshal(data, &grpcError); err == nil && grpcError.Message != "" {
		errMessage = grpcError.Message
	}
	return errutil.NewWithMessage(errMessage)
}

func NormalizeError(err error) error {
	if err == nil {
		return nil
	}
	return err
}
