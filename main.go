package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	fileMap = make(map[string]string)
)

type certFile struct {
	fileName string
}

// POST /api/cert
func postCert(c echo.Context) error {
	var cFile certFile

	if err := c.Bind(&cFile); err != nil {
		fmt.Println("Error binding received data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process content"}
	}

	reg, err := regexp.Compile("([^a-zA-Z0-9\\.\\-]+)")
	if err != nil {
		fmt.Println(err)
	}
	cFile.fileName = reg.ReplaceAllString(cFile.fileName, "")

	// check if file exists in list of files found in shared EmptyDir vol
	if _, err := os.Stat("/cert/" + cFile.fileName); err == nil {
		return c.Attachment(cFile.fileName, cFile.fileName)
	}
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to find file"}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/cert", postCert)

	//e.Logger.Info(e.Start(":8080"))
	e.Logger.Info(e.StartTLS(":8443", "tls.crt", "tls.key"))
}
