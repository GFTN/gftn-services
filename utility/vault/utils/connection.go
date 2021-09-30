// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utils

import (
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"strings"
	"crypto/x509"
	"log"
	"github.com/op/go-logging"
)

var LOGGER = logging.MustGetLogger("vault")

func (vs *Session) Get(url string, session string) ([]byte, error) {

	if vs.CertPath != "" {
		LOGGER.Infof("GET query URL with client-side certificate: %s", url)
	} else {
		LOGGER.Infof("GET query URL without client-side certificate: %s", url)
	}
	req, err := http.NewRequest("GET", vs.BaseURL+url, nil)

	req.Header.Add("Authorization", session)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	if vs.CertPath != "" {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(vs.CertPath, vs.KeyPath)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(vs.CertPath)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
			Renegotiation:      tls.RenegotiateFreelyAsClient,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{Transport: transport}
	}

	res, err := client.Do(req)
	if err != nil {
		LOGGER.Error("%s", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func (vs *Session) Post(url string, session string, payload *strings.Reader) ([]byte, error) {
	LOGGER.Infof("POST query URL: %s", url)
	req, err := http.NewRequest("POST", vs.BaseURL+url, payload)

	req.Header.Add("Authorization", session)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func (vs *Session) Put(url string, session string, payload *strings.Reader, additionHeader string) ([]byte, error) {
	LOGGER.Infof("PUT Query URL: %s", url)
	req, err := http.NewRequest("PUT", vs.BaseURL+url, payload)

	req.Header.Add("Authorization", session)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if additionHeader != "" {
		req.Header.Add(additionHeader, "Yes")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		LOGGER.Error("%s", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}
