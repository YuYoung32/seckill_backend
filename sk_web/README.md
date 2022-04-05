# Sk_web Service

这是秒杀服务的web网关部分

主要功能如下

- 使用gin作为路由，进行权限控制（token）、数据格式验证与转换、使用grpc调用注册在consul服务
- 将下订单请求存入rabbitmq消息队列进行流量削峰
- 读取redis里的请求下订单的结果并返回结果

使用时需要配置`conf/config_example`并重命名为`config.json`