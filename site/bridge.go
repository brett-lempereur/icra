package main

import (
	"encoding/json"
	"log"
	"net/url"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

const (
	mqttQos     = 0   // Quality-of-service level for subscriptions
	mqttQuiesce = 500 // Graceful shutdown delay
)

// MessagingBridge provides a bridge between a message broker and websocket
// clients connected to the visit service.
type MessagingBridge struct {
	client *mqtt.Client  // Message broker client.
	topics []string      // List of topics to subscribe to.
	vs     *VisitService // Destination for messages.
}

// NewMessagingBridge returns an initialised messaging bridge that will forward
// messages from the given topics to the given visit service.
func NewMessagingBridge(topics []string, vs *VisitService) *MessagingBridge {
	return &MessagingBridge{nil, topics, vs}
}

// Connect to the given message broker using the given identifier.
func (mb *MessagingBridge) Connect(uri, id string) error {

	log.Printf("Connecting to message broker '%s' as client '%s'", uri, id)

	// Attempt to parse the address.
	u, err := url.Parse(uri)
	if err != nil {
		log.Printf("Could not parse broker address: %s\n", uri)
		return err
	}

	// Strip authentication details from the address.
	ui := u.User
	u.User = nil

	// Construct the connection parameters.
	opts := mqtt.NewClientOptions()
	opts.AddBroker(u.String())
	opts.SetClientID(id)
	if ui != nil {
		opts.SetUsername(ui.Username())
		if pw, ok := ui.Password(); ok {
			opts.SetPassword(pw)
		}
	}
	opts.SetDefaultPublishHandler(mb.onMessage)
	opts.SetOnConnectHandler(mb.onConnect)
	opts.SetConnectionLostHandler(mb.onConnectionLost)

	// Attempt to connect to the broker.
	mb.client = mqtt.NewClient(opts)
	token := mb.client.Connect()
	token.Wait()
	return token.Error()

}

// Close the connection to the message broker.
func (mb *MessagingBridge) Close() {
	mb.client.Disconnect(mqttQuiesce)
}

// OnMessage is called when a message is received on one of the subscribed
// topics.
func (mb *MessagingBridge) onMessage(c *mqtt.Client, msg mqtt.Message) {

	var v Visit

	if err := json.Unmarshal(msg.Payload(), &v); err != nil {
		log.Printf("Could not decode message '%s': %s\n", string(msg.Payload()), err)
		return
	}
	mb.vs.Send(v)

}

// OnConnect is called when the message broker connection is established.
func (mb *MessagingBridge) onConnect(c *mqtt.Client) {

	log.Print("Connected to message broker, subscribing to topics...\n")

	for _, t := range mb.topics {
		token := c.Subscribe(t, mqttQos, nil)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Could not subscribe to '%s': %s", t, token.Error())
		}
	}
	log.Print("Finished subscribing to topics\n")

}

// OnConnectionLost is called when the connection to the message broker is
// lost.
func (mb *MessagingBridge) onConnectionLost(c *mqtt.Client, err error) {
	log.Printf("Lost connected to message broker: %s\n", err)
}
