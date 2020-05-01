package gobwas

import (
	"net/http"

	"github.com/kataras/neffos"

	gobwas "github.com/gobwas/ws"
)

// DefaultUpgrader is a gobwas/ws HTTP Upgrader with all fields set to the default values.
var DefaultUpgrader = Upgrader(gobwas.HTTPUpgrader{})

// Upgrader is a `neffos.Upgrader` type for the gobwas/ws subprotocol implementation.
// Should be used on `neffos.New` to construct the neffos server.
func Upgrader(upgrader gobwas.HTTPUpgrader) neffos.Upgrader {
	return func(w http.ResponseWriter, r *http.Request) (neffos.Socket, error) {
		underline, _, _, err := upgrader.Upgrade(r, w)
		if err != nil {
			return nil, err
		}

		return newSocket(underline, r, false), nil
	}
}
