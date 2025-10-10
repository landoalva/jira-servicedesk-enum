// Copyright 2025 İrem Kuyucu
// Copyright 2025 Laurynas Četyrkinas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	cookie  string
}

func newClient(baseURL, cookie string) *Client {
	return &Client{
		baseURL: baseURL,
		cookie:  cookie,
	}
}

func (c *Client) post(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cookie != "" {
		req.AddCookie(&http.Cookie{
			Name:  "customer.account.session.token",
			Value: c.cookie,
		})
	}

	return http.DefaultClient.Do(req)
}

func (c *Client) get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if c.cookie != "" {
		req.AddCookie(&http.Cookie{
			Name:  "customer.account.session.token",
			Value: c.cookie,
		})
	}

	return http.DefaultClient.Do(req)
}

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func unmarshalJSON(resp *http.Response, v interface{}) error {
	data, err := readBody(resp)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}
