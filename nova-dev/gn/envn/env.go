package envn

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvType 环境类型
type EnvType = string

//goland:noinspection GoUnusedConst
const (
	Local EnvType = "local" // 本地环境
	Dev   EnvType = "dev"   // 开发环境
	Test  EnvType = "test"  // 测试环境
	Prev  EnvType = "prev"  // 预发布环境
	Prod  EnvType = "prod"  // 线上环境
)

const (
	podA string = "a"
	podB string = "b"
)

// 环境变量集中定义
//
//goland:noinspection GoUnusedGlobalVariable
var (
	Env = func() string {
		env := os.Getenv("API_ENV")
		if env == "" {
			env = Local
		}
		return env
	}()
	Id          = os.Getenv("API_ID") // 整套环境的总编号，区分多套服务器
	HostName, _ = os.Hostname()
	PodId       int
	IsPodB      bool
	// MockSign mock生效开关, 默认情况下, 只有 Local 环境下才可以 Mock, 其他环境有特殊需要可单独改变该变量
	MockSign = Env == Local
)

func init() {
	var pobAb string
	pobAb, PodId = parseNode(HostName)
	IsPodB = pobAb == podB
}

// EnvName 不同环境需要区分表名、库名等, dev_1_player
// 本地开发环境可以指定env前缀
// prod环境无前缀
func EnvName(name string, s ...EnvType) string {
	return fmt.Sprintf("%s%s%s",
		func() string {
			if Env == Prod {
				return ""
			}
			if env := func() EnvType {
				e := Env
				if e == Local && len(s) > 0 {
					e = s[0]
				}
				return e
			}(); env == "" {
				return ""
			} else {
				return fmt.Sprintf("%s_", env)
			}
		}(),
		func() string {
			if Env == Prod {
				return ""
			}
			if Id == "" || Id == "0" {
				return ""
			}
			return fmt.Sprintf("%s_", Id)
		}(),
		name,
	)
}

func parseNode(hostname string) (podAb string, podId int) {
	parts := strings.Split(hostname, "-")
	length := len(parts)

	// 1. 防御性检查：如果没有横杠，或者长度不够切分，直接返回默认值
	// 这一步保证了后续取 parts[length-2] 不会越界 panic
	if length < 2 {
		return
	}

	// 2. 解析最后一段：PodId
	lastPart := parts[length-1]
	id, err := strconv.Atoi(lastPart)
	if err != nil {
		// 结尾不是数字，视为无 ID，也就无所谓 Type 了
		return
	}
	podId = id

	// 3. 解析倒数第二段：PodType
	// 既然上面已经校验过 length >= 2，这里直接取下标安全
	candidate := parts[length-2]
	if candidate == podA || candidate == podB {
		podAb = candidate
	}
	return
}
