package main

import (
	"context"
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
	"fmt"
	"crypto/sha1"
	"time"
)

var ErrEncodeResponse = errors.New("encoded response with error")

func MakeHandler(sc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(before),
	}

	loadParamsHandler := kithttp.NewServer(
		makeLoadImageEndpoint(sc),
		decodeLoadImageRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/", loadParamsHandler).Methods("GET")

	return r
}

func parseQuery(r *http.Request) (*loadImageRequest, error) {
	vars := r.URL.Query()

	req := &loadImageRequest{}
	if s := vars.Get("width"); s != "" {
		w, err := strconv.ParseUint(s, 10, 64)

		if err != nil {
			return nil, imageError{http.StatusBadRequest, err.Error()}
		}
		req.Width = uint(w)
	}

	if s := vars.Get("height"); s != "" {
		h, err := strconv.ParseUint(s, 10, 64)

		if err != nil {
			return nil, imageError{http.StatusBadRequest, err.Error()}
		}
		req.Height = uint(h)
	}

	req.Url = vars.Get("url")

	return req, nil
}

func decodeLoadImageRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return parseQuery(r)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok {
		encodeError(ctx, e, w)
		return nil
	}

	if data, ok := response.([]byte); ok {
		w.Header().Set("Content-Type", "image/jpeg")

		generateEtag := generate(data, true)

		w.Header().Set("Etag", generateEtag)
		w.Header().Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", 3600))

		if val := ctx.Value("etag"); val != nil {
			etag := val.(string)

			if etag == generateEtag {
				w.WriteHeader(http.StatusNotModified)
				return nil
			}
		}

		w.Write(data)

		return nil
	}

	return ErrEncodeResponse
}

func getHash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func generate(data []byte, weak bool) string {
	tag := fmt.Sprintf("\"%d-%s\"", len(data), getHash(data))
	if weak {
		tag = "W/" + tag
	}

	return tag
}

type ImageError interface {
	Status() int
}

type imageError struct {
	StatusCode int
	Err        string
}

func (err imageError) Status() int {
	return err.StatusCode
}

func (err imageError) Error() string {
	return err.Err
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err, ok := err.(ImageError); ok {
		w.WriteHeader(err.Status())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}

func before(ctx context.Context, r *http.Request) context.Context {
	etag := r.Header.Get("If-None-Match")
	ctx = context.WithValue(ctx, "etag", etag)
	return ctx
}