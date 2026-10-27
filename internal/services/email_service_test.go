package services

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSmtpTLSConfig_VerifiesServerCertificate(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	defer server.Close()

	certificate := server.Certificate()
	trustedRoots := x509.NewCertPool()
	trustedRoots.AddCert(certificate)

	tests := []struct {
		name      string
		host      string
		roots     *x509.CertPool
		wantError string
	}{
		{
			name:  "trusted certificate with matching hostname",
			host:  certificate.DNSNames[0],
			roots: trustedRoots,
		},
		{
			name:      "untrusted certificate",
			host:      certificate.DNSNames[0],
			roots:     x509.NewCertPool(),
			wantError: "unknown authority",
		},
		{
			name:      "trusted certificate with mismatched hostname",
			host:      "wrong-host.invalid",
			roots:     trustedRoots,
			wantError: "hostname mismatch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := smtpTLSConfig(tt.host)
			// Keep trust deterministic without changing certificate verification.
			config.RootCAs = tt.roots
			conn, err := tls.Dial("tcp", server.Listener.Addr().String(), config)
			if conn != nil {
				defer conn.Close()
			}
			switch tt.wantError {
			case "":
				if err != nil {
					t.Fatalf("TLS handshake failed: %v", err)
				}
			case "unknown authority":
				var unknownAuthority x509.UnknownAuthorityError
				if !errors.As(err, &unknownAuthority) {
					t.Fatalf("expected unknown authority error, got %v", err)
				}
			case "hostname mismatch":
				var hostnameError x509.HostnameError
				if !errors.As(err, &hostnameError) {
					t.Fatalf("expected hostname verification error, got %v", err)
				}
			}
		})
	}
}
