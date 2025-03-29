#!/bin/bash

set -e

# Stash any changes not being committed
git stash -q --keep-index

# Format all Go files and capture changes
gofmt -w -s .
GOFMT_STATUS=$?

# Capture formatted files
CHANGED_GO_FILES=$(git diff --name-only -- '*.go')

# Add the formatted files back to the commit
if [ -n "$CHANGED_GO_FILES" ]; then
  git add $CHANGED_GO_FILES
  echo "Formatted and added the following Go files:"
  echo "$CHANGED_GO_FILES"
fi

# Run golangci-lint if available
if command -v golangci-lint >/dev/null 2>&1; then
  golangci-lint run
  LINT_STATUS=$?
else
  LINT_STATUS=0
fi

# Restore stashed changes
git stash pop -q 2>/dev/null || true

# Exit with error code if any of the commands failed
[ $GOFMT_STATUS -eq 0 ] && [ $LINT_STATUS -eq 0 ]