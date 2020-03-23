package main

import (
	"fmt"
	"rabbitmqUtils/utils"
)

func main() {
	rabbitmq := utils.NewRabbitMQPubSub("imoocSimple")
	rabbitmq.PublishPub("hello work")
	fmt.Println("发送成功")

}
