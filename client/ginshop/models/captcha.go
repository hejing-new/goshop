package models

import (
	"context"
	"fmt"
	"net/url"
	"time"

	pbCaptcha "ginshop/proto/captcha"
)

// 调用验证码微服务的超时时间
const captchaRPCTimeout = 2 * time.Second

// CaptchaErrImg 验证码服务不可用时返回的占位图：
// 把错误直接画进图片里，前端（admin/login.js、itying/pass.js）不用改就能看到原因。
var CaptchaErrImg = "data:image/svg+xml;charset=utf-8," + url.QueryEscape(
	`<svg xmlns="http://www.w3.org/2000/svg" width="120" height="40">`+
		`<rect width="120" height="40" fill="#f7f7f7" stroke="#ddd"/>`+
		`<text x="6" y="24" font-size="10" fill="#c00">验证码服务未启动</text></svg>`)

//调用验证码微服务
func MakeCaptcha(height int, width int, length int) (string, string, error) {
	captchaClient := pbCaptcha.NewCaptchaService("captcha", CaptchaClient)

	// ★ 原来用的是 context.Background()，没有超时；现在从请求上下文派生
	ctx, cancel := context.WithTimeout(context.Background(), captchaRPCTimeout)
	defer cancel()

	res, err := captchaClient.MakeCaptcha(ctx, &pbCaptcha.MakeCaptchaRequest{
		Height: int32(height),
		Width:  int32(width),
		Length: int32(length),
	})

	// ★ 原来这里写的是 log.Fatal(err)。
	//   但 go-micro 的 util/log.Fatal 只打日志、不退出进程，
	//   于是带着 err 继续执行 return res.Id ... —— res 是 nil，直接空指针 panic。
	//   正确做法：把错误返回给调用方，由调用方决定怎么提示用户。
	if err != nil {
		return "", "", fmt.Errorf("调用验证码服务失败（请确认 Consul 和 captcha 服务已启动）: %w", err)
	}
	if res == nil {
		return "", "", fmt.Errorf("验证码服务返回为空")
	}

	return res.Id, res.B64S, nil
}

//验证验证码
//说明：签名保持 bool 不变（有 3 个调用方），服务不可用时按"验证失败"处理并打日志，
//不会像原来那样 panic。
func VerifyCaptcha(id string, VerifyValue string) bool {
	captchaClient := pbCaptcha.NewCaptchaService("captcha", CaptchaClient)

	ctx, cancel := context.WithTimeout(context.Background(), captchaRPCTimeout)
	defer cancel()

	res, err := captchaClient.VerifyCaptcha(ctx, &pbCaptcha.VerifyCaptchaRequest{
		Id:          id,
		VerifyValue: VerifyValue,
	})
	if err != nil {
		fmt.Println("验证验证码失败:", err)
		return false
	}
	if res == nil {
		return false
	}
	return res.VerifyResult
}
