package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/streadway/amqp"
)

type StockService interface {
	ParseStock(reader io.ReadCloser, sender string) error
}

type StockHandler struct {
	brokerUpstream *amqp.Channel
}

func NewStockHandler(brokerConn *amqp.Channel) *StockHandler {
	return &StockHandler{
		brokerUpstream: brokerConn,
	}

}

func (h *StockHandler) ParseStock(data io.ReadCloser, sender string) error {
	reader := csv.NewReader(data)

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("parser was expected two rows, found less than two instead")
	}

	//We are assuming that data is in the same format every time:
	//Symbol     Date      Time    Open   High    Low   Close  Volume
	//AAPL.US 2022-02-11 22:00:08 172.33 173.08 168.04 168.64 98670687

	//We know data is in second row and column 7 so we can avoid a nested loop

	symbol := records[1][0]
	closedPrice := records[1][6]

	payload := fmt.Sprintf("%s quote is $%s per share", symbol, closedPrice)

	if closedPrice == "N/D" {
		payload = "unrecognized stock symbol"
	}

	fmt.Printf("\nCLOSED PRICE: %s\n", closedPrice)

	msg := RabbitMQMessage{
		MsgFor:  sender,
		Payload: payload,
	}

	h.PublishMessageToBroker(msg)

	return nil

}

func (h *StockHandler) PublishMessageToBroker(msg RabbitMQMessage) {

	encoded, err := json.Marshal(&msg)
	if err != nil {
		fmt.Errorf("unable to encode msg: %+v", msg)
		return
	}

	h.brokerUpstream.Publish(
		"",
		QUEUE_NAME,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encoded,
		},
	)
}
