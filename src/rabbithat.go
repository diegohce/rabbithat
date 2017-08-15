package main

import (
    "fmt"
	"io"
    "net/http"
    "encoding/json"
)

const (
	DEFAULT_VHOST = "%2F"
	GET_QUEUES_URL = "/api/queues/%s" // vhost
	GET_EXCHANGES_URL = "/api/exchanges/%s"  // vhost
	GET_BINDINGS_URL = "/api/bindings/%s"  // vhost

	POST_BINDING_URL = "/api/bindings/%s/e/%s/q/%s" // vhost, exchange, queue

	PUT_QUEUE_URL = "/api/queues/%s/%s" // vhost, queue
	PUT_EXCHANGE_URL = "/api/exchanges/%s/%s" // vhost, queue
)


type Queue struct {
	Name       string `json:"name"`
	VHost      string `json:"vhost"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`

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
	Vhost         string `json:"vhost"`
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


func (r *RabbitMQ) doGet(url string) (io.ReadCloser, error) {

	fmt.Println("Requesting", url)

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

    req.SetBasicAuth(r.User, r.Password)
    resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}


func (r *RabbitMQ) readQueues() error {

    url := fmt.Sprintf(r.BaseUrl + GET_QUEUES_URL, r.VHost)

	body, err := r.doGet(url)
	if err != nil {
		return err
	}

    json.NewDecoder(body).Decode(&r.Queues)

	return nil
}

func (r *RabbitMQ) readExchanges() error {

    url := fmt.Sprintf(r.BaseUrl + GET_EXCHANGES_URL, r.VHost)

	body, err := r.doGet(url)
	if err != nil {
		return err
	}

    json.NewDecoder(body).Decode(&r.Exchanges)

	return nil
}

func (r *RabbitMQ) readBindings() error {

    url := fmt.Sprintf(r.BaseUrl + GET_BINDINGS_URL, r.VHost)

	body, err := r.doGet(url)
	if err != nil {
		return err
	}

    json.NewDecoder(body).Decode(&r.Bindings)

	return nil
}

func (r *RabbitMQ) CloneTo(hostport, user, password, vhost string) error {
	if err := r.readQueues(); err != nil {
		return err
	}

	if err := r.readExchanges(); err != nil {
		return err
	}

	if err := r.readBindings(); err != nil {
		return err
	}

	if vhost == "/" {
		vhost = DEFAULT_VHOST
	}

	for _, queue := range r.Queues {
		url := fmt.Sprintf(hostport + PUT_QUEUE_URL, vhost, queue.Name)
		fmt.Println(url)
	}

	for _, exchange := range r.Exchanges {
		url := fmt.Sprintf(hostport + PUT_EXCHANGE_URL, vhost, exchange.Name)
		fmt.Println(url)
	}

	for _, binding := range r.Bindings {
		url := fmt.Sprintf(hostport + POST_BINDING_URL, vhost, binding.Source, binding.Destination)
		fmt.Println(url)
	}

	return nil
}

func main() {

	r := &RabbitMQ{
        BaseUrl : "http://127.0.0.1:9000",
		User    : "guest",
		Password: "guest",
		VHost   : "schedule2",
	}

	r.CloneTo("OtherRabbit:55672", "guest", "guest", "/")

	//fmt.Println(r.Queues)
	//fmt.Println(r.Exchanges)
	//fmt.Println(r.Bindings)


/*
    manager := "http://127.0.0.1:9000/api/queues/schedule2"
    client := &http.Client{}
    req, _ := http.NewRequest("GET", manager, nil)
    req.SetBasicAuth("guest", "guest")
    resp, _ := client.Do(req)

    value := make([]Queue, 0)
    json.NewDecoder(resp.Body).Decode(&value)
    fmt.Println(value)
*/

}

/*

rmq to rmq
rmq to dump file
dump file to rmq

http://127.0.0.1:9000/api/bindings/schedule2
http://127.0.0.1:9000/api/queues/schedule2
http://127.0.0.1:9000/api/exchanges/schedule2

Default vhost %2F -> /

POST to /api/bindings/vhost/e/exchange/q/queue 
Will set binding 

PUT /api/queues/vhost/name
To create queue
{"auto_delete":false,"durable":true,"arguments":[],"node":"rabbit@smacmullen"}

PUT /api/exchanges/vhost/name
To create exchange
{"type":"direct","auto_delete":false,"durable":true,"internal":false,"arguments":[]}

*/


