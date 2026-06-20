package a

// Thing is a Null-Object interface (it declares IsNull() bool).
type Thing interface {
	IsNull() bool
	Foo() string
}

type null struct{}

func (null) IsNull() bool { return true }
func (null) Foo() string  { return "" }

// NullThing is the canonical Null object.
func NullThing() Thing { return null{} }

// Embedded marker via an embedded interface is still detected.
type Base interface{ IsNull() bool }
type Derived interface {
	Base
	Bar() int
}

func okConstructor() Thing { return NullThing() } // ok
func okNamed() Thing {
	t := NullThing()
	return t // ok: not a bare nil
}

func badReturn() Thing         { return nil }      // want `return a Null object .* instead of nil for a\.Thing`
func badEmbedded() Derived     { return nil }      // want `return a Null object .* for a\.Derived`
func badMulti() (Thing, error) { return nil, nil } // want `return a Null object .* for a\.Thing`

func notNullObject() error { return nil } // ok: error has no IsNull()
func plainPointer() *int   { return nil } // ok: not an interface

func inLiteral() Thing {
	f := func() Thing { return nil } // want `return a Null object .* for a\.Thing`
	return f()
}
