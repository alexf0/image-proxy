package main

import (
	"testing"
	"encoding/base64"
)

func TestValidate(t *testing.T) {
	req := &loadImageRequest{}

	if ok, err := validateRequest(req); ok {
		t.Errorf("should be equal false")
	} else if err == nil {
		t.Errorf("should be not equal nil")
	}

	url := []byte("/example0/example-1/example_2/image.jpg")
	encodeUrl := base64.RawURLEncoding.EncodeToString(url)

	req = &loadImageRequest{encodeUrl, 240, 240}

	if ok, _ := validateRequest(req); !ok {
		t.Errorf("should be equal true")
	}

	url = []byte("/example0/example-1/example_2/image.jpeg")
	encodeUrl = base64.RawURLEncoding.EncodeToString(url)

	req = &loadImageRequest{encodeUrl, 240, 240}

	if ok, _ := validateRequest(req); !ok {
		t.Errorf("should be equal true")
	}

	url = []byte("/example0/example-1/example_2/image.png")
	encodeUrl = base64.RawURLEncoding.EncodeToString(url)

	req = &loadImageRequest{encodeUrl, 240, 240}

	if ok, err := validateRequest(req); ok {
		t.Errorf("should be equal false")
	} else if err == nil {
		t.Errorf("should be not equal nil")
	}

	url = []byte("/example0/example-1/example_2/image.jpg")
	encodeUrl = base64.RawURLEncoding.EncodeToString(url)

	req = &loadImageRequest{encodeUrl, 0, 240}

	if ok, err := validateRequest(req); ok {
		t.Errorf("should be equal false")
	} else if err == nil {
		t.Errorf("should be not equal nil")
	}

	url = []byte("/example0/example-1/example_2/image.jpg")
	encodeUrl = base64.RawURLEncoding.EncodeToString(url)

	req = &loadImageRequest{encodeUrl, 240, 0}

	if ok, err := validateRequest(req); ok {
		t.Errorf("should be equal false")
	} else if err == nil {
		t.Errorf("should be not equal nil")
	}

	req = &loadImageRequest{"", 240, 240}

	if ok, err := validateRequest(req); ok {
		t.Errorf("should be equal false")
	} else if err == nil {
		t.Errorf("should be not equal nil")
	}
}