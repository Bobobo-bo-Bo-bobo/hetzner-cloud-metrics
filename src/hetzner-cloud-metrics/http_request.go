package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func httpRequest(_url string, method string, header *map[string]string, reader io.Reader, proxy string, apitoken string) (HTTPResult, error) {
	var result HTTPResult
	var transp *http.Transport
	var proxyURL *url.URL
	var err error

	if proxy != "" {
		proxyURL, err = url.Parse(proxy)
		if err != nil {
			return result, err
		}
	}
	transp = &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	if proxy != "" {
		transp.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Transport: transp,
	}

	request, err := http.NewRequest(method, _url, reader)
	if err != nil {
		return result, err
	}

	defer func() {
		if request.Body != nil {
			ioutil.ReadAll(request.Body)
			request.Body.Close()
		}
	}()

	// add required headers
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apitoken))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	// set User-Agent
	request.Header.Set("User-Agent", userAgent)

	// close connection after response and prevent re-use of TCP connection because some implementations (e.g. HP iLO4)
	// don't like connection reuse and respond with EoF for the next connections
	request.Close = true

	// add supplied additional headers
	if header != nil {
		for key, value := range *header {
			request.Header.Add(key, value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}

	defer func() {
		ioutil.ReadAll(response.Body)
		response.Body.Close()
	}()

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.Header = response.Header
	result.Content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	return result, nil
}
