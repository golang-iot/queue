package queue

import (
	"log"
	"errors"
	"github.com/streadway/amqp"
)

func err(msg string){
	err := errors.New(msg)
	log.Fatal(err)
}
	
type queueManager interface {
	Init(url string)
    GetChannel(name string) *amqp.Channel
    GetQueue(name string) amqp.Queue
	Consume(name string)
	Send(name string, message string) error
	Close()
}

type Queue struct {
	Connection *amqp.Connection
    Channel map[string]*amqp.Channel
	Queues map[string]amqp.Queue
}

func (q *Queue) Init(url string){
	var err error
	q.Connection, err = amqp.Dial(url)
	//fmt.Printf("Init: %+v\n", q)
	failOnError(err, "Failed to connect to RabbitMQ")
	q.Channel = make(map[string]*amqp.Channel)
	q.Queues = make(map[string]amqp.Queue)
}

func (q *Queue) GetChannel(name string) *amqp.Channel{
	ch, err := q.Connection.Channel()
	q.Channel[name] = ch
	failOnError(err, "Failed to open a channel")
	return ch
}

func (q *Queue) GetQueue(name string) amqp.Queue{
	var err error
	q.Queues[name], err = q.GetChannel(name).QueueDeclare(
		name, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q.Queues[name]
}

func (q Queue) Consume(name string) <-chan amqp.Delivery{
	//log.Printf("Consume %+v\n", q)
	_, ok  := q.Queues[name]
	if !ok{
		err("Declare a queue before consuming "+name)
	}
	
	_, ok  = q.Channel[name]
	if !ok{
		err("Declare a queue before consuming"+name)
	}
	
	messages, err := q.Channel[name].Consume(
		name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to listen to messages")
	return messages
}

func (q Queue) Send(name string, message string) error {
	//log.Printf("Send %+v\n", q)
	_, ok  := q.Queues[name]
	if !ok{
		err("Declare a queue before sending"+name)
	}
	err := q.Channel[name].Publish(
		"",     // exchange
		name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
		
	return err
}

func (q Queue) Close(){
	for _, ch := range q.Channel{
		ch.Close()
	}
	q.Connection.Close()
}