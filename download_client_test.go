package main

import (
	"net/http/httptest"
	"testing"
	"net/http"
	"io/ioutil"
	"os"
	"image/draw"
	"bytes"
	"image/jpeg"
	"strconv"
	"image"
	"image/color"
	"log"
	"fmt"
)

func TestSuccessDownload(t *testing.T) {
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

	data, err := downloadImage(ts.URL, client)

	if err != nil {
		t.Error(err)
	}

	haveImg, err := jpeg.Decode(bytes.NewReader(data))

	expectedImgWidth := uint(expectedImg.Bounds().Max.X)
	expectedImgHeight := uint(expectedImg.Bounds().Max.Y)

	haveImgWidth := uint(haveImg.Bounds().Max.X)
	haveImgHeight := uint(haveImg.Bounds().Max.Y)

	fmt.Println("expectedImgWidth: ", expectedImgWidth, "expectedImgHeight: ", expectedImgHeight)
	fmt.Println("haveImgWidth: ", haveImgWidth, "haveImgHeight: ", haveImgHeight)

	if expectedImgWidth != haveImgWidth || expectedImgHeight != haveImgHeight {
		t.Errorf("bounds should be equal")
	}
}

func TestFailedDownload(t *testing.T) {
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

	_, err = downloadImage(ts.URL, client)

	if imgErr, ok := err.(ImageError); !ok {
		t.Errorf("error should be ImageError")
	} else if imgErr.Status() != http.StatusNotFound {
		t.Errorf("status code should be equal")
	}
}

func writeImageTest(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func createImgTest() image.Image {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
	return m
}
