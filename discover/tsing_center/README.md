# tsing_center

从`Tsing Center`中获得端点

## 配置参数说明
- `name` 中间件的名称，必须是`tsing_center`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `addr`，string，必需， `Tsing Center`的HTTP/HTTPS监听地址:端口
  - `secret`，string，可选，连接`Tsing Center`的API时的secret字符串

## `config`字段示列
```json
{
  "addr": "http://127.0.0.1:10080",
  "secret": "123456"
}
```

## 完整参数示例：
```json
{
  "name": "tsing_center",
  "config": "{\"addr\":\"http://127.0.0.1:10080\",\"secret\":\"123456\"}"
}
```
