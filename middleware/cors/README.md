# cors

跨域资源共享控制

## 配置参数说明
- `name` 中间件的名称，必须是`cors`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `allow_origins` 允许客户端的来源域，默认值`["*"]`
  - `allow_credentials` 允许客户端携带cookie，默认值`true`
  - `allow_methods` 允许客户端请求的方法，默认值`["*"]`
  - `allow_headers` 允许客户端请求的头信息，默认值`["*"]`
  - `expose_headers` 允许响应的头信息，默认值`["*"]`

## `config`字段示列
```json
{
  "allow_origins": ["*"],
  "allow_headers": ["*"],
  "allow_credentials": true,
  "allow_methods": ["*"],
  "expose_headers": ["*"]
}
```

## 完整参数示例：
```json
{
  "name": "cors",
  "config": "{\"allow_origins\":[\"*\"],\"allow_headers\":[\"*\"],\"allow_credentials\":true,\"allow_methods\":[\"*\"],\"expose_headers\":[\"*\"]}"
}
```
