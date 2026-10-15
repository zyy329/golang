package converymock

//func funcA() int {
//	return 0
//}
//
// 测试多层嵌套关系下 mock 的影响
// 结果, 不能这么使用, 会出错, 多次运行结果不一致
//func TestMockWithConveyS(t *testing.T) {
//	mockey.PatchConvey("MockWithConveyS", t, func() {
//		convey.So(funcA(), convey.ShouldEqual, 0)
//
//		mockey.PatchConvey("mock 1", func() {
//			mockey.Mock(funcA).Return(1).Build()
//			convey.So(funcA(), convey.ShouldEqual, 1)
//
//			mockey.PatchConvey("mock 11", func() {
//				mockey.Mock(funcA).Return(11).Build()
//				convey.So(funcA(), convey.ShouldEqual, 11)
//			})
//			convey.So(funcA(), convey.ShouldEqual, 1)
//		})
//
//		mockey.PatchConvey("mock 2", func() {
//			mockey.Mock(funcA).Return(2).Build()
//			convey.So(funcA(), convey.ShouldEqual, 2)
//		})
//
//		mockey.PatchConvey("mock not", func() {
//			convey.So(funcA(), convey.ShouldEqual, 0)
//		})
//	})
//}

// 测试多个test 的 mock 是否会相互影响
// 结论, 多个测试函数一起运行时, mock会产生相互影响

//func TestMock1(t *testing.T) {
//	mockey.PatchConvey("mock 1", t, func() {
//		convey.So(funcA(), convey.ShouldEqual, 0)
//		mockey.PatchConvey("mock 1 in", func() {
//			mockey.Mock(funcA).Return(1).Build()
//			convey.So(funcA(), convey.ShouldEqual, 1)
//		})
//	})
//}
//
//func TestMock2(t *testing.T) {
//	mockey.PatchConvey("mock 2", t, func() {
//		convey.So(funcA(), convey.ShouldEqual, 0)
//		mockey.PatchConvey("mock 2 in", func() {
//			mockey.Mock(funcA).Return(2).Build()
//			convey.So(funcA(), convey.ShouldEqual, 2)
//		})
//	})
//}
