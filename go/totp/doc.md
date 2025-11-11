# 基于时间的一次性密码

Time-based One-Time Password

是一种根据预共享的密钥与当前时间，计算带有时效性的密码的算法。  
常用于2FA（Two-Factor Authenticate，二重因素验证）等场景。

## 过程

> 括号内为通常的默认参数

- 生成密钥K，服务端与客户端共享
- 记录协商时间T0（0）和时间步长TI（30秒）
- 协商加密hash算法（HMAC-SHA-1）
- 协商密码长度（6位）

## 对比

OTP -> HOTP -> TOTP

几种算法的区别在于hash函数的原文，OTP的hash原文是随机的、HOTP使用计数器、TOTP使用时间步数

## 生成密码

> 示例代码：`./demo.go`，dart代码见`github.com/mats0319/totp/lib/dart/totp.dart`

- OTP(K,C) = Truncate(HMAC-SHA-1(K,C)) -> C是随机数，T是截取函数
- step 1：sha-1得到20 Bytes，取最后一个Byte的后4 bits作为起始索引（起始索引属于0～15）
- step 2：从起始索引开始，往后获取总计4个Bytes，将其后31 bits转换为十进制数字作为随机长密码
- step 3：根据需要的位数d，取长密码的后d位，长度不足的在高位补0

## 问题

- 因为网络延迟、时钟偏移、用户延误等原因，有些服务器也接受使用上一个时间步长或下一个时间步长生成的密码
- 实际上TOTP密码的有效期通常是2倍或更多，这是基于上一条提到的原因的妥协
- TOTP没有限制重试次数，所以服务端程序需要处理用户多次失败的问题
- 2的31次方转换成十进制数是21开头的10位数，也就是说现有算法只能计算不多于10位的密码，且计算10位密码时，首位不会出现5～9

## 参考资料

[OTP算法步骤详解](https://notes.mengxin.science/2017/05/30/hotp-totp-algorithm-analysis/)  
[RFC 4226（HOTP）](https://datatracker.ietf.org/doc/html/rfc4226)  
[RFC 6238（TOTP）](https://datatracker.ietf.org/doc/html/rfc6238)
