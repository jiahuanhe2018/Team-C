# Lecture-1 summary

### ring signature和zk-snark进行比较
&emsp;&emsp;简单理解环形签名，有n个用户，他们都有一对公私钥，有一个用户有一个信息m，需要给其他所有用户证明自己拥有信息m。

- 第一阶段：签名阶段，这个用户保存其他所有用户的公钥，这个用户用所有人的公钥加上自己的私钥对信息m进行签名得到签名指纹S。
- 第二阶段：广播阶段，用户将签名指纹S和信息m广播给所有用户。
- 第三阶段：验证阶段，每个用户都可以用自己的公钥对指纹S和信息m进行验证。

decoy output是什么意思呢，简单理解，就是混入其他多少个公钥进来混淆加密，这样混的多了，任何一个验证方都难以知道是谁发起的这个签名。

项目：[https://github.com/t-bast/ring-signatures](https://github.com/t-bast/ring-signatures)

为了验证decoy多与少的时间消耗，写了一个脚本：

	func sign_time(){
        nums := []int{2, 3,5,7,11,57,111}
        PK, SK := ring.Generate(crand.Reader)
        message := "wtf"
        for _, num := range nums {
                var ringKeys []ring.PublicKey
                ringKeys = append(ringKeys, PK)
                for i :=0;i<num;i++ {
                        pk, _ := ring.Generate(crand.Reader)
                        ringKeys = append(ringKeys, pk)
                }
                t0 := time.Now().UnixNano()
                SK.Sign(crand.Reader,[]byte(message),ringKeys,0)
                t1 := time.Now().UnixNano()
                run_time := t1-t0
                fmt.Printf("decoy number:%d nanosecond",num)
                fmt.Printf("sign run time: %d nanosecond",run_time)
                fmt.Println()
        }
	}

	
结果如下：

	decoy number:2 nanosecondsign   run time: 23310644 nanosecond
	decoy number:3 nanosecondsign   run time: 32150456 nanosecond
	decoy number:5 nanosecondsign   run time: 49989242 nanosecond
	decoy number:7 nanosecondsign   run time: 66954814 nanosecond
	decoy number:11 nanosecondsign  run time: 103006057 nanosecond
	decoy number:57 nanosecondsign  run time: 512016500 nanosecond
	decoy number:111 nanosecondsign run time: 991204405 nanosecond

接下来看看zkSNARK，如果要想简单理解，可以参考我这博客：[零知识证明中的超级新星：zk-SNARKs](https://www.jianshu.com/p/1df33c10fd22)

当然要想更加深入理解zk-SNARKs，绝非这么容易，了解下基本的密码学概念：

- Zero knowledge：零知识，证明者向验证者证明某个信息的情况下，不会给验证者提供任何有用的信息。
- Succinctness：证据信息较短，方便验证。
- Non-interactivity：没有交互，这对于区块链来说就非常完美，直接把签名后的“无用”的信息放在公链让验证者可以公开验证。
- Arguments：证明过程是计算完好（computationally soundness）的，证明者无法在合理的时间内造出伪证。
- of knowledge：对一个证明者来说，在不知晓特定证明（witness）的前提下，构建一个有效的零知识是不可能的。

有一个C++的zk-SNARKs的项目[https://github.com/scipr-lab/libsnark](https://github.com/scipr-lab/libsnark)

zk-SNARKs占用空间大约40M左右，环形签名和zk-SNARKs比较，后者的运行更有效率。

### bitcoin，Ethereum，menero，zcash和EOS的交易，交易属性，块大小以及填入多少交易

<table>
	<tr>
		<td>币种</td>
		<td>出块时间</td>
		<td>交易属性</td>
		<td>共识算法</td>
		<td>块大小</td>
		<td>填入交易数</td>
	</tr>
	<tr>
		<td>Bitcoin</td>
		<td>10min</td>
		<td>最终一致性的确认需要多个出块来保证</td>
		<td>PoW</td>
		<td>1M</td>
		<td>2k-3k</td>
	</tr>
	<tr>
		<td>Ethereum</td>
		<td>14s</td>
		<td>gas消耗机制，支持图灵完备的智能合约</td>
		<td>PoW</td>
		<td>最大1500000Gas，每笔转账消耗21000Gas</td>
		<td>70</td>
	</tr>
	<tr>
		<td>Menero</td>
		<td>600s</td>
		<td>使用环签名，隐形地址保证匿名</td>
		<td>防ASIC的CryptoNote的核心算法</td>
		<td>不固定，最大为前100个块大小的中位数的2倍，最小块300KB</td>
		<td>交易数量不固定，能接近1000</td>
	</tr>
	<tr>
		<td>ZCash</td>
		<td>2.5min</td>
		<td>zk-SNARKs实现匿名</td>
		<td>抗ASIC的EquiHash算法</td>
		<td>2M</td>
		<td>3k左右</td>
	</tr>
	<tr>
		<td>EOS</td>
		<td>0.5s</td>
		<td>快速，地址记忆友好，权限体系</td>
		<td>DPoS+PBFT</td>
		<td>不固定</td>
		<td>1500以上</td>
	</tr>
</table>