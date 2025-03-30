package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// This is an example script that could be used to automatically replace fmt.Print statements
// with logger calls throughout the codebase.

// ReplacementRule defines a pattern to match and its replacement
type ReplacementRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// CreateReplacementRules creates the rules for replacing fmt.Print statements
func CreateReplacementRules() []ReplacementRule {
	return []ReplacementRule{
		// Replace logger.Error("Error occurred", logger.WithError(err)) with logger.Error("Operation failed", logger.WithError(err))
		{
			Pattern:     regexp.MustCompile(`fmt\.Println\(err\)`),
			Replacement: `logger.Error("Operation failed", logger.WithError(err))`,
		},
		// Replace logger.Info(fmt.Sprintf("Error: %v", err)) with logger.Error("Error occurred", logger.WithError(err))
		{
			Pattern:     regexp.MustCompile(`fmt\.Printf\("Error: %v", err\)`),
			Replacement: `logger.Error("Error occurred", logger.WithError(err))`,
		},
		// Replace logger.Info("Warning message") with logger.Warn("Warning message")
		{
			Pattern:     regexp.MustCompile(`fmt\.Println\("([^"]+)"\)`),
			Replacement: `logger.Info("$1")`,
		},
		// Replace logger.Info(fmt.Sprintf("Debug: %v\n", val)) with logger.Debug("Debug information", logger.WithField("value", val))
		{
			Pattern:     regexp.MustCompile(`fmt\.Printf\("Debug: %v\\n", ([a-zA-Z0-9_]+)\)`),
			Replacement: `logger.Debug("Debug information", logger.WithField("value", $1))`,
		},
		// Replace logger.Info(fmt.Sprintf("%s: %v\n", name, val)) with logger.Info(name, logger.WithField("value", val))
		{
			Pattern:     regexp.MustCompile(`fmt\.Printf\("%s: %v\\n", ([a-zA-Z0-9_]+), ([a-zA-Z0-9_]+)\)`),
			Replacement: `logger.Info($1, logger.WithField("value", $2))`,
		},
	}
}

// ProcessFile processes a single file and applies the replacement rules
func ProcessFile(filePath string, rules []ReplacementRule) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	// Convert to string for easier manipulation
	fileContent := string(content)
	
	// Check if the file imports fmt
	if !strings.Contains(fileContent, `import "fmt"`) && !strings.Contains(fileContent, `"fmt"`) {
		return nil // Skip files that don't import fmt
	}
	
	// Check if the file already imports logger
	hasLoggerImport := strings.Contains(fileContent, `import "spurt-cms/logger"`) || 
		strings.Contains(fileContent, `"spurt-cms/logger"`)
	
	// Apply each replacement rule
	modified := false
	for _, rule := range rules {
		if rule.Pattern.MatchString(fileContent) {
			fileContent = rule.Pattern.ReplaceAllString(fileContent, rule.Replacement)
			modified = true
		}
	}
	
	// If modifications were made but logger is not imported, add the import
	if modified && !hasLoggerImport {
		// Find the import block
		importRegex := regexp.MustCompile(`import\s*\(([\s\S]*?)\)`)
		if importRegex.MatchString(fileContent) {
			// Add logger to the import block
			fileContent = importRegex.ReplaceAllString(fileContent, 
				`import (
	"spurt-cms/logger"$1)`)
		} else {
			// Handle single line imports
			singleImportRegex := regexp.MustCompile(`import\s+"([^"]+)"`)
			if singleImportRegex.MatchString(fileContent) {
				fileContent = singleImportRegex.ReplaceAllString(fileContent, 
					`import (
	"spurt-cms/logger"
	"$1"
)`)
			}
		}
	}
	
	// Write the modified content back to the file if changes were made
	if modified {
		return os.WriteFile(filePath, []byte(fileContent), 0644)
	}
	
	return nil
}

// ProcessDirectory processes all Go files in a directory and its subdirectories
func ProcessDirectory(dirPath string, rules []ReplacementRule) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip vendor and node_modules directories
		if info.IsDir() && (info.Name() == "vendor" || info.Name() == "node_modules") {
			return filepath.SkipDir
		}
		
		// Process only Go files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			if err := ProcessFile(path, rules); err != nil {
				logger.Info(fmt.Sprintf("Error processing file %s: %v\n", path, err))
			}
		}
		
		return nil
	})
}

// Example of how to use the migration script
func ExampleUsage() {
	// Create replacement rules
	rules := CreateReplacementRules()
	
	// Process the entire project
	err := ProcessDirectory(".", rules)
	if err != nil {
		logger.Info(fmt.Sprintf("Error processing directory: %v\n", err))
	}
}
