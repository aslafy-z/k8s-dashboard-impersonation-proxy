// Application which greets you.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	// "github.com/caarlos0/env/v6"
)

var (
	ServiceAccountPath string = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	TargetHost         string = "http://localhost:8081"
)

// Header mapping
//
// References:
// - OAuth2-Proxy https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/overview/
// - Kubernetes dashboard https://user-images.githubusercontent.com/8249283/54769658-03c87a00-4bd8-11e9-86ea-ea9165bb82da.png
// OAuth2-Proxy                                                                 Kubernetes dashboard
// X-Auth-Request-User, X-Auth-Request-Preferred-Username                       Impersonate-User
// X-Auth-Request-Groups (format: A, B, C)                                      Impersonate-Group (one per group)

var serviceAccountToken string

func injectImpersonationHeaders(req *http.Request) {
	// Add Authorization header that holds the service account token to the request
	req.Header.Set("Authorization", "Bearer "+serviceAccountToken)

	// Add the Impersonate-User header to the request
	impersonateUser := ""
	if value := req.Header.Get("X-Auth-Request-User"); len(value) > 0 {
		impersonateUser = value
	}
	if value := req.Header.Get("X-Auth-Request-Preferred-Username"); len(value) > 0 {
		impersonateUser = value
	}
	if len(impersonateUser) > 0 {
		req.Header.Set("Impersonate-User", impersonateUser)
	}

	// Add the Impersonate-Group headers to the request
	impersonateGroups := strings.Split(req.Header.Get("X-Auth-Request-Groups"), ",")
	for _, impersonateGroup := range impersonateGroups {
		req.Header.Add("Impersonate-Group", impersonateGroup)
	}
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	targetUrl, _ := url.Parse(TargetHost)

	injectImpersonationHeaders(req)

	dump, _ := httputil.DumpRequest(req, false)
	fmt.Printf("%q\n", dump)

	proxy := httputil.ReverseProxy{Director: func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.Host = targetUrl.Host
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}}

	proxy.ServeHTTP(res, req)
}

func main() {
	if _, err := os.Stat(ServiceAccountPath); err == nil {
		if token, err := ioutil.ReadFile(ServiceAccountPath); err == nil {
			serviceAccountToken = string(token)
		} else {
			panic(err)
		}
	}

	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
