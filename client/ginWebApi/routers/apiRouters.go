package routers

import (
	"ginWebApi/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiRoutersInit(r *gin.Engine) {
	apiRouters := r.Group("/api")
	{
		apiRouters.GET("/", api.ApiController{}.Index)
		apiRouters.GET("/makeCaptcha", api.CaptchaController{}.MakeCaptcha)
		apiRouters.POST("/verifyCaptcha", api.CaptchaController{}.VerifyCaptcha)

	}

}
