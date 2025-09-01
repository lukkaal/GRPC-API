swag 工具负责生成 docs（OpenAPI 规范）

###### **gin-swagger** 负责把这些 docs 渲染成 Swagger UI 页面



swag init

会生成一个 docs/ 目录 里面有 docs.go 保存了 OpenAPI 规范的内容

同时 应该在每一个 handler 写好注释的 tag 比如:



// @Summary      Get all tasks

// @Description  获取所有任务列表

// @Tags         tasks

// @Produce      json

// @Success      200  {array}   Task

// @Router       /tasks \[get]

func GetTasks(c \*gin.Context) {

&nbsp;   tasks := \[]Task{

&nbsp;       {ID: 1, Name: "task1"},

&nbsp;       {ID: 2, Name: "task2"},

&nbsp;   }

&nbsp;   c.JSON(http.StatusOK, tasks)

}



以及 router 的最开始



// @title           Task API

// @version         1.0

// @description     这是一个示例的任务 API

// @host            localhost:8080

// @BasePath        /

func main() {

&nbsp;   r := gin.Default()

&nbsp;   r.GET("/tasks", GetTasks)



&nbsp;   // 注册 Swagger

&nbsp;   r.GET("/swagger/\*any", 

ginSwagger.WrapHandler(swaggerFiles.Handler))



&nbsp;   r.Run(":8080")

}



其中 

1）docs.go

Go 语言源码文件。

把 OpenAPI 规范转成 Go 代码形式的文档描述 

包含 swagger 的元数据（title、version、host 等）和所有 API 信息



2）swagger.json

纯 JSON 格式的 OpenAPI 规范文档。

内容和 docs.go 表达的一致



3）swagger.yaml

纯 YAML 格式的 OpenAPI 规范文档。

和 swagger.json 等价



执行步骤是：引入 gin-swagger 包
1）首先新建 docs 文件夹 最好在 gateway 下和 router 同级的位置

2）在 router 中引入：

r.GET("/swagger/\*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

注册一个 swagger 的访问网关 

3）在 router.go 中 引入 docs 同时在 package main 当中的 main 函数之前写上 tag

以及在每一个 handler 函数头顶写上 tag

4）然后在有 main.go (拥有全局 tag 的文件所在目录) 执行 swag init

就会将说明文档 yaml/json 写入到指定的 docs

5）使用 web 网页即可访问

