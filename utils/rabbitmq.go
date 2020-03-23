package utils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//url 格式
const MQURL = "amqp://ubantu:root@192.168.0.103:5672/"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//连接信息
	Mqurl string
}

//创建结构体实例
func NewRabbitMq(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	if err != nil {
		rabbitmq.failOnErr(err, "创建连接错误")
	}
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		rabbitmq.failOnErr(err, "获取channel 失败")
	}
	return rabbitmq
}

//断开连接
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建简单模式
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMq(queueName, "", "")
}

//简单模式生产
func (r *RabbitMQ) PublicSimple(message string) {
	//1首先申请队列如果不存在就自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//发送消息到队列中
	r.channel.Publish(
		r.Exchange, //交换机
		r.QueueName,
		false, //true 找不到队列回退发送者
		false, //如果没有绑定消费者则会发给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

}

//消费消息
func (r *RabbitMQ) ConsumeSimple() {
	//1首先申请队列如果不存在就自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//接收消息
	msg, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否应答自动删除消息
		true,
		//是否排他性
		false,
		//如果设置true表示发送消息
		false,
		//是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//消费消息
	forever := make(chan bool)
	go func() {
		for d := range msg {
			fmt.Println(string(d.Body))
		}
	}()

	<-forever
}

//创建订阅模式
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	return NewRabbitMq("", exchangeName, "")
}

//订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	//尝试连接交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"fanout",   //交换机类型
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare exchage")
	r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

//订阅模式消费
func (r *RabbitMQ) ConsumePub() {
	//尝试连接交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, //交换机名称
		"fanout",   //交换机类型
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare exchage")
	//尝试创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机产生对列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a queue")
	//绑定对列到exchange
	r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range message {
			fmt.Println(string(d.Body))
		}
	}()
	<-forever
}
