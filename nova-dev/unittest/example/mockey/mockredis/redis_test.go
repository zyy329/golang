package mockredis

import (
	"context"
	"testing"
	"time"

	"nova/framen"
	"nova/redisn"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
)

// 以redis 为例, 展示如何使用三方现成 mock库, 对redis操作进行注入
//(mysql 等通用数据库也有对应现成mock库)

// 测试启动, 做统一的初始化化工作, 同包下每个Test都会自动调用该函数
func TestMain(m *testing.M) {
	framen.TestMainRun(m)
}

func TestRedis(t *testing.T) {
	redisn.MockRedis(t)

	ctx := context.Background()
	client := redisn.NewClient("common")
	key := "s"
	t.Logf("cmd: %+v", client.Set(ctx, key, "2", 0))
	t.Logf("cmd: %+v", client.Get(ctx, key))
}

func TestRedisByConvey(t *testing.T) {
	// PatchConvey mockey 对 Convey的包装扩展
	mockey.PatchConvey("Redis Mock", t, func() {
		mr := redisn.MockRedis(t)
		ctx := context.Background()
		client := redisn.NewClient("common")

		// 测试其他单元测试用例 TestRedis() 的操作不会影响本次测试结果
		mockey.PatchConvey("TestRedis() not effect this", func() {
			key := "s"
			convey.So(client.Get(ctx, key).Val(), convey.ShouldEqual, "")
		})

		// list 测试
		mockey.PatchConvey("Test List 1", func() {
			key := "lst"
			client.LPush(ctx, key, "1")
			client.LPush(ctx, key, "2")
			client.LPush(ctx, key, "3")
			client.LPush(ctx, key, "4")

			convey.So(client.LPop(ctx, key).Val(), convey.ShouldEqual, "4")
			convey.So(client.LLen(ctx, key).Val(), convey.ShouldEqual, 3)
		})

		// 测试 Test List 1 的操作 不会跨 PatchConvey 影响测试结果
		// 相关说明可从参考: https://juejin.cn/post/7130533292310265886  执行顺序一节
		mockey.PatchConvey("Test List 1 not effect this", func() {
			key := "lst"
			convey.So(client.LLen(ctx, key).Val(), convey.ShouldEqual, 0)
		})

		// map 测试
		mockey.PatchConvey("Test Map", func() {
			key := "map"
			field := "f1"
			val := "fv2"
			client.HSet(ctx, key, field, val)
			convey.So(client.HGet(ctx, key, field).Val(), convey.ShouldEqual, val)
		})

		// 测试 TTL 过期
		mockey.PatchConvey("Test TTL", func() {
			key := "s"
			val := "2"
			client.Set(ctx, key, "2", 60*time.Minute)
			convey.So(client.Get(ctx, key).Val(), convey.ShouldEqual, val)
			convey.So(client.TTL(ctx, key).Val(), convey.ShouldBeBetween, 59*time.Minute, 61*time.Minute)

			// 快速过时间
			mr.FastForward(10 * time.Minute)
			convey.So(client.Get(ctx, key).Val(), convey.ShouldEqual, val)
			convey.So(client.TTL(ctx, key).Val(), convey.ShouldBeBetween, 49*time.Minute, 51*time.Minute)

			// 快速过期
			mr.FastForward(time.Hour)
			convey.So(client.Get(ctx, key).Val(), convey.ShouldEqual, "")
		})
	})
}
