// Web server that hosts the static web-browsing monitor and provides a bridge
// between the message broker and websockets.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
)

// Specification defines the environment variables required to configure
// the application.
type Specification struct {
	Port       string `default:"8080"`  // Web server port.
	Broker     string `required:"true"` // Address of the message broker.
	Topic      string `required:"true"` // Message broker topic.
	BufferSize int    `default:"100"`   // Number of visits to buffer.
	CacheSize  int    `default:"50"`    // Number of visits to cache.
}

// Main entry point that starts the messaging bridge, visit service, and web
// server.
func main() {

	var s Specification
	var port string

	// Load environment variables.
	if err := envconfig.Process("icra", &s); err != nil {
		log.Fatalf("Could not read environment variables: %s", err.Error())
	}
	if port = os.Getenv("VCAP_APP_PORT"); len(port) == 0 {
		log.Printf("Warning, VCAP_APP_PORT not set. Defaulting to %+v\n", s.Port)
		port = s.Port
	}

	// Generate a client identifier.
	id := fmt.Sprintf("icra-client-%d", os.Getpid())

	// Configure the visit service and messaging bridge.
	vs := NewVisitService(s.BufferSize, s.CacheSize)
	mb := NewMessagingBridge([]string{s.Topic}, vs)

	// Start the visit service.
	go vs.Serve()

	// Connect to the message broker.
	mb.Connect(s.Broker, id)

	// Configure the web server.
	http.Handle("/ws/browsing", vs.MakeHandler())
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Launch the web server.
	http.ListenAndServe(":"+port, nil)

}
