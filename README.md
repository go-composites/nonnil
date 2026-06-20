# golang-cop/nonnil

[![ci](https://github.com/golang-cop/nonnil/actions/workflows/ci.yml/badge.svg)](https://github.com/golang-cop/nonnil/actions/workflows/ci.yml)

A `go vet` analyzer that enforces the [golang-cop](https://github.com/golang-cop)
**Null-Object invariant**: a function whose result is a *Null-Object interface*
(any interface declaring `IsNull() bool`) must never return a bare `nil` — it must
hand back a real Null object (`Null.New()`, `NullError.New()`, …).

This turns "always return a Null object, never `nil`" from a matter of human
discipline into a check that **fails CI** — closing the one gap that lets the
nil-dereference the pattern is designed to prevent slip back in.

## Why

The composition style avoids `nil` dereferences because constructors initialise
every field to a Null object and callers send messages instead of testing
`!= nil`. But Go's type system still lets any function return a bare `nil` for an
interface result; one stray `return nil` and a caller panics. `nonnil` makes that
mistake un-mergeable.

## Install & use

```sh
go install github.com/golang-cop/nonnil/cmd/nonnil@latest
go vet -vettool=$(which nonnil) ./...
```

It reports, for example:

```
src/error.go:43:46: return a Null object (e.g. Null.New()) instead of nil for
  github.com/golang-cop/error/src.Interface: it is a Null-Object interface
  (IsNull() bool), and nil reintroduces the nil-dereference the pattern prevents
```

## What it flags

- A **bare `return nil`** in any result position whose declared type is an
  interface with an `IsNull() bool` method (directly or embedded).
- Works for single and multi-value returns, and inside function literals.

It deliberately does **not** flag interfaces without `IsNull()` (e.g. the builtin
`error`, or `Result.Interface`), so false positives are minimal. A deliberate
typed conversion `return (Thing)(nil)` or returning a nil interface variable is
also left alone — those are explicit and uncommon.

## CI

Add to each repo's workflow:

```yaml
  - run: go install github.com/golang-cop/nonnil/cmd/nonnil@latest
  - run: go vet -vettool=$(which nonnil) ./...
```

## License

BSD-3-Clause © the golang-cop/nonnil authors.
