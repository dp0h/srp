package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/dp0h/srp/app/pool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// SRPServer - reverse proxy server
type SRPServer struct {
	port          int
	sslMode       string
	host          string
	certFile      string
	keyFile       string
	autoCertPath  string
	validateCerts bool
	rwp           *pool.RandomWeightedPool
}

// NewReverseProxyServer creates a new reverse proxy server
func NewReverseProxyServer(port int, sslMode string, host string, certFile string, keyFile string, autoCertPath string, validateCerts bool, rwp *pool.RandomWeightedPool) *SRPServer {
	res := SRPServer{
		port:          port,
		sslMode:       sslMode,
		host:          host,
		certFile:      certFile,
		keyFile:       keyFile,
		autoCertPath:  autoCertPath,
		validateCerts: validateCerts,
		rwp:           rwp,
	}
	return &res
}

// Run reverse proxy server
func (s *SRPServer) Run() {
	log.Info().Int("port", s.port).Str("ssl-mode", s.sslMode).Msg("starting reverse proxy server")

	if !s.validateCerts {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	switch s.sslMode {
	case "none":
		s.runHttp()
		break
	case "static":
		s.runStatic()
		break
	case "auto":
		s.runAuto()
		break
	default:
		log.Fatal().Str("ssl-mode", s.sslMode).Msg("unrecognized ssl mode")
	}
}

func (s *SRPServer) runHttp() {
	addr := fmt.Sprintf(":%d", s.port)
	http.HandleFunc("/", s.handle)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func (s *SRPServer) runStatic() {
	if s.certFile == "" {
		log.Fatal().Msg("path to cert.pem is required")
	}
	if s.keyFile == "" {
		log.Fatal().Msg("path to key.pem is required")
	}

	addr := fmt.Sprintf(":%d", s.port)
	http.HandleFunc("/", s.handle)
	if err := http.ListenAndServeTLS(addr, s.certFile, s.keyFile, nil); err != nil {
		panic(err)
	}
}

func (s *SRPServer) runAuto() {
	if s.autoCertPath == "" {
		log.Fatal().Msg("autocert-path is required")
	}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(s.host),
		Cache:      autocert.DirCache(s.autoCertPath),
	}

	addr := fmt.Sprintf(":%d", s.port)
	http.HandleFunc("/", s.handle)
	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go func() {
		err := http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		if err != nil {
			log.Error().Err(err).Msg("failed to listen and server on :http")
		}
	}()

	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}

func (s *SRPServer) handle(res http.ResponseWriter, req *http.Request) {
	target, err := s.rwp.Next()
	if err != nil {
		log.Warn().Err(err).Msg("no live services available")
		res.WriteHeader(http.StatusGatewayTimeout)
		return
	}
	serveReverseProxy(target, res, req)
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse target")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	log.Debug().Str("host", req.Host).Str("url-host", req.URL.Host).Str("url-path", req.URL.Path).Str("url-scheme", req.URL.Scheme).Str("target", target).Msg("handle")

	req.URL.Host = targetUrl.Host
	req.URL.Scheme = targetUrl.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = targetUrl.Host

	proxy.ServeHTTP(res, req)
}
