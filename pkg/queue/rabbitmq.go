package queue

import (
    "encoding/json"
    "github.com/streadway/amqp"
)

type ImageProcessingMessage struct {
    ProductID uint     `json:"product_id"`
    Images    []string `json:"images"`
}

type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    // Declare the queue
    _, err = ch.QueueDeclare(
        "image_processing", // queue name
        true,              // durable
        false,             // delete when unused
        false,             // exclusive
        false,             // no-wait
        nil,              // arguments
    )
    if err != nil {
        return nil, err
    }

    return &RabbitMQ{
        conn:    conn,
        channel: ch,
    }, nil
}

func (r *RabbitMQ) PublishMessage(msg ImageProcessingMessage) error {
    body, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    return r.channel.Publish(
        "",                // exchange
        "image_processing", // routing key
        false,             // mandatory
        false,             // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}

func (r *RabbitMQ) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
    return r.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
}

func (r *RabbitMQ) Close() error {
    if err := r.channel.Close(); err != nil {
        return err
    }
    return r.conn.Close()
}