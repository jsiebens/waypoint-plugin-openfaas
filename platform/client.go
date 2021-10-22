package platform

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/openfaas/faas-cli/proxy"
)

var (
	defaultTimeout = 60 * time.Second
)

func (p *Platform) newClient() (*proxy.Client, error) {
	authChain, gateway, err := p.newAuthChain()
	if err != nil {
		return nil, err
	}

	transport := getDefaultCLITransport(p.config.TlsInsecure, &defaultTimeout)
	return proxy.NewClient(authChain, gateway, transport, &defaultTimeout)
}

func (p *Platform) newAuthChain() (proxy.ClientAuth, string, error) {
	username := selectValue(p.config.Username, "OPENFAAS_USERNAME")
	password := selectValue(p.config.Password, "OPENFAAS_PASSWORD")
	gateway := selectValue(p.config.Gateway, "OPENFAAS_URL")

	if username != "" && password != "" {
		return &BasicAuth{
			username: username,
			password: password,
		}, gateway, nil
	}

	token := selectValue(p.config.Token, "OPENFAAS_TOKEN")
	auth, err := proxy.NewCLIAuth(token, gateway)

	return auth, gateway, err
}

func selectValue(value, envVar string) string {
	if value != "" {
		return value
	}
	fromEnv := os.Getenv(envVar)
	if fromEnv != "" {
		return fromEnv
	}
	return ""
}

//BasicAuth basic authentication type
type BasicAuth struct {
	username string
	password string
}

func (auth *BasicAuth) Set(req *http.Request) error {
	req.SetBasicAuth(auth.username, auth.password)
	return nil
}

func getDefaultCLITransport(tlsInsecure bool, timeout *time.Duration) *http.Transport {
	if timeout != nil || tlsInsecure {
		tr := &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: false,
		}

		if timeout != nil {
			tr.DialContext = (&net.Dialer{
				Timeout: *timeout,
			}).DialContext

			tr.IdleConnTimeout = 120 * time.Millisecond
			tr.ExpectContinueTimeout = 1500 * time.Millisecond
		}

		if tlsInsecure {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: tlsInsecure}
		}
		tr.DisableKeepAlives = false

		return tr
	}
	return nil
}

func isFunctionNotFound(err error) bool {
	return strings.Contains(err.Error(), "404") ||
		strings.Contains(err.Error(), "No such function") ||
		strings.Contains(err.Error(), "No existing function")
}
