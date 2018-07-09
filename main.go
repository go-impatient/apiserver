package main

import (
	"encoding/json"
	"fmt"
	logger "github.com/lexkong/log"
	"github.com/moocss/apiserver/src"
	"github.com/moocss/apiserver/src/config"
	v "github.com/moocss/apiserver/src/pkg/version"
	"github.com/spf13/pflag"
	"os"
	"golang.org/x/sync/errgroup"
)

var usageStr = `
              .__                                        
_____  ______ |__| ______ ______________  __ ___________ 
\__  \ \____ \|  |/  ___// __ \_  __ \  \/ // __ \_  __ \
 / __ \|  |_> >  |\___ \\  ___/|  | \/\   /\  ___/|  | \/
(____  /   __/|__/____  >\___  >__|    \_/  \___  >__|   
     \/|__|           \/     \/                 \/       

Usage: apiserver [options]

Server Options:
	-c, --config <file>              Configuration file path
	-a, --address <address>          Address to bind (default: any)
	-p, --port <port>                Use port for clients (default: 9090)
Common Options:
	-h, --help                       Show this message
	-v, --version                    Show version
`

func main() {
	opts := config.ConfYaml{}

	var (
		showVersion bool
		configFile  string
	)

	pflag.StringVar(&configFile, "c", "", "Configuration file path.")
	pflag.StringVar(&configFile, "config", "", "Configuration file path.")
	pflag.BoolVar(&showVersion, "v", false, "Print version information.")
	pflag.BoolVar(&showVersion, "version", false, "Print version information.")
	pflag.StringVar(&opts.Core.Address, "a", "", "address to bind")
	pflag.StringVar(&opts.Core.Address, "address", "", "address to bind")
	pflag.StringVar(&opts.Core.Port, "p", "", "port number for gorush")
	pflag.StringVar(&opts.Core.Port, "port", "", "port number for gorush")
	pflag.Usage = usage
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

	var g errgroup.Group
	g.Go(func() error {
		// 启动服务
		return src.RunHTTPServer()
	})
	g.Go(func() error {
		// 健康检查
		return src.PingServer()
	})

	if err = g.Wait(); err != nil {
		logger.Error("接口服务出错了：", err)
	}

	logger.Infof("Start to listening the incoming requests on http address: %s", src.Conf.Core.Port)
}

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
