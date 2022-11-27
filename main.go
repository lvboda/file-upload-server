package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

var conf config
var executeDir = func() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex)
}()

func init() {
	initConfig()
}

func main() {
	gin.SetMode(conf.Server.Mode)
	router := gin.New()
	router.Use(authMiddleware)
	router.Static("/images", conf.Server.StoragePath)
	router.POST("/upload", upload)
	router.DELETE("/delete/:name", delete)
	if conf.Server.IsHttps {
		certFilePath := path.Join(executeDir, "../config/tls.pem")
		keyFilePath := path.Join(executeDir, "../config/tls.key")
		router.RunTLS(conf.Server.Port, certFilePath, keyFilePath)
	} else {
		router.Run(conf.Server.Port)
	}
}

type config struct {
	Server struct {
		Mode          string
		Port          string
		Token         string
		StoragePath   string
		IsHttps       bool
		ReturnBaseUrl string
	}
}

func initConfig() {
	file, err := os.ReadFile(path.Join(executeDir, "../config/config.toml"))
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
}

func authMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != conf.Server.Token {
		c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
		return
	}
	c.Next()
}

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	err = c.SaveUploadedFile(file, path.Join(executeDir, conf.Server.StoragePath, file.Filename))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": conf.Server.ReturnBaseUrl + file.Filename,
	})
}

func delete(c *gin.Context) {
	fileName := c.Param("name")
	err := os.Remove(path.Join(executeDir, conf.Server.StoragePath, fileName))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "ok",
	})
}
