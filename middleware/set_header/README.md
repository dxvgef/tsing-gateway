# set_header

用于在HTTP会话的请求或响应头部中设置Header参数

## 配置参数说明
- `name` 中间件的名称，必须是`set_header`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `request` key/value格式，可选，要写入的请求参数，key是参数名，value是参数值
  - `response` key/value格式，可选，要写入的请求参数，key是参数名，value是参数值

## `config`字段示列
```json
{
  "request": {
    "X-TEST-REQ-1": "test-req-1",
    "X-TEST-REQ-2": "test-req-2"
  },
  "response": {
    "X-TEST-RESP-1": "test-resp-1",
    "X-TEST-RESP-2": "test-resp-2"
  }
}
```

## 完整参数示例：

```json
{
  "name": "set_header",
  "config": "{\"request\":{\"X-TEST-REQ-1\":\"test-req-1\",\"X-TEST-REQ-2\":\"test-req-2\"},\"response\":{\"X-TEST-RESP-1\":\"test-resp-1\",\"X-TEST-RESP-2\":\"test-resp-2\"}}"
}
```
