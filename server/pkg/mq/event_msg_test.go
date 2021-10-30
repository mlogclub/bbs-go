package mq

import (
	"fmt"
	"testing"

	"github.com/gookit/event"
)

func TestEvent(t *testing.T) {
	// Register event listener
	event.On("evt1", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handle event: %s\n", e.Name())
		return nil
	}), event.Normal)

	// Register multiple listeners
	event.On("evt1", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handle event: %s\n", e.Name())
		return nil
	}), event.High)

	_, _ = event.Fire("evt1", map[string]interface{}{
		"a": "b",
	})

	_, _ = event.Fire("evt", nil)
}
