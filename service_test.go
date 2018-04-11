package main

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"os"
	"image/jpeg"
	"bytes"
	"net/http"
	"fmt"
)

func TestSuccessLoadImage(t *testing.T) {
	expectedImg := createImgTest()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeImageTest(w, &expectedImg)
	}))

	defer ts.Close()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(dir)

	client := NewDownloadClient(1, dir)

	service := NewService(client)

	expectedImgWidth := uint(expectedImg.Bounds().Max.X)
	expectedImgHeight := uint(expectedImg.Bounds().Max.Y)

	data, err := service.LoadImage(ts.URL, expectedImgWidth, expectedImgHeight)

	if err != nil {
		t.Error(err)
	}

	haveImg, err := jpeg.Decode(bytes.NewReader(data))

	haveImgWidth := uint(haveImg.Bounds().Max.X)
	haveImgHeight := uint(haveImg.Bounds().Max.Y)

	fmt.Println("expectedImgWidth: ", expectedImgWidth, "expectedImgHeight: ", expectedImgHeight)
	fmt.Println("haveImgWidth: ", haveImgWidth, "haveImgHeight: ", haveImgHeight)

	if expectedImgWidth != haveImgWidth || expectedImgHeight != haveImgHeight {
		t.Errorf("bounds should be equal")
	}
}

func TestFailedLoadImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer ts.Close()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(dir)

	client := NewDownloadClient(1, dir)

	service := NewService(client)

	_, err = service.LoadImage(ts.URL, 240, 240)

	if imgErr, ok := err.(ImageError); !ok {
		t.Errorf("error should be ImageError")
	} else if imgErr.Status() != http.StatusNotFound {
		t.Errorf("status code should be equal")
	}
}

func TestResizeLoadImage(t *testing.T) {
	expectedImg := createImgTest()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeImageTest(w, &expectedImg)
	}))

	defer ts.Close()

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(dir)

	client := NewDownloadClient(1, dir)

	service := NewService(client)

	expectedImgWidth := uint(expectedImg.Bounds().Max.X)
	expectedImgHeight := uint(expectedImg.Bounds().Max.Y)

	data, err := service.LoadImage(ts.URL, expectedImgWidth * 2, expectedImgHeight * 2)

	if err != nil {
		t.Error(err)
	}

	haveImg, err := jpeg.Decode(bytes.NewReader(data))

	haveImgWidth := uint(haveImg.Bounds().Max.X)
	haveImgHeight := uint(haveImg.Bounds().Max.Y)

	fmt.Println("expectedImgWidth: ", expectedImgWidth, "expectedImgHeight: ", expectedImgHeight)
	fmt.Println("haveImgWidth x2: ", haveImgWidth, "haveImgHeight x2: ", haveImgHeight)

	if expectedImgWidth * 2 != haveImgWidth || expectedImgHeight * 2 != haveImgHeight {
		t.Errorf("bounds should be equal")
	}
}

func TestResize(t *testing.T) {
	expectedImg := createImgTest()

	expectedImgWidth := uint(expectedImg.Bounds().Max.X)
	expectedImgHeight := uint(expectedImg.Bounds().Max.Y)

	buf := new(bytes.Buffer)

	err := jpeg.Encode(buf, expectedImg, nil)

	if err != nil {
		t.Error(err)
	}

	data, err := resizeImg(buf.Bytes(), expectedImgWidth * 2, expectedImgHeight * 2)

	if err != nil {
		t.Error(err)
	}

	haveImg, err := jpeg.Decode(bytes.NewReader(data))

	haveImgWidth := uint(haveImg.Bounds().Max.X)
	haveImgHeight := uint(haveImg.Bounds().Max.Y)

	fmt.Println("expectedImgWidth: ", expectedImgWidth, "expectedImgHeight: ", expectedImgHeight)
	fmt.Println("haveImgWidth x2: ", haveImgWidth, "haveImgHeight x2: ", haveImgHeight)

	if expectedImgWidth * 2 != haveImgWidth || expectedImgHeight * 2 != haveImgHeight {
		t.Errorf("bounds should be equal")
	}
}

func TestNotResize(t *testing.T) {
	expectedImg := createImgTest()

	expectedImgWidth := uint(expectedImg.Bounds().Max.X)
	expectedImgHeight := uint(expectedImg.Bounds().Max.Y)

	buf := new(bytes.Buffer)

	err := jpeg.Encode(buf, expectedImg, nil)

	if err != nil {
		t.Error(err)
	}

	data, err := resizeImg(buf.Bytes(), expectedImgWidth, expectedImgHeight)

	if err != nil {
		t.Error(err)
	}

	haveImg, err := jpeg.Decode(bytes.NewReader(data))

	haveImgWidth := uint(haveImg.Bounds().Max.X)
	haveImgHeight := uint(haveImg.Bounds().Max.Y)

	fmt.Println("expectedImgWidth: ", expectedImgWidth, "expectedImgHeight: ", expectedImgHeight)
	fmt.Println("haveImgWidth: ", haveImgWidth, "haveImgHeight: ", haveImgHeight)

	if expectedImgWidth != haveImgWidth || expectedImgHeight != haveImgHeight {
		t.Errorf("bounds should be equal")
	}
}