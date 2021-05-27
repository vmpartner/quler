package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vmpartner/go-mq/v2"
	"github.com/vmpartner/logs"
	"gopkg.in/ini.v1"
	"os"
	"strconv"
	"strings"
)

type Message struct {
	Delivery   amqp.Delivery
	BodyString string
}

func Run(config *ini.File) {

	var ch *amqp.Channel
	var m *mq.RabbitMQ
	var ms Message
	var f *os.File
	var err error
	var messages <-chan amqp.Delivery

	cfgApp := config.Section("app")
	cfgMQ := config.Section("mq")
	cfgFile := config.Section("file")

	// Create output file
	if !cfgFile.Key("message_per_file").MustBool() {
		fileName := cfgFile.Key("path").String()
		err = os.RemoveAll(fileName)
		if err != nil {
			panic(err)
		}
		fileName = strings.ReplaceAll(fileName, "_%", "")
		f, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	// Connect to MQ
	mq.User = cfgMQ.Key("user").String()
	mq.Pass = cfgMQ.Key("pass").String()
	mq.Host = cfgMQ.Key("host").String()
	mq.Port = cfgMQ.Key("port").String()
	mq.PingEachMinute, _ = cfgMQ.Key("ping_each_minute").Int()
	m, err = mq.New()
	if err != nil {
		panic(err)
	}
	ch, err = m.Conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Get messages
	messages, err = ch.Consume(cfgMQ.Key("queue_source").String(), cfgApp.Key("name").String(), false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	j := 0
	limit := cfgMQ.Key("limit_messages").MustInt()
	for message := range messages {

		j++

		// Convert message to read format
		ms = Message{}
		ms.Delivery = message
		ms.BodyString = string(ms.Delivery.Body)
		ms.Delivery.Body = nil

		if cfgFile.Key("message_per_file").MustBool(true) {
			fileName := cfgFile.Key("path").String()
			fileName = strings.Replace(fileName, "%", strconv.Itoa(j), 1)
			_ = os.RemoveAll(fileName)
			f, err = os.Create(fileName)
			if err != nil {
				panic(err)
			}
		}

		// Write to file
		_, err = f.WriteString(fmt.Sprintf("%+v", ms))
		if err != nil {
			panic(err)
		}

		// Close file
		if cfgFile.Key("message_per_file").MustBool(true) {
			_ = f.Close()
		} else {
			_, err = f.WriteString("\n===========================================================================\n")
			if err != nil {
				panic(err)
			}
		}

		// Sync file
		if !cfgFile.Key("message_per_file").MustBool(true) && cfgApp.Key("sync_each_message").MustBool(false) {
			err = f.Sync()
			if err != nil {
				panic(err)
			}
		}

		// Ack message
		if cfgMQ.Key("ack_message").MustBool(false) {
			err = message.Ack(false)
			if err != nil {
				panic(err)
			}
		}

		logs.InfoF("Message %d", message.Timestamp.Unix())

		// Limit
		if limit > 0 && j >= limit {
			return
		}
	}
}
