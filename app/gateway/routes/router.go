package routes

import (
	// get pprof
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/lukkaal/GRPC-API/app/gateway/internal"
	"github.com/lukkaal/GRPC-API/app/gateway/middleware"

	// swagger
	docs "github.com/lukkaal/GRPC-API/app/gateway/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()

	// middleware: Cors/ recover
	ginRouter.Use(middleware.Cors(), middleware.ErrorMiddleware())

	// swagger 配置
	docs.SwaggerInfo.BasePath = "/api/v1"
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := ginRouter.Group("/api/v1")
	{
		{
			v1.GET("ping", func(context *gin.Context) {
				context.JSON(200, "success")
			})
			v1.POST("/user/register", internal.UserRegister)
			v1.POST("/user/login", internal.UserLogin)
			v1.POST("/user/logout", internal.UserLogout)

			task := v1.Group("/")
			task.Use(middleware.JWT())
			{
				// 任务模块
				task.GET("task", internal.GetTaskList)
				task.POST("task", internal.CreateTask)
				task.PUT("task", internal.UpdateTask)
				task.DELETE("task", internal.DeleteTask)
			}
		}
	}

	pprofGroup := ginRouter.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile)) // cpu
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine"))) // goroutine
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))           // memory
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	return ginRouter
}

// GET /api/v1/ping
// POST /api/v1/user/register
// POST /api/v1/user/login

// DELETE PUT POST GET
// /api/v1/task
