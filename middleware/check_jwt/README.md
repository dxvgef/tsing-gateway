# check_jwt

网关自动检查JWT，支持Header、GET、POST、Cookie中获取Token，支持本地校验及HTTP远程校验

示例：
```json
{
    "middleware": [
        {
            "name": "check_jwt",
            "config": "{}"
        }
    ]
}
```
