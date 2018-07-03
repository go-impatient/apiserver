package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/lexkong/log"
	"github.com/moocss/apiserver/src"
	"github.com/moocss/apiserver/src/config"
	v "github.com/moocss/apiserver/src/pkg/version"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

var (
	cfg     = pflag.StringP("config", "c", "", "apiserver config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()

	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	g := src.New()
	srv := &http.Server{
		Addr:    viper.GetString("addr"),
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
		startBrowser(viper.GetString("url"))
	}()

	logger.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))

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
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
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
