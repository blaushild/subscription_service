package controller

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (c *controller) Swagger(w http.ResponseWriter, r *http.Request) {
	swaggerProxyURL, _ := url.Parse("http://swagger-ui:8080")
	proxy := httputil.NewSingleHostReverseProxy(swaggerProxyURL)
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(rw, "Bad Gateway", http.StatusBadGateway)
	}

	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/swagger")

	proxy.ServeHTTP(w, r)
}
