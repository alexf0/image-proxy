package main

import (
	"net/http"
	"github.com/nfnt/resize"
	"image/jpeg"
	"bytes"
)

type Service interface {
	LoadImage(url string, width, height uint) ([]byte, error)
}

type service struct {
	client *http.Client
}

func (s *service) LoadImage(url string, width, height uint) ([]byte, error) {
	data, err := downloadImage(url, s.client)

	if err != nil {
		return nil, err
	}

	return resizeImg(data, width, height)
}

func resizeImg(data []byte, width, height uint) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	w := uint(img.Bounds().Max.X)
	h := uint(img.Bounds().Max.Y)

	if w == width && h == height {
		return data, nil
	}

	newImg := resize.Resize(width, height, img, resize.Lanczos3)

	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, newImg, nil)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewService(client *http.Client) Service {
	return &service{client: client}
}