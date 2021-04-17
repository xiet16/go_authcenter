1 授权码模式
获取code:
前端: GET /authorize  

参数请求               类型         说明
client_id           string      移动端app_client,管理网站 manage_web
response_type       string      固定值code
scope               string      请求的权限范围 server，deviceservice,manageweb
state               string      验证请求的标志字段
redirect_uri        string      回调url

测试url:
http://localhost:9096/authorize?client_id=app_client
&response_type=code&scope=all&redirect_uri=http://localhost:3300

获取access_token:
后端 POST /token
参数请求               类型         说明
grant_type          string      固定值 authorization_code
code                string      上一步获取的code
redirect_uri        string      填写的重定向uri


2 password 模式
后端 POST /token   
grant_type string   固定值password
username   string   用户名
password   string   用户密码
scope      string   权限范围


遗留问题:
accss_token 再 redis 中的过期问题
redis中access_token 的管理问题
生成token，验证token过程中的异常问题
redis 集群崩溃问题

单点登录问题