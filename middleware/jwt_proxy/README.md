# jwt_proxy

JWT反向代理，网关自动转发JWT到上游并返回校验结果，支持从header、query、form、cookie中获取JWT，并以同样四种方式转发给上游。

## 配置参数说明
- `name` 中间件的名称，必须是`jwt_proxy`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `key` 为触发中间件的URL路径，`*`匹配所有路径
    - `source_type`，string，必需，来源类型，值限定：`header`,`query`,`form`,`cookie`
    - `source_name`，string，必需，来源参数
    - `upstream_url`，string,，必需，上游URL
    - `send_type`，string，必需，往上游发送的类型，值限定：`header`,`query`,`cookie`
    - `send_method`，string，必需，往上游发送的HTTP方法，值限定：`GET`,`HEAD`,`OPTIONS`,`TRACE`
    - `send_name`，string，必需，往上游发送使用的参数名
    - `upstream_success_body`，string，可选，上游JWT校验成功的响应HTTP.Body的字符串，用于处理只返回200状态码的API，如果留空则不校验body

## `config`字段示列
```json
{
  "source_type": "header",
  "source_name": "Authorization",
  "upstream_url": "/auth",
  "send_type": "header",
  "send_method": "GET",
  "send_name": "token",
  "upstream_success_body": "success"
}
```

## 完整参数示例：

```json
{
  "name": "jwt_proxy",
  "config": "{\"source_type\":\"header\",\"source_name\":\"Authorization\",\"upstream_url\":\"/auth\",\"send_type\":\"header\",\"send_method\":\"GET\",\"send_name\":\"token\",\"upstream_success_body\":\"success\"}"
}
```


