# PoS简单实现
小项目的执行命令：

	go run PoS/main.go
	./nc.exe localhost 9000

需要注意的点：

- 将验证者也就是出块者直接写到区块里。
- server看待验证的区块通道chan有没有区块，如果有，将其加入到临时区块队列中TempBlocks。
- 然后从待验证者中选出验证者胜利者。
	- 遍历临时区块列表，将临时区块中所有不在乐透池中的验证者都放到乐透池中。
	- 根据Token押注大小k，放入k份验证者到乐透池中。
	- 随机从乐透池中选出一个胜利者。
	- 让胜利者出块，然后往announcements通道广播自己活得出块的胜利。
	- 清空临时区块队列TempBlocks。
- 用server和client模拟，新的挖矿节点接入的时候，需要做如下的事情：
	- 用户输入抵押Token。
	- 实际上用户没有输入自己的address，这里是直接根据时间生成的address。
	- 然后扫描所有输入的抵押Token结果，每个验证者的抵押都有了。
	- 然后扫描素有输入的result结果，获取oldLastIndex。
	- 根据oldLastIndex，result和address生成新的区块newBlock。
	- 验证区块合法性，将新区块放到待验证的区块通道chan。