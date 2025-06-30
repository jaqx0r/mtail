// Copyright 2019 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package mtail_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/glog"
	"github.com/jaqx0r/mtail/internal/mtail"
	"github.com/jaqx0r/mtail/internal/testutil"
)

func TestLogMapping(t *testing.T) {
	testutil.SkipIfShort(t)
	tmpDir := testutil.TestTempDir(t)

	logDir := filepath.Join(tmpDir, "logs")
	progDir := filepath.Join(tmpDir, "progs")
	err := os.Mkdir(logDir, 0o700)
	testutil.FatalIfErr(t, err)
	err = os.Mkdir(progDir, 0o700)
	testutil.FatalIfErr(t, err)

	readLogFile := filepath.Join(logDir, "readlog")
	unreadLogFile := filepath.Join(logDir, "ignorelog")

	readFile := testutil.TestOpenFile(t, readLogFile)
	unreadFile := testutil.TestOpenFile(t, unreadLogFile)
	os.WriteFile(filepath.Join(progDir, "mapping.mtail"), []byte("logmapping \""+progDir+"/readlog\"\n"), 0o666)
	os.WriteFile(filepath.Join(progDir, "notmapping.mtail"), []byte("logmapping \""+progDir+"/ignorelog\"\n"), 0o666)
	defer readFile.Close()
	defer unreadFile.Close()

	m, stopM := mtail.TestStartServer(t, 1, 1, mtail.ProgramPath(progDir), mtail.LogPathPatterns(logDir+"/readlog", logDir+"/ignorelog"))

	defer stopM()

	m.AwakenPatternPollers(1, 1)
	m.AwakenLogStreams(1, 1) // Force read to EOF

	{
		lineCountCheck := m.ExpectMapExpvarDeltaWithDeadline("prog_lines_total", "mapping.mtail", 1)
		n, err := readFile.WriteString("line 1\n")
		testutil.FatalIfErr(t, err)
		glog.Infof("Wrote %d bytes", n)
		m.AwakenLogStreams(1, 1)
		lineCountCheck()
	}

	{
		lineCountCheck := m.ExpectMapExpvarDeltaWithDeadline("prog_lines_total", "notmapping.mtail", 0)
		n, err := unreadFile.WriteString("line 2\n")
		testutil.FatalIfErr(t, err)
		glog.Infof("Wrote %d bytes", n)
		m.AwakenLogStreams(1, 1)
		lineCountCheck()
	}
}
