package platform

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/multitemplate"
	"github.com/grt1st/wdproxy/g"
	"time"
	"fmt"
	"encoding/json"
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
	var numDomains, numRequests, numTags int
	g.DB.Model(g.WdproxyRecord{}).Count(&numRequests)
	g.DB.Model(g.WdproxyDomain{}).Count(&numDomains)
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "heatmap",
		"home": "active",
		"num_domains": numDomains,
		"num_requests": numRequests,
		"num_tags": numTags,
	})
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "heatmap",
		"home": "active",
	})
}

func Settings(c *gin.Context) {
	var num24h, numAll int
	t := time.Now()
	g.DB.Model(g.WdproxyRecord{}).Count(&numAll)
	g.DB.Model(g.WdproxyRecord{}).Where("created_at > ?", t.AddDate(0,0,-1)).Count(&num24h)
	c.HTML(http.StatusOK, "settings", gin.H{
		"title": "heatmap",
		"settings": "active",
		"check": "yes",
		"num_24h": num24h,
		"num_all": numAll,
	})
}

func Dashboard(c *gin.Context) {
	var records []g.WdproxyRecord
	g.DB.Model(g.WdproxyRecord{}).Find(&records)
	c.HTML(http.StatusOK, "dashboard", gin.H{
		"title": "heatmap",
		"dashboard": "active",
		"records": records,
	})
}

type WdproxyRecord struct {
	g.WdproxyRecord
	HeaderRequestJson   map[string]interface{}
	HeaderResponseJson  map[string]interface{}
}

func Detail(c *gin.Context) {
	id := c.Query("id")
	var record WdproxyRecord
	err := g.DB.Where("id = ?", id).First(&record).Error
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(record.HeaderRequest), &record.HeaderRequestJson)
	err = json.Unmarshal([]byte(record.HeaderResponse), &record.HeaderResponseJson)
	c.HTML(http.StatusOK, "detail", gin.H{
		"title": "heatmap",
		"dashboard": "active",
		"record": record,
	})
}