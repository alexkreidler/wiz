package client

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"log"

	//"encoding/json"

	"fmt"
	"github.com/alexkreidler/wiz/api/processors"
	"io/ioutil"
	"net/http"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
func (c Client) GetAllProcessors() (*processors.Processors, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(processors.ProcessorAPIEndpoints["GetAllProcessors"]))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result processors.Processors
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetProcessor(procID string) (*processors.Processor, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(processors.ProcessorAPIEndpoints["GetProcessor"], procID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result processors.Processor
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetAllRuns(procID string) (*processors.Runs, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(processors.ProcessorAPIEndpoints["GetAllRuns"], procID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result processors.Runs
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetRun(procID, runID string) (*processors.Run, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(processors.ProcessorAPIEndpoints["GetRun"], procID, runID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result processors.Run
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) GetConfig(procID, runID string) (*processors.Configuration, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(processors.ProcessorAPIEndpoints["GetConfig"], procID, runID))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result processors.Configuration
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) Configure(procID, runID string, config processors.Configuration) error {
	log.Println("Configuring", procID, runID)
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// TODO: maybe check response header type?
	_, err = c.Client.Post(c.baseUrl+fmt.Sprintf(processors.ProcessorAPIEndpoints["Configure"], procID, runID), "application/json", bytes.NewReader(body))

	return err
}

func (c Client) GetData(procID, runID string) (*processors.DataSpec, error) {
	resp, err := c.Client.Get(c.baseUrl + fmt.Sprintf(processors.ProcessorAPIEndpoints["GetData"]))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result processors.DataSpec
	json.Unmarshal(data, &result)

	return &result, nil
}

func (c Client) AddData(procID, runID string, data processors.Data) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Println("Preparing to send add data request", procID, runID, string(body))

	// TODO: maybe check response header type?
	_, err = c.Client.Post(c.baseUrl+fmt.Sprintf(processors.ProcessorAPIEndpoints["AddData"], procID, runID), "application/json", bytes.NewReader(body))

	return err
}
