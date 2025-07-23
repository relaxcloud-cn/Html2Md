// Package api provides HTTP API for HTML to Markdown conversion service
// @title HTML2Markdown API
// @version 1.0
// @description HTML转Markdown转换服务的REST API
// @termsOfService https://github.com/relaxcloud-cn/html2md
// @contact.name API Support
// @contact.url https://github.com/relaxcloud-cn/html2md/issues
// @contact.email support@example.com
// @license.name MIT
// @license.url https://github.com/relaxcloud-cn/html2md/blob/main/LICENSE
// @host localhost:8080
// @BasePath /
// @schemes http https
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/relaxcloud-cn/html2md/api/http/handler"
	"github.com/relaxcloud-cn/html2md/api/http/middleware"
	"github.com/relaxcloud-cn/html2md/internal/model"
	"github.com/relaxcloud-cn/html2md/internal/service"
)

// NewRouter 创建HTTP路由器
func NewRouter() *gin.Engine {
	// 在生产环境中设置为release模式
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 创建服务
	convertService := service.NewConvertService()
	convertHandler := handler.NewConvertHandler(convertService)

	// 首页重定向到Swagger文档
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	// Swagger文档
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 转换相关接口
		v1.POST("/convert", convertHandler.Convert)
		v1.POST("/convert/batch", convertHandler.ConvertBatch)
		v1.GET("/convert/simple", convertHandler.ConvertSimple)

		// 系统接口
		v1.GET("/health", convertHandler.Health)
		v1.GET("/info", convertHandler.GetConverterInfo)

		// 演示接口
		v1.GET("/demo", demoHandler)
	}

	// 错误处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.NewErrorResponse(
			model.CodeNotFound,
			"接口不存在",
			nil,
		))
	})

	return r
}

// demoHandler 演示处理器
// @Summary 演示接口
// @Description 提供一个简单的演示页面，展示API的使用方法
// @Tags 演示
// @Produce html
// @Success 200 {string} string "HTML演示页面"
// @Router /api/v1/demo [get]
func demoHandler(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HTML转Markdown API演示</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .container { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        textarea { width: 100%; height: 150px; margin: 10px 0; }
        button { background: #007cba; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #005a8b; }
        .result { background: white; padding: 15px; border-left: 4px solid #007cba; margin: 10px 0; }
        .error { border-left-color: #dc3545; background: #f8d7da; }
    </style>
</head>
<body>
    <h1>HTML转Markdown API演示</h1>
    
    <div class="container">
        <h3>在线转换</h3>
        <textarea id="htmlInput" placeholder="请输入HTML内容...">
<h1>示例标题</h1>
<p>这是一个<strong>粗体</strong>和<em>斜体</em>的示例段落。</p>
<ul>
<li>列表项 1</li>
<li>列表项 2</li>
</ul>
<a href="https://example.com">示例链接</a>
        </textarea>
        <br>
        <button onclick="convertHtml()">转换为Markdown</button>
        <div id="result"></div>
    </div>
    
    <div class="container">
        <h3>API端点</h3>
        <ul>
            <li><strong>POST /api/v1/convert</strong> - 转换HTML为Markdown</li>
            <li><strong>GET /api/v1/convert/simple</strong> - 简单转换（GET方式）</li>
            <li><strong>POST /api/v1/convert/batch</strong> - 批量转换</li>
            <li><strong>GET /api/v1/health</strong> - 健康检查</li>
            <li><strong>GET /api/v1/info</strong> - 转换器信息</li>
            <li><strong>GET /docs/index.html</strong> - Swagger API文档</li>
        </ul>
    </div>

    <script>
        async function convertHtml() {
            const html = document.getElementById('htmlInput').value;
            const resultDiv = document.getElementById('result');
            
            if (!html.trim()) {
                resultDiv.innerHTML = '<div class="result error">请输入HTML内容</div>';
                return;
            }
            
            try {
                const response = await fetch('/api/v1/convert', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        html: html,
                        plugins: ['commonmark']
                    })
                });
                
                const data = await response.json();
                
                if (data.code === 200) {
                    resultDiv.innerHTML = '<div class="result"><strong>转换结果：</strong><br><pre>' + 
                        escapeHtml(data.data.markdown) + '</pre></div>';
                } else {
                    resultDiv.innerHTML = '<div class="result error"><strong>转换失败：</strong>' + 
                        data.msg + '</div>';
                }
            } catch (error) {
                resultDiv.innerHTML = '<div class="result error"><strong>请求失败：</strong>' + 
                    error.message + '</div>';
            }
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    </script>
</body>
</html>
	`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
