#!/bin/bash

# Script to fix unused imports in the codebase
# This script will remove unused imports from Go files

echo "Starting unused import fixes..."

# Process specific files with known issues
files=(
  "./controllers/channelscontroller.go"
  "./controllers/emailtemplates.go"
  "./controllers/media_settingscontroller.go"
  "./controllers/menucontroller.go"
  "./controllers/templates.go"
)

for file in "${files[@]}"; do
  echo "Fixing unused imports in $file"
  
  # Check if fmt is imported but not used
  if grep -q "\"fmt\"" "$file" && ! grep -q "fmt\\." "$file"; then
    echo "  Removing unused fmt import from $file"
    # Remove the fmt import
    sed -i '' '/^[[:space:]]*"fmt"[[:space:]]*$/d' "$file"
    # Also handle fmt in import blocks
    sed -i '' '/import (/,/)/ { /[[:space:]]*"fmt"[[:space:]]*$/d }' "$file"
  fi
  
  # Check if logger is imported but not used
  if grep -q "\"spurt-cms/logger\"" "$file" && ! grep -q "logger\\." "$file"; then
    echo "  Removing unused logger import from $file"
    # Remove the logger import
    sed -i '' '/^[[:space:]]*"spurt-cms\/logger"[[:space:]]*$/d' "$file"
    # Also handle logger in import blocks
    sed -i '' '/import (/,/)/ { /[[:space:]]*"spurt-cms\/logger"[[:space:]]*$/d }' "$file"
  fi
done

echo "Unused import fixes complete."
echo "Please verify the changes and ensure all imports are correctly resolved."
