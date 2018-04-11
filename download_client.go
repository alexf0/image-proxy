package main

import (
	"net/http"
	"net"
	"time"
	"io/ioutil"
	"fmt"
)

func getUnknownErr(status int, body []byte) string {
	return fmt.Sprintf("Can't download image; Status: %d; %s", status, string(body))
}

func NewDownloadClient(downloadTimeout int, dir string) *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if dir != "" {
		transport.RegisterProtocol("local", http.NewFileTransport(http.Dir(dir)))
	}

	return &http.Client{
		Timeout:   time.Duration(downloadTimeout) * time.Second,
		Transport: transport,
	}
}

func downloadImage(url string, client *http.Client) ([]byte, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, imageError{res.StatusCode,getUnknownErr(res.StatusCode, body)}
	}

	return body, nil
}