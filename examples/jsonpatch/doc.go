// Package jsonpatch is an example implementation of machinery.Policy
// that targets Gateway API resources and assigns colors to them via a
// map of rules.
//
// ColorPolicy demonstrates a merge strategy based on JSON Merge Patch
// (RFC 7396): when two policies target the same resource, their specs are
// marshaled to JSON and combined with jsonpatch.MergePatch, with the
// Defaults or Overrides spec determining which policy acts as the patch
// and which as the base.
package jsonpatch
