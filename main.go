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
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Req is the request structure
type Req struct {
	Email string `json:"email"`
}

// Res is the response structure
type Res struct {
	r         Req
	Signature string `json:"signature"`
	Pubkey    string `json:"pubkey"`
}

func main() {
	// 	fmt.Println(`{
	//     "message":"theAnswerIs42",
	//     "signature":"MGUCMCDwlFyVdD620p0hRLtABoJTR7UNgwj8g2r0ipNbWPi4Us57YfxtSQJ3dAkHslyBbwIxAKorQmpWl9QdlBUtACcZm4kEXfL37lJ+gZ/hANcTyuiTgmwcEC0FvEXY35u2bKFwhA==",
	//     "pubkey":"-----BEGIN PUBLIC KEY-----\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEI5/0zKsIzou9hL3ZdjkvBeVZFKpDwxTb\nfiDVjHpJdu3+qOuaKYgsLLiO9TFfupMYHLa20IqgbJSIv/wjxANH68aewV1q2Wn6\nvLA3yg2mOTa/OHAZEiEf7bVEbnAov+6D\n-----END PUBLIC KEY-----\n"
	// }`)

	r := mux.NewRouter()

	r.HandleFunc("/", Auth).Methods("POST")
	r.HandleFunc("/", Home).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Home is the landing page of the api
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome To Certificate Service!")
}

// Auth serves as the authenticator entrypoint for the service
func Auth(w http.ResponseWriter, r *http.Request) {

	ca := CreateCertObject()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	cab, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	PublicKey(cab)

	PrivateKey(priv)

}

// CreateCertObject creates the x509 certificate object
func CreateCertObject() *x509.Certificate {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"Smart Edge"},
			Country:       []string{"USA"},
			Province:      []string{"TX"},
			Locality:      []string{"Plano"},
			StreetAddress: []string{"1234 Awesome Ave"},
			PostalCode:    []string{"12345"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	return ca
}

// PublicKey creates takes in a certificate slice and returns the public key
func PublicKey(cab []byte) *os.File {
	// Public Key
	certOut, err := os.Create("ca.crt")
	if err != nil {
		log.Panic("Public Key Creation failed", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cab})
	certOut.Close()
	log.Print("written cert.pem\n")
	return certOut
}

// PrivateKey takes in a private key and creates a private key file
func PrivateKey(priv *rsa.PrivateKey) *os.File {
	// Private Key
	keyOut, err := os.OpenFile("ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Panic("Private Key Creation failed", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("written key.pem\n")
	return keyOut
}
