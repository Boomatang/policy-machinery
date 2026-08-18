// Package machinery provides types for modeling topologies of targetable
// network resources (Namespaces, Services, Gateway API Gateways, Routes,
// etc.) and the policies attached to them.
//
// A Topology is a directed acyclic graph built with NewTopology (or the
// Gateway API-specific NewGatewayAPITopology), from targetables, policies,
// generic objects, and link functions describing how they relate. Once
// built, a Topology can be queried for targetables, policies, and objects,
// walked for parents/children/paths, and rendered as a Graphviz dot graph.
//
// Objects participate in a topology by implementing Object, Targetable,
// or Policy, and policies are associated with the targetables they
// reference via GetTargetRefs and, when they overlap, combined using a
// MergeStrategy.
package machinery
