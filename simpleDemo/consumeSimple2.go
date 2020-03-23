package main

import (
	"fmt"
	"rabbitmqUtils/utils"
)

func main() {
	rabbitmq := utils.NewRabbitMQPubSub("imoocSimple")
	rabbitmq.ConsumePub()
	fmt.Println("消费成功")
}
