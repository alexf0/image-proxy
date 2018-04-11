package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"fmt"
	"bytes"
	"net/http"
	"encoding/base64"
	"strings"
)

type loadImageRequest struct {
	Url string
	Width uint
	Height uint
}

func getInvalidErr(fields ... string) string {
	var buf bytes.Buffer
	for idx, field := range fields {
		buf.WriteString(field + " is blank")
		if idx != len(fields) - 1 {
			buf.WriteString(", ")
		}
	}

	return fmt.Sprintf("Can't validate request; %s", buf.String())
}

func validateRequest(req *loadImageRequest) (bool, error) {
	var fields []string
	if req.Url == "" {
		fields = append(fields, "url")
	} else {
		URL, err := base64.RawURLEncoding.DecodeString(req.Url)

		if err != nil {
			return false, imageError{http.StatusBadRequest,"invalid filename encoding"}
		}

		parts := strings.Split(string(URL), ".")
		last := parts[len(parts) - 1]

		if last != "jpg" && last != "jpeg" {
			return false, imageError{http.StatusBadRequest,"does not support image type"}
		}

		req.Url = string(URL)
	}

	if req.Width == 0 {
		fields = append(fields, "width")
	}

	if req.Height == 0 {
		fields = append(fields, "height")
	}

	if len(fields) == 0 {
		return true, nil
	}

	return false, imageError{http.StatusBadRequest, getInvalidErr(fields...)}
}

func makeLoadImageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*loadImageRequest)

		if ok, err := validateRequest(req); !ok  {
			return nil, err
		}

		data, err := s.LoadImage(req.Url, req.Height, req.Width)

		if err != nil {
			return nil, err
		}

		return data, nil
	}
}