#!/bin/bash

# Script to improve the log level assignments in the codebase
# This script will refine the automatic replacements with more appropriate log levels

echo "Starting refinement of logger levels..."

# Process files with logger calls
find . -name "*.go" -type f -exec grep -l "logger\.Info" {} \; > /tmp/logger_files.txt

# Process each file
total_files=$(wc -l < /tmp/logger_files.txt)
processed=0

while IFS= read -r file; do
  processed=$((processed + 1))
  echo "[$processed/$total_files] Refining log levels in $file"
  
  # Convert error-related logs to Error level
  sed -i '' 's/logger\.Info(fmt\.Sprintf("%v", err))/logger.Error("Error occurred", logger.WithError(err))/g' "$file"
  sed -i '' 's/logger\.Info(fmt\.Sprintf(".*error.*", .*)/logger.Error(fmt.Sprintf("\1", \2))/g' "$file"
  sed -i '' 's/logger\.Info(fmt\.Sprintf(".*failed.*", .*)/logger.Error(fmt.Sprintf("\1", \2))/g' "$file"
  sed -i '' 's/logger\.Info(fmt\.Sprintf(".*invalid.*", .*)/logger.Warn(fmt.Sprintf("\1", \2))/g' "$file"
  
  # Convert debug-related logs to Debug level
  sed -i '' 's/logger\.Info("Debug:/logger.Debug("Debug:/g' "$file"
  sed -i '' 's/logger\.Info("Starting /logger.Debug("Starting /g' "$file"
  sed -i '' 's/logger\.Info("Processing /logger.Debug("Processing /g' "$file"
  sed -i '' 's/logger\.Info("Loaded /logger.Debug("Loaded /g' "$file"
  
  # Convert warning-related logs to Warn level
  sed -i '' 's/logger\.Info("Warning:/logger.Warn("Warning:/g' "$file"
  sed -i '' 's/logger\.Info("Caution:/logger.Warn("Caution:/g' "$file"
  
  echo "Completed refining $file"
done < /tmp/logger_files.txt

echo "Refinement complete. Processed $total_files files."
echo "Please review the changes and make any necessary manual adjustments."
