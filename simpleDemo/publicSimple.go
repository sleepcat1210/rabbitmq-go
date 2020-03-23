package main

import (
	"fmt"
	"rabbitmqUtils/utils"
)

func main() {
	rabbitmq := utils.NewRabbitMQSimple("imoocSimple")
	rabbitmq.PublicSimple("hello work")
	fmt.Println("发送成功")

}
