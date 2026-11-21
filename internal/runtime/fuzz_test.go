// Copyright 2011 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package runtime

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaqx0r/mtail/internal/logline"
	"github.com/jaqx0r/mtail/internal/runtime/compiler"
	"github.com/jaqx0r/mtail/internal/runtime/vm"
)

// U+2424 SYMBOL FOR NEWLINE
const SEP = "␤"

// Enable this when debugging with a fuzz crash artifact; it slows the fuzzer down when enabled.
const dumpDebug = false

func FuzzE2E(f *testing.F) {
	// We need to successfully parse flags to initialize the glog logger used
	// by the compiler, but the fuzzer gets called with flags captured by the
	// libfuzzer main, which we don't want to intercept here.
	flag.CommandLine.Parse([]string{})

	entries, err := os.ReadDir("fuzz")
	if err != nil {
		if !os.IsNotExist(err) {
			f.Fatalf("os.ReadDir(fuzz): unexpected error: %v", err)
		} else {
			f.Logf("no `fuzz` corpus found, assuming in oss-fuzz")
		}
	} else {
		for i, entry := range entries {
			if entry.IsDir() {
				continue
			}
			path := filepath.Join("fuzz", entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				f.Fatalf("os.ReadFile(%s): %v", path, err)
			}
			f.Logf("Adding %s as seed %d", path, i)
			f.Add(data)
		}
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// Data contains the program and sample input, separated by SEP.
		offset := bytes.Index(data, []byte(SEP))
		if offset < 0 {
			// If no SEP, then append one and an empty line of input.
			offset = len(data)
			data = append(data, []byte(SEP+"\n")...)
		}
		t.Logf("data len %d, offset is %d, input starts at %d\n", len(data), offset, offset+len(SEP))

		cOpts := []compiler.Option{}
		if dumpDebug {
			cOpts = append(cOpts, compiler.EmitAst(), compiler.EmitAstTypes())
		}
		c, err := compiler.New(cOpts...)
		if err != nil {
			t.Error(err)
		}
		obj, err := c.Compile("fuzz", bytes.NewReader(data[:offset]))
		if err != nil {
			t.Log(err)
			return
		}
		v := vm.New("fuzz", obj, false, nil, dumpDebug, dumpDebug)
		if dumpDebug {
			t.Logf("%s\n", v.DumpByteCode())
		}
		v.HardCrash = true
		scanner := bufio.NewScanner(bytes.NewBuffer(data[offset+len(SEP):]))
		hash := logline.GetHash("fuzz")
		for scanner.Scan() {
			v.ProcessLogLine(context.Background(), logline.New(context.Background(), "fuzz", hash, scanner.Text()))
		}
	})
}
