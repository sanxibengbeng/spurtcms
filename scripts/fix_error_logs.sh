#!/bin/bash

# Script to fix error logs in the codebase
# This script will convert simple error logs to use the WithError helper

echo "Starting error log refinement..."

# Process files with logger calls
find . -name "*.go" -type f -exec grep -l "logger\.Info" {} \; > /tmp/error_files.txt

# Process each file
total_files=$(wc -l < /tmp/error_files.txt)
processed=0

while IFS= read -r file; do
  processed=$((processed + 1))
  echo "[$processed/$total_files] Fixing error logs in $file"
  
  # Find lines with err variable being logged
  grep -n "logger\.Info.*err" "$file" | while read -r line; do
    line_num=$(echo "$line" | cut -d':' -f1)
    
    # Extract the line content
    line_content=$(sed -n "${line_num}p" "$file")
    
    # Replace with Error level and WithError helper
    if [[ "$line_content" =~ logger\.Info\(fmt\.Sprintf\("%v",\ err\)\) ]]; then
      sed -i '' "${line_num}s/logger\.Info(fmt\.Sprintf(\"%v\", err))/logger.Error(\"Error occurred\", logger.WithError(err))/g" "$file"
    elif [[ "$line_content" =~ logger\.Info\(fmt\.Sprintf\(".*",\ err.*\)\) ]]; then
      # Extract the message part
      message=$(echo "$line_content" | sed -n 's/.*logger\.Info(fmt\.Sprintf("\(.*\)", err.*/\1/p')
      sed -i '' "${line_num}s/logger\.Info(fmt\.Sprintf(\"$message\", err))/logger.Error(\"$message\", logger.WithError(err))/g" "$file"
    fi
  done
  
  echo "Completed fixing $file"
done < /tmp/error_files.txt

echo "Error log refinement complete. Processed $total_files files."
echo "Please review the changes and make any necessary manual adjustments."
