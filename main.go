package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/susinda/constants"
	"github.com/susinda/controllers"
	amqp "github.com/susinda/pkg/active-mq"
	"github.com/susinda/transports"
)

func initializeViper() {
	viper.SetConfigName("appConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't open config file.")
	}
}

func activeMQConfig() map[string]string {
	var address, port, user, password, name, chanone, chantwo, domainkey string
	conf := make(map[string]string)

	address = viper.GetString(constants.DATABASE + "." + constants.ADDRESS)
	if address == "" {
		panic("ADDRESS variable required but not set")
	}
	port = viper.GetString(constants.DATABASE + "." + constants.PORT)
	if port == "" {
		panic("PORT variable required but not set")
	}
	user = viper.GetString(constants.DATABASE + "." + constants.USER)
	if user == "" {
		panic("USER variable required but not set")
	}
	password = viper.GetString(constants.DATABASE + "." + constants.PASS)
	if password == "" {
		panic("PASS variable required but not set")
	}
	name = viper.GetString(constants.DATABASE + "." + constants.NAME)
	if name == "" {
		panic("NAME variable required but not set")
	}

	chanone = viper.GetString(constants.DATABASE + "." + constants.CHANNELONE)
	if chanone == "" {
		panic("NAME variable required but not set")
	}
	chantwo = viper.GetString(constants.DATABASE + "." + constants.CHANNELTWO)
	if chantwo == "" {
		panic("NAME variable required but not set")
	}

	domainkey = viper.GetString(constants.DATABASE + "." + constants.DOMAINKEY)
	if chantwo == "" {
		panic("DOMAINKEY variable required but not set")
	}


	conf[constants.ADDRESS] = address
	conf[constants.PORT] = port
	conf[constants.USER] = user
	conf[constants.PASS] = password
	conf[constants.NAME] = name
	conf[constants.CHANNELONE] = chanone
	conf[constants.CHANNELTWO] = chantwo
	conf[constants.DOMAINKEY] = domainkey

	return conf
}

func Process(inp string) (result string) {
	for _, v := range inp {
		result = string(v) + result
	}
	return
}

func Send(inp string) (result string) {
	return "here i called for external inject backender " + inp
}

func initializeAmqp() {

	fmt.Println("\n1). here um in initializeAmqp ")

	mq := activeMQConfig()
	conntection := fmt.Sprintf("%s:%s", mq[constants.ADDRESS], mq[constants.PORT])

	controllers.AMQP = &amqp.AmqpClient{
		Username: mq[constants.USER],
		Password: mq[constants.PASS], Address: conntection,
		Channels:   []string{mq[constants.CHANNELONE], mq[constants.CHANNELTWO]},
		Custom:     true,
		Processor:  Process,
		Subscriber: Send,
	}

	controllers.AMQP.Connect()
}

func initializeLogrus(logfile string) {
	Formatter := new(log.JSONFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	log.SetFormatter(Formatter)

	var filename = logfile
	dir, errWD := os.Getwd()
	if errWD != nil {
		log.Error("Error occurred while getting present working directory, logging to stderror" + " - " + errWD.Error())
	}
	log.Info(dir)
	f, err := os.OpenFile(dir+"/"+filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error("Error occurred while opening log file, logging to stderror")
	} else {
		multiWriter := io.MultiWriter(os.Stdout, f)
		log.SetOutput(multiWriter)
	}
}

func init() {
	initializeLogrus("amqp.log")
	initializeViper()
	initializeAmqp()
}

func main() {
	// Here i have the string representation of the req
	httpTransport := transports.HttpTransport{}
	httpTransport.Init()

}
