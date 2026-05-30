// Copyright 2026 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package datum

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Duration describes a duration value at a given timestamp.
type Duration struct {
	BaseDatum
	mu    sync.RWMutex
	Value time.Duration
}

// Set sets the value of the Duration to the value at timestamp.
func (d *Duration) Set(value time.Duration, timestamp time.Time) {
	d.mu.Lock()
	d.Value = value
	d.stamp(timestamp)
	d.mu.Unlock()
}

// Get returns the value of the Duration.
func (d *Duration) Get() time.Duration {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Value
}

// ValueString returns the value of the Duration as a string.
func (d *Duration) ValueString() string {
	return fmt.Sprintf("%g", d.Get().Seconds())
}

// MarshalJSON returns a JSON encoding of the Duration.
func (d *Duration) MarshalJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	j := struct {
		Value int64
		Time  int64
	}{int64(d.Value), atomic.LoadInt64(&d.Time)}
	return json.Marshal(j)
}
