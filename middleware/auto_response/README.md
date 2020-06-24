# auto_response

根据客户端请求路径自动响应状态码以及内容，可用于自动处理`OPTIONS`方法以及`/favicon.ico`路径的请求

示例：
```json
{
    "middleware": [
        {
            "name": "auto_response",
            "config": "{\"*\":{\"method\":\"OPTIONS\",\"status\":204},\"/favicon.ico\":{\"method\":\"GET\",\"status\":204}}"
        }
    ]
}
```
