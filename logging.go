package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) LoadImage(url string, width, height uint) (bytes []byte, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "LoadImage",
			"url", url,
			"width", width,
			"height", height,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LoadImage(url, width, height)
}