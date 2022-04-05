# Sk_seckill_srv Service

这是秒杀服务的订单服务部分，主要功能如下

- 添加活动服务，将管理员给到的活动信息写入events表
- 活动信息查询服务：查询数据库里的events表，返回活动信息供管理员或用户查看
- 活动信息编辑服务：将管理员给到的新的活动信息更新events表
- 创建订单服务：从rabbitmq里取出创建订单请求，验证是否符合下订单条件并创建订单写入orders表
- 返回结果服务：根据创建订单结果存入redis供查询订单结果时使用

使用时需要更改`conf/config_example.json`文件，并重命名为`conf/config.json`