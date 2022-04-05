# Sk_user_srv Service

这是秒杀服务的用户服务部分，主要功能如下：

- 用户与管理员登陆服务，查询数据库里的users和admins表，返回是否验证成功
- 用户注册服务，通过beego发送验证码邮件，校验验证码与用户信息写入
- 用户信息查询服务：查询数据库里的users和admins表，返回用户信息供管理员查看

使用时需要更改`conf/config_example.json`文件，并重命名为`conf/config.json`