package main

import (
	"flag"
	"fmt"
	"github.com/cfloress-gb-cl/final-project-bootcamp/cmd/rest/user"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	md "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/rest/user/middlewares"
	"github.com/go-kit/log"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {

	// if we crash the go code, we get the file name and line number
	golog.SetFlags(golog.LstdFlags | golog.Lshortfile)

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var h http.Handler
	{
		h = user.MakeHTTPHandler(user.UserProxy{}, log.With(logger, "component", "HTTP"))
	}

	handler := md.UUIDContextMiddleware(h)
	handler = md.AuthenticationMiddleware(h)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	logger.Log("exit", <-errs)

}
