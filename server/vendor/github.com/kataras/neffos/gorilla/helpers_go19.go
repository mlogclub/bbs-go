// +build go1.9

package gorilla

import gorilla "github.com/gorilla/websocket"

// Options is just an alias for the `gorilla/websocket.Dialer` struct type.
type Options = gorilla.Dialer
