package cmd

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type EasylogVariable struct {
	Name         string
	Value        string
	ParseValueAs string
}

// Easylog message structure, seperate publish for each device: true
// Example json payload
//
//	{
//	    "date": "2022-08-10T07:56:55.0Z",
//	    "id": "AAAAAA010000000002100466",
//	    "device": "ele",
//	    "vars": [
//	        { "name": "mels_S", "value": "1.299999952316284"},
//	        {"name": "mels_N", "value": "0.200000002980232"},
//	        {"name": "lig_S", "value": "7.400000095367432"}
//	    ]
//	}
type EasylogPublishSeparate struct {
	Date   string
	Id     string
	Device string
	Vars   []EasylogVariable
}

// Easylog message structure, seperate publish for each device: false
// Example json payload
//
//	{
//	    "date": "2022-08-10T08:09:45.0Z",
//	    "id": "AAAAAA010000000002100466",
//	    "vars": [
//	        { "name": "ele > mels_S", "value": "7.150000095367432"},
//	        { "name": "ele > mels_N", "value": "4.644999980926514"},
//	        { "name": "ele > lig_S", "value": "18.100000381469727"},
//	    ]
//	}
type EasylogPublishTogether struct {
	Date   string
	Id     string
	Device string
	Vars   []EasylogVariable
}

func onMsg(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
