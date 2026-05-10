package middlewares

// import (
// 	"auth-service/internal/utils"
// 	"context"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// type ctxKey string

// const userCtxKey ctxKey = "claims"

// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
// 			return
// 		}
// 		parts := strings.SplitN(authHeader, " ", 2)
// 		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
// 			return
// 		}
// 		tokenStr := parts[1]
// 		claims, err := utils.ParseAccessToken(tokenStr)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
// 			return
// 		}
// 		ctx := context.WithValue(c.Request.Context(), userCtxKey, claims)
// 		c.Request = c.Request.WithContext(ctx)
// 		c.Set("claims", claims)
// 		c.Next()
// 	}
// }

// func GetClaimsFromContext(c *gin.Context) (*utils.Claims, bool) {
// 	v := c.Request.Context().Value(userCtxKey)
// 	if v == nil {
// 		return nil, false
// 	}
// 	claims, ok := v.(*utils.Claims)
// 	return claims, ok
// }

// func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		claims, ok := GetClaimsFromContext(c)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no claims"})
// 			return
// 		}
// 		for _, r := range allowedRoles {
// 			if claims.Role == r {
// 				c.Next()
// 				return
// 			}
// 		}
// 		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
// 	}
// }
