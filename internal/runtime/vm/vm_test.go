// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package vm

import (
	"context"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/jaqx0r/mtail/internal/logline"
	"github.com/jaqx0r/mtail/internal/metrics"
	"github.com/jaqx0r/mtail/internal/metrics/datum"
	"github.com/jaqx0r/mtail/internal/runtime/code"
	"github.com/jaqx0r/mtail/internal/testutil"
)

var instructions = []struct {
	name          string
	i             code.Instr
	re            []*regexp.Regexp
	str           []string
	reversedStack []interface{} // stack is inverted to be pushed onto vm stack

	expectedStack  []interface{}
	expectedThread thread
}{
	{
		"match",
		code.Instr{Opcode: code.Match, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{regexp.MustCompile("a*b")},
		[]string{},
		[]interface{}{},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{{text: "aaaab", indices: []int{0, 5}}}},
	},
	{
		"cmp lt",
		code.Instr{Opcode: code.Cmp, Operand: -1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1, "2"},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp eq",
		code.Instr{Opcode: code.Cmp, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"2", "2"},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp le",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, "2"},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp ne",
		code.Instr{Opcode: code.Cmp, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1", "2"},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp ge",
		code.Instr{Opcode: code.Cmp, Operand: -1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 2},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt float float",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"2.0", "1.0"},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt float int",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1.0", "2"},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt int float",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1", "2.0"},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp eq string string false",
		code.Instr{Opcode: code.Cmp, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"abc", "def"},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp eq string string true",
		code.Instr{Opcode: code.Cmp, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"abc", "abc"},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt float float",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2.0, 1.0},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt float int",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cmp gt int float",
		code.Instr{Opcode: code.Cmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1, 2.0},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"jnm",
		code.Instr{Opcode: code.Jnm, Operand: 37, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{false},
		[]interface{}{},
		thread{pc: 37, matches: []matchResult{}},
	},
	{
		"jm",
		code.Instr{Opcode: code.Jm, Operand: 37, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{false},
		[]interface{}{},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"jmp",
		code.Instr{Opcode: code.Jmp, Operand: 37, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{},
		thread{pc: 37, matches: []matchResult{}},
	},
	{
		"strptime",
		code.Instr{Opcode: code.Strptime, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"2012/01/18 06:25:00", "2006/01/02 15:04:05"},
		[]interface{}{},
		thread{
			pc: 0, time: time.Date(2012, 1, 18, 6, 25, 0, 0, time.UTC),
			matches: []matchResult{},
		},
	},
	{
		"iadd",
		code.Instr{Opcode: code.Iadd, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(3)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"isub",
		code.Instr{Opcode: code.Isub, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"imul",
		code.Instr{Opcode: code.Imul, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(2)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"idiv",
		code.Instr{Opcode: code.Idiv, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{4, 2},
		[]interface{}{int64(2)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"imod",
		code.Instr{Opcode: code.Imod, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{4, 2},
		[]interface{}{int64(0)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"imod 2",
		code.Instr{Opcode: code.Imod, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{3, 2},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"tolower",
		code.Instr{Opcode: code.Tolower, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"mIxeDCasE"},
		[]interface{}{"mixedcase"},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"length",
		code.Instr{Opcode: code.Length, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1234"},
		[]interface{}{4},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"length 0",
		code.Instr{Opcode: code.Length, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{""},
		[]interface{}{0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"shl",
		code.Instr{Opcode: code.Shl, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(4)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"shr",
		code.Instr{Opcode: code.Shr, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"and",
		code.Instr{Opcode: code.And, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(0)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"or",
		code.Instr{Opcode: code.Or, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(3)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"xor",
		code.Instr{Opcode: code.Xor, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 1},
		[]interface{}{int64(3)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"xor 2",
		code.Instr{Opcode: code.Xor, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 3},
		[]interface{}{int64(1)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"xor 3",
		code.Instr{Opcode: code.Xor, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{-1, 3},
		[]interface{}{int64(^3)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"neg",
		code.Instr{Opcode: code.Neg, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{0},
		[]interface{}{int64(-1)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"not",
		code.Instr{Opcode: code.Not, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{false},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"pow",
		code.Instr{Opcode: code.Ipow, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2, 2},
		[]interface{}{int64(4)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"s2i pop",
		code.Instr{Opcode: code.S2i, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"ff", 16},
		[]interface{}{int64(255)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"s2i",
		code.Instr{Opcode: code.S2i, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"190"},
		[]interface{}{int64(190)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"s2f",
		code.Instr{Opcode: code.S2f, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"1.0"},
		[]interface{}{float64(1.0)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"i2f",
		code.Instr{Opcode: code.I2f, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1},
		[]interface{}{float64(1.0)},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"settime",
		code.Instr{Opcode: code.Settime, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{int64(0)},
		[]interface{}{},
		thread{pc: 0, time: time.Unix(0, 0).UTC(), matches: []matchResult{}},
	},
	{
		"push int",
		code.Instr{Opcode: code.Push, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{1},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"push float",
		code.Instr{Opcode: code.Push, Operand: 1.0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{1.0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"setmatched false",
		code.Instr{Opcode: code.Setmatched, Operand: false, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{},
		thread{matched: false, pc: 0, matches: []matchResult{}},
	},
	{
		"setmatched true",
		code.Instr{Opcode: code.Setmatched, Operand: true, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{},
		thread{matched: true, pc: 0, matches: []matchResult{}},
	},
	{
		"otherwise",
		code.Instr{Opcode: code.Otherwise, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{true},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fadd",
		code.Instr{Opcode: code.Fadd, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2.0},
		[]interface{}{3.0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fsub",
		code.Instr{Opcode: code.Fsub, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2.0},
		[]interface{}{-1.0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fmul",
		code.Instr{Opcode: code.Fmul, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2.0},
		[]interface{}{2.0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fdiv",
		code.Instr{Opcode: code.Fdiv, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2.0},
		[]interface{}{0.5},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fmod",
		code.Instr{Opcode: code.Fmod, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2.0},
		[]interface{}{1.0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fpow",
		code.Instr{Opcode: code.Fpow, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{2.0, 2.0},
		[]interface{}{4.0},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"getfilename",
		code.Instr{Opcode: code.Getfilename, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{},
		[]interface{}{testFilename},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"i2s",
		code.Instr{Opcode: code.I2s, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1},
		[]interface{}{"1"},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"f2s",
		code.Instr{Opcode: code.F2s, Operand: nil, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{3.1},
		[]interface{}{"3.1"},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"cat",
		code.Instr{Opcode: code.Cat, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"first", "second"},
		[]interface{}{"firstsecond"},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"icmp gt false",
		code.Instr{Opcode: code.Icmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1, 2},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"fcmp gt false",
		code.Instr{Opcode: code.Fcmp, Operand: 1, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{1.0, 2.0},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"scmp eq false",
		code.Instr{Opcode: code.Scmp, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"abc", "def"},
		[]interface{}{false},
		thread{pc: 0, matches: []matchResult{}},
	},
	{
		"subst",
		code.Instr{Opcode: code.Subst, Operand: 0, SourceLine: 0},
		[]*regexp.Regexp{},
		[]string{},
		[]interface{}{"aa" /*old*/, "a" /*new*/, "caat"},
		[]interface{}{"cat"},
		thread{pc: 0, matches: []matchResult{}},
	},
}

const testFilename = "test"

// Testcode.Instrs tests that each instruction behaves as expected through one
// instruction cycle.
func TestInstrs(t *testing.T) {
	for _, tc := range instructions {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var m []*metrics.Metric
			m = append(m,
				metrics.NewMetric("foo", "test", metrics.Counter, metrics.Int),
				metrics.NewMetric("bar", "test", metrics.Counter, metrics.Int),
				metrics.NewMetric("quux", "test", metrics.Gauge, metrics.Float))
			obj := &code.Object{Regexps: tc.re, Strings: tc.str, Metrics: m, Program: []code.Instr{tc.i}}
			v := New(tc.name, obj, true, nil, false, false)
			v.t = new(thread)
			v.t.stack = make([]interface{}, 0)
			for _, item := range tc.reversedStack {
				v.t.Push(item)
			}
			v.t.matches = make([]matchResult, len(obj.Regexps))
			v.input = logline.New(context.Background(), testFilename, logline.GetHash(testFilename), "aaaab")
			v.execute(v.t, tc.i)
			if v.terminate {
				t.Fatalf("Execution failed, see info log.")
			}

			testutil.ExpectNoDiff(t, tc.expectedStack, v.t.stack)

			tc.expectedThread.stack = tc.expectedStack

			testutil.ExpectNoDiff(t, &tc.expectedThread, v.t, testutil.AllowUnexported(thread{}), testutil.AllowUnexported(matchResult{}))
		})
	}
}

// makeVM is a helper method for construction a single-instruction VM.
func makeVM(i code.Instr, m []*metrics.Metric) *VM {
	obj := &code.Object{Metrics: m, Program: []code.Instr{i}}
	v := New("test", obj, true, nil, false, false)
	v.t = new(thread)
	v.t.stack = make([]interface{}, 0)
	v.t.matches = make([]matchResult, 0)
	v.t.keysBuf = make([]string, 64)
	v.input = logline.New(context.Background(), testFilename, logline.GetHash(testFilename), "aaaab")
	return v
}

// makeMetrics returns a few useful metrics for observing under test.
func makeMetrics() []*metrics.Metric {
	var m []*metrics.Metric
	m = append(m,
		metrics.NewMetric("a", "tst", metrics.Counter, metrics.Int),
		metrics.NewMetric("b", "tst", metrics.Counter, metrics.Float),
		metrics.NewMetric("c", "tst", metrics.Gauge, metrics.String),
		metrics.NewMetric("d", "tst", metrics.Histogram, metrics.Float),
	)
	return m
}

type datumStoreTests struct {
	name     string
	i        code.Instr
	d        int // index of a metric in makeMetrics
	setup    func(t *thread, d datum.Datum)
	expected string
}

// code.Instructions with datum store side effects.
func TestDatumSetInstrs(t *testing.T) {
	tests := []datumStoreTests{
		{
			name: "simple inc",
			i:    code.Instr{Opcode: code.Inc, Operand: nil, SourceLine: 0},
			d:    0,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
			},
			expected: "1",
		},
		{
			name: "inc by int",
			i:    code.Instr{Opcode: code.Inc, Operand: 0, SourceLine: 0},
			d:    0,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push(2)
			},
			expected: "2",
		},
		{
			name: "inc by str",
			i:    code.Instr{Opcode: code.Inc, Operand: 0, SourceLine: 0},
			d:    0,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push("4")
			},
			expected: "4",
		},
		{
			name: "iset",
			i:    code.Instr{Opcode: code.Iset, Operand: nil, SourceLine: 0},
			d:    0,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push(2)
			},
			expected: "2",
		},
		{
			name: "iset str",
			i:    code.Instr{Opcode: code.Iset, Operand: nil, SourceLine: 0},
			d:    0,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push("3")
			},
			expected: "3",
		},
		{
			name: "fset",
			i:    code.Instr{Opcode: code.Fset, Operand: nil, SourceLine: 0},
			d:    1,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push(3.1)
			},
			expected: "3.1",
		},
		{
			name: "fset str",
			i:    code.Instr{Opcode: code.Fset, Operand: nil, SourceLine: 0},
			d:    1,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push("4.1")
			},
			expected: "4.1",
		},
		{
			name: "sset",
			i:    code.Instr{Opcode: code.Sset, Operand: nil, SourceLine: 0},
			d:    2,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push("4.1")
			},
			expected: "4.1",
		},
		{
			name: "dec",
			i:    code.Instr{Opcode: code.Dec, Operand: nil, SourceLine: 0},
			d:    0,
			setup: func(t *thread, d datum.Datum) {
				datum.SetInt(d, 1, time.Now())
				t.Push(d)
			},
			expected: "0",
		},
		{
			name: "set hist",
			i:    code.Instr{Opcode: code.Fset, Operand: nil, SourceLine: 0},
			d:    3,
			setup: func(t *thread, d datum.Datum) {
				t.Push(d)
				t.Push(3.1)
			},
			expected: "3.1",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			m := makeMetrics()
			v := makeVM(test.i, m)
			d, err := m[test.d].GetDatum()
			testutil.FatalIfErr(t, err)
			test.setup(v.t, d)
			v.execute(v.t, v.prog[0])
			if v.terminate {
				t.Fatalf("Execution failed, see INFO log.")
			}
			d, err = m[test.d].GetDatum()
			testutil.FatalIfErr(t, err)
			if d.ValueString() != test.expected {
				t.Errorf("unexpected value for datum %#v, want: %s, got %s", d, test.expected, d.ValueString())
			}
		})
	}
}

func TestStrptimeWithTimezone(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Skip("Skipping, timezone database not available:", err)
	}
	obj := &code.Object{Program: []code.Instr{{Opcode: code.Strptime, Operand: 0, SourceLine: 0}}}
	vm := New("strptimezone", obj, true, loc, false, false)
	vm.t = new(thread)
	vm.t.stack = make([]interface{}, 0)
	vm.t.Push("2012/01/18 06:25:00")
	vm.t.Push("2006/01/02 15:04:05")
	vm.execute(vm.t, obj.Program[0])
	if vm.t.time != time.Date(2012, 1, 18, 6, 25, 0, 0, loc) {
		t.Errorf("Time didn't parse with location: %s received", vm.t.time)
	}
}

func TestStrptimeWithoutTimezone(t *testing.T) {
	obj := &code.Object{Program: []code.Instr{{Opcode: code.Strptime, Operand: 0, SourceLine: 0}}}
	vm := New("strptimezone", obj, true, nil, false, false)
	vm.t = new(thread)
	vm.t.stack = make([]interface{}, 0)
	vm.t.Push("2012/01/18 06:25:00")
	vm.t.Push("2006/01/02 15:04:05")
	vm.execute(vm.t, obj.Program[0])
	if vm.t.time != time.Date(2012, 1, 18, 6, 25, 0, 0, time.UTC) {
		t.Errorf("Time didn't parse with location: %s received", vm.t.time)
	}
}

// code.Instructions with datum retrieve.
func TestDatumFetchInstrs(t *testing.T) {
	var m []*metrics.Metric
	m = append(m,
		metrics.NewMetric("a", "tst", metrics.Counter, metrics.Int),
		metrics.NewMetric("b", "tst", metrics.Counter, metrics.Float),
		metrics.NewMetric("c", "tst", metrics.Text, metrics.String))

	{
		// iget
		v := makeVM(code.Instr{Opcode: code.Iget, Operand: nil, SourceLine: 0}, m)
		d, err := m[0].GetDatum()
		testutil.FatalIfErr(t, err)
		datum.SetInt(d, 37, time.Now())
		v.t.Push(d)
		v.execute(v.t, v.prog[0])
		if v.terminate {
			t.Fatalf("Execution failed, see info log.")
		}
		i, err := v.t.PopInt()
		if err != nil {
			t.Fatalf("Execution failed, see info; %v", err)
		}
		if i != 37 {
			t.Errorf("unexpected value %d", i)
		}
	}

	{
		// fget
		v := makeVM(code.Instr{Opcode: code.Fget, Operand: nil, SourceLine: 0}, m)
		d, err := m[1].GetDatum()
		testutil.FatalIfErr(t, err)
		datum.SetFloat(d, 12.1, time.Now())
		v.t.Push(d)
		v.execute(v.t, v.prog[0])
		if v.terminate {
			t.Fatalf("Execution failed, see info log.")
		}
		i, err := v.t.PopFloat()
		if err != nil {
			t.Fatalf("Execution failed, see info: %v", err)
		}
		if i != 12.1 {
			t.Errorf("unexpected value %f", i)
		}
	}

	{
		// sget
		v := makeVM(code.Instr{Opcode: code.Sget, Operand: nil, SourceLine: 0}, m)
		d, err := m[2].GetDatum()
		testutil.FatalIfErr(t, err)
		datum.SetString(d, "aba", time.Now())
		v.t.Push(d)
		v.execute(v.t, v.prog[0])
		if v.terminate {
			t.Fatalf("Execution failed, see info log.")
		}
		i, err := v.t.PopString()
		if err != nil {
			t.Fatalf("Execution failed, see info log: %v", err)
		}
		if i != "aba" {
			t.Errorf("unexpected value %q", i)
		}
	}
}

func TestDeleteInstrs(t *testing.T) {
	var m []*metrics.Metric
	m = append(m,
		metrics.NewMetric("a", "tst", metrics.Counter, metrics.Int, "a"),
	)

	_, err := m[0].GetDatum("z")
	testutil.FatalIfErr(t, err)

	v := makeVM(code.Instr{Opcode: code.Expire, Operand: 1, SourceLine: 0}, m)
	v.t.Push(time.Hour)
	v.t.Push("z")
	v.t.Push(m[0])
	v.execute(v.t, v.prog[0])
	if v.terminate {
		t.Fatal("execution failed, see info log")
	}
	lv := m[0].FindLabelValueOrNil([]string{"z"})
	if lv == nil {
		t.Fatalf("couldn;t find label value in metric %#v", m[0])
	}
	if lv.Expiry != time.Hour {
		t.Fatalf("Expiry not correct, is %v", lv.Expiry)
	}
}

func TestTimestampInstr(t *testing.T) {
	var m []*metrics.Metric
	now := time.Now().UTC()
	v := makeVM(code.Instr{Opcode: code.Timestamp, Operand: nil, SourceLine: 0}, m)
	v.execute(v.t, v.prog[0])
	if v.terminate {
		t.Fatal("execution failed, see info log")
	}
	tos := time.Unix(v.t.Pop().(int64), 0).UTC()
	if now.Before(tos) {
		t.Errorf("Expecting timestamp to be after %s, was %s", now, tos)
	}

	newT := time.Unix(37, 0).UTC()
	v.t.time = newT
	v.execute(v.t, v.prog[0])

	if v.terminate {
		t.Fatal("execution failed, see info log")
	}
	tos = time.Unix(v.t.Pop().(int64), 0).UTC()
	if tos != newT {
		t.Errorf("Expecting timestamp to be %s, was %s", newT, tos)
	}
}

func TestProcessLogLineDoesNotPinLargeCaptureGroup(t *testing.T) {
	// Build a program that matches a regex, captures a substring, and stores it
	// as a metric label key via Dload.  The stored key must not pin the original
	// log line's backing array.
	re := regexp.MustCompile(`.(KEY).`)

	m := []*metrics.Metric{
		metrics.NewMetric("m", "tst", metrics.Counter, metrics.Int, "label"),
	}
	prog := []code.Instr{
		{Opcode: code.Match, Operand: 0, SourceLine: 0},
		{Opcode: code.Push, Operand: 0, SourceLine: 0},
		{Opcode: code.Capref, Operand: 1, SourceLine: 0},
		{Opcode: code.Mload, Operand: 0, SourceLine: 0},
		{Opcode: code.Dload, Operand: 1, SourceLine: 0},
	}
	obj := &code.Object{
		Metrics: m,
		Program: prog,
		Regexps: []*regexp.Regexp{re},
	}
	v := New("leaktest", obj, true, nil, false, false)

	// Allocate a large log line (~200KB).  The regex `.(KEY).` matches
	// "aKEYb" and captures "KEY".  Save the data pointer of the capture to
	// compare against the stored label value later.
	const offset = 100000
	large := strings.Repeat("x", offset) + "aKEYb" + strings.Repeat("y", 100000)
	capture := large[offset+1 : offset+4] // "KEY" — shares backing with large
	origPtr := unsafe.StringData(capture)

	v.ProcessLogLine(context.Background(), logline.New(context.Background(), "test", 0, large))
	runtime.KeepAlive(large) // ensure large survives until ProcessLogLine completes

	// Retrieve the stored label value.
	lv := m[0].FindLabelValueOrNil([]string{"KEY"})
	if lv == nil {
		t.Fatal("label value not found")
	}
	stored := lv.Labels[0]

	// Drop all references to the large backing array.  If strings.Clone was not
	// used, the stored key still pins the original backing array via the shared
	// data pointer.
	large = ""
	runtime.GC()

	// If the stored key shares the original backing array, its data pointer
	// matches.  If strings.Clone allocated a fresh copy, they differ.
	if unsafe.StringData(stored) == origPtr {
		t.Error("stored label key shares backing array with original log line; Dload should use strings.Clone to detach it")
	}
}

func TestProcessLogLineClearsInput(t *testing.T) {
	obj := &code.Object{Program: []code.Instr{}} // empty program returns immediately
	v := New("leaktest", obj, true, nil, false, false)
	line := logline.New(context.Background(), "test", 0, "hello world")

	v.ProcessLogLine(context.Background(), line)
	if v.input != nil {
		t.Error("v.input was not cleared after ProcessLogLine returned")
	}
	if v.t != nil {
		t.Error("v.t was not cleared after ProcessLogLine returned")
	}
}
