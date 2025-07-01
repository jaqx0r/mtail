// Copyright 2015 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package runtime

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/golang/glog"
	"github.com/jaqx0r/mtail/internal/logline"
	"github.com/jaqx0r/mtail/internal/metrics"
	"github.com/jaqx0r/mtail/internal/testutil"
)

func TestNewRuntime(t *testing.T) {
	store := metrics.NewStore()
	lines := make(chan *logline.LogLine)
	var wg sync.WaitGroup
	_, err := New(lines, &wg, "", store)
	testutil.FatalIfErr(t, err)
	close(lines)
	wg.Wait()
}

func TestNewRuntimeErrors(t *testing.T) {
	store := metrics.NewStore()
	lines := make(chan *logline.LogLine)
	var wg sync.WaitGroup
	_, err := New(lines, nil, "", store)
	if err == nil {
		t.Error("New(..., nil) expecting error, got nil")
	}
	_, err = New(lines, &wg, "", nil)
	if err == nil {
		t.Error("New(..., nil) expecting error, got nil")
	}
}

func TestCompileAndRun(t *testing.T) {
	testProgram := "/$/ {}\n"
	store := metrics.NewStore()
	lines := make(chan *logline.LogLine)
	var wg sync.WaitGroup
	l, err := New(lines, &wg, "", store)
	testutil.FatalIfErr(t, err)
	if err := l.CompileAndRun("Test", strings.NewReader(testProgram)); err != nil {
		t.Errorf("CompileAndRun returned error: %s", err)
	}
	l.handleMu.Lock()
	if len(l.handles) < 1 {
		t.Errorf("no vm handles: %v", l.handles)
	}
	l.handleMu.Unlock()
	l.handleMu.Lock()
	h := l.handles["Test"]
	if h == nil {
		t.Errorf("No handle for Test: %v", l.handles)
	}
	l.handleMu.Unlock()
	close(lines)
	wg.Wait()
}

var testProgram = "/$/ {}\n"

var testProgFiles = []string{
	"test.wrongext",
	"test.mtail",
	".test",
}

func TestLoadProg(t *testing.T) {
	store := metrics.NewStore()
	tmpDir := testutil.TestTempDir(t)

	lines := make(chan *logline.LogLine)
	var wg sync.WaitGroup
	l, err := New(lines, &wg, tmpDir, store)
	testutil.FatalIfErr(t, err)

	for _, name := range testProgFiles {
		f := testutil.TestOpenFile(t, filepath.Join(tmpDir, name))
		n, err := f.WriteString(testProgram)
		testutil.FatalIfErr(t, err)
		glog.Infof("Wrote %d bytes", n)
		err = l.LoadProgram(filepath.Join(tmpDir, name))
		testutil.FatalIfErr(t, err)
		f.Close()
	}
	close(lines)
	wg.Wait()
}

func TestLogLineFilter(t *testing.T) {

	tests := []struct {
		lines       []*logline.LogLine
		logmappings map[uint32]struct{}
		expected    []*logline.LogLine
	}{
		{
			// Test case where one file is processed and one not.
			lines: []*logline.LogLine{
				{
					Filenamehash: 12345,
					Line:         "This is a valid log line",
				},
				{
					Filenamehash: 67890,
					Line:         "This log line should be ignored",
				},
			},
			expected: []*logline.LogLine{
				{
					Filenamehash: 12345,
					Line:         "This is a valid log line",
				},
			},
			logmappings: map[uint32]struct{}{
				12345: {}, // This maps to the first line.
			},
		},
		{
			// Test case where both file are processed.
			lines: []*logline.LogLine{
				{
					Filenamehash: 12345,
					Line:         "This is a valid log line",
				},
				{
					Filenamehash: 67890,
					Line:         "This is a valid log line",
				},
			},
			expected: []*logline.LogLine{
				{
					Filenamehash: 12345,
					Line:         "This is a valid log line",
				},
				{
					Filenamehash: 67890,
					Line:         "This is a valid log line",
				},
			},
			logmappings: map[uint32]struct{}{
				12345: {}, // This maps to the first line.
				67890: {}, // This maps to the first line.
			},
		},
		{
			// Test case where file process because no mapping.
			lines: []*logline.LogLine{
				{
					Filenamehash: 12345,
					Line:         "This is a valid log line",
				},
			},
			expected: []*logline.LogLine{
				{
					Filenamehash: 12345,
					Line:         "This is a valid log line",
				},
			},
			logmappings: map[uint32]struct{}{}, // empty to all logs match
		},
		{
			// empty test case
			lines:       []*logline.LogLine{},
			expected:    []*logline.LogLine{},
			logmappings: map[uint32]struct{}{}, // empty to all logs match
		},
	}

	for _, tc := range tests {
		// Create a channel for log lines.
		lines := make(chan *logline.LogLine, len(tc.lines))
		for _, line := range tc.lines {
			lines <- line
		}
		close(lines)

		var wg sync.WaitGroup

		store := metrics.NewStore()
		// Create a Runtime instance with a logmapping for the first line.
		r, err := New(lines, &wg, "", store)
		testutil.FatalIfErr(t, err)

		// Add a logmapping for the first line.
		r.logmappings["test_program"] = tc.logmappings

		// Add a mock program handle for "test_program".
		linesReceived := make(chan *logline.LogLine, len(tc.expected))

		// Start processing lines.
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				if _, ok := r.logmappings["test_program"][line.Filenamehash]; ok {
					linesReceived <- line
				}
			}
		}()

		// Wait for the Runtime to finish processing.
		wg.Wait()

		// Validate the lines received by the "test_program" handle.

		testutil.ExpectLinesReceivedNoDiff(t, tc.expected, linesReceived)
	}

}
