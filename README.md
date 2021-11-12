# techtrainingcamp-red-envelope

## 接口文档介绍

### 消息状态码

| code |        状态        |
| :--: | :----------------: |
|  0   |      请求成功      |
|  1   |     无效的参数     |
|  2   |       被限流       |
|  3   |   抢到达最大限制   |
|  4   |      请求失败      |
|  5   | 同一个用户请求太快 |



### 抢红包接口

- 每个用户最多只能抢到max_count次，次数可配置
- 一定概率能抢到，概率可配置
- 同一个uid最多2s调用一次该接口

请求

POST	.../v0/snatch

```json
{
    "uid": "123" 				// 用户id string类型
}
```

响应

```json
{
    "code": 0,
    "data": {
        "cur_count": 1,				// 当前第几次抢到
        "enveloped_id": "c67479bbu3iato50giu0", // 红包id 
        "max_count": 10 			// 最多抢到次数
    },
    "msg": "success"
}
```



### 拆红包接口

- 当uid与envelope_id不匹配或对应的红包已被拆，会返回msg为无效的参数

请求

POST	.../v0/open

```json
{
    "uid": "123",				// 用户id，string类型
    "envelope_id":"c67479bbu3iato50giu0"	// 红包id，string类型
}
```

响应

```json
{
    "code": 0,
    "data": {
        "value": 3	    			// 红包金额，以"分"为单位
    },
    "msg": "success"
}
```



### 钱包列表接口

- 显示总余额
- 显示抢到的红包，有已拆红包和未拆红包
- 按红包获取时间排序，从最新获取开始
- 为应对高并发，该接口有10s延迟

请求

POST	.../v0/get_wallet_list

```
{
    "uid": "123" // 用户id string类型
}
```

响应

```json
{
    "code": 0,
    "data": {
        "amount": 2,
        "envelope_list": [
            {
                "envelope_id": "c676vv2849te1d3rfsbg",
                "opened": false,
                "snatch_time": 1636724856
            },
            {
                "envelope_id": "c676vv2849te1d3rding",
                "value": 1,
                "opened": true,
                "snatch_time": 1636724851
            },
            {
                "envelope_id": "c676vuq849te1d3qqpug",
                "value": 1,
                "opened": true,
                "snatch_time": 1636724821
            },
            {
                "envelope_id": "c676vuq849te1d3qo52g",
                "opened": false,
                "snatch_time": 1636724819
            }
        ]
    },
    "msg": "success"
}
```

### 显示当前配置接口

- 显示当前已发出去的金额、红包数量，总共的金额、红包数以及每个用户最大抢次数和概率

请求

POST	.../get_config

响应

```json
{
    "code": 0,
    "data": {
        "cur_amount": 1, 		// 当前发出去的总金额
        "cur_size": 1,			// 当前发出去的总红包数
        "max_amount": 1000,		// 当前设置的总金额
        "max_count": 10,		// 用户最多可抢的次数
        "max_size": 1000,		// 当前设置的总红包数
        "snatched_pr": 100		// 用户抢到的概率
    },
    "msg": "success"
}
```

### 更改当前配置接口

- 可在活动期间热更新相关配置

请求

POST .../config

```json
{
    "max_amount": 1000000,	// 多添加的金额（注意是多添加）
    "max_count": 10,		// 用户最多可抢的次数
    "max_size": 1000000,	// 多添加的红包数（注意是多添加）
    "snatched_pr": 80		// 用户抢到的概率
}
```

响应

```
{
    "code": 0,
    "data": {
        "max_amount": 4000000,
        "max_count": 20,
        "max_size": 3000000,
        "snatched_pr": 80
    },
    "msg": "success"
}
```

