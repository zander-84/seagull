package think

// BenchmarkToGrpcString-16    	 1856049	       650.7 ns/op
//func BenchmarkToGrpcString(b *testing.B) {
//	err := ErrSystemSpace("this is a  system error")
//	thinkErr := FromError(err)
//
//	for i := 0; i < b.N; i++ {
//	}
//
//}
//
//// BenchmarkToGrpcStringParallel-16    	 4887981	       268.8 ns/op
//func BenchmarkToGrpcStringParallel(b *testing.B) {
//
//	err := ErrSystemSpace("this is a  system error")
//	thinkErr := FromError(err)
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//		}
//	})
//}
