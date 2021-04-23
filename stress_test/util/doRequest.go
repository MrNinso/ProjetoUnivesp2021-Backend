package util

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func DoRequest(metodo, url, json, email, token string) (time.Duration, int, []byte) {
	req, err := http.NewRequest(metodo, url, strings.NewReader(json))

	if err != nil {
		panic(err)
	}

	req.Header.Add("token_api", token)
	req.Header.Add("email_api", email)

	client := http.Client{}

	initTime := time.Now()

	resp, err := client.Do(req)

	requetsTime := time.Since(initTime)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return requetsTime, resp.StatusCode, body
}
