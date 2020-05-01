package neffos

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	uuid "github.com/iris-contrib/go.uuid"
)

// Upgrader is the definition type of a protocol upgrader, gorilla or gobwas or custom.
// It is the first parameter of the `New` function which constructs a neffos server.
type Upgrader func(w http.ResponseWriter, r *http.Request) (Socket, error)

// IDGenerator is the type of function that it is used
// to generate unique identifiers for new connections.
//
// See `Server.IDGenerator`.
type IDGenerator func(w http.ResponseWriter, r *http.Request) string

// DefaultIDGenerator returns a universal unique identifier for a new connection.
// It's the default `IDGenerator` for `Server`.
var DefaultIDGenerator IDGenerator = func(http.ResponseWriter, *http.Request) string {
	id, err := uuid.NewV4()
	if err != nil {
		return strconv.FormatInt(time.Now().Unix(), 10)
	}
	return id.String()
}

// Server is the neffos server.
// Keeps the `IDGenerator` which can be customized, by default it's the `DefaultIDGenerator`  which
// generates connections unique identifiers using the uuid/v4.
//
// Callers can optionally register callbacks for connection, disconnection and errored.
// Its most important methods are `ServeHTTP` which is used to register the server on a specific endpoint
// and `Broadcast` and `Close`.
// Use the `New` function to create a new server, server starts automatically, no further action is required.
type Server struct {
	uuid string

	upgrader      Upgrader
	IDGenerator   IDGenerator
	StackExchange StackExchange

	// If `StackExchange` is set then this field is ignored.
	//
	// It overrides the default behavior(when no StackExchange is not used)
	// which publishes a message independently.
	// In short the default behavior doesn't wait for a message to be published to all clients
	// before any next broadcast call.
	//
	// Therefore, if set to true,
	// each broadcast call will publish its own message(s) by order.
	SyncBroadcaster bool

	mu         sync.RWMutex
	namespaces Namespaces

	// connection read/write timeouts.
	readTimeout  time.Duration
	writeTimeout time.Duration

	count uint64

	connections       map[*Conn]struct{}
	connect           chan *Conn
	disconnect        chan *Conn
	actions           chan action
	broadcastMessages chan []Message

	broadcaster *broadcaster

	// messages that this server must waits
	// for a reply from one of its own connections(see `waitMessages`).
	waitingMessages      map[string]chan Message
	waitingMessagesMutex sync.RWMutex

	closed uint32

	// OnUpgradeError can be optionally registered to catch upgrade errors.
	OnUpgradeError func(err error)
	// OnConnect can be optionally registered to be notified for any new neffos client connection,
	// it can be used to force-connect a client to a specific namespace(s) or to send data immediately or
	// even to cancel a client connection and dissalow its connection when its return error value is not nil.
	// Don't confuse it with the `OnNamespaceConnect`, this callback is for the entire client side connection.
	OnConnect func(c *Conn) error
	// OnDisconnect can be optionally registered to notify about a connection's disconnect.
	// Don't confuse it with the `OnNamespaceDisconnect`, this callback is for the entire client side connection.
	OnDisconnect func(c *Conn)
}

// New constructs and returns a new neffos server.
// Listens to incoming connections automatically, no further action is required from the caller.
// The second parameter is the "connHandler", it can be
// filled as `Namespaces`, `Events` or `WithTimeout`, same namespaces and events can be used on the client-side as well,
// Use the `Conn#IsClient` on any event callback to determinate if it's a client-side connection or a server-side one.
//
// See examples for more.
func New(upgrader Upgrader, connHandler ConnHandler) *Server {
	readTimeout, writeTimeout := getTimeouts(connHandler)
	namespaces := connHandler.GetNamespaces()
	s := &Server{
		uuid:              uuid.Must(uuid.NewV4()).String(),
		upgrader:          upgrader,
		namespaces:        namespaces,
		readTimeout:       readTimeout,
		writeTimeout:      writeTimeout,
		connections:       make(map[*Conn]struct{}),
		connect:           make(chan *Conn, 1),
		disconnect:        make(chan *Conn),
		actions:           make(chan action),
		broadcastMessages: make(chan []Message),
		broadcaster:       newBroadcaster(),
		waitingMessages:   make(map[string]chan Message),
		IDGenerator:       DefaultIDGenerator,
	}

	go s.start()

	return s
}

// UseStackExchange can be used to add one or more StackExchange
// to the server.
// Returns a non-nil error when "exc"
// completes the `StackExchangeInitializer` interface and its `Init` failed.
//
// Read more at the `StackExchange` type's docs.
func (s *Server) UseStackExchange(exc StackExchange) error {
	if exc == nil {
		return nil
	}

	if err := stackExchangeInit(exc, s.namespaces); err != nil {
		return err
	}

	if s.usesStackExchange() {
		s.StackExchange = wrapStackExchanges(s.StackExchange, exc)
	} else {
		s.StackExchange = exc
	}

	return nil
}

// usesStackExchange reports whether this server
// uses one or more `StackExchange`s.
func (s *Server) usesStackExchange() bool {
	return s.StackExchange != nil
}

func (s *Server) start() {
	atomic.StoreUint32(&s.closed, 0)

	for {
		select {
		case c := <-s.connect:
			s.connections[c] = struct{}{}
			atomic.AddUint64(&s.count, 1)
		case c := <-s.disconnect:
			if _, ok := s.connections[c]; ok {
				// close(c.out)
				delete(s.connections, c)
				atomic.AddUint64(&s.count, ^uint64(0))
				// println("disconnect...")
				if s.OnDisconnect != nil {
					// don't fire disconnect if was immediately closed on the `OnConnect` server event.
					if !c.readiness.isReady() || (c.readiness.err != nil) {
						continue
					}
					s.OnDisconnect(c)
				}

				if s.usesStackExchange() {
					s.StackExchange.OnDisconnect(c)
				}
			}
		case msgs := <-s.broadcastMessages:
			for c := range s.connections {
				publishMessages(c, msgs)
			}
		case act := <-s.actions:
			for c := range s.connections {
				act.call(c)
			}

			if act.done != nil {
				act.done <- struct{}{}
			}
		}
	}
}

// Close terminates the server and all of its connections, client connections are getting notified.
func (s *Server) Close() {
	if atomic.CompareAndSwapUint32(&s.closed, 0, 1) {
		s.Do(func(c *Conn) {
			c.Close()
		}, false)
	}
}

var (
	errServerClosed  = errors.New("server closed")
	errInvalidMethod = errors.New("no valid request method")
)

// URLParamAsHeaderPrefix is the prefix that server parses the url parameters as request headers.
// The client's `URLParamAsHeaderPrefix` must match.
// Note that this is mostly useful for javascript browser-side clients, nodejs and go client support custom headers by default.
// No action required from end-developer, exported only for chance to a custom parsing.
const URLParamAsHeaderPrefix = "X-Websocket-Header-"

func tryParseURLParamsToHeaders(r *http.Request) {
	q := r.URL.Query()
	for k, values := range q {
		if len(k) <= len(URLParamAsHeaderPrefix) {
			continue
		}

		k = http.CanonicalHeaderKey(k) // canonical, so no X-WebSocket thing.

		idx := strings.Index(k, URLParamAsHeaderPrefix)
		if idx != 0 { // must be prefix.
			continue
		}

		if r.Header == nil {
			r.Header = make(http.Header)
		}

		k = k[len(URLParamAsHeaderPrefix):]

		for _, v := range values {
			r.Header.Add(k, v)
		}
	}
}

var errUpgradeOnRetry = errors.New("check status")

// IsTryingToReconnect reports whether the returning "err" from the `Server#Upgrade`
// is from a client that was trying to reconnect to the websocket server.
//
// Look the `Conn#WasReconnected` and `Conn#ReconnectTries` too.
func IsTryingToReconnect(err error) (ok bool) {
	return err != nil && err == errUpgradeOnRetry
}

// This header key should match with that browser-client's `whenResourceOnline->re-dial` uses.
const websocketReconectHeaderKey = "X-Websocket-Reconnect"

func isServerConnID(s string) bool {
	return strings.HasPrefix(s, "neffos(0x")
}

func genServerConnID(s *Server, c *Conn) string {
	return fmt.Sprintf("neffos(0x%s(%s%p))", s.uuid, c.id, c)
}

// Upgrade handles the connection, same as `ServeHTTP` but it can accept
// a socket wrapper and a "customIDGen" that overrides the server's IDGenerator
// and it does return the connection or any errors.
func (s *Server) Upgrade(
	w http.ResponseWriter,
	r *http.Request,
	socketWrapper func(Socket) Socket,
	customIDGen IDGenerator,
) (*Conn, error) {
	if atomic.LoadUint32(&s.closed) > 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, errServerClosed
	}

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusFound)
		return nil, errUpgradeOnRetry
	}

	if r.Method != http.MethodGet {
		// RCF rfc2616 https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
		// The response MUST include an Allow header containing a list of valid methods for the requested resource.
		//
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Allow#Examples
		w.Header().Set("Allow", http.MethodGet)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return nil, errInvalidMethod
	}

	tryParseURLParamsToHeaders(r)

	socket, err := s.upgrader(w, r)
	if err != nil {
		if s.OnUpgradeError != nil {
			s.OnUpgradeError(err)
		}
		return nil, err
	}

	if socketWrapper != nil {
		socket = socketWrapper(socket)
	}

	c := newConn(socket, s.namespaces)
	if customIDGen != nil {
		c.id = customIDGen(w, r)
	} else {
		c.id = s.IDGenerator(w, r)
	}
	c.serverConnID = genServerConnID(s, c)

	c.readTimeout = s.readTimeout
	c.writeTimeout = s.writeTimeout
	c.server = s

	retriesHeaderValue := r.Header.Get(websocketReconectHeaderKey)
	if retriesHeaderValue != "" {
		c.ReconnectTries, _ = strconv.Atoi(retriesHeaderValue)
	}

	if !s.usesStackExchange() && !s.SyncBroadcaster {
		go func(c *Conn) {
			for s.waitMessages(c) {
			}
		}(c)
	}

	s.connect <- c

	go c.startReader()

	// Before `OnConnect` in order to be able
	// to Broadcast inside the `OnConnect` custom func.
	if s.usesStackExchange() {
		if err := s.StackExchange.OnConnect(c); err != nil {
			c.readiness.unwait(err)
			return nil, err
		}
	}

	// Start the reader before `OnConnect`, remember clients may remotely connect to namespace before `Server#OnConnect`
	// therefore any `Server:NSConn#OnNamespaceConnected` can write immediately to the client too.
	// Note also that the `Server#OnConnect` itself can do that as well but if the written Message's Namespace is not locally connected
	// it, correctly, can't pass the write checks. Also, and most important, the `OnConnect` is ready to connect a client to a namespace (locally and remotely).
	//
	// This has a downside:
	// We need a way to check if the `OnConnect` returns an non-nil error which means that the connection should terminate before namespace connect or anything.
	// The solution is to still accept reading messages but add them to the queue(like we already do for any case messages came before ack),
	// the problem to that is that the queue handler is fired when ack is done but `OnConnect` may not even return yet, so we introduce a `mark ready` atomic scope
	// and a channel which will wait for that `mark ready` if handle queue is called before ready.
	// Also make the same check before emit the connection's disconnect event (if defined),
	// which will be always ready to be called because we added the connections via the connect channel;
	// we still need the connection to be available for any broadcasting on connected events.
	// ^ All these only when server-side connection in order to correctly handle the end-developer's `OnConnect`.
	//
	// Look `Conn.serverReadyWaiter#startReader##handleQueue.serverReadyWaiter.unwait`(to hold the events until no error returned or)
	// `#Write:serverReadyWaiter.unwait` (for things like server connect).
	// All cases tested & worked perfectly.
	if s.OnConnect != nil {
		if err = s.OnConnect(c); err != nil {
			// TODO: Do something with that error.
			// The most suitable thing we can do is to somehow send this to the client's `Dial` return statement.
			// This can be done if client waits for "OK" signal or a failure with an error before return the websocket connection,
			// as for today we have the ack process which does NOT block and end-developer can send messages and server will handle them when both sides are ready.
			// So, maybe it's a better solution to transform that process into a blocking state which can handle any `Server#OnConnect` error and return it at client's `Dial`.
			// Think more later today.
			// Done but with a lot of code.... will try to cleanup some things.
			//println("OnConnect error: " + err.Error())
			c.readiness.unwait(err)
			// No need to disconnect here, connection's .Close will be called on readiness ch errored.

			// c.Close()
			return nil, err
		}
	}

	//println("OnConnect does not exist or no error, fire unwait")
	c.readiness.unwait(nil)

	return c, nil
}

// ServeHTTP completes the `http.Handler` interface, it should be passed on a http server's router
// to serve this neffos server on a specific endpoint.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Upgrade(w, r, nil, nil)
}

// GetTotalConnections returns the total amount of the connected connections to the server, it's fast
// and can be used as frequently as needed.
func (s *Server) GetTotalConnections() uint64 {
	return atomic.LoadUint64(&s.count)
}

type action struct {
	call func(*Conn)
	done chan struct{}
}

// Do loops through all connected connections and fires the "fn", with this method
// callers can do whatever they want on a connection outside of a event's callback,
// but make sure that these operations are not taking long time to complete because it delays the
// new incoming connections.
// If "async" is true then this method does not block the flow of the program.
func (s *Server) Do(fn func(*Conn), async bool) {
	act := action{call: fn}
	if !async {
		act.done = make(chan struct{})
		// go func() { s.actions <- act }()
		// <-act.done
	}

	s.actions <- act
	if !async {
		<-act.done
	}
}

func publishMessages(c *Conn, msgs []Message) bool {
	for _, msg := range msgs {
		if msg.from == c.ID() {
			// if the message is not supposed to return back to any connection with this ID.

			return true
		}

		// if "To" field is given then send to a specific connection.
		if msg.To != "" && msg.To != c.ID() {
			return true
		}

		// c.Write may fail if the message is not supposed to end to this client
		// but the connection should be still open in order to continue.
		if !c.Write(msg) && c.IsClosed() {
			return false
		}
	}

	return true
}

func (s *Server) waitMessages(c *Conn) bool {
	s.broadcaster.mu.Lock()
	defer s.broadcaster.mu.Unlock()

	msgs, ok := s.broadcaster.waitUntilClosed(c.closeCh)
	if !ok {
		return false
	}

	return publishMessages(c, msgs)
}

type stringerValue struct{ v string }

func (s stringerValue) String() string { return s.v }

// Exclude can be passed on `Server#Broadcast` when
// caller does not have access to the `Conn`, `NSConn` or a `Room` value but
// has access to a string variable which is a connection's ID instead.
//
// Example Code:
// nsConn.Conn.Server().Broadcast(
//	neffos.Exclude("connection_id_here"),
//  neffos.Message{Namespace: "default", Room: "roomName or empty", Event: "chat", Body: [...]})
func Exclude(connID string) fmt.Stringer { return stringerValue{connID} }

// Broadcast method is fast and does not block any new incoming connection by-default,
// it can be used as frequently as needed. Use the "msg"'s Namespace, or/and Event or/and Room to broadcast
// to a specific type of connection collectives.
//
// If first "exceptSender" parameter is not nil then the message "msg" will be
// broadcasted to all connected clients except the given connection's ID,
// any value that completes the `fmt.Stringer` interface is valid. Keep note that
// `Conn`, `NSConn`, `Room` and `Exclude(connID) global function` are valid values.
//
// Example Code:
// nsConn.Conn.Server().Broadcast(
//	nsConn OR nil,
//  neffos.Message{Namespace: "default", Room: "roomName or empty", Event: "chat", Body: [...]})
//
// Note that it if `StackExchange` is nil then its default behavior
// doesn't wait for a publish to complete to all clients before any
// next broadcast call. To change that behavior set the `Server.SyncBroadcaster` to true
// before server start.
func (s *Server) Broadcast(exceptSender fmt.Stringer, msgs ...Message) {

	if exceptSender != nil {
		var fromExplicit, from string

		switch c := exceptSender.(type) {
		case *Conn:
			fromExplicit = c.serverConnID
		case *NSConn:
			fromExplicit = c.Conn.serverConnID
		default:
			from = exceptSender.String()
		}

		for i := range msgs {
			if from != "" {
				msgs[i].from = from
			} else {
				msgs[i].FromExplicit = fromExplicit
			}
		}
	}

	if s.usesStackExchange() {
		s.StackExchange.Publish(msgs)
		return
	}

	if s.SyncBroadcaster {
		s.broadcastMessages <- msgs
		return
	}

	s.broadcaster.broadcast(msgs)
}

// Ask is like `Broadcast` but it blocks until a response
// from a specific connection if "msg.To" is filled otherwise
// from the first connection which will reply to this "msg".
//
// Accepts a context for deadline as its first input argument.
// The second argument is the request message
// which should be sent to a specific namespace:event
// like the `Conn.Ask`.
func (s *Server) Ask(ctx context.Context, msg Message) (Message, error) {
	if ctx == nil {
		ctx = context.TODO()
	}

	msg.wait = genWait(false)

	if s.usesStackExchange() {
		msg.wait = genWaitStackExchange(msg.wait)
		return s.StackExchange.Ask(ctx, msg, msg.wait)
	}

	ch := make(chan Message)
	s.waitingMessagesMutex.Lock()
	s.waitingMessages[msg.wait] = ch
	s.waitingMessagesMutex.Unlock()

	s.Broadcast(nil, msg)

	select {
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case receive := <-ch:
		s.waitingMessagesMutex.Lock()
		delete(s.waitingMessages, msg.wait)
		s.waitingMessagesMutex.Unlock()

		return receive, receive.Err
	}
}

// GetConnectionsByNamespace can be used as an alternative way to retrieve
// all connected connections to a specific "namespace" on a specific time point.
// Do not use this function frequently, it is not designed to be fast or cheap, use it for debugging or logging every 'x' time.
// Users should work with the event's callbacks alone, the usability is enough for all type of operations. See `Do` too.
//
// Not thread safe.
func (s *Server) GetConnectionsByNamespace(namespace string) map[string]*NSConn {
	conns := make(map[string]*NSConn)

	s.mu.RLock()
	for c := range s.connections {
		if ns := c.Namespace(namespace); ns != nil {
			conns[ns.Conn.ID()] = ns
		}
	}
	s.mu.RUnlock()

	return conns
}

// GetConnections can be used as an alternative way to retrieve
// all connected connections to the server on a specific time point.
// Do not use this function frequently, it is not designed to be fast or cheap, use it for debugging or logging every 'x' time.
//
// Not thread safe.
func (s *Server) GetConnections() map[string]*Conn {
	conns := make(map[string]*Conn)

	s.mu.RLock()
	for c := range s.connections {
		conns[c.ID()] = c
	}
	s.mu.RUnlock()

	return conns
}

var (
	// ErrBadNamespace may return from a `Conn#Connect` method when the remote side does not declare the given namespace.
	ErrBadNamespace = errors.New("bad namespace")
	// ErrBadRoom may return from a `Room#Leave` method when trying to leave from a not joined room.
	ErrBadRoom = errors.New("bad room")
	// ErrWrite may return from any connection's method when the underline connection is closed (unexpectedly).
	ErrWrite = errors.New("write closed")
)
