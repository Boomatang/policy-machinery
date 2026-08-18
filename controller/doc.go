// Package controller provides a simplified controller abstraction, built on
// k8s.io/apimachinery and sigs.k8s.io/controller-runtime, for watching
// Gateway API, policy, and other linked resources and reacting to changes
// via a machinery.Topology.
//
// A Controller is configured with functional options (WithClient,
// WithLogger, WithPolicyKinds, WithObjectKinds, WithRunnable,
// WithReconcile, ...) and started with Controller.Start. As watched
// resources change, the controller rebuilds its topology and invokes a
// ReconcileFunc with the resulting events.
//
// Subscription and Workflow can be used to filter events and compose
// multiple reconciliation functions together.
package controller
