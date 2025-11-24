# JWT学习笔记

JWT：json web token，是一种客户端与服务端之间确认身份的方案

服务端在通过验证后（例如验证密码），向客户端发送一组信息，其形如：

```JWT payload
{
  "name": "mario",
  "identify": "admin",
  "expire_time": "xxx"
}
```

客户端每次向服务端发送请求时携带该信息，服务端认为该信息中的内容是可信的

## JWT的结构

e.g.
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

- Header：base64url编码，通常包含算法（S部分的加密算法）、token类型等信息
- Payload：base64url编码，包含想要传输的实际数据
- Signature：使用H中的算法得到的签名，服务端信任payload的基础

## JWT的使用

- 服务端生成P，根据H、P计算S（服务端持有密钥），整体拼好传给客户端
- 客户端每次请求时带上该token（放在HTTP请求的头信息Authorization字段）
- 服务端拿到token并验证，通过验证后可以认为P部分的内容是可信的

## base64url编码

与base64编码相似，只是将其中几个在url中有特殊含义的字符替换成其他字符：

- '+' -> '-'
- '/' -> '_'
- 不使用'='在末尾填充
