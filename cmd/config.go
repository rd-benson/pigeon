package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rd-benson/pigeon/common"
)

// Device configuration
type Device struct {
	Topic     string // mqtt topic (suffix) for this device
	ValueType string // value type
}

// Check device configuration value type
func (d *Device) validateValueType() error {
	validValueTypes := []string{"bool", "int8", "int16", "uint32", "uint8", "uint16", "uint32", "float32", "float64"}
	if !common.Contains(d.ValueType, validValueTypes) {
		return fmt.Errorf("invalid value type in configuration")
	}
	return nil
}

// Configuration per site (assumed Easylog transfers data)
type Site struct {
	Name            string
	PublishSeparate bool
	QOS             uint8
	Devices         []Device
}

// Check QOS value
func (s *Site) validateQOS() error {
	validQOS := []uint8{0, 1, 2}
	if !common.Contains(s.QOS, validQOS) {
		return fmt.Errorf("invalid QOS in configuration")
	}
	return nil
}

// Check device configuration
func (s *Site) validateDevices() {
	for _, device := range s.Devices {
		if err := device.validateValueType(); err != nil {
			log.Fatalf("site: %s, contains device configuration errors (invalid value type)", s.Name)
		}
	}
}

// Get mqtt topic to subscribe to
// Easylog topics of the form: site/(device)
func (s *Site) GetFullTopicStrings() []string {
	var topics []string
	if !s.PublishSeparate {
		return []string{s.Name}
	}
	for _, device := range s.Devices {
		topics = append(topics, fmt.Sprintf("%v/%v", s.Name, device.Topic))
	}
	return topics
}

// Mqtt broker configuration
type Mqtt struct {
	Url        string
	Port       uint16
	Encryption string
}

// Get the mqtt broker URI
// Example tcp://broker.hivemq.com:1883
func (b *Mqtt) GetURI() string {
	return fmt.Sprintf("%v://%v:%v", b.Encryption, b.Url, b.Port)
}

// Influxdb cloud configuration
type Database struct {
	Url   string // influxdb cloud url
	Token string // ENV variable containing token (case sensitive)
}

// Get influxdb token value from the ENV variable in configuration file
func (d *Database) GetToken() {
	d.Token = os.Getenv(d.Token)
}

// Total configuration
type Config struct {
	Mqtt     Mqtt
	Database Database
	Sites    []Site
}

// Make sure configuration is sane
func (c *Config) validateConfig() {
	// Encryption: tcp/tls
	mqtt := c.Mqtt
	if !(mqtt.Encryption == "tls" || mqtt.Encryption == "tcp") {
		log.Fatal("supported encryption protocols: tcp/tls")
	}
	// Sites
	sites := c.Sites
	for _, site := range sites {
		site.validateQOS()
		site.validateDevices()
	}

}
