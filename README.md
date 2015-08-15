# cofile
支持客户端断点上传的文件服务器


上传请求:
---------

* header

|   字段名   |   类型   |   说明   |
| ---------- | -------- | -------- |
| ver | string | 请求版本号 |
| filetype | string | 文件类型 |
| filemd5  | string | 文件MD5(16进制，小写，32B) |
| filesize | int64 | 文件总大小 |
| chunksize | int64 |   数据块大小 |
| offset | int64 | 该数据块在整个文件中的偏移位置(zero-based) |


* body -- `chunk`


上传回复:
---------
* header

|   字段名   |   类型   |   说明   |
| ---------- | -------- | -------- |
| ver | string | 回复版本号 |
| offset | int64 | 期望下次收到的数据块的偏移位置(zero-based) |
| retcode | int64 | 0：继续上传，1：上传成功，2：上传失败 |
| errmsg | string | 错误描述信息 |

