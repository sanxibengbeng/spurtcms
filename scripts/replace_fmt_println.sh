#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始替换 fmt.Println 为 logger 调用...${NC}"

# 统计原始的 fmt.Println 数量
ORIGINAL_COUNT=$(grep -r "fmt.Println" --include="*.go" . | grep -v "vendor" | grep -v "scripts" | wc -l)
echo -e "${YELLOW}原始 fmt.Println 调用数量: $ORIGINAL_COUNT${NC}"

# 替换错误日志
echo "1. 替换错误日志模式..."
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Println(err)/logger.Error("Operation failed", logger.WithError(err))/g' {} \;

# 替换带有 Error 的打印
echo "2. 替换错误打印模式..."
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Printf("Error \(.*\): %v", \(.*\))/logger.Error("\1", logger.WithError(\2))/g' {} \;
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Printf("Error \(.*\): %s", \(.*\))/logger.Error("\1", logger.WithError(\2))/g' {} \;

# 替换普通信息日志
echo "3. 替换信息日志模式..."
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Println("\(.*\)")/logger.Info("\1")/g' {} \;

# 替换带有变量的打印
echo "4. 替换调试日志模式..."
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Printf("\(.*\): %v", \(.*\))/logger.Debug("\1", logger.WithField("value", \2))/g' {} \;
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Printf("\(.*\): %s", \(.*\))/logger.Debug("\1", logger.WithField("value", \2))/g' {} \;

# 替换简单变量打印
echo "5. 替换简单变量打印模式..."
find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./scripts/*" -exec sed -i '' 's/fmt.Println(\(.*\))/logger.Debug("Debug value", logger.WithField("value", \1))/g' {} \;

# 添加 logger 导入
echo "6. 添加 logger 导入..."
FILES_WITH_CHANGES=$(grep -l "logger\." --include="*.go" . | grep -v "vendor" | grep -v "scripts")

for file in $FILES_WITH_CHANGES; do
  # 检查文件是否已经导入了 logger 包
  if ! grep -q "\"spurt-cms/logger\"" "$file"; then
    # 检查是否有 import 块
    if grep -q "import (" "$file"; then
      # 在 import 块中添加 logger 导入
      sed -i '' '/import (/a\'$'\n''\t"spurt-cms/logger"' "$file"
    else
      # 检查是否有单行导入
      if grep -q "import " "$file"; then
        # 在第一个 import 后添加 logger 导入
        sed -i '' '/import /a\'$'\n''import "spurt-cms/logger"' "$file"
      fi
    fi
  fi
done

# 统计替换后的 fmt.Println 数量
REMAINING_COUNT=$(grep -r "fmt.Println" --include="*.go" . | grep -v "vendor" | grep -v "scripts" | wc -l)
echo -e "${GREEN}替换完成!${NC}"
echo -e "${YELLOW}替换前 fmt.Println 调用数量: $ORIGINAL_COUNT${NC}"
echo -e "${YELLOW}替换后剩余 fmt.Println 调用数量: $REMAINING_COUNT${NC}"
echo -e "${GREEN}成功替换: $(($ORIGINAL_COUNT - $REMAINING_COUNT)) 处${NC}"

if [ $REMAINING_COUNT -gt 0 ]; then
  echo -e "${RED}警告: 仍有 $REMAINING_COUNT 处 fmt.Println 调用未被替换，可能需要手动处理。${NC}"
  echo "以下是剩余的 fmt.Println 调用:"
  grep -r "fmt.Println" --include="*.go" . | grep -v "vendor" | grep -v "scripts" | head -10
  if [ $REMAINING_COUNT -gt 10 ]; then
    echo "... 以及更多 $(($REMAINING_COUNT - 10)) 处调用"
  fi
fi
