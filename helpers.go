package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

func rmCerts(seedName string) {
	for _, t := range []string{"cert", "key"} {
		filePath := fmt.Sprintf(`%s.%s`, seedName, t)
		err := os.Remove(filePath)
		if err != nil {
			log.Printf("error revoming cert file %s err=%v", filePath, err)
		}
	}
}

func generateCertFiles(seedName string) {
	// Generate a private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Create a certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year validity

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Big Corp"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create a self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic(err)
	}

	// Encode and save the certificate to a .crt file
	certOut, err := os.Create(fmt.Sprintf(`%s.cert`, seedName))
	if err != nil {
		panic(err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	// Encode and save the private key to a .key file
	keyOut, err := os.Create(fmt.Sprintf(`%s.key`, seedName))
	if err != nil {
		panic(err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
}
