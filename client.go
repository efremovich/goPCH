package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

//ErrTooManyRedirect - Too many redirects
//ErrHTTPRedirect - Redirect to non-https server
var (
	ErrTooManyRedirect = errors.New("Too many redirects")
	ErrHTTPRedirect    = errors.New("Redirect to non-https server")
)

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) > 10 {
		return ErrTooManyRedirect
	}
	return nil
}

//GetClient - инициализация клиента
func GetClient() http.Client {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	redir := redirectPolicyFunc

	client := http.Client{
		Jar:           jar,
		CheckRedirect: redir,
	}
	return client

}

//GetClient - инициализация клиента
func GetClientTLS() http.Client {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	redir := redirectPolicyFunc

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Jar:           jar,
		CheckRedirect: redir,
		Transport:     tr,
	}
	return client

}
