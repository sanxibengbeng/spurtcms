#!/bin/bash

# Script to fix logger import paths in the codebase
# This script will replace the full GitHub path with relative imports

echo "Starting import path fixes..."

# Find all files with the problematic import
find . -name "*.go" -type f -exec grep -l "\"github.com/spurtcms/spurtcms/logger\"" {} \; > /tmp/import_files.txt

# Process each file
total_files=$(wc -l < /tmp/import_files.txt)
processed=0

while IFS= read -r file; do
  processed=$((processed + 1))
  echo "[$processed/$total_files] Fixing imports in $file"
  
  # Replace the full GitHub path with a relative path
  # First check if the file is in the root directory
  if [[ $(dirname "$file") == "." ]]; then
    # For files in the root directory
    sed -i '' 's|"github.com/spurtcms/spurtcms/logger"|"./logger"|g' "$file"
  else
    # For files in subdirectories, calculate the relative path
    dir_depth=$(echo "$file" | tr -cd '/' | wc -c)
    rel_path=""
    
    # Build the relative path based on directory depth
    for ((i=1; i<dir_depth; i++)); do
      rel_path="../$rel_path"
    done
    
    sed -i '' "s|\"github.com/spurtcms/spurtcms/logger\"|\"${rel_path}logger\"|g" "$file"
  fi
  
  echo "Completed fixing $file"
done < /tmp/import_files.txt

echo "Import path fixes complete. Processed $total_files files."
echo "Please verify the changes and ensure all imports are correctly resolved."
