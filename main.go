package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	ServiceAccountPath string = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	TargetURL          string = os.Getenv("TARGET_URL")
)

// Header mapping
//
// References:
// - OAuth2-Proxy https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/overview/
// - Kubernetes dashboard https://user-images.githubusercontent.com/8249283/54769658-03c87a00-4bd8-11e9-86ea-ea9165bb82da.png
// OAuth2-Proxy                                                                 Kubernetes dashboard
// X-Forwarded-Preferred-Username                                               Impersonate-User
// X-Forwarded-Groups (format: A, B, C)                                         Impersonate-Group (one per group)

var serviceAccountToken string

func injectImpersonationHeaders(req *http.Request) {
	// Add Authorization header that holds the service account token to the request
	req.Header.Set("Authorization", "Bearer "+serviceAccountToken)

	// Add the Impersonate-User header to the request
	if value := req.Header.Get("X-Forwarded-Preferred-Username"); len(value) > 0 {
		req.Header.Set("Impersonate-User", value)
	}

	// Add the Impersonate-Group headers to the request
	if value := req.Header.Get("X-Forwarded-Groups"); len(value) > 0 {
		// Add the Impersonate-Group headers to the request
		groups := strings.Split(value, ",")
		for _, group := range groups {
			req.Header.Add("Impersonate-Group", group)
		}
	}
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	targetUrl, _ := url.Parse(TargetURL)

	dump, _ := httputil.DumpRequest(req, false)
	log.Printf("Before injection: %q\n", dump)

	injectImpersonationHeaders(req)

	dump, _ = httputil.DumpRequest(req, false)
	log.Printf("After injection: %q\n", dump)

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
