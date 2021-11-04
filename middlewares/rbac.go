package middlewares

//Rbac  验证权限
//func Rbac() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		//用户id
//		userInfo, _ := ctx.Get("UserInfo")
//		userId := userInfo.(*models.AdminUserClaim).UserId
//
//		//rbac的租户
//		domain := "admin"
//
//		//获取请求的URI
//		obj := ctx.Request.URL.RequestURI()
//
//		//获取请求方法
//		act := ctx.Request.Method
//
//		//检测是否具有该执行权限
//		isCan, _ := core.Enforcer.Enforce(userId, domain, obj, act)
//		if isCan {
//			ctx.Next()
//		} else {
//			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "权限不足~", "code": 0, "data": ""})
//		}
//	}
//}
