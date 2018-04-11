package main

import (
	"flag"
	"os"
	"syscall"
	"fmt"
	"net/http"
	"os/signal"
	"github.com/go-kit/kit/log"
	"io/ioutil"
)

const defaultPort = "8080"

func main() {
	var (
		addr     = envString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":" + addr, "HTTP listen address")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	dir, err := ioutil.TempDir("", "images")

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir)

	client := NewDownloadClient(1, dir)

	sc := NewService(client)

	sc = NewLoggingService(log.With(logger, "component", "proxy_image"), sc)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/", MakeHandler(sc, httpLogger))

	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}