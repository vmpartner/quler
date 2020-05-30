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

var ch *amqp.Channel
var m *mq.RabbitMQ
var ms Message
var f *os.File
var config *ini.File
var err error
var messages <-chan amqp.Delivery

func main() {

	// Parse config
	config, err = ini.Load("./app.conf")
	if err != nil {
		panic(err)
	}
	cfgApp := config.Section("app")
	cfgMQ := config.Section("mq")
	cfgFile := config.Section("file")

	// Create output file
	if !cfgFile.Key("message_per_file").MustBool() {
		os.RemoveAll(cfgFile.Key("path").String())
		f, err = os.Create(cfgFile.Key("path").String())
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

	// Get messages
	messages, err = ch.Consume(cfgMQ.Key("queue_source").String(), cfgApp.Key("name").String(), false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	k := 0
	j := 0
	for message := range messages {

		k++
		j++

		// Convert message to read format
		ms = Message{}
		ms.Delivery = message
		ms.BodyString = string(ms.Delivery.Body)
		ms.Delivery.Body = nil

		if cfgFile.Key("message_per_file").MustBool(true) {
			fileName := cfgFile.Key("path").String()
			fileName = strings.Replace(fileName, "%", strconv.Itoa(j), 1)
			os.RemoveAll(fileName)
			f, err = os.Create(fileName)
		}

		// Write to file
		_, err = f.WriteString(fmt.Sprintf("%+v", ms))
		if err != nil {
			panic(err)
		}

		// Close file
		if cfgFile.Key("message_per_file").MustBool(true) {
			f.Close()
		} else {
			_, err = f.WriteString("\n===========================================================================\n")
			if err != nil {
				panic(err)
			}
		}

		// Sync file
		if k >= cfgApp.Key("sync_each_message").MustInt(1) {
			k = 0
			f.Sync()
		}

		// Ack message
		if cfgMQ.Key("ack_message").MustBool(false) {
			err = message.Ack(false)
		}

		logs.InfoF("Message %d", message.Timestamp.Unix())
	}
}
