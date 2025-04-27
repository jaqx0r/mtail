#!/usr/bin/env bash

echo STABLE_GIT_BRANCH $(git --no-pager rev-parse --abbrev-ref HEAD)
echo STABLE_GIT_VERSION $(git --no-pager describe --tags --always --dirty)
echo STABLE_GIT_REVISION $(git --no-pager rev-parse HEAD)
