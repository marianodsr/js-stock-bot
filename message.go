package main

type RabbitMQMessage struct {
	MsgFor  string `json:"msg_for"`
	Payload string `json:"payload"`
}
