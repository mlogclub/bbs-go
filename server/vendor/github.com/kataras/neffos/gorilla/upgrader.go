package gorilla

import (
	"net/http"

	"github.com/kataras/neffos"

	gorilla "github.com/gorilla/websocket"
)

// DefaultUpgrader is a gorilla/websocket Upgrader with all fields set to the default values.
var DefaultUpgrader = Upgrader(gorilla.Upgrader{})

// Upgrader is a `neffos.Upgrader` type for the gorilla/websocket subprotocol implementation.
// Should be used on `New` to construct the neffos server.
func Upgrader(upgrader gorilla.Upgrader) neffos.Upgrader {
	return func(w http.ResponseWriter, r *http.Request) (neffos.Socket, error) {
		underline, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			return nil, err
		}

		return newSocket(underline, r, false), nil
	}
}
