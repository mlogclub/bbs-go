// +build go1.9

package gobwas

import gobwas "github.com/gobwas/ws"

// Options is just an alias for the `gobwas/ws.Dialer` struct type.
type Options = gobwas.Dialer

// Header is an alias to the adapter that allows the use of `http.Header` as
// `gobwas/ws.Dialer.HandshakeHeader`.
type Header = gobwas.HandshakeHeaderHTTP
