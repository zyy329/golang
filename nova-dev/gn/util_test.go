package gn

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
)

//func TestMain(m *testing.M) {
//	framen.TestMainRun(m)
//}

type TS struct {
	a int
	b string
}

func TestArrFirst(t *testing.T) {
	mockey.PatchConvey("ArrFirst", t, func() {
		mockey.PatchConvey("int", func() {
			mockey.PatchConvey("zero", func() {
				var nArr []int
				ret := ArrFirst(nArr)
				convey.So(ret, convey.ShouldEqual, 0)
				ret = ArrFirstOrDefault(nArr, 0)
				convey.So(ret, convey.ShouldEqual, 0)
				ret = ArrFirstOrDefault(nArr, 1)
				convey.So(ret, convey.ShouldEqual, 1)
			})
			mockey.PatchConvey("normal", func() {
				nArr := []int{2, 3, 4, 5, 6}
				ret := ArrFirst(nArr)
				convey.So(ret, convey.ShouldEqual, 2)
				ret = ArrFirstOrDefault(nArr, 0)
				convey.So(ret, convey.ShouldEqual, 2)
				ret = ArrFirstOrDefault(nArr, 1)
				convey.So(ret, convey.ShouldEqual, 2)
			})
		})

		mockey.PatchConvey("struct", func() {
			var def = TS{-1, "-1"}
			mockey.PatchConvey("zero", func() {
				var tArr []TS
				ret2 := ArrFirst(tArr)
				convey.So(ret2.a, convey.ShouldEqual, 0)
				ret2 = ArrFirstOrDefault(tArr, def)
				convey.So(ret2.a, convey.ShouldEqual, -1)
			})
			mockey.PatchConvey("normal", func() {
				tArr := []TS{{2, "2"}, {3, "3"}}
				ret2 := ArrFirst(tArr)
				convey.So(ret2.a, convey.ShouldEqual, 2)
				ret2 = ArrFirstOrDefault(tArr, def)
				convey.So(ret2.a, convey.ShouldEqual, 2)
			})
		})

		mockey.PatchConvey("point int", func() {
			var def = -1
			mockey.PatchConvey("zero", func() {
				var nArr []*int
				ret := ArrFirst(nArr)
				convey.So(ret, convey.ShouldBeNil)
				ret = ArrFirstOrDefault(nArr, &def)
				convey.So(*ret, convey.ShouldEqual, def)
			})
			mockey.PatchConvey("normal", func() {
				nArr := []int{2, 3, 4, 5, 6}
				var nPArr []*int
				for _, n := range nArr {
					nPArr = append(nPArr, &n)
				}
				ret2 := ArrFirst(nPArr)
				convey.So(*ret2, convey.ShouldEqual, 2)
				ret2 = ArrFirstOrDefault(nPArr, &def)
				convey.So(*ret2, convey.ShouldEqual, 2)
			})
		})

		mockey.PatchConvey("point struct", func() {
			var def = TS{-1, "-1"}
			mockey.PatchConvey("zero", func() {
				var tArr []*TS
				ret := ArrFirst(tArr)
				convey.So(ret, convey.ShouldBeNil)
				ret = ArrFirstOrDefault(tArr, &def)
				convey.So(ret.a, convey.ShouldEqual, def.a)
			})
			mockey.PatchConvey("normal", func() {
				tArr := []*TS{{2, "2"}, {3, "3"}}
				ret2 := ArrFirst(tArr)
				convey.So(ret2.a, convey.ShouldEqual, 2)
				ret2 = ArrFirstOrDefault(tArr, &def)
				convey.So(ret2.a, convey.ShouldEqual, 2)
			})
		})
	})
}
