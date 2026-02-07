package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// utils.ReturnErrorResponse(c, http.StatusUnauthorized, "Token Error", "Token Required")
			c.Abort()
			return
		}

		// Check Bearer scheme, tokenString := parts[1]
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// utils.ReturnErrorResponse(c, http.StatusUnauthorized, "Token Error", "Invalid authorization format")
			c.Abort()
			return
		}
		// actualToken := parts[1]

		// Parse and validate token
		// claims, err := pkg.ParseToken(actualToken, &pkg.JWTAccessClaims{})
		// if err != nil {
		// 	utils.ReturnErrorResponse(c, http.StatusUnauthorized, "Token Error", "Invalid token claims")
		// 	c.Abort()
		// 	return
		// }

		// validate token
		// if err = pkg.ValidateClaim(claims); err != nil {
		// 	utils.ReturnErrorResponse(c, http.StatusUnauthorized, "Token Error", "Invalid token claims")
		// 	c.Abort()
		// 	return
		// }

		// Set user information in context
		// if ac, ok := claims.(*pkg.JWTAccessClaims); ok {
		// c.Set(constant.CtxUserID, ac.UserID)
		// c.Set(constant.CtxUserEmail, ac.Email)
		// c.Set(constant.CtxUsername, ac.Username)
		// c.Set(constant.CtxJWTAccessToken, actualToken)
		// decline if can't find they key in redis (expired or logout)
		// var clientInfo string
		// clientInfoTmp, exists := c.Get(constant.CtxDeviceInfo)
		// clientInfo, ok := clientInfoTmp.(string)
		// if !exists || !ok {
		// 	clientInfo = utils.RedisJWTTokenKey(constant.UnknownIpStr, constant.UnknownStr)
		// }
		// query := redisCache.NewRedisQuerier(q)
		// if _, err := query.GetUserToken(c.Request.Context(), utils.RedisJWTTokenKey(clientInfo, ac.UserID)); err != nil {
		// 	q.Log.Warn().Msgf("can't find user: %v, with client info: %v", ac.UserID, clientInfo)
		// 	utils.ReturnErrorResponse(c, http.StatusUnauthorized, "Token Error", "Token not found")
		// 	c.Abort()
		// 	return
		// }
		// }

		c.Next()
	}
}
