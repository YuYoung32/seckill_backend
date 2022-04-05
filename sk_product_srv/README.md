# Sk_product_srv Service

这是秒杀服务的商品服务部分，主要功能如下

- 添加产品服务，将管理员给到的商品信息写入products表
- 产品信息查询服务：查询数据库里的products表，返回商品信息供管理员或用户查看
- 产品信息编辑服务：将管理员给到的新的商品信息更新products表


使用时需要更改`conf/config_example.json`文件，并重命名为`conf/config.json`