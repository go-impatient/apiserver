package src

import (
	"github.com/gin-gonic/gin"
	"github.com/moocss/apiserver/src/router"
	"github.com/moocss/apiserver/src/router/middleware"
	"golang.org/x/crypto/acme/autocert"
	"github.com/lexkong/log"
	"net/http"
	"crypto/tls"
	"time"
	"errors"
	"github.com/seccom/kpass/src/logger"
	"github.com/moocss/apiserver/src/service"
)

// New returns a app instance
func New() *gin.Engine {
	// init db
	Storage.Init(Conf)
	defer service.DB.Close()

	// Set gin mode.
	gin.SetMode(Conf.Core.Mode)

	// Create the Gin engine.
	g := gin.New()

	// Routes
	router.Load(
		// Cores
		g,
		// Middlwares
		middleware.VersionMiddleware(),
		// middleware.Logging(),
		middleware.RequestId(),
	)
	return g
}

func autoTLSServer() *http.Server {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(Conf.Core.AutoTLS.Host),
		Cache:      autocert.DirCache(Conf.Core.AutoTLS.Folder),
	}
	return &http.Server{
		Addr:      	":https",
		TLSConfig: 	&tls.Config{GetCertificate: m.GetCertificate},
		Handler:  	New(),
	}
}

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer() (err error) {
	if !Conf.Core.Enabled {
		log.Debug("httpd server is disabled.")
		return nil
	}

	if Conf.Core.AutoTLS.Enabled {
		s := autoTLSServer()
		err = s.ListenAndServeTLS("", "")
	} else if Conf.Core.TLS.CertPath != "" && Conf.Core.TLS.KeyPath != "" {
		err = http.ListenAndServeTLS(Conf.Core.Address+":"+Conf.Core.TLS.Port, Conf.Core.TLS.CertPath, Conf.Core.TLS.KeyPath, New())
	} else {
		err = http.ListenAndServe(Conf.Core.Address+":"+Conf.Core.Port, New())
	}

	return
}

// PingServer
func PingServer() (err error) {
	maxPingConf := Conf.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + Conf.Core.Port + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		logger.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	err = errors.New("Cannot connect to the router.")
	return err
}
