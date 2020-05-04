package rabbit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	DEFAULT_VHOST     = "%2F"
	GET_QUEUES_URL    = "/api/queues/%s"    // vhost
	GET_EXCHANGES_URL = "/api/exchanges/%s" // vhost
	GET_BINDINGS_URL  = "/api/bindings/%s"  // vhost

	POST_BINDING_URL = "/api/bindings/%s/e/%s/q/%s" // vhost, exchange, queue

	PUT_QUEUE_URL    = "/api/queues/%s/%s"    // vhost, queue
	PUT_EXCHANGE_URL = "/api/exchanges/%s/%s" // vhost, queue
)

type Queue struct {
	Name       string                 `json:"name"`
	VHost      string                 `json:"vhost"`
	Durable    bool                   `json:"durable"`
	AutoDelete bool                   `json:"auto_delete"`
	Arguments  map[string]interface{} `json:"arguments"`
}

type Exchange struct {
	Name       string `json:"name"`
	VHost      string `json:"vhost"`
	Type       string `json:"type"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
}

type Binding struct {
	Source        string `json:"source"`
	VHost         string `json:"vhost"`
	Destination   string `json:"destination"`
	DestType      string `json:"destination_type"`
	Routingkey    string `json:"routing_key"`
	PropertiesKey string `json:"properties_key"`
}

type RabbitMQ struct {
	BaseUrl   string
	User      string
	Password  string
	VHost     string
	Queues    []Queue
	Exchanges []Exchange
	Bindings  []Binding
}

func (r *RabbitMQ) doCall(method, url string, payload io.Reader) (io.ReadCloser, error) {

	fmt.Println(method, "->", url)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(r.User, r.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (r *RabbitMQ) doGet(url string) (io.ReadCloser, error) {

	return r.doCall("GET", url, nil)
}

func (r *RabbitMQ) doPost(url string, body []byte) (io.ReadCloser, error) {

	ior_body := bytes.NewReader(body)

	return r.doCall("POST", url, ior_body)
}

func (r *RabbitMQ) doPut(url string, body []byte) (io.ReadCloser, error) {

	ior_body := bytes.NewReader(body)

	return r.doCall("PUT", url, ior_body)
}

func (r *RabbitMQ) readQueues() error {

	url := fmt.Sprintf(r.BaseUrl+GET_QUEUES_URL, r.VHost)

	body, err := r.doGet(url)
	if err != nil {
		return err
	}

	return json.NewDecoder(body).Decode(&r.Queues)
}

func (r *RabbitMQ) readExchanges() error {

	url := fmt.Sprintf(r.BaseUrl+GET_EXCHANGES_URL, r.VHost)

	body, err := r.doGet(url)
	if err != nil {
		return err
	}

	return json.NewDecoder(body).Decode(&r.Exchanges)
}

func (r *RabbitMQ) readBindings() error {

	url := fmt.Sprintf(r.BaseUrl+GET_BINDINGS_URL, r.VHost)

	body, err := r.doGet(url)
	if err != nil {
		return err
	}

	return json.NewDecoder(body).Decode(&r.Bindings)
}

func (r *RabbitMQ) Collect() error {

	if err := r.readQueues(); err != nil {
		return err
	}

	if err := r.readExchanges(); err != nil {
		return err
	}

	if err := r.readBindings(); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) CloneTo(hostport, user, password, vhost string) error {

	if vhost == "/" {
		vhost = DEFAULT_VHOST
	}

	baseurl := "http://" + hostport

	for _, queue := range r.Queues {
		url := fmt.Sprintf(baseurl+PUT_QUEUE_URL, vhost, queue.Name)

		queue.VHost = vhost

		body, _ := json.Marshal(queue)

		fmt.Println(string(body))

		if res, err := r.doPut(url, body); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}

	}

	for _, exchange := range r.Exchanges {

		if strings.HasPrefix(exchange.Name, "amq.") {
			continue
		}

		url := fmt.Sprintf(baseurl+PUT_EXCHANGE_URL, vhost, exchange.Name)

		exchange.VHost = vhost

		body, _ := json.Marshal(exchange)

		fmt.Println(string(body))

		if res, err := r.doPut(url, body); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}

	for _, binding := range r.Bindings {

		if binding.Source == "" {
			//fmt.Println("Skipping binding", binding, "empty source")
			continue
		}

		if binding.DestType == "queue" {
			url := fmt.Sprintf(baseurl+POST_BINDING_URL, vhost, binding.Source, binding.Destination)

			binding.VHost = vhost

			body, _ := json.Marshal(binding)

			fmt.Println(string(body))

			if res, err := r.doPost(url, body); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}

		} else {
			fmt.Println("Unknown binding destination type '", binding.DestType, "'")
		}
	}

	return nil
}

func (r *RabbitMQ) DumpTo(filename string) error {

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	json.NewEncoder(f).Encode(r)

	f.Close()

	return nil
}

func CollectFromFile(filename string) (*RabbitMQ, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	r := &RabbitMQ{}

	if err := json.NewDecoder(f).Decode(r); err != nil {
		return nil, err
	}

	f.Close()

	return r, nil
}
