#!/bin/bash

# Script to replace fmt.Println statements with appropriate logger calls
# This script will process all files identified by find_fmt_prints.sh

echo "Starting replacement of fmt.Println statements with logger calls..."

# Check if the file list exists
if [ ! -f /tmp/fmt_files.txt ]; then
  echo "File list not found. Please run ./scripts/find_fmt_prints.sh first."
  exit 1
fi

# Make sure the logger package is imported
add_logger_import() {
  local file=$1
  
  # Check if logger is already imported
  if ! grep -q "\"github.com/spurtcms/spurtcms/logger\"" "$file" && ! grep -q "\"./logger\"" "$file" && ! grep -q "\".*/logger\"" "$file"; then
    # Find the import block and add logger import
    if grep -q "^import (" "$file"; then
      # Multi-line import block
      sed -i '' '/^import (/,/^)/ s/)$/\t"github.com\/spurtcms\/spurtcms\/logger"\n)/' "$file"
    elif grep -q "^import \"" "$file"; then
      # Single import, convert to block and add logger
      sed -i '' 's/^import "/import (\n\t"/g' "$file"
      sed -i '' '/^import (/,/^$/ s/^$/\t"github.com\/spurtcms\/spurtcms\/logger"\n)/' "$file"
    else
      # No imports, add import block
      sed -i '' '/^package/a\\nimport (\n\t"github.com/spurtcms/spurtcms/logger"\n)' "$file"
    fi
  fi
}

# Process each file
total_files=$(wc -l < /tmp/fmt_files.txt)
processed=0

while IFS= read -r file; do
  processed=$((processed + 1))
  echo "[$processed/$total_files] Processing $file"
  
  # Add logger import if needed
  add_logger_import "$file"
  
  # Replace different types of fmt.Println statements
  
  # 1. Simple fmt.Println("message")
  sed -i '' 's/fmt\.Println("\([^"]*\)")/logger.Info("\1")/g' "$file"
  
  # 2. fmt.Println with variables
  sed -i '' 's/fmt\.Println(\([^)]*\))/logger.Info(fmt.Sprintf("%v", \1))/g' "$file"
  
  # 3. fmt.Printf statements - convert to appropriate logger calls
  sed -i '' 's/fmt\.Printf("\([^"]*\)", \([^)]*\))/logger.Info(fmt.Sprintf("\1", \2))/g' "$file"
  
  # 4. Error-related prints - use Error level
  sed -i '' 's/fmt\.Println("Error:/logger.Error("Error:/g' "$file"
  sed -i '' 's/fmt\.Println("Failed/logger.Error("Failed/g' "$file"
  sed -i '' 's/fmt\.Println(".*error.*")/logger.Error("\0")/g' "$file"
  
  # 5. Warning-related prints - use Warn level
  sed -i '' 's/fmt\.Println("Warning:/logger.Warn("Warning:/g' "$file"
  sed -i '' 's/fmt\.Println(".*warning.*")/logger.Warn("\0")/g' "$file"
  
  # 6. Debug-related prints - use Debug level
  sed -i '' 's/fmt\.Println("Debug:/logger.Debug("Debug:/g' "$file"
  
  echo "Completed processing $file"
done < /tmp/fmt_files.txt

echo "Replacement complete. Processed $total_files files."
echo "Please review the changes and make any necessary adjustments."
echo "You may need to manually fix some complex fmt.Println statements."
