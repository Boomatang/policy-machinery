//go:build unit

package controller_test

import (
	"testing"

	"github.com/kuadrant/policy-machinery/controller"
)

func TestEventType_String(t *testing.T) {
	tests := []struct {
		name  string
		event controller.EventType
		want  string
	}{
		{"create event", controller.CreateEvent, "create"},
		{"update event", controller.UpdateEvent, "update"},
		{"delete event", controller.DeleteEvent, "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := tt.event
			got := event.String()
			if got != tt.want {
				t.Errorf("EventType(%d).String() = %q, want %q", tt.event, got, tt.want)
			}
		})
	}
}

func TestEventType_Values(t *testing.T) {
	if controller.CreateEvent != 0 {
		t.Errorf("CreateEvent = %d, want 0", controller.CreateEvent)
	}
	if controller.UpdateEvent != 1 {
		t.Errorf("UpdateEvent = %d, want 1", controller.UpdateEvent)
	}
	if controller.DeleteEvent != 2 {
		t.Errorf("DeleteEvent = %d, want 2", controller.DeleteEvent)
	}
}
