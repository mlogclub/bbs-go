package neffos

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// The Message is the structure which describes the incoming and outcoming data.
// Emitter's "body" argument is the `Message.Body` field.
// Emitter's return non-nil error is the `Message.Err` field.
// If native message sent then the `Message.Body` is filled with the body and
// when incoming native message then the `Message.Event` is the `OnNativeMessage`,
// native messages are allowed only when an empty namespace("") and its `OnNativeMessage` callback are present.
//
// The the raw data received/sent structured following this order:
// <wait()>;
// <namespace>;
// <room>;
// <event>;
// <isError(0-1)>;
// <isNoOp(0-1)>;
// <body||error_message>
//
// Internal `serializeMessage` and
// exported `DeserializeMessage` functions
// do the job on `Conn#Write`, `NSConn#Emit` and `Room#Emit` calls.
type Message struct {
	wait string

	// The Namespace that this message sent to/received from.
	Namespace string
	// The Room that this message sent to/received from.
	Room string
	// The Event that this message sent to/received from.
	Event string
	// The actual body of the incoming/outcoming data.
	Body []byte
	// The Err contains any message's error, if any.
	// Note that server-side and client-side connections can return an error instead of a message from each event callbacks,
	// except the clients's force Disconnect which its local event doesn't matter when disconnected manually.
	Err error

	// if true then `Err` is filled by the error message and
	// the last segment of incoming/outcoming serialized message is the error message instead of the body.
	isError bool
	isNoOp  bool

	isInvalid bool

	// the CONN ID, filled automatically if `Server#Broadcast` first parameter of sender connection's ID is not empty,
	// not exposed to the subscribers (rest of the clients).
	// This is the ID across neffos servers when scale.
	from string
	// When sent by the same connection of the current running server instance.
	// This field is serialized/deserialized but it's clean on sending or receiving from a client
	// and it's only used on StackExchange feature.
	// It's serialized as the first parameter, instead of wait signal, if incoming starts with 0x.
	FromExplicit string // the exact Conn's pointer in this server instance.
	// Reports whether this message is coming from a stackexchange.
	// This field is not exposed and it's not serialized at all, ~local-use only~.
	//
	// The "wait" field can determinate if this message is coming from a stackexchange using its second char,
	// This value set based on "wait" on deserialization when coming from remote side.
	// Only server-side can actually set it.
	FromStackExchange bool

	// To is the connection ID of the receiver, used only when `Server#Broadcast` is called, indeed when we only need to send a message to a single connection.
	// The Namespace, Room are still respected at all.
	//
	// However, sending messages to a group of connections is done by the `Room` field for groups inside a namespace or just `Namespace` field as usual.
	// This field is not filled on sending/receiving.
	To string

	// True when event came from local (i.e client if running client) on force disconnection,
	// i.e OnNamespaceDisconnect and OnRoomLeave when closing a conn.
	// This field is not filled on sending/receiving.
	// Err does not matter and never sent to the other side.
	IsForced bool
	// True when asking the other side and fire the respond's event (which matches the sent for connect/disconnect/join/leave),
	// i.e if a client (or server) onnection want to connect
	// to a namespace or join to a room.
	// Should be used rarely, state can be checked by `Conn#IsClient() bool`.
	// This field is not filled on sending/receiving.
	IsLocal bool

	// True when user define it for writing, only its body is written as raw native websocket message, namespace, event and all other fields are empty.
	// The receiver should accept it on the `OnNativeMessage` event.
	// This field is not filled on sending/receiving.
	IsNative bool

	// Useful rarely internally on `Conn#Write` namespace and rooms checks, i.e `Conn#DisconnectAll` and `NSConn#RemoveAll`.
	// If true then the writer's checks will not lock connectedNamespacesMutex or roomsMutex again. May be useful in the future, keep that solution.
	locked bool

	// if server or client should write using Binary message.
	// This field is not filled on sending/receiving.
	SetBinary bool
}

func (m *Message) isConnect() bool {
	return m.Event == OnNamespaceConnect
}

func (m *Message) isDisconnect() bool {
	return m.Event == OnNamespaceDisconnect
}

func (m *Message) isRoomJoin() bool {
	return m.Event == OnRoomJoin
}

func (m *Message) isRoomLeft() bool {
	return m.Event == OnRoomLeft
}

// Serialize returns this message's transport format.
func (m Message) Serialize() []byte {
	return serializeMessage(nil, m)
}

type (
	// MessageObjectMarshaler is an optional interface that "objects"
	// can implement to customize their byte representation, see `Object` package-level function.
	MessageObjectMarshaler interface {
		Marshal() ([]byte, error)
	}

	// MessageObjectUnmarshaler is an optional interface that "objects"
	// can implement to customize their structure, see `Message.Object` method.
	MessageObjectUnmarshaler interface {
		Unmarshal(body []byte) error
	}
)

var (
	// DefaultMarshaler is a global, package-level alternative for `MessageObjectMarshaler`.
	// It's used when the `Marshal.v` parameter is not a `MessageObjectMarshaler`.
	DefaultMarshaler = json.Marshal
	// DefaultUnmarshaler is a global, package-level alternative for `MessageObjectMarshaler`.
	// It's used when the `Message.Unmarshal.outPtr` parameter is not a `MessageObjectUnmarshaler`.
	DefaultUnmarshaler = json.Unmarshal
)

// Marshal marshals the "v" value and returns a Message's Body.
// If the "v" value is `MessageObjectMarshaler` then it returns the result of its `Marshal` method,
// otherwise the DefaultMarshaler will be used instead.
// Errors are pushed to the result, use the object's Marshal method to catch those when necessary.
func Marshal(v interface{}) []byte {
	if v == nil {
		panic("nil assigment")
	}

	var (
		body []byte
		err  error
	)

	if marshaler, ok := v.(MessageObjectMarshaler); ok {
		body, err = marshaler.Marshal()
	} else {
		body, err = DefaultMarshaler(v)
	}

	if err != nil {
		return []byte(err.Error())
	}
	return body
}

// Unmarshal unmarshals this Message's body to the "outPtr".
// The "outPtr" must be a pointer to a value that can customize its decoded value
// by implementing the `MessageObjectUnmarshaler`, otherwise the `DefaultUnmarshaler` will be used instead.
func (m *Message) Unmarshal(outPtr interface{}) error {
	if outPtr == nil {
		panic("nil assigment")
	}

	if unmarshaler, ok := outPtr.(MessageObjectUnmarshaler); ok {
		return unmarshaler.Unmarshal(m.Body)
	}

	return DefaultUnmarshaler(m.Body, outPtr)
}

const (
	waitIsConfirmationPrefix   = '#'
	waitComesFromClientPrefix  = '$'
	waitComesFromStackExchange = '!'
)

// IsWait reports whether this message waits for a response back.
func (m *Message) IsWait(isClientConn bool) bool {
	if m.wait == "" {
		return false
	}

	if m.wait[0] == waitIsConfirmationPrefix {
		// true even if it's not client-client but it's a confirmation message.
		return true
	}

	if m.wait[0] == waitComesFromClientPrefix {
		if isClientConn {
			return true
		}
		return false
	}

	return true
}

// ClearWait clears the wait token, rarely used.
func (m *Message) ClearWait() bool {
	if m.FromExplicit == "" && m.wait != "" {
		m.wait = ""
		return true
	}

	return false
}

func genWait(isClientConn bool) string {
	now := time.Now().UnixNano()
	wait := strconv.FormatInt(now, 10)

	if isClientConn {
		wait = string(waitComesFromClientPrefix) + wait
	}

	return wait
}

func genWaitConfirmation(wait string) string {
	return string(waitIsConfirmationPrefix) + wait
}

func genWaitStackExchange(wait string) string {
	if len(wait) < 2 {
		return ""
	}

	// This is the second special character.
	// If found, it is removed on the deserialization
	// and Message.FromStackExchange is set to true.
	return string(wait[0]+waitComesFromStackExchange) + wait[1:]
}

type (
	// MessageEncrypt type kept for future use when serializing a message.
	MessageEncrypt func(out []byte) []byte
	// MessageDecrypt type kept for future use when deserializing a message.
	MessageDecrypt func(in []byte) []byte
)

var (
	trueByte  = []byte{'1'}
	falseByte = []byte{'0'}

	messageSeparatorString = ";"
	messageSeparator       = []byte(messageSeparatorString)
	// we use this because has zero chance to be part of end-developer's Message.Namespace, Room, Event, To and Err fields,
	// semicolon has higher probability to exists on those values. See `escape` and `unescape`.
	messageFieldSeparatorReplacement = "@%!semicolon@%!"
)

// called on `serializeMessage` to all message's fields except the body (and error).
func escape(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.Replace(s, messageSeparatorString, messageFieldSeparatorReplacement, -1)
}

// called on `DeserializeMessage` to all message's fields except the body (and error).
func unescape(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.Replace(s, messageFieldSeparatorReplacement, messageSeparatorString, -1)
}

func serializeMessage(encrypt MessageEncrypt, msg Message) (out []byte) {
	if msg.IsNative && msg.wait == "" {
		out = msg.Body
	} else {
		if msg.FromExplicit != "" {
			if msg.wait != "" {
				// this should never happen unless manual set of FromExplicit by end-developer which is forbidden by the higher level calls.
				panic("msg.wait and msg.FromExplicit cannot work together")
			}

			msg.wait = msg.FromExplicit
		}
		out = serializeOutput(msg.wait, escape(msg.Namespace), escape(msg.Room), escape(msg.Event), msg.Body, msg.Err, msg.isNoOp)
	}

	if encrypt != nil {
		out = encrypt(out)
	}

	return out
}

func serializeOutput(wait, namespace, room, event string,
	body []byte,
	err error,
	isNoOp bool,
) []byte {

	var (
		isErrorByte = falseByte
		isNoOpByte  = falseByte
		waitByte    = []byte{}
	)

	if err != nil {
		if b, ok := isReply(err); ok {
			body = b
		} else {
			body = []byte(err.Error())
			isErrorByte = trueByte
		}
	}

	if isNoOp {
		isNoOpByte = trueByte
	}

	if wait != "" {
		waitByte = []byte(wait)
	}

	msg := bytes.Join([][]byte{ // this number of fields should match the deserializer's, see `validMessageSepCount`.
		waitByte,
		[]byte(namespace),
		[]byte(room),
		[]byte(event),
		isErrorByte,
		isNoOpByte,
		body,
	}, messageSeparator)

	return msg
}

// DeserializeMessage accepts a serialized message []byte
// and returns a neffos Message.
// When allowNativeMessages only Body is filled and check about message format is skipped.
func DeserializeMessage(decrypt MessageDecrypt, b []byte, allowNativeMessages, shouldHandleOnlyNativeMessages bool) Message {
	if decrypt != nil {
		b = decrypt(b)
	}

	wait, namespace, room, event, body, err, isNoOp, isInvalid := deserializeInput(b, allowNativeMessages, shouldHandleOnlyNativeMessages)

	fromExplicit := ""
	if isServerConnID(wait) {
		fromExplicit = wait
		wait = ""
	}

	fromStackExchange := len(wait) > 2 && wait[1] == waitComesFromStackExchange
	if fromStackExchange {
		// remove the second special char, we need to reform it,
		// this wait token is compared to the waiter side as it's without the information about stackexchnage.
		wait = string(wait[0]) + wait[2:]
	}

	return Message{
		wait:              wait,
		Namespace:         unescape(namespace),
		Room:              unescape(room),
		Event:             unescape(event),
		Body:              body,
		Err:               err,
		isError:           err != nil,
		isNoOp:            isNoOp,
		isInvalid:         isInvalid,
		from:              "",
		FromExplicit:      fromExplicit,
		FromStackExchange: fromStackExchange,
		To:                "",
		IsForced:          false,
		IsLocal:           false,
		IsNative:          allowNativeMessages && event == OnNativeMessage,
		locked:            false,
		SetBinary:         false,
	}
}

const validMessageSepCount = 7

var knownErrors = []error{ErrBadNamespace, ErrBadRoom, ErrWrite, ErrInvalidPayload}

// RegisterKnownError registers an error that it's "known" to both server and client sides.
// This simply adds an error to a list which, if its static text matches
// an incoming error text then its value is set to the `Message.Error` field on the events callbacks.
//
// For dynamic text error, there is a special case which if
// the error "err" contains
// a `ResolveError(errorText string) bool` method then,
// it is used to report whether this "err" is match to the incoming error text.
func RegisterKnownError(err error) {
	for _, knownErr := range knownErrors {
		if err == knownErr {
			return
		}
	}

	knownErrors = append(knownErrors, err)
}

func resolveError(errorText string) error {
	for _, knownErr := range knownErrors {
		if resolver, ok := knownErr.(interface {
			ResolveError(errorText string) bool
		}); ok {
			if resolver.ResolveError(errorText) {
				return knownErr
			}
		}

		if knownErr.Error() == errorText {
			return knownErr
		}
	}

	return errors.New(errorText)
}

func deserializeInput(b []byte, allowNativeMessages, shouldHandleOnlyNativeMessages bool) ( // go-lint: ignore line
	wait,
	namespace,
	room,
	event string,
	body []byte,
	err error,
	isNoOp bool,
	isInvalid bool,
) {

	if len(b) == 0 {
		isInvalid = true
		return
	}

	if shouldHandleOnlyNativeMessages {
		event = OnNativeMessage
		body = b
		return
	}

	// Note: Go's SplitN returns the remainder in[6] but JavasSript's string.split behaves differently.
	dts := bytes.SplitN(b, messageSeparator, validMessageSepCount)
	if len(dts) != validMessageSepCount {
		if !allowNativeMessages {
			isInvalid = true
			return
		}

		event = OnNativeMessage
		body = b
		return
	}

	wait = string(dts[0])
	namespace = string(dts[1])
	room = string(dts[2])
	event = string(dts[3])
	isError := bytes.Equal(dts[4], trueByte)
	isNoOp = bytes.Equal(dts[5], trueByte)
	if b := dts[6]; len(b) > 0 {
		if isError {
			errorText := string(b)
			err = resolveError(errorText)
		} else {
			body = b // keep it like that.
		}
	}

	return
}

func genEmptyReplyToWait(wait string) []byte {
	return append([]byte(wait), bytes.Repeat(messageSeparator, validMessageSepCount-1)...)
}
