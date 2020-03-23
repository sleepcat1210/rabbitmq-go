package main

import (
	"fmt"
	"rabbitmqUtils/utils"
)

func main() {
	rabbitmq := utils.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
	fmt.Println("消费成功")
}
