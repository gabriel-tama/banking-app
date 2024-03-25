package middleware

import (
	"net/http"

	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gabriel-tama/banking-app/common/response"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		res := &response.ResponseBody{
			Message: "unauthorized",
		}

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		const BEARER_SCHEMA = "BEARER "
		tokenString := authHeader[len(BEARER_SCHEMA):]
		_, err := jwtService.ValidateToken(tokenString)

		if err != nil {
			res.Message = "invalid token"
			c.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}
	}
}
