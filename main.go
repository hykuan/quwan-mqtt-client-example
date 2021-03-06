package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cisco/senml"
	"github.com/eclipse/paho.mqtt.golang"
)

const (
	defTopic    = "channels/1/messages"
	defBroker   = "ssl://[YOUR_NGINX_HOST]:8883"
	defUser     = "1"
	defPassword = "YOUR_THING_TOKEN"

	envTopic    = "TOPIC"
	envBroker   = "BROKER"
	envUser     = "USER"
	envPassword = "PASSWORD"
)

type config struct {
	topic    string
	broker   string
	user     string
	password string
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func loadConfig() config {
	return config{
		topic:    getEnv(envTopic, defTopic),
		broker:   getEnv(envBroker, defBroker),
		user:     getEnv(envUser, defUser),
		password: getEnv(envPassword, defPassword),
	}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	raw, err := senml.Decode(msg.Payload(), senml.JSON)
	if err != nil {
		fmt.Println("got error: ", err)
		return
	}

	normalizedData := senml.Normalize(raw)
	fmt.Println("Normalized data: ", normalizedData)
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	config := loadConfig()

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.broker)
	opts.SetUsername(config.user)
	opts.SetPassword(config.password)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(config.topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for i := 0; i < 1; i++ {
		// the message should be senml format, the following json string is example from mainflux
		// [{"bn":"some-base-name:","bt":1.276020076001e+09, "bu":"A","bver":5, "n":"voltage","u":"V","v":125}, {"n":"current","t":-5,"v":1.2}, {"n":"current","t":-4,"v":0.382131}]
		record := []senml.SenMLRecord{
			{
				BaseName: "Device_1",
				BaseTime: float64(time.Now().UnixNano()),
				Name:     "voltage",
				Unit:     "V",
				Value:    func(value float64) *float64 { return &value }(110),
			},
			{
				BaseName: "Device_1",
				BaseTime: float64(time.Now().UnixNano()),
				Name:     "temperature",
				Unit:     "celsius",
				Value:    func(value float64) *float64 { return &value }(60),
			},
		}

		text, _ := json.Marshal(record)

		token := c.Publish(config.topic, 0, false, string(text))
		token.Wait()
		time.Sleep(1 * time.Second)
	}

	<-ch
}
