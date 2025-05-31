// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package symbol

import (
	"bytes"
	"fmt"

	"github.com/golang/glog"
	"github.com/jaqx0r/mtail/internal/runtime/compiler/errors"
	"github.com/jaqx0r/mtail/internal/runtime/compiler/position"
	"github.com/jaqx0r/mtail/internal/runtime/compiler/types"
)

// Kind enumerates the kind of a Symbol.
type Kind int

// Kind enumerates the kinds of symbols found in the program text.
const (
	VarSymbol     Kind = iota // Variables
	CaprefSymbol              // Capture group references
	DecoSymbol                // Decorators
	PatternSymbol             // Named pattern constants
	endSymbol                 // for testing
)

func (k Kind) String() string {
	switch k {
	case VarSymbol:
		return "variable"
	case CaprefSymbol:
		return "capture group reference"
	case DecoSymbol:
		return "decorator"
	case PatternSymbol:
		return "named pattern constant"
	default:
		panic("unexpected symbolkind")
	}
}

// Symbol describes a named program object.
type Symbol struct {
	Name    string             // identifier name
	Kind    Kind               // kind of program object
	Type    types.Type         // object's type
	Pos     *position.Position // Source file position of definition
	Binding interface{}        // binding to storage allocated in runtime
	Addr    int                // Address offset in another structure, object specific
	Used    bool               // Optional marker that this symbol is used after declaration.
}

// NewSymbol creates a record of a given symbol kind, named name, found at loc.
func NewSymbol(name string, kind Kind, pos *position.Position) (sym *Symbol) {
	return &Symbol{name, kind, types.Undef, pos, nil, 0, false}
}

// Scope maintains a record of the identifiers declared in the current program
// scope, and a link to the parent scope.  A program can insert multiple
// symbols with the same identifier into the symbol table; multiple definition
// errors are detected by `Check`, below.
type Scope struct {
	Parent  *Scope
	Symbols map[string][]*Symbol
}

// NewScope creates a new scope within the parent scope.
func NewScope(parent *Scope) *Scope {
	return &Scope{parent, make(map[string][]*Symbol)}
}

// Insert attempts to insert a symbol into the scope.  If the scope already
// contains a symbol with the same name, the new symbol is appended to the
// list.
func (s *Scope) Insert(sym *Symbol) (alt *Symbol) {
	if len(s.Symbols[sym.Name]) > 0 {
		alt = s.Symbols[sym.Name][0]
	}
	s.Symbols[sym.Name] = append(s.Symbols[sym.Name], sym)
	return
}

// InsertAlias attempts to insert a duplicate name for an existing symbol into
// the scope.
func (s *Scope) InsertAlias(sym *Symbol, alias string) (alt *Symbol) {
	if len(s.Symbols[alias]) > 0 {
		alt = s.Symbols[alias][0]
	}
	s.Symbols[alias] = append(s.Symbols[alias], sym)
	return
}

// Lookup returns the symbol with the given name if it is found in this or any
// parent scope, otherwise nil.  If the symbol has more than one definition,
// the first registered symbol is returned.
func (s *Scope) Lookup(name string, kind Kind) *Symbol {
	for scope := s; scope != nil; scope = scope.Parent {
		symList := scope.Symbols[name]
		if len(symList) > 0 && symList[0].Kind == kind {
			return symList[0]
		}
	}
	return nil
}

// String prints the current scope and all parents to a string, recursing up to
// the root scope.  This method is only used for debugging.
func (s *Scope) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "scope %p {", s)
	if s != nil {
		fmt.Fprintln(&buf)
		if len(s.Symbols) > 0 {
			for name, sym := range s.Symbols {
				for _, s := range sym {
					fmt.Fprintf(&buf, "\t%q: %v %q %v\n", name, s.Kind, s.Name, s.Used)
				}
			}
		}
		if s.Parent != nil {
			fmt.Fprintf(&buf, "%s", s.Parent.String())
		}
	}
	fmt.Fprintf(&buf, "}\n")
	return buf.String()
}

// CopyFrom copies all the symbols from another scope object into this one.
// It recurses up the input scope copying all visible symbols into one.
func (s *Scope) CopyFrom(o *Scope) {
	for _, syms := range o.Symbols {
		for _, sym := range syms {
			s.Insert(sym)
		}
	}
	if o.Parent != nil {
		s.CopyFrom(o.Parent)
	}
}

// Check checks a symbol table for validity and emits errors into the given error list if any are found.
func (s *Scope) Check(errors *errors.ErrorList) {
	for _, symList := range s.Symbols {
		multiple := len(symList) > 1
		for i, sym := range symList {
			if multiple && i > 0{
				errors.Add(sym.Pos, fmt.Sprintf("Redeclaration of %s `%s' previously declared at %s", sym.Kind, sym.Name, symList[0].Pos))
				continue
			}
			if !sym.Used {
				// Users don't have control over the patterns given from decorators
				// so this should never be an error; but it can be useful to know
				// if a program is doing unnecessary work.
				if sym.Kind == CaprefSymbol {
					if sym.Addr == 0 {
						// Don't warn about the zeroth capture group; it's not user-defined.
						continue
					}
					glog.Infof("capture group reference `%s' at %s appears to be unused", sym.Name, sym.Pos)
					continue
				}
				errors.Add(sym.Pos, fmt.Sprintf("Declaration of %s `%s' here is never used.", sym.Kind, sym.Name))
			}
		}
	}
}
