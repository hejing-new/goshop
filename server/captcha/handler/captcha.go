package handler

import (
	// 注意：如果你的 proto 目录下没有 captcha 文件夹，这里改成 "captcha/proto"
	pb "captcha/proto/captcha"
	"context"
	"fmt"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

// 创建store
var store = base64Captcha.DefaultMemStore

type Captcha struct{}

// 获取验证码的方法
func (e *Captcha) MakeCaptcha(ctx context.Context, req *pb.MakeCaptchaRequest, res *pb.MakeCaptchaResponse) error {
	fmt.Println(req)
	var driver base64Captcha.Driver
	driverString := base64Captcha.DriverString{
		Height:          int(req.Height),
		Width:           int(req.Width),
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          int(req.Length),
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 102,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driver = driverString.ConvertFonts()

	c := base64Captcha.NewCaptcha(driver, store)

	// 修改这里：v1.3.5 版本返回 4 个值，用 _ 忽略第三个 answer
	id, b64s, _, err := c.Generate()

	res.Id = id
	res.B64S = b64s

	return err
}

// 验证码的方法
func (e *Captcha) VerifyCaptcha(ctx context.Context, req *pb.VerifyCaptchaRequest, res *pb.VerifyCaptchaResponse) error {
	fmt.Println(req)
	if store.Verify(req.Id, req.VerifyValue, true) {
		res.VerifyResult = true
	} else {
		res.VerifyResult = false
	}
	return nil
}
