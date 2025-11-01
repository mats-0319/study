package main

import (
	"log"
	"time"

	"github.com/xlzd/gotp"
)

// 程序没有保持开启并一直循环生成验证码，这是为了防止忘记关闭。
// 我们认为生成的第一个验证码剩余有效期过短导致来不及输入是可能的，
// 但如果连续两个验证码都没能成功输入，那这就属于用户操作失误了，程序不予处理。
func main() {
	otp := gotp.NewDefaultTOTP("<< Your github 2FA secret >>")

	log.Printf("> TOTP: %s, Remaining validity time: %d\n", otp.Now(), 30-(time.Now().Unix())%30)

	time.Sleep(time.Second * time.Duration(30-(time.Now().Unix())%30))

	log.Printf("> TOTP: %s, Remaining validity time: %d\n", otp.Now(), 30-(time.Now().Unix())%30)

	time.Sleep(time.Second * 30)
}
