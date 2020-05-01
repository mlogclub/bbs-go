package neffos

import (
	"reflect"
	"strings"
	"time"
)

// ConnHandler is the interface which namespaces and events can be retrieved through.
// Built-in ConnHandlers are the`Events`, `Namespaces`, `WithTimeout` and `NewStruct`.
// Users of this are the `Dial`(client) and `New` (server) functions.
type ConnHandler interface {
	GetNamespaces() Namespaces
}

var (
	_ ConnHandler = (Events)(nil)
	_ ConnHandler = (Namespaces)(nil)
	_ ConnHandler = WithTimeout{}
	_ ConnHandler = (*Struct)(nil)
)

// Events completes the `ConnHandler` interface.
// It is a map which its key is the event name
// and its value the event's callback.
//
// Events type completes the `ConnHandler` itself therefore,
// can be used as standalone value on the `New` and `Dial` functions
// to register events on empty namespace as well.
//
// See `Namespaces`, `New` and `Dial` too.
type Events map[string]MessageHandlerFunc

// GetNamespaces returns an empty namespace with the "e" Events.
func (e Events) GetNamespaces() Namespaces {
	return Namespaces{"": e}
}

func (e Events) fireEvent(c *NSConn, msg Message) error {
	if h, ok := e[msg.Event]; ok {
		return h(c, msg)
	}

	if h, ok := e[OnAnyEvent]; ok {
		return h(c, msg)
	}

	return nil
}

// On is a shortcut of Events { eventName: msgHandler }.
// It registers a callback "msgHandler" for an event "eventName".
func (e Events) On(eventName string, msgHandler MessageHandlerFunc) {
	e[eventName] = msgHandler
}

// Namespaces completes the `ConnHandler` interface.
// Can be used to register one or more namespaces on the `New` and `Dial` functions.
// The key is the namespace literal and the value is the `Events`,
// a map with event names and their callbacks.
//
// See `WithTimeout`, `New` and `Dial` too.
type Namespaces map[string]Events

// GetNamespaces just returns the "nss" namespaces.
func (nss Namespaces) GetNamespaces() Namespaces { return nss }

// On is a shortcut of Namespaces { namespace: Events: { eventName: msgHandler } }.
// It registers a callback "msgHandler" for an event "eventName" of the particular "namespace".
func (nss Namespaces) On(namespace, eventName string, msgHandler MessageHandlerFunc) Events {
	if nss[namespace] == nil {
		nss[namespace] = make(Events)
	}
	nss[namespace][eventName] = msgHandler

	return nss[namespace]
}

// WithTimeout completes the `ConnHandler` interface.
// Can be used to register namespaces and events or just events on an empty namespace
// with Read and Write timeouts.
//
// See `New` and `Dial`.
type WithTimeout struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Namespaces Namespaces
	Events     Events
}

// GetNamespaces returns combined namespaces from "Namespaces" and "Events" fields
// with read and write timeouts from "ReadTimeout" and "WriteTimeout" fields of "t".
func (t WithTimeout) GetNamespaces() Namespaces {
	return JoinConnHandlers(t.Namespaces, t.Events).GetNamespaces()
}

func getTimeouts(h ConnHandler) (readTimeout time.Duration, writeTimeout time.Duration) {
	if t, ok := h.(WithTimeout); ok {
		readTimeout = t.ReadTimeout
		writeTimeout = t.WriteTimeout
	}

	if s, ok := h.(*Struct); ok {
		readTimeout = s.readTimeout
		writeTimeout = s.writeTimeout
	}

	return
}

// EventMatcherFunc is a type of which a Struct matches the methods with neffos events.
type EventMatcherFunc = func(methodName string) (string, bool)

// Struct is a ConnHandler. All fields are unexported, use `NewStruct` instead.
// It converts any pointer to a struct value to `neffos.Namespaces` using reflection.
type Struct struct {
	ptr reflect.Value

	// defaults to empty and tries to get it through `Struct.Namespace() string` method.
	namespace string
	// defaults to nil, if specified
	// then it matches the events based on the result string or false if this method shouldn't register as event.
	eventMatcher              EventMatcherFunc
	readTimeout, writeTimeout time.Duration

	// This field is set when external dependency injection system is used.
	injector StructInjector

	events Events
}

// SetNamespace sets a namespace that this Struct is responsible for,
// Alterinatively create a method on the controller named `Namespace() string`
// to retrieve this namespace at build time.
func (s *Struct) SetNamespace(namespace string) *Struct {
	s.namespace = namespace
	return s
}

var (
	// EventPrefixMatcher matches methods to events based on the "prefix".
	EventPrefixMatcher = func(prefix string) EventMatcherFunc {
		return func(methodName string) (string, bool) {
			if strings.HasPrefix(methodName, prefix) {
				return methodName, true
			}

			return "", false
		}
	}

	// EventTrimPrefixMatcher matches methods based on the "prefixToTrim"
	// and events are registered without this prefix.
	EventTrimPrefixMatcher = func(prefixToTrim string) EventMatcherFunc {
		return func(methodName string) (string, bool) {
			if strings.HasPrefix(methodName, prefixToTrim) {
				return methodName[len(prefixToTrim):], true
			}

			return "", false
		}
	}
)

// SetEventMatcher sets an event method matcher which applies to every
// event except the system events (OnNamespaceConnected, and so on).
func (s *Struct) SetEventMatcher(matcher EventMatcherFunc) *Struct {
	s.eventMatcher = matcher
	return s
}

// SetTimeouts sets read and write deadlines on the underlying network connection.
// After a read or write have timed out, the websocket connection is closed.
//
// Defaults to 0, no timeout except an `Upgrader` or `Dialer` specifies its own values.
func (s *Struct) SetTimeouts(read, write time.Duration) *Struct {
	s.readTimeout = read
	s.writeTimeout = write

	return s
}

// SetInjector sets a custom injector and overrides the neffos default behavior
// on dynamic structs.
// The "fn" should handle to fill static fields and the NSConn.
// This "fn" will only be called when dynamic struct "ptr" is passed
// on the `NewStruct`.
// The caller should return a
// valid type of "ptr" reflect.Value.
func (s *Struct) SetInjector(fn StructInjector) *Struct {
	s.injector = fn
	return s
}

// NewStruct returns a new Struct value instance type of ConnHandler.
// The "ptr" should be a pointer to a struct.
// This function is used when you want to convert a structure to
// neffos.ConnHandler based on the struct's methods.
// The methods if "ptr" structure value
// can be func(msg neffos.Message) error if the structure contains a *neffos.NSConn field,
// otherwise they should be like any event callback: func(nsConn *neffos.NSConn, msg neffos.Message) error.
// If contains a field of type *neffos.NSConn then on each new connection to the namespace a new controller is created
// and static fields(if any) are set on runtime with the NSConn itself.
// If it's a static controller (does not contain a NSConn field)
// then it just registers its functions as regular events without performance cost.
//
// Users of this method is `New` and `Dial`.
//
// Note that this method has a tiny performance cost when an event's callback's logic has small footprint.
func NewStruct(ptr interface{}) *Struct {
	if ptr == nil {
		panic("NewStruct: value is nil")
	}

	if s, ok := ptr.(*Struct); ok { // if it's already a *Struct then just return it.
		return s
	}

	var v reflect.Value // use for methods with receiver Ptr.
	if rValue, ok := ptr.(reflect.Value); ok {
		v = rValue
	} else {
		v = reflect.ValueOf(ptr)
	}

	if !v.IsValid() {
		panic("NewStruct: value is not a valid one")
	}

	typ := v.Type() // use for methods with receiver Ptr.

	if typ.Kind() != reflect.Ptr {
		panic("NewStruct: value should be a pointer")
	}

	if typ.ConvertibleTo(nsConnType) {
		panic("NewStruct: conversion for type" + typ.String() + " NSConn is not allowed.")
	}

	if indirectType(typ).Kind() != reflect.Struct {
		panic("NewStruct: value does not points to a struct")
	}

	n := typ.NumMethod()
	_, hasNamespaceMethod := typ.MethodByName("Namespace")
	if n == 0 || (n == 1 && hasNamespaceMethod) {
		panic("NewStruct: value does not contain any exported methods")
	}

	return &Struct{
		ptr: v,
	}
}

// Events builds and returns the Events.
// Callers of this method is users that want to add Structs to different namespaces
// in the same application.
// When a single namespace is used then this call is unnecessary,
// the `Struct` is already a fully featured `ConnHandler` by itself.
func (s *Struct) Events() Events {
	if s.events != nil {
		return s.events
	}

	s.events = makeEventsFromStruct(s.ptr, s.eventMatcher, s.injector)
	return s.events
}

// GetNamespaces creates and returns Namespaces based on the
// pointer to struct value provided by the "s".
func (s *Struct) GetNamespaces() Namespaces { // completes the `ConnHandler` interface.
	if s.namespace == "" {
		s.namespace, _ = resolveStructNamespace(s.ptr)
	}

	return Namespaces{
		s.namespace: s.Events(),
	}
}

// JoinConnHandlers combines two or more "connHandlers"
// and returns a result of a single `ConnHandler` that
// can be passed on the `New` and `Dial` functions.
func JoinConnHandlers(connHandlers ...ConnHandler) ConnHandler {
	namespaces := Namespaces{}

	for _, h := range connHandlers {
		nss := h.GetNamespaces()
		if len(nss) > 0 {
			for namespace, events := range nss {
				if events == nil {
					continue
				}
				clonedEvents := make(Events, len(events))
				for evt, cb := range events {
					clonedEvents[evt] = cb
				}

				if curEvents, exists := namespaces[namespace]; exists {
					// fill missing events.
					for evt, cb := range clonedEvents {
						curEvents[evt] = cb
					}

				} else {
					namespaces[namespace] = clonedEvents
				}
			}
		}
	}

	return namespaces
}
