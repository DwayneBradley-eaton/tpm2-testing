package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := flag.String("f", "", "Filename to attemp to parse")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatal("Filename is required")
	}

	certCerFile, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	derBytes := make([]byte, 2000)

	count, err := certCerFile.Read(derBytes)
	if err != nil {
		log.Fatal(err)
	}

	certCerFile.Close()

	// trim the bytes to actual length in call
	cert, err := x509.ParseCertificate(derBytes[0:count])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name %s\n", cert.Subject.CommonName)
	fmt.Printf("Not before %s\n", cert.NotBefore.String())
	fmt.Printf("Not after %s\n", cert.NotAfter.String())
}
