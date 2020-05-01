package neffos

import (
	"context"
	"reflect"
	"sync"
)

// NSConn describes a connection connected to a specific namespace,
// it emits with the `Message.Namespace` filled and it can join to multiple rooms.
// A single `Conn` can be connected to one or more namespaces,
// each connected namespace is described by this structure.
type NSConn struct {
	Conn *Conn
	// Static from server, client can select which to use or not.
	// Client and server can ask to connect.
	// Server can forcely disconnect.
	namespace string
	// Static from server, client can select which to use or not.
	events Events

	// Dynamically channels/rooms for each connected namespace.
	// Client can ask to join, server can forcely join a connection to a room.
	// Namespace(room(fire event)).
	rooms      map[string]*Room
	roomsMutex sync.RWMutex

	// value is just a temporarily value.
	// Storage across event callbacks for this namespace.
	value reflect.Value
}

func newNSConn(c *Conn, namespace string, events Events) *NSConn {
	return &NSConn{
		Conn:      c,
		namespace: namespace,
		events:    events,
		rooms:     make(map[string]*Room),
	}
}

// String method simply returns the Conn's ID().
// Useful method to this connected to a namespace connection to be passed on `Server#Broadcast` method
// to exclude itself from the broadcasted message's receivers.
func (ns *NSConn) String() string {
	return ns.Conn.String()
}

// Emit method sends a message to the remote side
// with its `Message.Namespace` filled to this specific namespace.
func (ns *NSConn) Emit(event string, body []byte) bool {
	if ns == nil { // if for any reason Namespace() called without be available.
		return false
	}

	return ns.Conn.Write(Message{Namespace: ns.namespace, Event: event, Body: body})
}

// Ask method writes a message to the remote side and blocks until a response or an error received.
func (ns *NSConn) Ask(ctx context.Context, event string, body []byte) (Message, error) {
	if ns == nil {
		return Message{}, ErrWrite
	}

	return ns.Conn.Ask(ctx, Message{Namespace: ns.namespace, Event: event, Body: body})
}

// JoinRoom method can be used to join a connection to a specific room, rooms are dynamic.
// Returns the joined `Room`.
func (ns *NSConn) JoinRoom(ctx context.Context, roomName string) (*Room, error) {
	if ns == nil {
		return nil, ErrWrite
	}

	return ns.askRoomJoin(ctx, roomName)
}

// Room method returns a joined `Room`.
func (ns *NSConn) Room(roomName string) *Room {
	if ns == nil {
		return nil
	}

	ns.roomsMutex.RLock()
	room := ns.rooms[roomName]
	ns.roomsMutex.RUnlock()

	return room
}

// Rooms returns a slice copy of the joined rooms.
func (ns *NSConn) Rooms() []*Room {
	ns.roomsMutex.RLock()
	rooms := make([]*Room, len(ns.rooms))
	i := 0
	for _, room := range ns.rooms {
		rooms[i] = room
		i++
	}
	ns.roomsMutex.RUnlock()

	return rooms
}

// LeaveAll method sends a remote and local leave room signal `OnRoomLeave` to and for all rooms
// and fires the `OnRoomLeft` event if succeed.
func (ns *NSConn) LeaveAll(ctx context.Context) error {
	if ns == nil {
		return nil
	}

	ns.roomsMutex.Lock()
	defer ns.roomsMutex.Unlock()

	leaveMsg := Message{Namespace: ns.namespace, Event: OnRoomLeave, IsLocal: true, locked: true}
	for room := range ns.rooms {
		leaveMsg.Room = room
		if err := ns.askRoomLeave(ctx, leaveMsg, false); err != nil {
			return err
		}
	}

	return nil
}

func (ns *NSConn) forceLeaveAll(isLocal bool) {
	ns.roomsMutex.Lock()
	defer ns.roomsMutex.Unlock()

	leaveMsg := Message{Namespace: ns.namespace, Event: OnRoomLeave, IsForced: true, IsLocal: isLocal}
	for room := range ns.rooms {
		leaveMsg.Room = room
		ns.events.fireEvent(ns, leaveMsg)

		delete(ns.rooms, room)

		leaveMsg.Event = OnRoomLeft
		ns.events.fireEvent(ns, leaveMsg)

		leaveMsg.Event = OnRoomLeave
	}
}

// Disconnect method sends a disconnect signal to the remote side and fires the local `OnNamespaceDisconnect` event.
func (ns *NSConn) Disconnect(ctx context.Context) error {
	if ns == nil {
		return nil
	}

	return ns.Conn.askDisconnect(ctx, Message{
		Namespace: ns.namespace,
		Event:     OnNamespaceDisconnect,
	}, true)
}

func (ns *NSConn) askRoomJoin(ctx context.Context, roomName string) (*Room, error) {
	ns.roomsMutex.RLock()
	room, ok := ns.rooms[roomName]
	ns.roomsMutex.RUnlock()
	if ok {
		return room, nil
	}

	joinMsg := Message{
		Namespace: ns.namespace,
		Room:      roomName,
		Event:     OnRoomJoin,
		IsLocal:   true,
	}

	_, err := ns.Conn.Ask(ctx, joinMsg)
	if err != nil {
		return nil, err
	}

	err = ns.events.fireEvent(ns, joinMsg)
	if err != nil {
		return nil, err
	}

	room = newRoom(ns, roomName)
	ns.roomsMutex.Lock()
	ns.rooms[roomName] = room
	ns.roomsMutex.Unlock()

	joinMsg.Event = OnRoomJoined
	ns.events.fireEvent(ns, joinMsg)
	return room, nil
}

func (ns *NSConn) replyRoomJoin(msg Message) {
	if ns == nil || msg.wait == "" || msg.isNoOp {
		return
	}

	ns.roomsMutex.RLock()
	_, ok := ns.rooms[msg.Room]
	ns.roomsMutex.RUnlock()
	if !ok {
		err := ns.events.fireEvent(ns, msg)
		if err != nil {
			msg.Err = err
			ns.Conn.Write(msg)
			return
		}
		ns.roomsMutex.Lock()
		ns.rooms[msg.Room] = newRoom(ns, msg.Room)
		ns.roomsMutex.Unlock()

		msg.Event = OnRoomJoined
		ns.events.fireEvent(ns, msg)
	}

	ns.Conn.writeEmptyReply(msg.wait)
}

func (ns *NSConn) askRoomLeave(ctx context.Context, msg Message, lock bool) error {
	if ns == nil {
		return nil
	}

	if lock {
		ns.roomsMutex.RLock()
	}
	_, ok := ns.rooms[msg.Room]
	if lock {
		ns.roomsMutex.RUnlock()
	}

	if !ok {
		return ErrBadRoom
	}

	_, err := ns.Conn.Ask(ctx, msg)
	if err != nil {
		return err
	}

	// msg.IsLocal = true
	err = ns.events.fireEvent(ns, msg)
	if err != nil {
		return err
	}

	if lock {
		ns.roomsMutex.Lock()
	}

	delete(ns.rooms, msg.Room)

	if lock {
		ns.roomsMutex.Unlock()
	}

	msg.Event = OnRoomLeft
	ns.events.fireEvent(ns, msg)

	return nil
}

func (ns *NSConn) replyRoomLeave(msg Message) {
	if ns == nil || msg.wait == "" || msg.isNoOp {
		return
	}

	room := ns.Room(msg.Room)
	if room == nil {
		ns.Conn.writeEmptyReply(msg.wait)
		return
	}

	// if client then we need to respond to server and delete the room without ask the local event.
	if ns.Conn.IsClient() {
		ns.events.fireEvent(ns, msg)

		ns.roomsMutex.Lock()
		delete(ns.rooms, msg.Room)
		ns.roomsMutex.Unlock()

		ns.Conn.writeEmptyReply(msg.wait)

		msg.Event = OnRoomLeft
		ns.events.fireEvent(ns, msg)
		return
	}

	// server-side, check for error on the local event first.
	err := ns.events.fireEvent(ns, msg)
	if err != nil {
		msg.Err = err
		ns.Conn.Write(msg)
		return
	}

	ns.roomsMutex.Lock()
	delete(ns.rooms, msg.Room)
	ns.roomsMutex.Unlock()

	msg.Event = OnRoomLeft
	ns.events.fireEvent(ns, msg)

	ns.Conn.writeEmptyReply(msg.wait)
}
