package push

import (
	"AlarmPawServer/config"
	"crypto/tls"
	"crypto/x509"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"runtime"
)

func init() {
	authKey, err := token.AuthKeyFromBytes([]byte(config.LocalConfig.Apple.ApnsPrivateKey))
	if err != nil {
		log.Printf("failed to create APNS auth key: %v\n", err)
	}

	var rootCAs *x509.CertPool

	system := func() string { return runtime.GOOS }()

	if system == "windows" {
		rootCAs = x509.NewCertPool()
	} else {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			log.Printf("failed to get rootCAs: %v\n", err)
		}
	}

	for _, ca := range config.ApnsCAs {
		rootCAs.AppendCertsFromPEM([]byte(ca))
	}

	CLI = &apns2.Client{
		Token: &token.Token{
			AuthKey: authKey,
			KeyID:   config.LocalConfig.Apple.KeyID,
			TeamID:  config.LocalConfig.Apple.TeamID,
		},
		HTTPClient: &http.Client{
			Transport: &http2.Transport{
				DialTLS: apns2.DialTLS,
				TLSClientConfig: &tls.Config{
					RootCAs: rootCAs,
				},
			},
			Timeout: apns2.HTTPClientTimeout,
		},
		Host: apns2.HostProduction,
	}
	log.Printf("init apns client success...\n")
}
