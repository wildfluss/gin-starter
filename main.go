package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	if os.Getenv("DEBUG_ENABLED") == "true" {
		// https://gin-gonic.com/docs/examples/html-rendering/
		r.LoadHTMLFiles("html/index.html")
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
		// https://gin-gonic.com/docs/examples/serving-static-files/
		r.Static("/assets", "./html/assets")
	} else {
		t, err := loadTemplate()
		if err != nil {
			panic(err)
		}
		r.SetHTMLTemplate(t)
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "/html/index.html", nil)
		})
		r.GET("/assets/*filepath", StaticHandler)
	}

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	// log.Println(os.Getenv("DEBUG_ENABLED"))

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		// log.Println("loadTemplate:", name)

		// if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
		// 	continue
		// }

		if strings.HasSuffix(name, ".svg") {
			// log.Println("loadTemplate: skipped:", name)
			continue
		}

		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}

		// log.Println("loadTemplate: parsed:", name)
	}
	return t, nil
}

// XXX
var assetsData = make(map[string][]byte)

func StaticHandler(c *gin.Context) {
	p := c.Param("filepath")
	fixed := strings.Join([]string{"/html/assets", p}, "")
	// log.Println("StaticHandler:", fixed)
	_, ok := assetsData[fixed]
	if !ok {
		_, ok := Assets.Files[fixed]
		if !ok {
			return
		}
		data, err := ioutil.ReadAll(Assets.Files[fixed])
		if err != nil {
			return
		}
		assetsData[fixed] = data
	}
	data := assetsData[fixed]
	// log.Println("StaticHandler: data len:", len(data))
	c.Writer.Header().Set("Content-Type", "image/svg+xml")
	c.Writer.Write(data)
}
