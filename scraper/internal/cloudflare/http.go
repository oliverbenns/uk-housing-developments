package cloudflare

import (
	"crypto/tls"
	"net/http"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"

func NewTransport() http.Transport {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		},
	}

	return http.Transport{
		TLSClientConfig: tlsConfig,
	}
}

// Bypass bot detection...
func CreateHttpClient() *http.Client {
	transport := NewTransport()
	cfTransport := transportWithUA{
		T: &transport,
	}

	client := &http.Client{
		Transport: &cfTransport,
	}

	return client
}

type transportWithUA struct {
	T http.RoundTripper
}

func (cft *transportWithUA) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", ua)
	return cft.T.RoundTrip(req)
}
