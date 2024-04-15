package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	env "github.com/caarlos0/env/v8"
)

type config struct {
	ServiceAccountToken string  `env:"SERVICE_ACCOUNT_TOKEN,unset"`
	ServiceAccountPath  string  `env:"SERVICE_ACCOUNT_PATH" envDefault:"/var/run/secrets/kubernetes.io/serviceaccount/token"`
	TargetURL           url.URL `env:"TARGET_URL,required"`
	HeaderUsername      string  `env:"HEADER_USERNAME,required" envDefault:"X-Auth-Request-Preferred-Username"`
	HeaderGroups        string  `env:"HEADER_GROUPS,required" envDefault:"X-Auth-Request-Groups"`
	ListenAddress       string  `env:"LISTEN_ADDRESS,required" envDefault:":8080"`
	InsecureTLSVerify   bool    `env:"INSECURE_TLS_VERIFY" envDefault:"false"`
	Debug               bool    `env:"DEBUG" envDefault:"false"`
}

var (
	cfg   config
	proxy httputil.ReverseProxy
)

// IsUrl checks if a string is a valid URL
func IsValidUrl(u *url.URL) bool {
	return u.Scheme != "" && u.Host != ""
}

// injectHeaders injects Authorization, User and Groups headers to the request
func injectHeaders(req *http.Request) {
	// inject service account token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+cfg.ServiceAccountToken)

	// inject user to the Impersonate-User header
	if value := req.Header.Get(cfg.HeaderUsername); len(value) > 0 {
		req.Header.Del(cfg.HeaderUsername)
		req.Header.Set("Impersonate-User", value)
	}

	// inject groups to Impersonate-Group headers
	if value := req.Header.Get(cfg.HeaderGroups); len(value) > 0 {
		req.Header.Del(cfg.HeaderGroups)
		groups := strings.Split(value, ",")
		for _, group := range groups {
			group = strings.TrimSpace(group)
			if len(group) > 0 {
				req.Header.Add("Impersonate-Group", group)
			}
		}
	}
}

// handleRequest handles incoming requests
func handleRequest(res http.ResponseWriter, req *http.Request) {
	if cfg.Debug {
		// log original request
		dump, _ := json.Marshal(req.Header)
		log.Printf("debug: before headers injection: %+v\n", string(dump))
	}

	// mutate request
	injectHeaders(req)

	if cfg.Debug {
		// log mutated request
		dump, _ := json.Marshal(req.Header)
		rDump := strings.Replace(string(dump), "Bearer "+cfg.ServiceAccountToken, "[REDACTED]", 1)
		log.Printf("debug: after headers injection: %+v\n", rDump)
	}

	// forward request
	proxy.ServeHTTP(res, req)
}

// handleReadinessRequest handles incoming readiness requests
func handleReadinessRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL)
		handler.ServeHTTP(res, req)
	})
}

func NewReverseProxy() *httputil.ReverseProxy {
	p := &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.URL.Scheme = cfg.TargetUrl.Scheme
		req.URL.Host = cfg.TargetUrl.Host
		req.Host = cfg.TargetUrl.Host
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}}
	p.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: cfg.InsecureTLSVerify},
	}
	return p
}

func main() {
	// retrieve configuration
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error: configuration parsing: %+v\n", err)
	}

	// validate configuration
	if ok := IsValidUrl(&cfg.TargetURL); !ok {
		log.Fatalln("error: target URL is not valid")
	}
	if len(cfg.ServiceAccountToken) == 0 && len(cfg.ServiceAccountPath) == 0 {
		log.Fatalln("error: at least one of service account token or path is required")
	}

	// show configuration
	rCfg := cfg
	if len(rCfg.ServiceAccountToken) > 0 {
		rCfg.ServiceAccountToken = "[REDACTED]"
	}
	log.Printf("Configuration: %+v\n", rCfg)

	// read service account token
	if len(cfg.ServiceAccountToken) > 0 {
		log.Printf("info: using service account token from environment variable\n")
	} else {
		token, err := ioutil.ReadFile(cfg.ServiceAccountPath)
		if err != nil {
			log.Fatalf("error: read token: %+v\n", err)
		}
		cfg.ServiceAccountToken = string(token)
		log.Printf("info: using service account token from '%s'\n", cfg.ServiceAccountPath)
	}

	// initialize reverse proxy
	proxy := NewReverseProxy()

	// listen and serve
	http.HandleFunc("/-/ready", handleReadinessRequest)
	http.HandleFunc("/", handleRequest)
	log.Printf("info: listening on %s\n", cfg.ListenAddress)
	if err := http.ListenAndServe(cfg.ListenAddress, logRequest(http.DefaultServeMux)); err != nil {
		log.Fatalf("error: listen: %+v\n", err)
	}
}
