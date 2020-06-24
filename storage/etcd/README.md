# etcd

基于etcd的数据存储引擎

## 配置参数说明
- `name` 中间件的名称，必须是`url_rewrite`
- `config`字段是`JSON (Object)`格式并压缩并转义后的`string`类型
  - `key_prefix`，string 类型，必需，etcd的键名前缀
  - `endpoints`，[]string 类型，必需，etcd的连接地址
  - `dial_timeout`，uint 类型，可选，etcd的`dial_timeout`参数
  - `username`，string 类型，可选，etcd的`username`参数
  - `password`，string 类型，可选，etcd的`password`参数
  - `auto_sync_interval`，uint 类型，可选，etcd的`auto_sync_interval`参数
  - `dial_keep_alive_time`，uint 类型，可选，etcd的`dial_keep_alive_time`参数
  - `dial_keep_alive_timeout`，uint 类型，可选，etcd的`dial_keep_alive_timeout`参数
  - `max_call_send_msg_size`，uint 类型，可选，etcd的`max_call_send_msg_size`参数
  - `max_call_recv_msg_size`，uint 类型，可选，etcd的`max_call_recv_msg_size`参数
  - `reject_old_cluster`，bool 类型，可选，etcd的`reject_old_cluster`参数
  - `permit_without_stream`，bool 类型，可选，etcd的`permit_without_stream`参数

## `config`字段示列
```json
{
  "key_prefix": "/tsing-gateway",
  "endpoints": ["http://127.0.0.1:2379"]
}
```

## 完整参数示例：

```json
{
  "name": "etcd",
  "config": "{\"key_prefix\":\"/tsing-gateway\",\"endpoints\":[\"http://127.0.0.1:2379\"]}"
}
```
