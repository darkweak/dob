package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"encoding/base64"
)

// DomainAcmeFile structure is composed of DNS Main and SANs
type DomainAcmeFile struct {
	Main string   `json:"Main"`
	SANs []string `json:"SANs"`
}

// CertificateAcmeFile structure is composed of Certificate, DomainAcmeFile and Key
type CertificateAcmeFile struct {
	Certificate string         `json:"Certificate"`
	Domain      DomainAcmeFile `json:"Domain"`
	Key         string         `json:"Key"`
}

// BodyAcmeFile structure is composed of Contact, Status
type BodyAcmeFile struct {
	Contact []string `json:"contact"`
	Status  string   `json:"status"`
}

// RegistrationAcmeFile structure is composed of Body, URI
type RegistrationAcmeFile struct {
	Body BodyAcmeFile `json:"body"`
	URI  string       `json:"uri"`
}

// AccountAcmeFile structure is composed of Email, KeyType, PrivateKey, Registration
type AccountAcmeFile struct {
	Email        string               `json:"Email"`
	KeyType      string               `json:"KeyType"`
	PrivateKey   string               `json:"PrivateKey"`
	Registration RegistrationAcmeFile `json:"Registration"`
}

// AcmeFile structure is Træfik acme.json hierachy
type AcmeFile struct {
	Account      AccountAcmeFile       `json:"Account"`
	Certificates []CertificateAcmeFile `json:"Certificates"`
}

func (a *AcmeFile) writeTofile () {
	data, _ := json.MarshalIndent(a, "", "\t")
	ioutil.WriteFile("test.json", data, 777)
}

func acmeFileBuilder (email string) *AcmeFile {
	acmeFile := &AcmeFile{
		Account: AccountAcmeFile{
			Email: fmt.Sprintf("%s", email),
			KeyType: "4096",
			PrivateKey: "",
			Registration: RegistrationAcmeFile{
				Body: BodyAcmeFile{
					Contact: []string{fmt.Sprintf("%s", email)},
					Status: "active",
				},
			},
		},
		Certificates: []CertificateAcmeFile{},
	}
	acmeFile.writeTofile()

	return acmeFile
}

func (a *AcmeFile) contains (domain string) (bool, int) {
	val := false
	i := 0
	for index, cert := range a.Certificates {
		if cert.Domain.Main == domain {
			val = true
			i = index
			break
		}
	}
	return val, i
}

func (a *AcmeFile) addCertificate (c *certificate) {
	v, i := a.contains(c.Domain)
	certif := readFile(c.Cert)
	key := readFile(c.Key)
	cert := CertificateAcmeFile{
		Certificate: base64.StdEncoding.EncodeToString(certif),
		Domain: DomainAcmeFile{
			Main: c.Domain,
			SANs: c.Sans,
		},
		Key: base64.StdEncoding.EncodeToString(key),
	}
	if v {
		a.Certificates[i] = cert
	} else {
		a.Certificates = append(a.Certificates, cert)
	}
	a.writeTofile()
}
