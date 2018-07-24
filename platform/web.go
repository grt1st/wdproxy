package platform

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/multitemplate"
)

func RunWebPlatform() {
	router := gin.Default()

	// gin模版配置
	router.Static("/static", "platform/static")
	router.LoadHTMLGlob("platform/templates/*")
	router.HTMLRender = createRender()

	router.GET("/", Index)
	router.GET("/login", Login)
	router.POST("/login", Login)
	router.GET("/settings", Settings)
	router.GET("/dashboard", Dashboard)
	router.GET("/dashboard/detail", Detail)
	router.Run("127.0.0.1:3000")

}

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "platform/templates/base.html", "platform/templates/index.html")
	r.AddFromFiles("settings", "platform/templates/base.html", "platform/templates/settings.html")
	r.AddFromFiles("login", "platform/templates/base.html", "platform/templates/login.html")
	r.AddFromFiles("dashboard", "platform/templates/base.html", "platform/templates/dashboard.html")
	r.AddFromFiles("detail", "platform/templates/base.html", "platform/templates/detail.html")
	return r
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "heatmap",
		"home": "active",
	})
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "heatmap",
		"home": "active",
	})
}

func Settings(c *gin.Context) {
	c.HTML(http.StatusOK, "settings", gin.H{
		"title": "heatmap",
		"settings": "active",
	})
}

func Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard", gin.H{
		"title": "heatmap",
		"dashboard": "active",
	})
}

func Detail(c *gin.Context) {
	c.HTML(http.StatusOK, "detail", gin.H{
		"title": "heatmap",
		"dashboard": "active",
	})
}