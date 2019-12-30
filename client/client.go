package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Scheme  string
	Host    string
	baseUrl string
	Client  http.Client
}

func NewClient(host string) Client {
	h := "http://"
	return Client{
		Host:    host,
		Scheme:  h,
		baseUrl: h + host,
		Client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// below is so repetitive --> codegen?
func (c Client) GetAllProcessors() (*api.Processors, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(api.ProcessorAPIEndpoints["GetAllProcessors"]))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.Processors
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetProcessor(procID string) (*api.Processor, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(api.ProcessorAPIEndpoints["GetProcessor"], procID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.Processor
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetAllRuns(procID string) (*api.Runs, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(api.ProcessorAPIEndpoints["GetAllRuns"], procID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.Runs
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetRun(procID, runID string) (*api.Run, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(api.ProcessorAPIEndpoints["GetRun"], procID, runID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.Run
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetConfig(procID, runID string) (*api.Configuration, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(api.ProcessorAPIEndpoints["GetConfig"], procID, runID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.Configuration
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) Configure(procID, runID string, config api.Configuration) error {
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// TODO: maybe check return type?
	_, err = c.Client.Post(c.baseUrl+fmt.Sprintf(api.ProcessorAPIEndpoints["Configure"], procID, runID), "application/json", bytes.NewReader(body))

	return err
}

func (c Client) GetData(procID, runID string) (*api.DataSpec, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(api.ProcessorAPIEndpoints["GetData"]))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.DataSpec
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) AddData(procID, runID string, data api.Data) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// TODO: maybe check return type?
	_, err = c.Client.Post(c.baseUrl+fmt.Sprintf(api.ProcessorAPIEndpoints["AddData"], procID, runID), "application/json", bytes.NewReader(body))

	return err
}
