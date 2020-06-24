# auto_response

根据客户端请求路径自动响应状态码以及内容，可用于自动处理`OPTIONS`方法以及`/favicon.ico`路径的请求

## 配置参数说明
- `name` 中间件的名称，必须是`auto_response`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `key` 为触发中间件的URL路径，`*`匹配所有路径
    - `method`，string，必需，触发中间件的HTTP方式，值限定：所有HTTP Method，`ANY`匹配所有方法
    - `status`，int，必需，为中间件要输出的HTTP状态码，值限定：所有HTTP Status Code
    - `data`，string,可选，为中间件要输出的HTTP Body，如果是空，则不输出HTTP Body

## `config`字段示列
```json
{
  "*": {
    "method": "OPTIONS",
    "status": 204
  },
  "/favicon.ico": {
    "method": "GET",
    "status": 204
  }
}
```

## 完整参数示例：
```json
{
  "name": "auto_response",
  "config": "{\"*\":{\"method\":\"OPTIONS\",\"status\":204},\"/favicon.ico\":{\"method\":\"GET\",\"status\":204}}"
}
```
