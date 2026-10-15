#!/usr/bin/env bash

##################################################
# 单元测试 和 测试覆盖率检测脚本
# 常规使用时由 pre-commit 运行是调用
# 此时是无参调用, 将驱动有修改的包(已经git add 的包)运行单元测试，并检测测试覆盖率
# 需要时可通过 ./run-go-unit-test-coverage.sh all 命令执行所有包的单元测试并检测全文覆盖率
##################################################

# 定义需要跳过的包名（根据实际情况修改，支持完整路径或部分匹配）
SKIP_PACKAGES=(
#  "data-center/api"
)

test_fail() {
  echo "❌ unit tests failed"
  exit 1
}

TEST_ALL=$1
# 设置最低覆盖率阈值
THRESHOLD=80.0
PACKAGES=""
BASE_PATH=$(go list -m)

if [ "${TEST_ALL}" == "all" ]; then
  # 附带 all 参数, 测试所有包
  PACKAGES=$(go list ./... | grep -v /vendor/) || fail
else
  # 只测试本次已经add待提交的包
  # 获取本次提交的变更文件，筛选出 .go 文件
  changed_files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

  # 如果没有 Go 文件变更，直接退出
  if [ -z "$changed_files" ]; then
      echo "✅ No Go files changed, skipping tests."
      exit 0
  fi

  # 提取变更文件对应的包路径,去重,添加 BASE_PATH 前缀(和go list结果格式保持一致)
  PACKAGES=$(echo "$changed_files" | xargs -n1 dirname | sort -u | sed "s|^|${BASE_PATH}/|" | sed "s|^${BASE_PATH}/.$|${BASE_PATH}|")
fi

# 对每个变更的包执行 Go 单元测试
for pkg in $PACKAGES; do
  # 检查当前包是否在跳过列表中
  for skip in "${SKIP_PACKAGES[@]}"; do
    # 若包路径包含跳过的名字（支持部分匹配），则跳过
    if [[ "$pkg" == *"$skip"* ]]; then
      echo "skip pkg: $pkg"
      continue 2  # 跳出内层循环，继续外层循环的下一个包
    fi
  done

  echo "test pkg: $pkg"

  # -gcflags="all=-l -N" 禁用内联和编译优化,支持 mockey 注入
  go test -gcflags="all=-l -N" -coverprofile=coverage.out -timeout 30s -short -v "${pkg}" || test_fail

  # 覆盖率检查忽略根目录; 因为根目录主要是 main.go 没有覆盖率检查的意义
  if [ "$pkg" = "${BASE_PATH}" ]; then
    continue
  fi

  # 提取覆盖率
  COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
  # 检查是否低于阈值
  if (( $(echo "${COVERAGE} < ${THRESHOLD}" | bc -l) )); then
    echo "❌ Code coverage (${COVERAGE}%) is below the required threshold (${THRESHOLD}%). Increase test coverage before committing."

    # 打印不合格的包覆盖率细节
    echo "---coverage detail--------------------------"
    echo "---pkg: ${pkg}"
    go tool cover -func=coverage.out
    echo "--------------------------------------------"
    exit 1
  fi
done

# 单元测试通过, 删除临时的 coverage.out
rm coverage.out || true

#####################覆盖率相关命令说明#####################
## 输出覆盖率基础信息
#go test -cover
## 生成覆盖率报告
#go test -coverprofile=coverage.out
## 查看每个函数的覆盖率信息
#go tool cover -func=coverage.out
## 生成一个可视化的 HTML 报告
#go tool cover -html=coverage.out
## 统计某个包及其子包的整体覆盖率
#go test -coverpkg=./... -coverprofile=coverage.out ./...
#########################################################