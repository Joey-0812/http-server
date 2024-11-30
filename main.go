package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

func generateCertificate() {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"My Organization"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certOut, err := os.Create("cert.pem")
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.Create("key.pem")
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
	}
	key, _ := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: key})
	keyOut.Close()

	log.Println("Generated cert.pem and key.pem")
}

func handlerHttps(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, HTTPS with self-generated certificate!"))
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, HTTP request!"))
}

func main() {
	generateCertificate()

	go func() {
		http.HandleFunc("/hello", handlerHttp)
		log.Fatal(http.ListenAndServe(":80", nil))
	}()
	http.HandleFunc("/", handlerHttps)
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}
