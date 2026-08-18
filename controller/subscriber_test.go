//go:build unit

package controller_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kuadrant/policy-machinery/controller"
	"github.com/kuadrant/policy-machinery/machinery"
)

func newTestConfigMap(namespace, name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
}

func TestSubscriptionReconcile(t *testing.T) {
	fooKind := schema.GroupKind{Group: "test", Kind: "Foo"}
	barKind := schema.GroupKind{Group: "test", Kind: "Bar"}

	createEvent := controller.CreateEvent
	deleteEvent := controller.DeleteEvent

	fooCreated := controller.ResourceEvent{
		Kind:      fooKind,
		EventType: controller.CreateEvent,
		NewObject: newTestConfigMap("ns1", "foo"),
	}
	fooUpdated := controller.ResourceEvent{
		Kind:      fooKind,
		EventType: controller.UpdateEvent,
		OldObject: newTestConfigMap("ns1", "foo"),
		NewObject: newTestConfigMap("ns1", "foo"),
	}
	barDeleted := controller.ResourceEvent{
		Kind:      barKind,
		EventType: controller.DeleteEvent,
		OldObject: newTestConfigMap("ns2", "bar"),
	}

	testCases := []struct {
		name           string
		events         []controller.ResourceEventMatcher
		resourceEvents []controller.ResourceEvent
		reconcileFunc  controller.ReconcileFunc // nil if not expecting a call
		expectedEvents []controller.ResourceEvent
		expectedErr    error
	}{
		{
			name:           "no resource events",
			events:         []controller.ResourceEventMatcher{{Kind: &fooKind}},
			resourceEvents: nil,
		},
		{
			name:           "no matchers",
			events:         nil,
			resourceEvents: []controller.ResourceEvent{fooCreated},
		},
		{
			name:           "matches by kind",
			events:         []controller.ResourceEventMatcher{{Kind: &fooKind}},
			resourceEvents: []controller.ResourceEvent{fooCreated, barDeleted},
			expectedEvents: []controller.ResourceEvent{fooCreated},
		},
		{
			name:           "matches by event type",
			events:         []controller.ResourceEventMatcher{{EventType: &deleteEvent}},
			resourceEvents: []controller.ResourceEvent{fooCreated, fooUpdated, barDeleted},
			expectedEvents: []controller.ResourceEvent{barDeleted},
		},
		{
			name:           "matches by namespace",
			events:         []controller.ResourceEventMatcher{{ObjectNamespace: "ns2"}},
			resourceEvents: []controller.ResourceEvent{fooCreated, barDeleted},
			expectedEvents: []controller.ResourceEvent{barDeleted},
		},
		{
			name:           "matches by name",
			events:         []controller.ResourceEventMatcher{{ObjectName: "bar"}},
			resourceEvents: []controller.ResourceEvent{fooCreated, barDeleted},
			expectedEvents: []controller.ResourceEvent{barDeleted},
		},
		{
			name: "matches by kind, event type, namespace and name combined",
			events: []controller.ResourceEventMatcher{
				{Kind: &fooKind, EventType: &createEvent, ObjectNamespace: "ns1", ObjectName: "foo"},
			},
			resourceEvents: []controller.ResourceEvent{fooCreated, fooUpdated, barDeleted},
			expectedEvents: []controller.ResourceEvent{fooCreated},
		},
		{
			name: "does not match when one field of a matcher differs",
			events: []controller.ResourceEventMatcher{
				{Kind: &fooKind, EventType: &deleteEvent},
			},
			resourceEvents: []controller.ResourceEvent{fooCreated, fooUpdated},
			expectedEvents: nil,
		},
		{
			name: "matches on any of multiple matchers (OR semantics)",
			events: []controller.ResourceEventMatcher{
				{Kind: &barKind},
				{EventType: &createEvent},
			},
			resourceEvents: []controller.ResourceEvent{fooCreated, fooUpdated, barDeleted},
			expectedEvents: []controller.ResourceEvent{fooCreated, barDeleted},
		},
		{
			name:           "empty matcher matches everything",
			events:         []controller.ResourceEventMatcher{{}},
			resourceEvents: []controller.ResourceEvent{fooCreated, fooUpdated, barDeleted},
			expectedEvents: []controller.ResourceEvent{fooCreated, fooUpdated, barDeleted},
		},
		{
			name:           "falls back to NewObject for namespace/name when OldObject is nil",
			events:         []controller.ResourceEventMatcher{{ObjectNamespace: "ns1", ObjectName: "foo"}},
			resourceEvents: []controller.ResourceEvent{fooCreated},
			expectedEvents: []controller.ResourceEvent{fooCreated},
		},
		{
			name:           "no reconcile func set",
			events:         []controller.ResourceEventMatcher{{Kind: &fooKind}},
			resourceEvents: []controller.ResourceEvent{fooCreated},
			reconcileFunc:  nil,
		},
		{
			name:           "reconcile func error is propagated",
			events:         []controller.ResourceEventMatcher{{Kind: &fooKind}},
			resourceEvents: []controller.ResourceEvent{fooCreated},
			expectedEvents: []controller.ResourceEvent{fooCreated},
			expectedErr:    fmt.Errorf("boom"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var calledWith []controller.ResourceEvent
			called := false

			reconcileFunc := controller.ReconcileFunc(func(_ context.Context, events []controller.ResourceEvent, _ *machinery.Topology, _ error, _ *sync.Map) error {
				called = true
				calledWith = events
				return tc.expectedErr
			})

			// Special case: explicitly test the nil ReconcileFunc path.
			if tc.name == "no reconcile func set" {
				reconcileFunc = nil
			}

			s := controller.Subscription{
				ReconcileFunc: reconcileFunc,
				Events:        tc.events,
			}

			err := s.Reconcile(t.Context(), tc.resourceEvents, nil, nil, nil)

			expectCall := len(tc.expectedEvents) > 0 && reconcileFunc != nil
			if called != expectCall {
				t.Errorf("expected ReconcileFunc called: %t, got %t", expectCall, called)
			}

			if expectCall {
				if len(calledWith) != len(tc.expectedEvents) {
					t.Fatalf("expected %d events passed to ReconcileFunc, got %d", len(tc.expectedEvents), len(calledWith))
				}
				for i, expected := range tc.expectedEvents {
					if calledWith[i].Kind != expected.Kind || calledWith[i].EventType != expected.EventType {
						t.Errorf("event %d: expected kind %v/eventType %v, got kind %v/eventType %v",
							i, expected.Kind, expected.EventType, calledWith[i].Kind, calledWith[i].EventType)
					}
				}
			}

			if expectCall && tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
