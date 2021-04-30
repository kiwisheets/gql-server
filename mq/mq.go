package mq

import (
	"fmt"

	c_mq "github.com/cheshir/go-mq"
	"github.com/kiwisheets/util"
	"github.com/sirupsen/logrus"
)

type MQ struct {
	mq           c_mq.MQ
	CreateClient c_mq.SyncProducer
	UpdateClient c_mq.SyncProducer
}

func (m *MQ) Close() {
	m.mq.Close()
}

func setupProducers(m *MQ) {
	var err error
	m.CreateClient, err = m.mq.SyncProducer("client_create")
	if err != nil {
		panic(fmt.Errorf("failed to create producer: client_create: %s", err))
	}

	m.UpdateClient, err = m.mq.SyncProducer("client_update")
	if err != nil {
		panic(fmt.Errorf("failed to create producer: client_update: %s", err))
	}
}

func setupConsumers(m *MQ) {

}

func Init() *MQ {
	m := MQ{
		mq: util.NewMQ(),
	}
	setupConsumers(&m)
	setupProducers(&m)

	go handleMQErrors(m.mq.Error())

	return &m
}

func handleMQErrors(errors <-chan error) {
	for err := range errors {
		logrus.Errorf("mq error: %s", err)
	}
}
