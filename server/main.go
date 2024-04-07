package main

import (
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	host := flag.String("h", "127.0.0.1:8123", "host string (can be an IP or a domain)")
	username := flag.String("u", "", "basic auth user")
	password := flag.String("p", "", "basic auth password")
	dest := flag.String("f", "pics/", " folder for pictures")
	flag.Parse()

	// check if
	ip, port, err := net.SplitHostPort(*host)
	runsLocal := err == nil && ip != "" && port != ""

	// enforce basic auth when not running locally
	if !runsLocal && (*username == "" || *password == "") {
		println("Missing username or password")
		flag.PrintDefaults()
		os.Exit(1)
	}

	r := gin.Default()

	if !runsLocal {
		r.Use(gin.BasicAuth(gin.Accounts{*username: *password}))
	}

	// server last picure on TLD
	r.GET("/", func(c *gin.Context) {
		c.File(*dest + "/last.jpg")
	})

	// list gives a list of the last 50 pictures
	r.GET("/list", func(c *gin.Context) {
		list := "<html><body><ul>"
		files, err := os.ReadDir(*dest)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading directory")
			return
		}
		for i := len(files) - 1; i >= 0 && i >= len(files)-50; i-- {
			if files[i].IsDir() {
				continue
			}
			list += "<li><a href=/pics/" + files[i].Name() + ">" + files[i].Name() + "</a></li>"
		}
		list += "</ul></body></html>"
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(list))
	})

	// serve the static pictures
	r.Static("/pics", *dest)

	// receive a picture
	r.POST("/send", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		c.SaveUploadedFile(file, *dest+file.Filename)
		c.SaveUploadedFile(file, *dest+"last.jpg")
		c.String(http.StatusOK, "ok")
	})

	if runsLocal {
		r.Run(*host)
	} else {
		autotls.Run(r, *host)
	}
}
