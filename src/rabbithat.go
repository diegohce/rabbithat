package main

import (
	"rabbit"
)




func main() {

	r := &rabbit.RabbitMQ{
        BaseUrl : "http://127.0.0.1:9000",
		User    : "guest",
		Password: "guest",
		VHost   : "schedule2",
	}

	r.Collect()

	r.DumpTo("dump.json")
	r.CloneTo("10.0.3.214:55672", "guest", "guest", "diego")

/*
	r, _ := rabbit.CollectFromFile("dump.json")
	r.CloneTo("10.0.3.214:55672", "guest", "guest", "diego")

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
To create binding 
{"routing_key":"my_routing_key","arguments":[]}


PUT /api/queues/vhost/name
To create queue
{"auto_delete":false,"durable":true,"arguments":[],"node":"rabbit@smacmullen"}

PUT /api/exchanges/vhost/name
To create exchange
{"type":"direct","auto_delete":false,"durable":true,"internal":false,"arguments":[]}

*/


