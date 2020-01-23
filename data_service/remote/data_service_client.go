// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package remote

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/evilsocket/islazy/tui"

	"github.com/maxlandon/wiregost/data_service/handlers"
)

// DbClient is a preconfigured HTTP client for query the data_service
type DbClient struct {
	*http.Client
	// Config
	BaseURL *url.URL
}

// NewClient instantiates a DbClient, loads Environment and configures the DbClient
// with data_service parameters, including TLS transport security.
func newClient() *DbClient {
	env := handlers.LoadEnv()

	client := &DbClient{
		&http.Client{},
		&url.URL{},
	}

	// Setup URL
	client.BaseURL.Scheme = "https"
	client.BaseURL.Host = env.Service.Address + ":" + strconv.Itoa(env.Service.Port)

	// Setup client TLS transport
	cert, err := tls.LoadX509KeyPair(env.Service.Certificate, env.Service.Key)
	if err != nil {
		log.Println(tui.Red("[!] Error loading X509 keypair: ") + err.Error())
	}

	conf := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	client.Transport = &http.Transport{TLSClientConfig: &conf}

	return client
}

// newRequest is a standard request that takes care of URL settings, request headers and body encoding.
func (c *DbClient) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	// Fill body
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Forge request
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Set workspace context header
	if ctx != nil {
		ws := ctx.Value("workspace_id")
		if ws != nil {
			ws := ws.(uint)
			ws64 := uint64(ws)
			wsID := strconv.FormatUint(ws64, 10)
			req.Header.Set("workspace_id", wsID)
		}
	}

	return req, nil
}

// do is a wrapper around http.Client.Do(), that uses the above custom newRequest function.
// It also unmarshals HTTP response's JSON body in the interface passed as a argument
func (c *DbClient) do(req *http.Request, v interface{}) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot parse response body")
		return err
	}

	if v != nil {
		err = json.Unmarshal(respBody, v)
		return err
	}

	return err
}
