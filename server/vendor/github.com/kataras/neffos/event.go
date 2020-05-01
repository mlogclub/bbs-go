package neffos

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// MessageHandlerFunc is the definition type of the events' callback.
// Its error can be written to the other side on specific events,
// i.e on `OnNamespaceConnect` it will abort a remote namespace connection.
// See examples for more.
type MessageHandlerFunc func(*NSConn, Message) error

var (
	// OnNamespaceConnect is the event name which its callback is fired right before namespace connect,
	// if non-nil error then the remote connection's `Conn.Connect` will fail and send that error text.
	// Connection is not ready to emit data to the namespace.
	OnNamespaceConnect = "_OnNamespaceConnect"
	// OnNamespaceConnected is the event name which its callback is fired after namespace successfully connected.
	// Connection is ready to emit data back to the namespace.
	OnNamespaceConnected = "_OnNamespaceConnected"
	// OnNamespaceDisconnect is the event name which its callback is fired when
	// remote namespace disconnection or local namespace disconnection is happening.
	// For server-side connections the reply matters, so if error returned then the client-side cannot disconnect yet,
	// for client-side the return value does not matter.
	OnNamespaceDisconnect = "_OnNamespaceDisconnect" // if allowed to connect then it's allowed to disconnect as well.
	// OnRoomJoin is the event name which its callback is fired right before room join.
	OnRoomJoin = "_OnRoomJoin" // able to check if allowed to join.
	// OnRoomJoined is the event name which its callback is fired after the connection has successfully joined to a room.
	OnRoomJoined = "_OnRoomJoined" // able to broadcast messages to room.
	// OnRoomLeave is the event name which its callback is fired right before room leave.
	OnRoomLeave = "_OnRoomLeave" // able to broadcast bye-bye messages to room.
	// OnRoomLeft is the event name which its callback is fired after the connection has successfully left from a room.
	OnRoomLeft = "_OnRoomLeft" // if allowed to join to a room, then its allowed to leave from it.
	// OnAnyEvent is the event name which its callback is fired when incoming message's event is not declared to the ConnHandler(`Events` or `Namespaces`).
	OnAnyEvent = "_OnAnyEvent" // when event no match.
	// OnNativeMessage is fired on incoming native/raw websocket messages.
	// If this event defined then an incoming message can pass the check (it's an invalid message format)
	// with just the Message's Body filled, the Event is "OnNativeMessage" and IsNative always true.
	// This event should be defined under an empty namespace in order this to work.
	OnNativeMessage = "_OnNativeMessage"
)

// IsSystemEvent reports whether the "event" is a system event,
// OnNamespaceConnect, OnNamespaceConnected, OnNamespaceDisconnect,
// OnRoomJoin, OnRoomJoined, OnRoomLeave and OnRoomLeft.
func IsSystemEvent(event string) bool {
	switch event {
	case OnNamespaceConnect, OnNamespaceConnected, OnNamespaceDisconnect,
		OnRoomJoin, OnRoomJoined, OnRoomLeave, OnRoomLeft:
		return true
	default:
		return false
	}
}

// CloseError can be used to send and close a remote connection in the event callback's return statement.
type CloseError struct {
	error
	Code int
}

func (err CloseError) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.error.Error())
}

// IsDisconnectError reports whether the "err" is a timeout or a closed connection error.
func IsDisconnectError(err error) bool {
	if err == nil {
		return false
	}

	return IsCloseError(err) || IsTimeoutError(err)
}

func isManualCloseError(err error) bool {
	if _, ok := err.(CloseError); ok {
		return true
	}

	return false
}

// IsCloseError reports whether the "err" is a "closed by the remote host" network connection error.
func IsCloseError(err error) bool {
	if err == nil {
		return false
	}

	if isManualCloseError(err) {
		return true
	}

	if err == io.ErrUnexpectedEOF || err == io.EOF {
		return true
	}

	if netErr, ok := err.(*net.OpError); ok {
		if netErr.Err == nil {
			return false
		}

		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			if sysErr.Err == nil {
				return false
			}
			// return strings.HasSuffix(sysErr.Err.Error(), "closed by the remote host.")
			return true
		}

		return strings.HasSuffix(err.Error(), "use of closed network connection")
	}

	return false
}

// IsTimeoutError reports whether the "err" is caused by a defined timeout.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	if netErr, ok := err.(*net.OpError); ok {
		// poll.TimeoutError is the /internal/poll of the go language itself, we can't use it directly.
		return netErr.Timeout()
	}

	return false
}

type reply struct {
	Body []byte
}

func (r reply) Error() string {
	return ""
}

func isReply(err error) ([]byte, bool) {
	if err != nil {
		if r, ok := err.(reply); ok {
			return r.Body, true
		}
	}
	return nil, false
}

// Reply is a special type of custom error which sends a message back to the other side
// with the exact same incoming Message's Namespace (and Room if specified)
// except its body which would be the given "body".
func Reply(body []byte) error {
	return reply{body}
}
