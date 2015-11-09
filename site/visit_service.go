package main

import (
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

// VisitService provides a websocket endpoint that streams visits as JSON
// messages.
type VisitService struct {
	clients    map[*connection]bool // Set of registered connections.
	visits     chan Visit           // Channel for incoming visits.
	register   chan *connection     // Register requests from new connections.
	unregister chan *connection     // Unregister requests from connections.
	shutdown   chan struct{}        // Shutdown channel that terminates service.
	cache      *VisitCache          // Cache of recent visits.
	wg         sync.WaitGroup       // Wait group for service and connections.
	bufferSize int                  // Connection visit buffer size.
}

// NewVisitService returns an initialised service with the given connection
// buffer and cache sizes.
func NewVisitService(bufferSize, cacheSize int) *VisitService {

	// Verify that the buffer size is at least twice as large as the cache
	// size.  When a client first connects, we fill its buffer with a replay
	// of cached messages; this avoids blocking the service but means that we
	// might accidentally disconnect clients if the buffer isn't sufficiently
	// large.
	if cacheSize > bufferSize/2 {
		panic("Buffer size must be twice as large as cache")
	}

	// Initialise and return the new visit service.
	return &VisitService{
		clients:    make(map[*connection]bool),
		visits:     make(chan Visit),
		register:   make(chan *connection),
		unregister: make(chan *connection),
		shutdown:   make(chan struct{}),
		cache:      NewVisitCache(cacheSize),
		bufferSize: bufferSize,
	}

}

// Send a visit to all connected clients.
func (vs *VisitService) Send(v Visit) {
	vs.visits <- v
}

// MakeHandler returns a handler that upgrades connections to the websocket
// protocol and registers them with the service.
func (vs *VisitService) MakeHandler() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {

		log.Printf("Client connected to visit service: %s\n", ws.RemoteAddr())

		c := &connection{
			ws:       ws,
			visits:   make(chan Visit, vs.bufferSize),
			shutdown: make(chan struct{}),
		}

		vs.register <- c
		defer func() { vs.unregister <- c }()
		vs.wg.Add(1)
		defer vs.wg.Done()
		c.serve()

		log.Printf("Client disconnected from visit service: %s\n", ws.RemoteAddr())

	})
}

// Serve accepts registration requests and sends clients cached and new
// visits.
func (vs *VisitService) Serve() {

	log.Printf("Visit service started")

	vs.wg.Add(1)
	defer vs.wg.Done()

	for {
		select {
		case c := <-vs.register:
			vs.clients[c] = true
			for v := range vs.cache.Iterator() {
				c.visits <- v
			}
		case c := <-vs.unregister:
			if _, ok := vs.clients[c]; ok {
				delete(vs.clients, c)
				close(c.shutdown)
			}
		case v := <-vs.visits:
			vs.cache.Add(v)
			for c := range vs.clients {
				select {
				case c.visits <- v:
				default:
					delete(vs.clients, c)
					close(c.shutdown)
				}
			}
		case <-vs.shutdown:
			for c := range vs.clients {
				close(c.shutdown)
			}
		}
	}

	log.Printf("Visit service stopped")

}

// Stop the service by closing the service's channel, blocking until the
// service has gracefully stopped.
func (vs *VisitService) Stop() {
	close(vs.shutdown)
	vs.wg.Wait()
}

// Connection is a one-direction flow of visits from the service to a remote
// client.
type connection struct {
	ws       *websocket.Conn
	visits   chan Visit
	shutdown chan struct{}
}

// Serve the connection by serialising and sending incoming visits.
func (c *connection) serve() {
	defer c.ws.Close()
	for {
		select {
		case v := <-c.visits:
			err := websocket.JSON.Send(c.ws, v)
			if err != nil {
				break
			}
		case <-c.shutdown:
			break
		}
	}
}
