### BlockDemo
hash算法 （任意输入）返回一个256位的哈希 (转化为16进制一共是64个字节)
block
    - 必要条件
        - Block 编号
        - Nonce nonce值
        - Data  数据
        - prev  上一个区块哈希（创世区块为0）
        - Hash  由所有必要条件组成进行哈希得到的哈希值

blockChain
    - 将每个区块序列化
    - 序列化后存入数组
    - 数据库除了存储每一个区块以外，还会存储最后区块的哈希


### POW（Proof of Work）工作量证明——多劳多得
    -PoW机制中根据矿工的工作量来执行货币的分配和记账权的确定。算力竞争的胜者将获得相应区块记账权和比特币奖励。因此,矿机芯片的算力越高,挖矿的时间更长,就可以获得更多的数字货币。

    优点：算法简单，容易实现；节点间无需交换额外的信息即可达成共识；破坏系统需要投入极大的成本。
    缺点：浪费能源；区块的确认时间难以缩短；新的区块链必须找到一种不同的散列算法，否则就会面临比特币的算力攻击；容易产生分叉，需要等待多个确认；永远没有最终性，需要检查点机制来弥补最终性。

### POS（Proof of Stake）股权证明算法——持有越多，获得越多
    - 股权证明算法——持有越多，获得越多

    POS 机制采用类似股权证明与投票的机制，选出记帐人，由它来创建区块。持有股权愈多则有较大的特权，且需负担更多的责任来产生区块，同时也获得更多收益的权力。
    POS 机制中一般用币龄来计算记账权，每个币持有一天算一个币龄,比如持有 100 个币，总共持有了 30 天,那么此时的币龄就为 3000。在 POS 机制下,如果记账人发现一个 POS 区块, 他的币龄就会被清空为 0,每被清空 365 币龄,将会从区块中获得 0.05 个币的利息(可理解为年利率 5%)。

    优点：在一定程度上缩短了共识达成的时间；不再需要大量消耗能源挖矿。
    缺点：还是需要挖矿，本质上没有解决商业应用的痛点；所有的确认都只是一个概率上的表达，而不是一个确定性的事情，理论上有可能存在其他攻击影响。

### DPOS（Delegated Proof-of-Stake）股份授权证明
    - DPOS 是在 POS 基础之上发展起来的。与PoS的主要区别在于持币者投出一定数量的节点，代理他们进行验证和记账。其合规监管、性能、资源消耗和容错性与PoS相似。
    DPoS的工作原理为：每个股东按其持股比例拥有影响力，51%股东投票的结果将是不可逆且有约束力的。其挑战是通过及时而高效的方法达到51%批准。为达到这个目标，每个股东可以将其投票权授予一名代表。获票数最多的前100位代表按既定时间表轮流产生区块。每名代表分配到一个时间段来生产区块。所有的代表将收到等同于一个平均水平的区块所含交易费的10%作为报酬。如果一个平均水平的区块含有100股作为交易费，一名代表将获得1股作为报酬。DPoS的投票模式可以每30秒产生一个新区块。

    优点：大幅缩小参与验证和记账节点的数量，可以达到秒级的共识验证。
    缺点：整个共识机制还是依赖于代币，很多商业应用是不需要代币存在的。

### 区块结构体定义的基本条件
1. 区块高度（编号）
2. 交易数据
3. 时间戳
4. 上一个区块的哈希
5. 哈希

```
type Block struct {
    Height      int64
    Data        []byte
    Timestamp   int64
    PrevHash    []byte
    Hash        []byte
}
```

### 哈希生成的思路
- 1.将区块高度转换为字节数组
- 2.将时间戳转换为字节数组
    - 2_1.使用strconv.FormatInt()转换为2进制的字符串
    - 2_2.使用[]byte()强制类型转换成byte数组
- 3.拼接Height,Data,Timestamp,PrevHash生成新的[]byte数组
- 4.将拼接后的字节数组进行256哈希

### 总结：创建区块链的过程

> 分两个文件，两个结构体
file 1 Block.go         用来生成区块
```
type Block struct {
    Height      int64
    Data        []byte
    Timestamp   int64
    PrevHash    []byte
    Hash        []byte
}
```

file 2 BlockChain.go    创建链，将区块加入到链中
```
type BlockChain struct {
    BlockChain   []*Block
}
```

- 1. mian.go => BlockChain.go 通过主函数BlockChain的方法，用来创建一条区块链，并且将创世块加入到里面

- 2. BlockChain.go => Block.go 调用Block的创建创世块的方法 =>  调用创建一个新区块的方法 => 调用生成Hash的方法

- 3. Block.go => BlockChain.go => 接收到生成Hash方法返回的Hash => 将Hash存入区块中并返回区块 => 创世区块创建完成并返回 => 区块链接收到区块并加入到链中  => 返回区块链对象

- 4. 添加区块：Block.go => BlockChain.go => 调用Block的NewBlock方法继续添加区块。


### 工作量证明算法的原理

区块中若干个属性拼接生成的Hash值如果小于当前系统给定的难度值，那么就挖矿成功。

比如说：
    区块中有当前属性：`Height, Data, Timestamp, PrevHash。`这些属性，
    由于时间戳不一样所以每次生成的Hash值也不一样。
    当然在其它区块链中除了时间戳不一样以外：还有一个被称作`nonce`的值，它每次循环都会自增1

首先我们有一个变量为int型1，然后将它进行左移位。如下：
`0000 0001 ... 0000 0000` (一共256位)

比如我们设置难度值为多少就往左移(256-难度值)多少，比如难度值我设为20那么对应的我应该往左移236位也就是：
`0000 0000 0000 0000 0001 ... 0000 0000`

如上可见，当难度值设置的越大的时候，这个数就越小。

而我们挖矿实际上就是不断的替换`nonce`值去生成`hash`如果小于当前的难度值则算挖矿成功。举个例子：
如果难度值设置为256，那么我们挖到矿的条件就只有一个，那就是`hash`为1


### 新增bolb数据库
区块存储逻辑：
- 创世块：
1. 将区块序列化：以`key=>value`的方式进行存储 （`key`为当前区块的`Hash，Value`为当前区块序列化后的字符串）
2. 还需要存储当前最新区块的Hash 用l表示 (l => 最新区块的哈希)

- 新增区块
1. 新增区块的上一个区块的Hash是通过数据库l字段获取
2. 当新增完成后需要更新l字段为当前最新区块的哈希值

### 遍历区块
1. 通过`for`循环去遍历数据库中的区块。需要设置一个变量用来接收查找的区块Hash，没次查找完毕后将上一个区块的哈希赋值给变量
2. 通过数据库获取最新区块的哈希，找到最新的区块。
3. 将最新区块的上一个区块哈希赋值给变量



