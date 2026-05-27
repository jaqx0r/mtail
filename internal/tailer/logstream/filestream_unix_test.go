// Copyright 2020 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

//go:build unix

package logstream_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/jaqx0r/mtail/internal/logline"
	"github.com/jaqx0r/mtail/internal/tailer/logstream"
	"github.com/jaqx0r/mtail/internal/testutil"
	"github.com/jaqx0r/mtail/internal/waker"
)

// TestFileStreamRotation is a unix-specific test because on Windows, files cannot be removed
// or renamed while there is an open read handle on them. Instead, log rotation would
// have to be implemented by copying and then truncating the original file. That test
// case is already covered by TestFileStreamTruncation.
func TestFileStreamRotation(t *testing.T) {
	var wg sync.WaitGroup

	tmpDir := testutil.TestTempDir(t)

	name := filepath.Join(tmpDir, "log")
	f := testutil.TestOpenFile(t, name)
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	waker, awaken := waker.NewTest(ctx, 1, "stream")

	// OneShotDisabled because we hit EOF and need to wait.
	fs, err := logstream.New(ctx, &wg, waker, name, logstream.OneShotDisabled)
	testutil.FatalIfErr(t, err)

	expected := []*logline.LogLine{
		{Context: context.TODO(), Filename: name, Line: "1", Filenamehash: logline.GetHash(name)},
		{Context: context.TODO(), Filename: name, Line: "2", Filenamehash: logline.GetHash(name)},
	}
	checkLineDiff := testutil.ExpectLinesReceivedNoDiff(t, expected, fs.Lines())

	awaken(1, 1) // sync to eof

	glog.Info("write 1")
	testutil.WriteString(t, f, "1\n")
	awaken(1, 1)

	glog.Info("rename")
	testutil.FatalIfErr(t, os.Rename(name, name+".1"))
	// filestream won't notice if there's a synchronisation point between
	// rename and create, that path relies on the tailer
	f = testutil.TestOpenFile(t, name)
	defer f.Close()

	awaken(1, 1)
	glog.Info("write 2")
	testutil.WriteString(t, f, "2\n")
	awaken(1, 1)

	cancel()
	wg.Wait()

	checkLineDiff()

	if v := <-fs.Lines(); v != nil {
		t.Errorf("expecting filestream to be complete because stopped")
	}
}

func TestFileStreamURL(t *testing.T) {
	var wg sync.WaitGroup

	tmpDir := testutil.TestTempDir(t)

	name := filepath.Join(tmpDir, "log")
	f := testutil.TestOpenFile(t, name)
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	waker, awaken := waker.NewTest(ctx, 1, "stream")

	fs, err := logstream.New(ctx, &wg, waker, "file://"+name, logstream.OneShotDisabled)
	testutil.FatalIfErr(t, err)

	expected := []*logline.LogLine{
		{Context: context.TODO(), Filename: name, Line: "yo", Filenamehash: logline.GetHash(name)},
	}
	checkLineDiff := testutil.ExpectLinesReceivedNoDiff(t, expected, fs.Lines())

	awaken(1, 1)

	testutil.WriteString(t, f, "yo\n")
	awaken(1, 1)

	cancel()
	wg.Wait()

	checkLineDiff()

	if v := <-fs.Lines(); v != nil {
		t.Errorf("expecting filestream to be complete because stopped")
	}
}

// manualWaker is a simple waker that lets the test directly signal wakeups
// by closing and replacing the wake channel.
type manualWaker struct {
	wake chan struct{}
}

func (w *manualWaker) Wake() <-chan struct{} {
	return w.wake
}

// wake signals the waker, closing the current channel and creating a new one.
func (w *manualWaker) wakeAndReset() {
	close(w.wake)
	w.wake = make(chan struct{})
}

// TestFileStreamRotationPermissionDenied tests that when a rotation is detected
// and the new file cannot be opened (permission denied), the filestream closes
// its Lines channel so the tailer can clean up and retry later.
func TestFileStreamRotationPermissionDenied(t *testing.T) {
	testutil.SkipIfRoot(t)
	var wg sync.WaitGroup

	tmpDir := testutil.TestTempDir(t)

	name := filepath.Join(tmpDir, "log")
	f := testutil.TestOpenFile(t, name)
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mw := &manualWaker{wake: make(chan struct{})}

	fs, err := logstream.New(ctx, &wg, mw, name, logstream.OneShotDisabled)
	testutil.FatalIfErr(t, err)

	expectedLines := []*logline.LogLine{
		{Context: context.TODO(), Filename: name, Line: "1", Filenamehash: logline.GetHash(name)},
	}
	checkLineDiff := testutil.ExpectLinesReceivedNoDiff(t, expectedLines, fs.Lines())

	mw.wakeAndReset() // sync to EOF

	testutil.WriteString(t, f, "1\n")
	mw.wakeAndReset()

	// Rotate: rename old file, create new file with no read permissions.
	testutil.FatalIfErr(t, os.Rename(name, name+".1"))
	f2, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
	testutil.FatalIfErr(t, err)
	f2.Close()

	// Wake the stream to detect rotation.  The goroutine will see a new
	// inode, try to open the new file, fail with permission denied, and
	// (with our fix) close fs.lines.
	mw.wakeAndReset()

	// Wait for the Lines channel to close, indicating the goroutine exited.
	ok, err := testutil.DoOrTimeout(func() (bool, error) {
		select {
		case _, stillOpen := <-fs.Lines():
			return !stillOpen, nil
		default:
			return false, nil
		}
	}, time.Second, 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("Lines channel not closed after rotation + permission denied")
	}

	// Wait for the goroutine to fully shut down.
	wg.Wait()

	// Fix permissions on the new file.
	testutil.FatalIfErr(t, os.Chmod(name, 0o666))

	// Now create a new stream; this should succeed.
	fs2, err := logstream.New(ctx, &wg, mw, name, logstream.OneShotDisabled)
	testutil.FatalIfErr(t, err)

	expectedLines2 := []*logline.LogLine{
		{Context: context.TODO(), Filename: name, Line: "2", Filenamehash: logline.GetHash(name)},
	}
	checkLineDiff2 := testutil.ExpectLinesReceivedNoDiff(t, expectedLines2, fs2.Lines())

	f3, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0o666)
	testutil.FatalIfErr(t, err)
	defer f3.Close()

	mw.wakeAndReset() // sync to EOF

	testutil.WriteString(t, f3, "2\n")
	mw.wakeAndReset()

	cancel()
	wg.Wait()

	checkLineDiff()
	checkLineDiff2()
}

// TestFileStreamOpenFailure is a unix-specific test because on Windows, it is not possible to create a file
// that you yourself cannot read (minimum permissions are 0222).
func TestFileStreamOpenFailure(t *testing.T) {
	// can't force a permission denied if run as root
	testutil.SkipIfRoot(t)
	var wg sync.WaitGroup

	tmpDir := testutil.TestTempDir(t)

	name := filepath.Join(tmpDir, "log")
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
	//nolint:staticcheck // test code
	defer f.Close()

	testutil.FatalIfErr(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	waker := waker.NewTestAlways()

	_, err = logstream.New(ctx, &wg, waker, name, logstream.OneShotEnabled)
	if err == nil || !os.IsPermission(err) {
		t.Errorf("Expected a permission denied error, got: %v", err)
	}
	cancel()
}
