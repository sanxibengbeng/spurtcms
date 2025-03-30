#!/bin/bash

# Script to perform a final check on the logger replacements
# This script will check for any remaining fmt.Println statements and verify logger imports

echo "Starting final verification..."

# Check for any remaining fmt.Println statements
remaining_fmt=$(grep -r "fmt\.Println" --include="*.go" . | grep -v "/vendor/" | grep -v "/scripts/" | wc -l)

echo "Remaining fmt.Println statements: $remaining_fmt"

if [ "$remaining_fmt" -gt 0 ]; then
  echo "Files with remaining fmt.Println statements:"
  grep -l "fmt\.Println" --include="*.go" . | grep -v "/vendor/" | grep -v "/scripts/"
fi

# Check for logger imports
missing_logger=$(grep -l "logger\." --include="*.go" . | grep -v "/vendor/" | grep -v "/scripts/" | xargs grep -L "\"github.com/spurtcms/spurtcms/logger\"" | xargs grep -L "\"./logger\"" | xargs grep -L "\".*/logger\"" | wc -l)

echo "Files using logger without proper import: $missing_logger"

if [ "$missing_logger" -gt 0 ]; then
  echo "Files missing logger import:"
  grep -l "logger\." --include="*.go" . | grep -v "/vendor/" | grep -v "/scripts/" | xargs grep -L "\"github.com/spurtcms/spurtcms/logger\"" | xargs grep -L "\"./logger\"" | xargs grep -L "\".*/logger\""
fi

echo "Final verification complete."
