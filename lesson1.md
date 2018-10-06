```
作业：
1、需要对ring signature，zk-snark进行比较，ring signature decoy的数量在多少的时候更消耗时间空间，zk-sanrk在使用上占用多少空间，计算时间相比哪个更快
2、将bitcoin，ethereum，monero，zcash，EOS 的交易、相关交易易属性、块大小以及填入多少交易写在report中
```

A1：
ring signature测试代码来至于https://github.com/t-bast/ring-signatures

```
func benchmarkSign(ringSize int, b *testing.B) {
	pubKeys, privKeys := GenerateKeys(ringSize)
	i := rand.Intn(ringSize)
	message := []byte("Benchmark me like the french people do.")
	for n := 0; n < b.N; n++ {
		_, err := privKeys[i].Sign(nil, message, pubKeys, i)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSign3(b *testing.B)   { benchmarkSign(3, b) }
func BenchmarkSign12(b *testing.B)  { benchmarkSign(12, b) }
func BenchmarkSign100(b *testing.B) { benchmarkSign(100, b) }
```

```
ubuntu@ubuntu:~/go/src/github.com/ring-signatures/ring$ go test -test.bench BenchmarkSign3
goos: linux
goarch: amd64
pkg: github.com/ring-signatures/ring
BenchmarkSign3-4   	      50	  21895635 ns/op
PASS
ok  	github.com/ring-signatures/ring	1.495s
ubuntu@ubuntu:~/go/src/github.com/ring-signatures/ring$ go test -test.bench BenchmarkSign12
goos: linux
goarch: amd64
pkg: github.com/ring-signatures/ring
BenchmarkSign12-4   	      10	 104238347 ns/op
PASS
ok  	github.com/ring-signatures/ring	1.569s
ubuntu@ubuntu:~/go/src/github.com/ring-signatures/ring$ go test -test.bench BenchmarkSign100
goos: linux
goarch: amd64
pkg: github.com/ring-signatures/ring
BenchmarkSign100-4   	       1	1340726054 ns/op
PASS
ok  	github.com/ring-signatures/ring	1.701s

```

从上面可以看到, ringSize(对应decoy)越大，消耗的时间越多。


A2：

calculateRunway() 函数每次执行消耗的Gas如下：  

| 项目     | 块大小   | 交易大小 | 交易数/块 |  交易属性  |
| -------- | -------- | -------- | --------  |   :----:   |
| bitcoin  | < 1M | 取决于输入和输出的个数 | 3000左右 | 
| ethereum | 26s | 不固定 |
| monero   | 2 * M100 | 取决于ring的大小  |  
| zcash    | < 2M |  
| EOS      | 没有限制 | 不固定  | 3000左右  |  1笔交易可以包含多个action  |




