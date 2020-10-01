package web_test

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"ticket/web"
	"time"

	"github.com/phayes/freeport"
)

const (
	certFile = "testdata/localhost.cert"
	keyFile = "testdata/localhost.key"
)

func startTestServer(t *testing.T) (addr string, client *http.Client) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	addr = net.JoinHostPort("localhost", strconv.Itoa(port))
	go web.ListenAndServeWithTLS(addr, certFile, keyFile)
	url := "https://" + addr + "/healthz"
	client = localhostTrustingClient(t)
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	for resp.StatusCode != http.StatusOK {
		log.Println("RETRYING...")
		time.Sleep(10 * time.Millisecond)
		// we know there won't be an error
		resp, _ = http.Get(url)
	}
	return addr, client
}

func localhostTrustingClient(t *testing.T) *http.Client {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := ioutil.ReadFile(certFile)
	if err != nil {
		t.Fatalf("Failed to append %q to RootCAs: %v", certFile, err)
	}
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		t.Fatal("No certs appended. That's bad")
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
	}}}
}
func TestWeb(t *testing.T) {
	t.Parallel()
	addr, client := startTestServer(t)
	resp, err := client.Get("https://" + addr + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}