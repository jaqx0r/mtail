// Copyright 2025 The mtail Authors. All Rights Reserved.
// This file is available under the Apache license.

package mtail_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/jaqx0r/mtail/internal/mtail"
	"github.com/jaqx0r/mtail/internal/testutil"
)

func TestIncludeLog(t *testing.T) {
	testutil.SkipIfShort(t)

	logDir := testutil.TestTempDir(t)

	// The logfile isn't known until now so we need to create out own Program Dir, copy in the program and update the log_filter
	progDir := testutil.TestTempDir(t)

	content, err := ioutil.ReadFile("../../examples/log_filter.mtail")
	testutil.FatalIfErr(t, err)

	logFile := filepath.Join(logDir, "log")

	updatedContent := strings.ReplaceAll(string(content), "<<ADD_LOGFILE_HERE>>", logFile)

	progFile := filepath.Join(progDir, "/test.mtail")

	// Write the updated content to the destination file
	err = ioutil.WriteFile(progFile, []byte(updatedContent), 0644)

	testutil.FatalIfErr(t, err)

	// Verify the file was written correctly
	writtenContent, err := ioutil.ReadFile(progFile)
	testutil.FatalIfErr(t, err)
	if string(writtenContent) != updatedContent {
		t.Fatalf("File content mismatch: expected %q, got %q", updatedContent, string(writtenContent))
	}

	m, stopM := mtail.TestStartServer(t, 1, 1, mtail.LogPathPatterns(logFile), mtail.ProgramPath(progFile))
	defer stopM()

	lineCountCheck := m.ExpectMapExpvarDeltaWithDeadline("log_lines_total", logFile, 3)
	progLineCountCheck := m.ExpectMapExpvarDeltaWithDeadline("prog_lines_total", "test.mtail", 3)
	logCountCheck := m.ExpectExpvarDeltaWithDeadline("log_count", 1)

	f := testutil.TestOpenFile(t, logFile)
	defer f.Close()
	m.AwakenPatternPollers(1, 1) // Find `logFile`
	m.AwakenLogStreams(1, 1)     // Force a sync to EOF

	for i := 1; i <= 3; i++ {
		testutil.WriteString(t, f, fmt.Sprintf("%d\n", i))
	}
	m.AwakenLogStreams(1, 1) // Expect to read 3 lines here.

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		lineCountCheck()
	}()
	go func() {
		defer wg.Done()
		logCountCheck()
	}()
	go func() {
		defer wg.Done()
		progLineCountCheck()
	}()
	wg.Wait()
}

func TestExcludeLog(t *testing.T) {
	testutil.SkipIfShort(t)

	logDir := testutil.TestTempDir(t)

	logFile := filepath.Join(logDir, "log")

	m, stopM := mtail.TestStartServer(t, 1, 1, mtail.LogPathPatterns(logFile), mtail.ProgramPath(mtail.ProgramPath("../../examples/log_filter_ignore.mtail")))
	defer stopM()

	lineCountCheck := m.ExpectMapExpvarDeltaWithDeadline("log_lines_total", logFile, 3)
	progLineCountCheck := m.ExpectMapExpvarDeltaWithDeadline("prog_lines_total", "log_filter_ignore.mtail", 0) // not read
	logCountCheck := m.ExpectExpvarDeltaWithDeadline("log_count", 1)

	f := testutil.TestOpenFile(t, logFile)
	defer f.Close()
	m.AwakenPatternPollers(1, 1) // Find `logFile`
	m.AwakenLogStreams(1, 1)     // Force a sync to EOF

	for i := 1; i <= 3; i++ {
		testutil.WriteString(t, f, fmt.Sprintf("%d\n", i))
	}
	m.AwakenLogStreams(1, 1) // Expect to read 3 lines here.

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		lineCountCheck()
	}()
	go func() {
		defer wg.Done()
		logCountCheck()
	}()
	go func() {
		defer wg.Done()
		progLineCountCheck()
	}()
	wg.Wait()
}
