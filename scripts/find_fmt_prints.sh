#!/bin/bash

# Script to find all fmt.Println statements in Go files
# This script will output the files containing fmt.Println statements

echo "Finding all Go files with fmt.Println statements..."

# Find all Go files containing fmt.Println statements
grep -l "fmt.Println" $(find . -name "*.go" | grep -v "/vendor/" | grep -v "/scripts/") > /tmp/fmt_files.txt

# Count the number of files
file_count=$(wc -l < /tmp/fmt_files.txt)

echo "Found $file_count files containing fmt.Println statements."
echo "Files list saved to /tmp/fmt_files.txt"

# Sample some of the fmt.Println statements to help determine appropriate log levels
echo "Sampling some fmt.Println statements for analysis:"
for file in $(head -5 /tmp/fmt_files.txt); do
  echo "File: $file"
  grep -n "fmt.Println" "$file" | head -3
  echo "---"
done

echo "Run ./scripts/replace_fmt_prints.sh to replace these with logger calls."
