package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
	"github.com/spf13/pflag"
	logger "github.com/lexkong/log"
	"github.com/moocss/apiserver/src"
	"github.com/moocss/apiserver/src/config"
	v "github.com/moocss/apiserver/src/pkg/version"
)



func main() {
	opts := config.ConfYaml{}

	var (
		showVersion bool
		configFile string
	)

	pflag.StringVar(&configFile, "c" , "", "Configuration file path.")
	pflag.StringVar(&configFile, "config", "", "Configuration file path.")
	pflag.BoolVar(&showVersion, "v", false, "Print version information.")
	pflag.BoolVar(&showVersion, "version", false, "Print version information.")
	pflag.StringVar(&opts.Core.Address, "A", "", "address to bind")
	pflag.StringVar(&opts.Core.Address, "address", "", "address to bind")
	pflag.StringVar(&opts.Core.Port, "p", "", "port number for gorush")
	pflag.StringVar(&opts.Core.Port, "port", "", "port number for gorush")

	pflag.Parse()

	v.SetVersion(src.Version)

	if showVersion {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}

	var err error
	// set default parameters.
	src.Conf, err = config.Init(configFile)
	if err != nil {
		fmt.Printf("Load yaml config file error: '%v'", err)
		return
	}

	// overwrite server port and address
	if opts.Core.Port != "" {
		src.Conf.Core.Port = opts.Core.Port
	}
	if opts.Core.Address != "" {
		src.Conf.Core.Address = opts.Core.Address
	}

	g := src.New()
	srv := &http.Server{
		Addr: src.Conf.Core.Address + ":" + src.Conf.Core.Port,
		Handler: g,
	}

	// 启动服务
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	// 健康检查
	go func() {
		// ping server
		if err := pingServer(); err != nil {
			logger.Fatal("The router has no response, or it might took too long to start up. ", err)
		}
		logger.Info("The router has been deployed successfully.")
	}()

	// 打开浏览器
	go func() {
		time.Sleep(time.Second * 6)
		startBrowser("http://localhost:" + src.Conf.Core.Port)
	}()

	logger.Infof("Start to listening the incoming requests on http address: %s", src.Conf.Core.Port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	maxPingConf := src.Conf.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + src.Conf.Core.Port + "/sd/health")
		// defer resp.Body.Close();
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		logger.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
