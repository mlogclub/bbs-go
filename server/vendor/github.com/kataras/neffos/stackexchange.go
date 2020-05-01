package neffos

import (
	"context"
)

// StackExchange is an optional interface
// that can be used to change the way neffos
// sends messages to its clients, i.e
// communication between multiple neffos servers.
//
// See the "kataras/neffos/stackexchange" subpackage for more details.
// Real-World example and usage documentation
// can be found at: "kataras/neffos/_examples/redis".
type StackExchange interface {
	// OnConnect should prepare the connection's subscriber.
	// It's called automatically after the neffos server's OnConnect (if any)
	// on incoming client connections.
	OnConnect(c *Conn) error
	// OnDisconnect should close the connection's subscriber that
	// created on the `OnConnect` method.
	// It's called automatically when a connection goes offline,
	// manually by server or client or by network failure.
	OnDisconnect(c *Conn)

	// Publish should publish messages through a stackexchange.
	// It's called automatically on neffos broadcasting.
	Publish(msgs []Message) bool
	// Subscribe should subscribe to a specific namespace,
	// it's called automatically on neffos namespace connected.
	Subscribe(c *Conn, namespace string)
	// Unsubscribe should unsubscribe from a specific namespace,
	// it's called automatically on neffos namespace disconnect.
	Unsubscribe(c *Conn, namespace string) // should close the subscriber.
	// Ask should be able to perform a server Ask to a specific client or to all clients
	// It blocks until response from a specific client if msg.To is filled,
	// otherwise will return on the first responder's reply.
	Ask(ctx context.Context, msg Message, token string) (Message, error)
	// NotifyAsk should notify and unblock a subscribed connection for this
	// specific message, "token" is the neffos wait signal for this message.
	NotifyAsk(msg Message, token string) error
}

// StackExchangeInitializer is an optional interface for a `StackExchange`.
// It contains a single `Init` method which accepts
// the registered server namespaces  and returns error.
// It does not called on manual `Server.StackExchange` field set,
// use the `Server.UseStackExchange` to make sure that this implementation is respected.
type StackExchangeInitializer interface {
	// Init should initialize a stackexchange, it's optional.
	Init(Namespaces) error
}

func stackExchangeInit(s StackExchange, namespaces Namespaces) error {
	if s != nil {
		if sinit, ok := s.(StackExchangeInitializer); ok {
			return sinit.Init(namespaces)
		}
	}

	return nil
}

// internal use only when more than one stack exchanges are registered.
type stackExchangeWrapper struct {
	// read-only fields.
	parent  StackExchange
	current StackExchange
}

func wrapStackExchanges(existingExc StackExchange, newExc StackExchange) StackExchange {
	return &stackExchangeWrapper{
		parent:  existingExc,
		current: newExc,
	}
}

func (s *stackExchangeWrapper) OnConnect(c *Conn) error {
	// return on first error, do not wrap errors,
	// the server should NOT run if at least one is errored.
	err := s.parent.OnConnect(c)
	if err != nil {
		return err
	}

	return s.current.OnConnect(c)
}

func (s *stackExchangeWrapper) OnDisconnect(c *Conn) {
	s.parent.OnDisconnect(c)
	s.current.OnDisconnect(c)
}

func (s *stackExchangeWrapper) Publish(msgs []Message) bool {
	// keep try on the next but return false on any failure.
	okParent := s.parent.Publish(msgs)
	okCurrent := s.current.Publish(msgs)

	return okParent && okCurrent
}

func (s *stackExchangeWrapper) Ask(ctx context.Context, msg Message, token string) (Message, error) {
	// we run Ask and if one is failing then we keep trying for all stackexchanges.
	msg, err := s.parent.Ask(ctx, msg, token)
	if err != nil {
		msg, err = s.current.Ask(ctx, msg, token)
	}

	return msg, err
}

func (s *stackExchangeWrapper) NotifyAsk(msg Message, token string) error {
	err := s.parent.NotifyAsk(msg, token)
	if err != nil {
		return s.current.NotifyAsk(msg, token)
	}

	return nil
}

func (s *stackExchangeWrapper) Subscribe(c *Conn, namespace string) {
	s.parent.Subscribe(c, namespace)
	s.current.Subscribe(c, namespace)
}

func (s *stackExchangeWrapper) Unsubscribe(c *Conn, namespace string) {
	s.parent.Unsubscribe(c, namespace)
	s.current.Unsubscribe(c, namespace)
}
