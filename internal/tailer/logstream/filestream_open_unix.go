// Copyright 2026 The mtail Authors. All Rights Reserved.
// This file is available under the Apache license.

//go:build !windows

package logstream

import (
	"os"
)

func openFile(pathname string) (*os.File, error) {
	return os.OpenFile(pathname, os.O_RDONLY, 0o600)
}
