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
	Filenamehash uint16 // stored for efficiency key lookup
	Line         string // The text of the log line itself up to the newline.
}

// New creates a new LogLine object.
func New(ctx context.Context, filename string, line string) *LogLine {
	hash := sha256.Sum256([]byte(filename))
	return &LogLine{ctx, filename, binary.BigEndian.Uint16(hash[:2]), line}
}
