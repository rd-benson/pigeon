package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rd-benson/pigeon/cmd"
)

var pigeon cmd.Pigeon

func init() {
	pigeon.InitConfig()
}

func main() {

	// Mqtt
	pigeon.MqttStart()

	// Handle exit cleanly
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		if s := <-sigs; s != nil {
			fmt.Printf("Received signal: %s\n", s)
			pigeon.MqttStop(250)
			exit()
		}
	}()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}

}

/* func main() {

	// Set MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.GetURI())
	opts.SetClientID("pigeon")
	opts.SetPingTimeout(1 * time.Second)

	// Create client
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe
	for i := 0; i < len(cfg.Sites); i++ {
		for j := 0; j < len(cfg.Sites[i].Devices); i++ {
			if token := c.Subscribe(cfg.Sites[i].Devices[j].Topic, 0, onMsg); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}
		}
	}

	// Handle exit cleanly
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		if s := <-sigs; s != nil {
			fmt.Printf("Received signal: %s\n", s)
			c.Disconnect(250)
			exit()
		}
	}()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}

} */

func exit() {
	fmt.Println("exiting ...")
	os.Exit(1)
}
