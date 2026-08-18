// Package colorpolicy is an example implementation of machinery.Policy
// that targets Gateway API resources and assigns colors to them via a
// list of rules.
//
// ColorPolicy demonstrates how to define custom merge strategies for
// policies with overlapping target references. A policy spec can be
// declared under Defaults or Overrides, each with a strategy ("atomic" or
// "merge") that determines how a policy combines with others when more
// than one applies to the same target: atomic strategies keep one policy's
// rules wholesale, while merge strategies combine rules from both
// policies, rule by rule, using rule Id.
package colorpolicy
