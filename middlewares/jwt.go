package middlewares

//TokenValid  验证token
//func TokenValid() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		token := ctx.GetHeader("Authorization")
//		jwtKey := viper.GetString("adminJwt.key")
//		userInfo := models.AdminUserClaim{}
//		getToken, err := jwt.ParseWithClaims(token, &userInfo, func(token *jwt.Token) (interface{}, error) {
//			return []byte(jwtKey), nil
//		})
//
//		if getToken != nil && getToken.Valid { //token验证通过
//			ctx.Set("AdminUser", getToken.Claims) //在具体业务中可以使用：c.MustGet("AdminUser").(*models.AdminUserClaim).UserName 获取具体的值
//			ctx.Next()
//		} else if ve, ok := err.(*jwt.ValidationError); ok {
//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "错误的token", "code": 100, "data": ""})
//			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "token过期或未启用", "code": 100, "data": ""})
//			} else {
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Couldn't handle this token:" + err.Error(), "code": 100, "data": ""})
//			}
//		} else {
//			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "无法解析此token" + err.Error(), "code": 100, "data": ""})
//		}
//	}
//}
