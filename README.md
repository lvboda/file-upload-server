# upload-server
一个简单的文件上传服务。

## 配置

``` toml
[server]
    mode          = "release"
    port          = ":3002"
    token         = ""
    storagePath   = "../files"
    isHttps       = false
    returnBaseUrl = "https://lvboda.cn/"
```

### mode

可选值：debug | release | test

对应gin框架的mode，没有特殊情况默认release就可以。 [了解更多](https://gin-gonic.com/zh-cn/)

### port

服务运行的端口，前面需要带:，不然会跑不起来。

### token

相当于一个固定密码，每次请求与header里的Authorization值进行对比，不一致会调用失败。

### storagePath

文件在服务器上的存储地址，这个路径相对于可执行文件main，如果路径里出现不存在的文件夹需手动创建。

### isHttps

是否开启https，如果为true则在config文件夹里准备tls.key与tls.pem私钥与公钥文件，文件名固定。

### returnBaseUrl

返回的基础url，接口返回值url格式为：returnBaseUrl + 文件名，returnBaseUrl结尾要有'/'。

## 接口

### 上传

请求地址：/upload

请求方法：POST

请求头：{ Authorization: 服务端配置文件里token的值, Content-Type: "multipart/form-data" }

请求体: formData：{ file: 上传的文件 }

返回值: { "code": 200, "data": 文件url }

### 删除

请求地址：/delete/文件名称

请求方法：DELETE

请求头：{ Authorization: 服务端配置文件里token的值 }

返回值: { "code": 200, "data": "ok" }

### 查看

请求地址：/static/文件名称

请求方法：GET

返回值: 文件

## 注意事项

1.发生错误会直接返回400状态码。

2.修改config文件后需重启服务才会生效。

3.文件最好在上传前做一下hash或者随机数重命名处理，服务端会直接用上传的文件名，如果发生名称相同的情况后者会覆盖前者。

## License

[MIT](https://github.com/lvboda/upload-server/blob/master/LICENSE)

Copyright (c) 2022 Boda Lü
