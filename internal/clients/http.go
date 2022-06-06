package clients

import (
	"crypto/tls"
	"net/http"
	"time"
)

var httpClient *http.Client

func InitHTTPClient(insecureSkipVerify bool) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureSkipVerify}
	httpClient = &http.Client{Timeout: 10 * time.Second, Transport: customTransport}
}

func GetHTTPClient() *http.Client {
	return httpClient
}
