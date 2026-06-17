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

// Time describes a timestamp value at a given timestamp.
type Time struct {
	BaseDatum
	mu    sync.RWMutex
	Value time.Time
}

// Set sets the value of the Time to the value at timestamp.
func (d *Time) Set(value time.Time, timestamp time.Time) {
	d.mu.Lock()
	d.Value = value
	d.stamp(timestamp)
	d.mu.Unlock()
}

// Get returns the value of the Time.
func (d *Time) Get() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Value
}

// ValueString returns the value of the Time as a string.
func (d *Time) ValueString() string {
	t := d.Get()
	return fmt.Sprintf("%d.%09d", t.Unix(), t.Nanosecond())
}

// MarshalJSON returns a JSON encoding of the Time.
func (d *Time) MarshalJSON() ([]byte, error) {
	t := d.Get()
	j := struct {
		Value time.Time `json:"Value"`
		Time  int64     `json:"Time"`
	}{t, atomic.LoadInt64(&d.Time)}
	return json.Marshal(j)
}
