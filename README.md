# tsing-gateway
Tsing Gateway是一个开源、跨平台、去中心化集群、动态配置的API网关。

## 功能特性
- 虚拟主机，按主机名隔离数据，可多项目共用API网关
- 路由组/路由，根据URL匹配上游
- 服务发现，从`Tsing Center`、`Nacos`等服务中心获取服务节点
- 中间件，在网关运行时，为主机或服务动态挂载或卸载的中间件机制
- 持久存储，网关的配置可持久存储，并支持`etcd`、`consul`、`redis`多种数据源
- 去中心化集群，轻松组建横向扩展的网关集群，并用任意网关节点做流量入口
- 动态配置，可通过Admin API对网关的配置进行动态变更，无需重启网关进程

### 中间件
- `auto_response`，根据客户端请求路径和方法，自动响应状态码以及内容
- `jwt_proxy`，JWT反向代理，转发JWT到上游并返回校验结果
- `set_header`，在HTTP会话的请求或响应头部中设置Header参数
- `url_rewrite`，将客户端的请求URL重写后转发给上游

### 存储数据源
- [x] etcd
- [ ] consul
- [ ] redis

### 服务发现
- [x] Tsing Center
- [ ] Nacos

## 相关资源

- [Tsing](https://github.com/dxvgef/tsing) 高性能、微核心的Go语言HTTP服务框架
- [Tsing Center](https://github.com/dxvgef/tsing-center) 开源、跨平台、去中心化集群、动态配置的服务中心

## 用户及案例

如果你在使用本项目，请通过[Issues](https://github.com/dxvgef/tsing-gateway/issues)告知我们项目的简介

## 帮助/说明

本项目处于开发初期阶段，API和数据存储结构可能会频繁变更，暂不建议在生产环境中使用，如有问题可在[Issues](https://github.com/dxvgef/tsing-gateway/issues)里提出。

诚邀更多的开发者为本项目开发管理面板和官方网站等资源，帮助这个开源项目更好的发展。
