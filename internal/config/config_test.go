package config

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	u, _ := url.ParseRequestURI("/etc/main.conf")
	fmt.Println(u.Scheme)
}
