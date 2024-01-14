// Copyright 2024 Google Inc.  ll Rights Reserved.
// This file is available under the Apache license.

package mtail_test

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestExecMtail(t *testing.T) {
	mtailPath, ok := bazel.FindBinary("cmd/mtail", "mtail")
	if !ok {
		t.Fatal("`mtail` not found in runfiles")
	}
	cs := []string{
		"-progs",
		"../../examples",
		"-logs", "testdata/rsyncd.log",
		"-one_shot",
		"-one_shot_format=prometheus",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, mtailPath, cs...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Error(err)
	}
}
