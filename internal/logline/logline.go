// Copyright 2017 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package logline

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
)

// LogLine contains all the information about a line just read from a log.
type LogLine struct {
	Context context.Context

	Filename     string // The log filename that this line was read from
	Filenamehash uint32 // stored for efficient key lookup
	Line         string // The text of the log line itself up to the newline.
}

// New creates a new LogLine object.
func New(ctx context.Context, filename string, filenamehash uint32, line string) *LogLine {
	return &LogLine{ctx, filename, filenamehash, line}
}

// External as unit tests need it
func GetHash(filename string) uint32 {
	hash := sha256.Sum256([]byte(filename))
	return binary.BigEndian.Uint32(hash[:8])
}
