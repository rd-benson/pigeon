package cmd

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type MqttClient struct {
	opts   *mqtt.ClientOptions
	client mqtt.Client
	topics map[string]mqtt.MessageHandler
}

type Pigeon struct {
	config Config
	mqtt   MqttClient
}

// Get config from viper
func (p *Pigeon) InitConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	viper.WatchConfig()
	viper.Unmarshal(&p.config)
	p.config.Database.GetToken()
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(&p.config)
		p.config.Database.GetToken()
	})
	p.config.validateConfig()
}

// Start mqtt client and subscriptions
func (p *Pigeon) MqttStart() {
	p.newMqttClientOptions()
	p.newMqttClient()
	fmt.Printf("connecting to mqtt broker at %v ...", p.config.Mqtt.GetURI())
	p.connectMqttBroker()
	fmt.Println(" connected!")
	p.subscribe()
}

// Stop mqtt client
func (p *Pigeon) MqttStop(quiesce uint) {
	p.mqtt.client.Disconnect(quiesce)
}

func (p *Pigeon) newMqttClientOptions() {
	// Set MQTT client options
	p.mqtt.opts = mqtt.NewClientOptions()
	p.mqtt.opts.AddBroker(p.config.Mqtt.GetURI())
	p.mqtt.opts.SetClientID("pigeon")
	p.mqtt.opts.SetPingTimeout(1 * time.Second)
}

func (p *Pigeon) newMqttClient() {
	p.mqtt.client = mqtt.NewClient(p.mqtt.opts)
}

func (p *Pigeon) connectMqttBroker() {
	token := p.mqtt.client.Connect()
	go func() {
		for {
			select {
			case <-token.Done():
				return
			default:
				fmt.Print(".")
				time.Sleep(2500 * time.Millisecond)
			}
		}
	}()

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (p *Pigeon) subscribe() {
	for _, site := range p.config.Sites {
		for _, topic := range site.GetFullTopicStrings() {
			p.mqtt.client.Subscribe(topic, site.QOS, onMsg)
		}
	}

}
