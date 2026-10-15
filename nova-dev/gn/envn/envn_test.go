package envn

import (
	"testing"
)

func TestEnvName(t *testing.T) {
	t.Log(EnvName("player"))

	t.Log(EnvName("player", Dev))

	t.Log(EnvName("player", ""))

	Env = ""
	t.Log(EnvName("player"))

	Env = Dev
	t.Log(EnvName("player"))

	Env = Prev
	Id = "1"
	t.Log(EnvName("player"))
}

func TestParsePod(t *testing.T) {
	// 1. 定义测试用例结构体
	tests := []struct {
		name     string // 用例描述
		input    string // 输入
		wantType string // 期望的 Type
		wantId   int    // 期望的 ID
	}{
		// 正常情况 - A组
		{
			name:     "Standard format type A",
			input:    "5117-a-2",
			wantType: "a",
			wantId:   2,
		},
		// 正常情况 - B组
		{
			name:     "Standard format type B",
			input:    "game-server-b-10",
			wantType: "b",
			wantId:   10,
		},
		// 只有 ID，没有 a/b
		{
			name:     "ID only",
			input:    "xxx-3",
			wantType: "",
			wantId:   3,
		},
		// 只有 ID，且中间段不是 a/b
		{
			name:     "ID with other suffix",
			input:    "server-c-5",
			wantType: "", // c 不是 a 或 b
			wantId:   5,
		},
		// 无效格式 - 无分隔符
		{
			name:     "No separator",
			input:    "xxxx",
			wantType: "",
			wantId:   0,
		},
		// 无效格式 - 结尾非数字
		{
			name:     "Invalid ID",
			input:    "5117-a-test",
			wantType: "",
			wantId:   0,
		},
		// 边界情况 - 只有数字
		{
			name:     "Only Number",
			input:    "100",
			wantType: "",
			wantId:   0, // 因为没有横杠分割，视作整体字符串，不解析ID
		},
	}

	// 2. 循环执行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotId := parseNode(tt.input)

			// 3. 断言：检查实际结果与期望值是否一致
			if gotType != tt.wantType {
				t.Errorf("ParsePodName() type = %v, want %v", gotType, tt.wantType)
			}
			if gotId != tt.wantId {
				t.Errorf("ParsePodName() id = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
