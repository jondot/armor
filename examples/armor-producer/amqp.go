package main

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Amqp struct {
	Uri        string
	channel    *amqp.Channel
	connection *amqp.Connection
	Exchange   string
}

func NewAmqp(uri string, exchange string) *Amqp {
	return &Amqp{Uri: uri, Exchange: exchange}
}

func (a *Amqp) Publish(queue string, body []byte) error {
	err := a.channel.Publish(
		a.Exchange, // exchange
		queue,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	return err
}
func (a *Amqp) Declare(queue string) error {
	if a.channel == nil {
		return errors.New("Please dial first")
	}

	_, err := a.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return err
}

func (a *Amqp) Dial() error {
	conn, err := amqp.Dial(a.Uri)
	if err != nil {
		log.Printf("AMQP: cannot connect")
		return err
	}
	a.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("AMQP: cannot open channel")
		return err
	}

	a.channel = ch

	return err
}

func (a *Amqp) Close() {
	if a.channel != nil {
		a.channel.Close()
	}
	if a.connection != nil {
		a.connection.Close()
	}
}

func (a *Amqp) AutoHeal() {
	notifyClose := a.connection.NotifyClose(make(chan *amqp.Error))
	for {
		select {
		case <-notifyClose:
			for {
				err := a.Dial()
				if err == nil {
					return
				}
				time.Sleep(2 * time.Second)
			}
		}
	}
	return
}
