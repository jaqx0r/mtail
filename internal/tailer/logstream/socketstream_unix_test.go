// Copyright 2020 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

//go:build unix

package logstream_test

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jaqx0r/mtail/internal/logline"
	"github.com/jaqx0r/mtail/internal/tailer/logstream"
	"github.com/jaqx0r/mtail/internal/testutil"
	"github.com/jaqx0r/mtail/internal/waker"
)

func TestSocketStreamReadCompletedBecauseSocketClosed(t *testing.T) {
	for _, scheme := range []string{"unix", "tcp"} {
		scheme := scheme
		t.Run(scheme, testutil.TimeoutTest(time.Second, func(t *testing.T) { //nolint:thelper
			var wg sync.WaitGroup

			var addr string
			switch scheme {
			case "unix":
				tmpDir := testutil.TestTempDir(t)
				addr = filepath.Join(tmpDir, "sock")
			case "tcp":
				addr = fmt.Sprintf("127.0.0.1:%d", testutil.FreePort(t))
			default:
				t.Fatalf("bad scheme %s", scheme)
			}

			ctx, cancel := context.WithCancel(context.Background())
			// The stream is not shut down with cancel in this test.
			defer cancel()
			waker := waker.NewTestAlways()

			sockName := scheme + "://" + addr
			ss, err := logstream.New(ctx, &wg, waker, sockName, logstream.OneShotEnabled)
			testutil.FatalIfErr(t, err)

			expected := []*logline.LogLine{
				{Context: context.TODO(), Filename: sockName, Line: "1", Filenamehash: logline.GetHash(sockName)},
			}
			checkLineDiff := testutil.ExpectLinesReceivedNoDiff(t, expected, ss.Lines())

			s, err := net.Dial(scheme, addr)
			testutil.FatalIfErr(t, err)

			_, err = s.Write([]byte("1\n"))
			testutil.FatalIfErr(t, err)

			// Close the socket to signal to the socketStream to shut down.
			testutil.FatalIfErr(t, s.Close())

			wg.Wait()

			checkLineDiff()

			if v := <-ss.Lines(); v != nil {
				t.Errorf("expecting socketstream to be complete because socket closed")
			}
		}))
	}
}

func TestSocketStreamReadCompletedBecauseCancel(t *testing.T) {
	for _, scheme := range []string{"unix", "tcp"} {
		scheme := scheme
		t.Run(scheme, testutil.TimeoutTest(time.Second, func(t *testing.T) { //nolint:thelper
			var wg sync.WaitGroup

			var addr string
			switch scheme {
			case "unix":
				tmpDir := testutil.TestTempDir(t)
				addr = filepath.Join(tmpDir, "sock")
			case "tcp":
				addr = fmt.Sprintf("127.0.0.1:%d", testutil.FreePort(t))
			default:
				t.Fatalf("bad scheme %s", scheme)
			}

			ctx, cancel := context.WithCancel(context.Background())
			waker := waker.NewTestAlways()

			sockName := scheme + "://" + addr
			ss, err := logstream.New(ctx, &wg, waker, sockName, logstream.OneShotDisabled)
			testutil.FatalIfErr(t, err)

			s, err := net.Dial(scheme, addr)
			testutil.FatalIfErr(t, err)

			_, err = s.Write([]byte("1\n"))
			testutil.FatalIfErr(t, err)

			// Read from Lines to synchronise with the stream goroutine: this
			// blocks until the stream has consumed the data, proving it read
			// "1\n" from the socket before we cancel.
			expected := &logline.LogLine{Context: context.TODO(), Filename: sockName, Line: "1", Filenamehash: logline.GetHash(sockName)}
			received := <-ss.Lines()
			testutil.ExpectNoDiff(t, expected, received, testutil.IgnoreFields(logline.LogLine{}, "Context"))

			cancel() // This cancellation should cause the stream to shut down.
			wg.Wait()

			if v := <-ss.Lines(); v != nil {
				t.Errorf("expecting socketstream to be complete because cancel")
			}
		}))
	}
}
