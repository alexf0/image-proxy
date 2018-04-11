package main

import (
	"testing"
	"net/http/httptest"
	"encoding/base64"
	"fmt"
)

func TestParseQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "http://test.com", nil)

	if _, err := parseQuery(req); err != nil {
		t.Errorf("should be not equal nil")
	}

	req = httptest.NewRequest("GET", "http://test.com?width=", nil)

	if _, err := parseQuery(req); err != nil {
		t.Errorf("should be not equal nil")
	}

	req = httptest.NewRequest("GET", "http://test.com?width=-100", nil)

	if _, err := parseQuery(req); err == nil {
		t.Errorf("should be not equal nil")
	}

	req = httptest.NewRequest("GET", "http://test.com?height=", nil)

	if _, err := parseQuery(req); err != nil {
		t.Errorf("should be not equal nil")
	}

	req = httptest.NewRequest("GET", "http://test.com?height=-100", nil)

	if _, err := parseQuery(req); err == nil {
		t.Errorf("should be not equal nil")
	}

	req = httptest.NewRequest("GET", "http://test.com?url=", nil)

	if _, err := parseQuery(req); err != nil {
		t.Errorf("should be not equal nil")
	}

	url := []byte("/example0/example-1/example_2/image.jpg")
	encodeUrl := base64.RawURLEncoding.EncodeToString(url)
	fullUrl := "http://test.com?url=" + encodeUrl + "&width=480&height=480"

	req = httptest.NewRequest("GET", fullUrl, nil)
	if imgReq, _ := parseQuery(req); imgReq == nil {
		t.Errorf("should be not equal nil")
	} else if imgReq.Url != encodeUrl {
		fmt.Println(imgReq.Url)
		t.Errorf("should be equal")
	} else if imgReq.Width != uint(480) {
		t.Errorf("should be equal")
	} else if imgReq.Height != uint(480) {
		t.Errorf("should be equal")
	}
}