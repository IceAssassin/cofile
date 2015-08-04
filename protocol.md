上传请求:
---------------------------------------------------
[header]

字段名       类型      说明
ver          int       请求版本号
filetype     string    文件类型
filemd5      string    文件MD5(16进制，小写，32B)
filesize     unsigned  文件总大小
rangestart   unsigned  上传文件起始位置下标
rangeend     unsigned  上传文件结束位置下标

---------------------------------------------------
[body]

数据块，长度为: rangeend - rangestart + 1
---------------------------------------------------



上传回复:
---------------------------------------------------
[header]

字段名       类型      说明
ver          int       回复版本号
flag         int       0：继续上传，1：上传成功，2：上传失败
rangestart   unsigned  上传文件起始位置下标
rangeend     unsigned  上传文件结束位置下标
errmsg       string    错误描述信息
---------------------------------------------------

