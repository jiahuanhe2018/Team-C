# PoW简单实现
小项目的执行命令：

	go run Pow/main.go
	curl -i -X POST -H 'Content-Type:application/json' -d '{"Result":45}' http://localhost:9000

需要注意的点：

- 主要是进行区块的PoW挖矿。
- 然后跑一个简单的server来模拟。
- 一个Get请求，当我们收到一个http请求的时候，写入blockchain。
- Client发起一个POST请求，然后开始生成新的区块，加入区块链，返回给Client。