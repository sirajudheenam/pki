// // server.go
// package main

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// )

// func main() {
// 	// Load server cert + key
// 	cert, err := tls.LoadX509KeyPair("certs/server/server.chain.pem", "certs/server/server.key.pem")
// 	if err != nil {
// 		log.Fatal("Failed loading server cert/key:", err)
// 	}

// 	// Load Root CA (trust chain for client certs)
// 	caCert, err := os.ReadFile("certs/server/root.cert.pem")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	interCert, err := os.ReadFile("certs/server/intermediate.cert.pem")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)
// 	caCertPool.AppendCertsFromPEM(interCert)

// 	// TLS config: require client cert
// 	tlsConfig := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 		ClientCAs:    caCertPool,
// 		ClientAuth:   tls.RequireAndVerifyClientCert,
// 		MinVersion:   tls.VersionTLS12,
// 	}

// 	server := &http.Server{
// 		Addr:      ":8443",
// 		TLSConfig: tlsConfig,
// 	}

// 	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		// Print info about the client cert
// 		if len(r.TLS.PeerCertificates) > 0 {
// 			clientCert := r.TLS.PeerCertificates[0]
// 			fmt.Fprintf(w, "Hello, %s!\n", clientCert.Subject.CommonName)
// 		} else {
// 			fmt.Fprintf(w, "Hello, unknown client!\n")
// 		}
// 	})

// 	log.Println("Server listening on https://localhost:8443")
// 	log.Fatal(server.ListenAndServeTLS("", "")) // certs provided in TLSConfig
// }
