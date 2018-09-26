# 第一题
libsnark 运行报出Floating Point Exception，还没解决。弄好了补上
# 第二题

币 | 相关交易属性（主要） | 块大小 | 块交易数
--- | --- | --- | ---
bitcoin | inputs, outputs, witness, lock_time | 平均约0.8M，理论极限接近4M（with Segwit） | 平均约1500
ethereum | nonce, gasPrice, gasLimit, to, value, data | 平均约20k，由block gas limit决定 | 平均约80，由block gas limit决定
monero | ? | 平均约100k | 平均约6.58（过去一年统计）
zcach | ? | 最大2M，平均约50k | 上限（估算）：6.67tx/s for Shielded tx, 26.67 tx/s for transparent tx
EOS | ? | 最大1M（可投票修改） | 没查到~
