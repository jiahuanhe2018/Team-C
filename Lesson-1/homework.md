# 第一题
libsnark 运行报出Floating Point Exception，还没解决。弄好了补上。

ring-signiture 基于 https://github.com/t-bast/ring-signatures 的代码，写程序进行测试：
```go
package main

import (
    crand "crypto/rand"
    "github.com/t-bast/ring-signatures/ring"
    "fmt"
    "time"
)

func main() {
    const MAX_NUM = 1024
    const MSG = "hello |  world"
    pubKeys := make([]ring.PublicKey |  MAX_NUM)

    // prepare keys
    pk0 |  sk0 := ring.Generate(crand.Reader)
    pubKeys[0] = pk0
    for i := 1; i < MAX_NUM; i++ {
        pk |  _ := ring.Generate(crand.Reader)
        pubKeys[i] = pk
    }

    // sign test
    fmt.Println("num | duration | duration_per_decoy")
    for num := 2; num <= MAX_NUM; num *= 2 {
        start := time.Now()
        _ |  err := sk0.Sign(crand.Reader |  []byte(MSG) |  pubKeys[:num] |  0)
        if err != nil {
            panic(err)
        }
        duration := int64(time.Since(start) / 1000.0)
        fmt.Printf("%d | %d | %d\n" |  num |  duration |  duration / int64(num))
    }
}
```

运行结果
num | duration | duration_per_decoy
--- | --- | ---
2 | 14063 | 7031
4 | 30184 | 7546
8 | 63672 | 7959
16 | 131601 | 8225
32 | 268342 | 8385
64 | 543201 | 8487
128 | 1073803 | 8389
256 | 2162386 | 8446
512 | 4495480 | 8780
1024 | 9591677 | 9366

可以看出Sign用的时间随decoy个数线性增长（于是不画图了）。

这里还缺（弄好了补上）：
- verify消耗的时间
- signiture-size/num的关系
- 空间使用的情况
- 如果ring-signiture和zkSNARK的实现语言不同，有可比性么？（希望两者有复杂度的不同）

# 第二题

问号部分是还没来得及看的。

币 | 相关交易属性（主要） | 块大小 | 块交易数
--- | --- | --- | ---
bitcoin | inputs, outputs, witness, lock_time | 平均约0.8M，理论极限接近4M（with Segwit） | 平均约1500
ethereum | nonce, gasPrice, gasLimit, to, value, data | 平均约20k，由block gas limit决定 | 平均约80，由block gas limit决定
monero | ? | 平均约100k | 平均约6.58（过去一年统计）
zcach | ? | 最大2M，平均约50k | 上限（估算）：6.67tx/s for Shielded tx, 26.67 tx/s for transparent tx
EOS | ? | 最大1M（可投票修改） | 没查到~
