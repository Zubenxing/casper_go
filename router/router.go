package router

import (
	"casper_go/config"
	"casper_go/core/middleware"
	"casper_go/features/auth"
	"casper_go/features/certificate"
	"casper_go/features/password"
	"casper_go/features/work"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Casper Platform API
// @version 1.0
// @description 自主化监控平台 API 文档
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer 你的token"（注意 Bearer 后有空格）

// Setup 设置路由
func Setup(mode string, cfg *config.Config) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()

	// 全局中间件
	r.Use(middleware.CORS())

	// 请求ID追踪（如果启用）
	if cfg.API.RequestIDEnabled {
		r.Use(middleware.RequestID())
	}

	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 限流中间件（如果启用）
	if cfg.API.RateLimitEnabled {
		middleware.InitRateLimiter(cfg.API.RateLimitRequests)
		r.Use(middleware.RateLimit())
	}

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Casper Platform is running",
		})
	})

	// API 路由组（支持版本控制）
	apiPath := "/api"
	if cfg.API.EnableVersion {
		apiPath = "/api/" + cfg.API.CurrentVersion
	}
	api := r.Group(apiPath)
	{
		// 认证相关路由（不需要认证）
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", auth.LoginAPI)
			authGroup.POST("/refresh", auth.RefreshTokenAPI)
		}

		// 证书工具路由（公开访问，不需要认证）
		certTools := api.Group("/certificates/tools")
		{
			certTools.POST("/generate-csr", certificate.GenerateCSRAPI)
			certTools.POST("/validate-csr", certificate.ValidateCSRAPI)
			certTools.POST("/validate-cert", certificate.ValidateCertificateAPI)
		}

		// 文件服务（公开访问，不需要认证 - 允许浏览器直接加载图片）
		api.GET("/files/*filename", work.ServeFile)

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// 用户信息
			authenticated.GET("/profile", auth.GetProfileAPI)
			authenticated.POST("/logout", auth.LogoutAPI)

			// 证书监控路由
			cert := authenticated.Group("/certificates")
			{
				cert.POST("", certificate.AddAPI)
				cert.GET("", certificate.GetListAPI)
				cert.GET("/:id", certificate.GetDetailAPI)
				cert.PUT("/:id", certificate.UpdateAPI)
				cert.PATCH("/:id/info", certificate.UpdateInfoAPI)
				cert.DELETE("/:id", certificate.DeleteAPI)
				cert.POST("/check-all", certificate.CheckAllAPI)
			}

			// 密码管理路由
			passwords := authenticated.Group("/passwords")
			{
				passwords.POST("", password.CreateAPI)
				passwords.GET("", password.GetListAPI)
				passwords.GET("/categories", password.GetCategoriesAPI)
				passwords.GET("/stats", password.GetStatsAPI)
				passwords.GET("/:id", password.GetDetailAPI)
				passwords.GET("/:id/password", password.GetPasswordAPI)
				passwords.PUT("/:id", password.UpdateAPI)
				passwords.DELETE("/:id", password.DeleteAPI)
			}

			// 工作记录路由
			works := authenticated.Group("/works")
			{
				works.POST("", work.CreateWorkAPI)
				works.GET("", work.GetWorkListAPI)
				works.GET("/stats", work.GetWorkStatsAPI)
				works.GET("/:id", work.GetWorkDetailAPI)
				works.PUT("/:id", work.UpdateWorkAPI)
				works.DELETE("/:id", work.DeleteWorkAPI)
			}

			// 问题记录路由
			issues := authenticated.Group("/work-issues")
			{
				issues.POST("", work.CreateWorkIssueAPI)
				issues.GET("", work.GetWorkIssueListAPI)
				issues.GET("/stats", work.GetWorkIssueStatsAPI)
				issues.GET("/:id", work.GetWorkIssueDetailAPI)
				issues.PUT("/:id", work.UpdateWorkIssueAPI)
				issues.DELETE("/:id", work.DeleteWorkIssueAPI)
				issues.POST("/upload", work.UploadImage) // 上传问题截图
			}
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			admin.GET("/status", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Admin area",
				})
			})
		}
	}

	return r
}
