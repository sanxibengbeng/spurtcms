package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 替换规则
type ReplaceRule struct {
	Pattern     string
	Replacement string
}

func main() {
	// 定义替换规则
	rules := []ReplaceRule{
		// 错误日志
		{Pattern: `fmt.Println(err)`, Replacement: `logger.Error("Operation failed", logger.WithError(err))`},
		{Pattern: `fmt.Printf("Error`, Replacement: `logger.Error(`},
		
		// 信息日志
		{Pattern: `fmt.Println("`, Replacement: `logger.Info("`},
		{Pattern: `fmt.Printf("`, Replacement: `logger.Info(`},
		
		// 调试日志 - 变量打印
		{Pattern: `fmt.Println(`, Replacement: `logger.Debug("Debug value", logger.WithField("value", `},
	}

	// 获取项目根目录
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	// 遍历所有 .go 文件
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过 vendor 目录
		if strings.Contains(path, "vendor") || strings.Contains(path, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理 .go 文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			processFile(path, rules)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", rootDir, err)
	}
}

func processFile(filePath string, rules []ReplaceRule) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	fileContent := string(content)
	originalContent := fileContent

	// 检查文件是否导入了 fmt 包
	if !strings.Contains(fileContent, `"fmt"`) {
		return
	}

	// 检查文件是否包含 fmt.Println
	if !strings.Contains(fileContent, "fmt.Println") && !strings.Contains(fileContent, "fmt.Printf") {
		return
	}

	// 解析 Go 源代码
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, fileContent, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", filePath, err)
		return
	}

	// 检查是否已经导入了 logger 包
	hasLoggerImport := false
	for _, imp := range f.Imports {
		if imp.Path != nil && strings.Contains(imp.Path.Value, "logger") {
			hasLoggerImport = true
			break
		}
	}

	// 应用替换规则
	for _, rule := range rules {
		fileContent = strings.Replace(fileContent, rule.Pattern, rule.Replacement, -1)
	}

	// 如果文件内容有变化，并且需要添加 logger 导入
	if fileContent != originalContent && !hasLoggerImport {
		// 添加 logger 导入
		importIndex := strings.Index(fileContent, "import (")
		if importIndex != -1 {
			// 找到 import 块的结束位置
			closingIndex := strings.Index(fileContent[importIndex:], ")")
			if closingIndex != -1 {
				insertPos := importIndex + closingIndex
				fileContent = fileContent[:insertPos] + "\n\t\"spurt-cms/logger\"" + fileContent[insertPos:]
			}
		} else {
			// 单行导入
			importIndex = strings.Index(fileContent, "import ")
			if importIndex != -1 {
				newlineIndex := strings.Index(fileContent[importIndex:], "\n")
				if newlineIndex != -1 {
					insertPos := importIndex + newlineIndex + 1
					fileContent = fileContent[:insertPos] + "import \"spurt-cms/logger\"\n" + fileContent[insertPos:]
				}
			}
		}

		// 写回文件
		err = ioutil.WriteFile(filePath, []byte(fileContent), 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			return
		}

		fmt.Printf("Updated file: %s\n", filePath)
	}
}
