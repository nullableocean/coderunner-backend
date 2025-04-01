package httpcommiter

import (
	"bytes"
	"codeproccesor/src/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiMethod = "POST"
)

type HttpCommiter struct {
	commitUrl    string
	accessToken  string
	accessHeader string
}

func NewHttpCommiter(commitUrl string, token string, tokenHeader string) *HttpCommiter {
	return &HttpCommiter{
		commitUrl:    commitUrl,
		accessToken:  token,
		accessHeader: tokenHeader,
	}
}

func (c *HttpCommiter) Commit(r domain.Result) error {
	req, err := c.createRequest(&r)
	if err != nil {
		return err
	}

	err = c.sendRequest(req)

	return err
}

func (c *HttpCommiter) createRequest(r *domain.Result) (*http.Request, error) {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(r)
	if err != nil {
		return nil, fmt.Errorf("result json marshal error: %s", err)
	}

	req, err := http.NewRequest(apiMethod, c.commitUrl, body)
	if err != nil {
		return nil, fmt.Errorf("commiter new request error: %s", err)
	}

	req.Header.Set(c.accessHeader, c.accessToken)

	return req, nil
}

func (c *HttpCommiter) sendRequest(req *http.Request) error {
	tries := 3
	var ok bool
	var err error
	var resp *http.Response

	client := http.DefaultClient

	for tries != 0 {
		err = nil

		resp, err = client.Do(req)
		if err != nil {
			err = fmt.Errorf("http request error: %s", err)
		}

		if resp.StatusCode == http.StatusOK {
			err = nil
			ok = true
			break
		}

		tries--
	}

	if !ok && err == nil {
		b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("commite result api error status: %s body: %s", resp.Status, b)
	}

	return err
}
