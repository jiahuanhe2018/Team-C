> 作业：
> 
> 1. 需要对ring signature，zk-snark进行比较，ring signature decoy的数量在多少的时候更消耗时间空间，zk-sanrk在使用上占用多少空间，计算时间相比哪个更快
> 
> 2. 将bitcoin，ethereum，monero，zcash，EOS 的交易、相关交易易属性、块大小以及填入多少交易写在report中

## A1

> [匿名币的核心技术环签名和zk-SNARK的对比](https://www.tuoluocaijing.cn/article/detail-11625.html)

### ring signature

> 参考项目 <https://github.com/fernandolobato/ecc_linkable_ring_signatures>

编写测试代码
```
from linkable_ring_signature import ring_signature, verify_ring_signature

from ecdsa.util import randrange
from ecdsa.curves import SECP256k1

import sys
import random

number_participants = int(sys.argv[1])

x = [ randrange(SECP256k1.order) for i in range(number_participants)]
y = list(map(lambda xi: SECP256k1.generator * xi, x))

message = "Every move we made was a kiss"

i = random.randrange(number_participants)
signature = ring_signature(x[i], i, message, y)

assert(verify_ring_signature(message, y, *signature))
```

输入不同数量的decoy，统计时间、空间
```

breeze@ubuntu:~/dev/blockchain/ecc_linkable_ring_signatures$ /usr/bin/time -f "%e %M" python3 test.py 1
0.22 12776
breeze@ubuntu:~/dev/blockchain/ecc_linkable_ring_signatures$ /usr/bin/time -f "%e %M" python3 test.py 10
1.90 12836
breeze@ubuntu:~/dev/blockchain/ecc_linkable_ring_signatures$ /usr/bin/time -f "%e %M" python3 test.py 10
1.88 12832
breeze@ubuntu:~/dev/blockchain/ecc_linkable_ring_signatures$ /usr/bin/time -f "%e %M" python3 test.py 100
18.83 12884
breeze@ubuntu:~/dev/blockchain/ecc_linkable_ring_signatures$ /usr/bin/time -f "%e %M" python3 test.py 200
38.30 13052
breeze@ubuntu:~/dev/blockchain/ecc_linkable_ring_signatures$ /usr/bin/time -f "%e %M" python3 test.py 1000
189.01 13760
```

可见，`ring signature`算法消耗的时间与`decoy`数量呈线性关系，**decoy越多，破解难度越高，消耗的时间和空间也越多**

### zk-snark

> 参考项目 <https://github.com/Charterhouse/pysnark>  

还没完全弄懂，等弄明白代码再补充答案

## A2


> Bitcoin <https://blockexplorer.com/>  
> Ethereum <https://etherscan.io/>  
> Monero <https://moneroblocks.info/>  
> Zcash <https://explorer.zcha.in/>  
> EOS <https://eostracker.io/>  

|   Name   | 块大小 | 出块时间 |   主要交易属性  |
|----------|--------|----------|-----------------|
| Bitcoin  | < 1M   | ~15 min  | Value, From, To |
| Ethereum | ~32K   | ~15 secs | Value, From, To |
| Monero   | ~      | ~5 secs  |                 |
| Zcash    | ~      | ~2 min   |                 |
| EOS      | ~      | ~0.5 s   |                 |


