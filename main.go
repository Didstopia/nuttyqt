package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
	nut "github.com/robbiet480/go.nut"
	"github.com/sirupsen/logrus"
)

// Config holds the configuration for the application.
type Config struct {
	// MQTT broker protocol. Defaults to "tcp".
	MQTTBrokerProtocol string

	// MQTT broker host. Defaults to "localhost".
	MQTTBrokerHost string

	// MQTT broker port. Defaults to 1883.
	MQTTBrokerPort int

	// MQTT client ID. Defaults to "nuttyqt".
	MQTTClient string

	// MQTT topic. Defaults to "nuttyqt".
	MQTTTopic string

	// MQTT username. Defaults to "".
	MQTTUser string

	// MQTT password. Defaults to "".
	MQTTPass string

	// NUT server host. Defaults to "localhost".
	NUTServerHost string

	// NUT server port. Defaults to 3493.
	NUTServerPort int

	// NUT username. Defaults to "".
	NUTUser string

	// NUT password. Defaults to "".
	NUTPass string

	// Update interval in seconds. Defaults to 60.
	UpdateInterval int

	// Verbose logging. Defaults to false.
	Verbose bool
}

var (
	config = Config{
		MQTTBrokerProtocol: "tcp",
		MQTTBrokerHost:     "localhost",
		MQTTBrokerPort:     1883,
		MQTTClient:         "nuttyqt",
		MQTTTopic:          "nuttyqt",
		MQTTUser:           "",
		MQTTPass:           "",

		NUTServerHost: "localhost",
		NUTServerPort: 3493,
		NUTUser:       "",
		NUTPass:       "",

		UpdateInterval: 60,
		Verbose:        false,
	}

	// MQTT
	// mqttClient *mqtt.Client
	mqttClient mqtt.Client

	// NUT
	nutClient *nut.Client
	upsDevice *nut.UPS
)

// Logger for the application.
var log = logrus.New()

// Get the value of an environment variable or return a default value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Load the configuration from environment variables.
func LoadConfig() {
	// MQTT
	config.MQTTBrokerProtocol = GetEnv("MQTT_BROKER_PROTOCOL", config.MQTTBrokerProtocol)
	config.MQTTBrokerHost = GetEnv("MQTT_BROKER_HOST", config.MQTTBrokerHost)
	config.MQTTBrokerPort, _ = strconv.Atoi(GetEnv("MQTT_BROKER_PORT", strconv.Itoa(config.MQTTBrokerPort)))
	config.MQTTClient = GetEnv("MQTT_CLIENT", config.MQTTClient)
	config.MQTTTopic = GetEnv("MQTT_TOPIC", config.MQTTTopic)
	config.MQTTUser = GetEnv("MQTT_USER", config.MQTTUser)
	config.MQTTPass = GetEnv("MQTT_PASS", config.MQTTPass)

	// NUT
	config.NUTServerHost = GetEnv("NUT_SERVER", config.NUTServerHost)
	config.NUTServerPort, _ = strconv.Atoi(GetEnv("NUT_PORT", strconv.Itoa(config.NUTServerPort)))
	config.NUTUser = GetEnv("NUT_USER", config.NUTUser)
	config.NUTPass = GetEnv("NUT_PASS", config.NUTPass)

	// Other
	config.UpdateInterval, _ = strconv.Atoi(GetEnv("UPDATE_INTERVAL", strconv.Itoa(config.UpdateInterval)))
	config.Verbose, _ = strconv.ParseBool(GetEnv("VERBOSE", strconv.FormatBool(config.Verbose)))
}

// FIXME: Do we even need this at this point?!
var mqttMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debug("TOPIC: %s\n", msg.Topic())
	log.Debug("MSG: %s\n", msg.Payload())
}

// Get the UPS device from NUT.
func GetUPS() *nut.UPS {
	// FIXME: Can we keep the connection open AND handle reconnects?
	// Create a new NUT client and connect to the server.
	var isNewClient bool
	if nutClient == nil {
		log.Info(fmt.Sprintf("Connecting to NUT server at %s:%d ...", config.NUTServerHost, config.NUTServerPort))
		client, connectErr := nut.Connect(config.NUTServerHost, config.NUTServerPort)
		if connectErr != nil {
			log.Fatal("Failed to connect to NUT server: ", connectErr)
		}
		nutClient = &client
		isNewClient = true
	} else {
		log.Debug("Reusing existing NUT client ...")
	}

	// Authenticate with the NUT server.
	if isNewClient {
		log.Debug("Authenticating with NUT server ...")
		if config.NUTUser != "" && config.NUTPass != "" {
			_, authErr := nutClient.Authenticate(config.NUTUser, config.NUTPass)
			if authErr != nil {
				log.Fatal("Failed to authenticate with NUT server: ", authErr)
			}
		} else {
			log.Debug("No NUT credentials provided. Skipping authentication ...")
		}
	}

	// Get a list of all available UPS devices.
	log.Debug("Getting a list of all UPS devices ...")
	upsList, listErr := nutClient.GetUPSList()
	if listErr != nil {
		log.Fatal("Failed to get a list of UPS devices: ", listErr)
	}

	// TODO: Return all UPS devices instead and send them all to MQTT?
	// Return the first UPS in the list.
	log.Debug("Return the first UPS device ...")
	return &upsList[0]
}

// Create a new MQTT client and connect to the MQTT broker.
func CreateMQTTClient() {
	//
	// NOTE: Usage examples for the Paho MQTT client:
	//
	// https://github.com/eclipse/paho.mqtt.golang/tree/master/cmd
	//

	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	// mqtt.ERROR = log.New(os.Stdout, "", 0)
	mqtt.ERROR = logrus.New() // FIXME: Ideally redirect these to our existing logger, instead of a new one.

	// TODO: Does this library handle reconnecting automatically?
	// TODO: Set up LWT (Last Will and Testament) message.
	log.Debug("Setting up MQTT client ...")
	opts := mqtt.NewClientOptions()
	opts.SetConnectRetry(false)
	opts.SetAutoReconnect(true)
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", config.MQTTBrokerProtocol, config.MQTTBrokerHost, config.MQTTBrokerPort))
	opts.SetClientID(config.MQTTClient)
	opts.SetDefaultPublishHandler(mqttMessageHandler)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(mqttMessageHandler)
	opts.SetPingTimeout(1 * time.Second)

	// FIXME: If mqttClient can't be a pointer, how can we check if it's nil and not recreate it?
	log.Debug("Creating MQTT client ...")
	mqttClient = mqtt.NewClient(opts)

	log.Info(fmt.Sprintf("Connecting to MQTT broker at %s://%s:%d ...", config.MQTTBrokerProtocol, config.MQTTBrokerHost, config.MQTTBrokerPort))
	if token := mqttClient.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		log.Fatal("Failed to connect to MQTT broker: ", token.Error())
	}
}

func main() {
	// Load the configuration from the environment.
	LoadConfig()

	// Setup logging.
	log.Out = os.Stdout
	if config.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	// Create the MQTT client.
	CreateMQTTClient()

	// Start the update loop in a goroutine.
	go Update()

	// Setup graceful shutdown and wait for SIGINT or SIGTERM.
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulShutdown

	// Create a context with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// TODO: Actually utilize both the context and cancel,
	//       eg. in Close() but maybe in the Update() loop too?

	// Disconnect from the NUT server and MQTT broker when we're done.
	defer Close(ctx, cancel)
}

// Update loop that runs at the configured interval, updating the UPS device
// and sending the UPS device data to the MQTT broker.
func Update() {
	for {
		if !mqttClient.IsConnected() || !mqttClient.IsConnectionOpen() {
			log.Debug("MQTT client is not connected, skipping update ...")
		} else {
			// Get the UPS device.
			log.Debug("Updating UPS device ...")
			upsDevice = GetUPS()

			// Serialize the UPS device to JSON.
			log.Debug("Serializing UPS device to JSON ...")
			upsDeviceJSON, jsonErr := json.Marshal(upsDevice)
			if jsonErr != nil {
				log.Fatal("Failed to serialize UPS device to JSON: ", jsonErr)
			}

			// Send the data to the MQTT broker.
			log.Debug("Sending data to MQTT broker ...")
			mqttMessageToken := mqttClient.Publish(config.MQTTTopic, 0, false, upsDeviceJSON)
			mqttMessageToken.WaitTimeout(5 * time.Second)
			if mqttMessageToken.Error() != nil {
				log.Fatal("Failed to send data to MQTT broker: ", mqttMessageToken.Error())
			}
		}

		// Wait 15 seconds before updating again.
		log.Debug("Sleeping for ", config.UpdateInterval, " seconds ...")
		time.Sleep(time.Duration(config.UpdateInterval) * time.Second)
	}
}

// Close the application.
func Close(ctx context.Context, cancel context.CancelFunc) {
	log.Info("Shutting down ...")

	// Disconnect from the NUT server and MQTT broker.
	if err := CloseNUT(); err != nil {
		log.Warn("Failed to close NUT connection: ", err)
	}
	if err := CloseMQTT(); err != nil {
		log.Warn("Failed to close MQTT connection: ", err)
	}

	// Cancel the context.
	cancel()

	log.Info("Shutdown complete, terminating ...")
	os.Exit(0)
}

// Close the NUT client.
func CloseNUT() error {
	if nutClient != nil {
		log.Debug("Disconnecting from NUT server ...")
		_, err := nutClient.Disconnect()
		if err != nil {
			return err
		}
	} else {
		log.Debug("No NUT client to disconnect from, skipping ...")
	}
	return nil
}

// Close the MQTT client.
func CloseMQTT() error {
	if mqttClient != nil {
		log.Debug("Disconnecting from MQTT broker ...")
		mqttClient.Disconnect(250)
	} else {
		log.Debug("No MQTT client to disconnect from, skipping ...")
	}
	return nil
}
