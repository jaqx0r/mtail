// Copyright 2026 The mtail Authors. All Rights Reserved.
// This file is available under the Apache license.

//go:build windows

package logstream

import (
	"os"
	"syscall"
)

func openFile(pathname string) (*os.File, error) {
	pathp, err := syscall.UTF16PtrFromString(pathname)
	if err != nil {
		return nil, &os.PathError{Op: "utf16", Path: pathname, Err: err}
	}
	handle, err := syscall.CreateFile(
		pathp,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: pathname, Err: err}
	}
	return os.NewFile(uintptr(handle), pathname), nil
}
