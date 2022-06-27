# simple-api-service
用Golang实现的简单API服务，可以用于一些定制化的场景。
比如，在一些内网的混合产品的环境里，某个服务没有提供对外的接口，而我们又想获取该服务的某些数据，就可以用这个简单API应用，通过连上该服务的数据库取数据。

## 简析
先通过getToken获取请求用的token，然后再调用接口取数据。示例代码中，服务监听了10000端口，提供了/getToken、/v1/get/test、/v1/custom/sql三个接口，其中sql接口传入自定义sql直接取数据。

## 用到的技术
- 简单JWT
- mysql连接与操作
- Gin API服务
- yaml配置文件读取
