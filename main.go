package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"
)

const (
	iniPath = "./app.conf"
)

func main() {

	// Parse config
	config, err := ini.Load(iniPath)
	if err != nil {
		panic(err)
	}

	a := app.New()
	w := a.NewWindow("QULER")
	w.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})

	// app name
	appName := config.Section("app").Key("name").String()
	bindAppName := binding.NewString()
	_ = bindAppName.Set(appName)

	// sync_each_message
	appSyncEachMessage := config.Section("app").Key("sync_each_message").MustBool()
	bindAppSyncEachMessage := binding.NewBool()
	_ = bindAppSyncEachMessage.Set(appSyncEachMessage)

	// user
	mqUser := config.Section("mq").Key("user").String()
	bindMqUser := binding.NewString()
	_ = bindMqUser.Set(mqUser)

	// pass
	mqPass := config.Section("mq").Key("pass").String()
	bindMqPass := binding.NewString()
	_ = bindMqPass.Set(mqPass)

	// host
	mqHost := config.Section("mq").Key("host").String()
	bindMqHost := binding.NewString()
	_ = bindMqHost.Set(mqHost)

	// port
	mqPort := config.Section("mq").Key("port").String()
	bindMqPort := binding.NewString()
	_ = bindMqPort.Set(mqPort)

	// queue_source
	mqQueueSource := config.Section("mq").Key("queue_source").String()
	bindMqQueueSource := binding.NewString()
	_ = bindMqQueueSource.Set(mqQueueSource)

	// ack_message
	mqAckMessage := config.Section("mq").Key("ack_message").MustBool()
	bindMqAckMessage := binding.NewBool()
	_ = bindMqAckMessage.Set(mqAckMessage)

	// limit_messages
	mqLimitMessages := config.Section("mq").Key("limit_messages").String()
	bindMqLimitMessages := binding.NewString()
	_ = bindMqLimitMessages.Set(mqLimitMessages)

	// message_per_file
	appMessagePerFile := config.Section("file").Key("message_per_file").MustBool()
	bindAppMessagePerFile := binding.NewBool()
	_ = bindAppMessagePerFile.Set(appMessagePerFile)

	// path
	mqPath := config.Section("file").Key("path").String()
	bindMqPath := binding.NewString()
	_ = bindMqPath.Set(mqPath)

	saveConfig := func() {
		// app name
		s, _ := bindAppName.Get()
		config.Section("app").Key("name").SetValue(s)

		// sync_each_message
		b, _ := bindAppSyncEachMessage.Get()
		if b {
			s = "true"
		} else {
			s = "false"
		}
		config.Section("app").Key("sync_each_message").SetValue(s)

		// user
		s, _ = bindMqUser.Get()
		config.Section("mq").Key("user").SetValue(s)

		// pass
		s, _ = bindMqPass.Get()
		config.Section("mq").Key("pass").SetValue(s)

		// host
		s, _ = bindMqHost.Get()
		config.Section("mq").Key("host").SetValue(s)

		// port
		s, _ = bindMqPort.Get()
		config.Section("mq").Key("port").SetValue(s)

		// queue_source
		s, _ = bindMqQueueSource.Get()
		config.Section("mq").Key("queue_source").SetValue(s)

		// ack_message
		b, _ = bindMqAckMessage.Get()
		if b {
			s = "true"
		} else {
			s = "false"
		}
		config.Section("mq").Key("ack_message").SetValue(s)

		// limit_messages
		s, _ = bindMqLimitMessages.Get()
		config.Section("mq").Key("limit_messages").SetValue(s)

		// message_per_file
		b, _ = bindAppMessagePerFile.Get()
		if b {
			s = "true"
		} else {
			s = "false"
		}
		config.Section("file").Key("message_per_file").SetValue(s)

		// limit_messages
		s, _ = bindMqPath.Get()
		config.Section("file").Key("path").SetValue(s)

		config.SaveTo(iniPath)

		// Parse config
		config, err = ini.Load(iniPath)
		if err != nil {
			panic(err)
		}
	}

	w.SetContent(container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Имя приложения"),
			widget.NewEntryWithData(bindAppName),
		),
		container.NewGridWithColumns(2,
			widget.NewCheckWithData("Скидывать на диск каждое сообщение", bindAppSyncEachMessage),
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Пользователь"),
			widget.NewEntryWithData(bindMqUser),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Пароль"),
			widget.NewEntryWithData(bindMqPass),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Хост"),
			widget.NewEntryWithData(bindMqHost),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Порт"),
			widget.NewEntryWithData(bindMqPort),
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Очередь"),
			widget.NewEntryWithData(bindMqQueueSource),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Лимит сообщений"),
			widget.NewEntryWithData(bindMqLimitMessages),
		),
		container.NewGridWithColumns(2,
			widget.NewCheckWithData("Подтверждать сообщение при получении", bindMqAckMessage),
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewCheckWithData("Отдельный файл для каждого сообщения", bindAppMessagePerFile),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Путь к файлам"),
			widget.NewEntryWithData(bindMqPath),
		),
		widget.NewSeparator(),
		widget.NewSeparator(),
		container.NewGridWithColumns(3,
			widget.NewLabel(""),
			widget.NewButton("Запустить", func() {
				saveConfig()
				Run(config)
			}),
			widget.NewLabel(""),
		),
	))

	w.ShowAndRun()
}
