package middleware

import (
	"log"
	"time"

	"github.com/Heath000/fzuSE2024/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var authMiddleware *jwt.GinJWTMiddleware
var identityKey = "email"

// Login struct
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

// Auth middleware 返回 JWT 中间件实例
func Auth() *jwt.GinJWTMiddleware {
	return authMiddleware
}

func init() {
	var err error
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "gin-skeleton",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		SendCookie:  true,
		PayloadFunc: func(data any) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
					"ID":        v.ID,
					"name":      v.Name,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) any {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				ID:    uint(claims["ID"].(float64)),
				Email: claims[identityKey].(string),
				Name:  claims["name"].(string),
			}
		},

		// Authenticator: 用户登录验证
		Authenticator: func(c *gin.Context) (any, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password
			return model.LoginByEmailAndPassword(email, password)
		},

		// Authorizator: 权限控制
		Authorizator: func(data any, c *gin.Context) bool {
			if user, ok := data.(*model.User); ok {
				// 如果是 /admin 路由或者以 /admin 开头的路由，检查是否为管理员
				if c.FullPath() == "/admin" || c.FullPath()[:6] == "/admin" {
					if user.Email != "admin" {
						return false // 非管理员用户不允许访问 /admin 路由
					}
				}
				// 其他路由，所有已登录用户都可以访问
				return true
			}
			return false
		},

		// 未授权的响应
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		// TokenLookup 配置：从请求中查找 Token 的位置
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}
