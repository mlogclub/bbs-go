package neffos

import (
	"context"
	"strings"
)

// Client is the neffos client. Contains the neffos client-side connection
// and the ID came from server on acknowledgement process of the `Dial` function.
// Use its `Connect` to connect to a namespace or
// `WaitServerConnect` to wait for server to force-connect this client to a namespace.
type Client struct {
	conn *Conn

	// ID comes from server, local changes are not reflected,
	// use the `Server#IDGenerator` if you want to set a custom logic for ID set.
	ID string

	// NotifyClose can be optionally registered to notify about the client's disconnect.
	// This callback is for the entire client side connection,
	// the channel is notified after namespace disconnected and any room left events.
	// Don't confuse it with the `OnNamespaceDisconnect` event.
	// Usage:
	// <- client.NotifyClose // blocks until local `Close` or remote close of connection.
	NotifyClose <-chan struct{}
}

// Close method terminates the client-side connection.
// Forces the client to disconnect from all connected namespaces and leave from all joined rooms,
// server gets notified.
func (c *Client) Close() {
	if c == nil || c.conn == nil {
		return
	}

	c.conn.Close()
}

// WaitServerConnect method blocks until server manually calls the connection's `Connect`
// on the `Server#OnConnected` event.
//
// See `Conn#WaitConnect` for more details.
func (c *Client) WaitServerConnect(ctx context.Context, namespace string) (*NSConn, error) {
	return c.conn.WaitConnect(ctx, namespace)
}

// Connect method returns a new connected to the specific "namespace" `NSConn` value.
// The "namespace" should be declared in the `connHandler` of both server and client sides.
// Returns error if server-side's `OnNamespaceConnect` event callback returns an error.
//
// See `Conn#Connect` for more details.
func (c *Client) Connect(ctx context.Context, namespace string) (*NSConn, error) {
	return c.conn.Connect(ctx, namespace)
}

// Dialer is the definition type of a dialer, gorilla or gobwas or custom.
// It is the second parameter of the `Dial` function.
type Dialer func(ctx context.Context, url string) (Socket, error)

// Dial establishes a new neffos client connection.
// Context "ctx" is used for handshake timeout.
// Dialer "dial" can be either `gobwas.Dialer/DefaultDialer` or `gorilla.Dialer/DefaultDialer`,
// custom dialers can be used as well when complete the `Socket` and `Dialer` interfaces for valid client.
// URL "url" is the endpoint of the neffos server, i.e "ws://localhost:8080/echo".
// The last parameter, and the most important one is the "connHandler", it can be
// filled as `Namespaces`, `Events` or `WithTimeout`, same namespaces and events can be used on the server-side as well.
//
// See examples for more.
func Dial(ctx context.Context, dial Dialer, url string, connHandler ConnHandler) (*Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if !strings.HasPrefix(url, "ws://") && !strings.HasPrefix(url, "wss://") {
		url = "ws://" + url
	}

	underline, err := dial(ctx, url)
	if err != nil {
		return nil, err
	}

	if connHandler == nil {
		connHandler = Namespaces{}
	}

	c := newConn(underline, connHandler.GetNamespaces())
	readTimeout, writeTimeout := getTimeouts(connHandler)
	c.readTimeout = readTimeout
	c.writeTimeout = writeTimeout

	go c.startReader()

	if err = c.sendClientACK(); err != nil {
		return nil, err
	}

	return &Client{conn: c, ID: c.id, NotifyClose: c.closeCh}, nil
}
