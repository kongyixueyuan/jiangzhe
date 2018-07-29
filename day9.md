1. 尽可能多的写出整个课程，我们学习并掌握到的知识点
    - 创世区块以及区块的生成过程
    - 区块链的创建并且如何将区块添加到区块链中
    - POW算法的实现原理和如何将POW算法集成到区块链中
    - 数据的持久化，使用数据库存储区块信息
    - 打印和输出所有的区块信息，并且进行优化（迭代输出）
    - CLI命令行工具的使用（通过命令行工具去执行某一功能）
    - UTXO（未花费交易输出模型）实现原理
    - 默克尔树实现原理和钱包地址生成的过程并且将钱包地址集成到项目中
    - 主节点，钱包节点，矿工节点的模拟

2. 描述我们课程中POW代码实现逻辑
    - 在新增区块的过程中，区块的哈希由POW计算而得来
    - 创建一个1的大数，并且左移256位。设置难度值如果为4的(x)倍则代表计算而来的数前面有x个零
    - 设置一个nonce值，不断的改变nonce值去尝试找到比我们创建出来的数小的数
    - 如果找到了那么就是挖矿成功并且返回hash和nonce值。区块得到这两个值后加入到自己的属性中

3. 课程中数据库如何实现增删改查
    - 使用到的是开源的key=>value的数据库 blotdb（https://github.com/boltdb/bolt）
    - 新增
    ```
    err := db.Update(func(tx *bolt.Tx) error {
        ...
        return nil
    })
    ```
    - 查询
    ```
    err := db.View(func(tx *bolt.Tx) error {
    	...
    	return nil
    })
    ```

    - 更新
    ```
    db.Update(func(tx *bolt.Tx) error {
    	b := tx.Bucket([]byte("MyBucket"))
    	err := b.Put([]byte("answer"), []byte("42"))
    	return err
    })
    ```
    - 删除

    ```
    db.Update(func(tx *bolt.Tx) error {
    	b := tx.DeleteBucket([]byte("MyBucket"))
    	err := b.Put([]byte("answer"), []byte("42"))
    	return err
    })
    ```

4. 图文并貌完整的描述钱包地址生成过程
    1. 椭圆曲线算法生成私钥并且通过私钥生成公钥
    2. 根据公钥生成地址（使用hash160对对公钥两次Hash）并拼接版本信息
    3. 将2再进行两次hash256加密
    4. 将2和3拼接在一起使用base58加密

5. 图文并貌描述据两个实例描述UTXO模型的巧妙设计
    - （交易输入）tx_input 和 （交易输出）tx_output
    - 创世区块A得到100个token的奖励那么就是 tx_input{[]byte{}, "-1", "FirstBlock"} tx_output{100, A}
    - 如果A给B转了50块钱
        - tx_input{100,0,A} tx_output{50, B}, tx_output{50, A}
        - 首先扣除了A的100个Token，然后给B输出了50个Token，又给A输出了50个Token

6. 私钥签名，公钥是如何验证签名的

7. 完整的描述节点区块同步的逻辑和过程
    - 启动节点，（主节点：）阻塞等待客户端发送命令（非主节点：）向主节点请求数据
    - 如果主节点的最高区块大于请求节点，那么主节点向请求节点发送version否则向请求节点请求数据

8. 钱包节点转账、主节点、矿工节点之间的完整交互逻辑

9. 怎么理解libp2p实现节点数据同步

10. 运行Otto，编写一个简单的合约，将合约提交到虚拟机进行编译运行，附属上相关截图