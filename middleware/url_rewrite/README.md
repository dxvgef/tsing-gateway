# url_rewrite

网关转发客户端的请求到上游时，重写请求的URL路径，目前支持以下重写规则：

- prefix 重写前缀
- suffix 重写后缀
- replace 替换重写

## 配置参数说明
- `name` 中间件的名称，必须是`url_rewrite`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `prefix` key/value类型，可选，要重写前缀的URL，key是目标字符串，value是重写后的字符串
  - `suffix` key/value类型，可选，要重写后缀的URL，key是目标字符串，value是重写后的字符串
  - `replace` key/value类型，可选，要替换重写的URL，key是目标字符串，value是重写后的字符串

## `config`字段示列
```json
{
  "prefix": {
    "/test/": "/"
  },
  "suffix": {
    "/test/": "/"
  },
  "replace": {
    "/test/": "/"
  }
}
```

## 完整参数示例：

```json
{
  "name": "url_rewrite",
  "config": "{\"prefix\":{\"/test/\":\"/\"},\"suffix\":{\"/test/\":\"/\"},\"replace\":{\"/test/\":\"/\"}}"
}
```
