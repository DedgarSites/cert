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
	port    = ":" + os.Getenv("TLS_PORT")
	cert    = os.Getenv("SSCS_CERT")
	key     = os.Getenv("SSCS_KEY")
)

type certFile struct {
	FileName string `json:"FileName"`
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
	cFile.FileName = reg.ReplaceAllString(cFile.FileName, "")

	// check if file exists in list of files found in shared EmptyDir vol
	if _, err := os.Stat("/cert/certificates/" + cFile.FileName); err == nil {
		return c.Attachment(cFile.FileName, cFile.FileName)
	}
	return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to find file"}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/cert", postCert)

	//e.Logger.Info(e.Start(":8080"))
	fmt.Println("printing port, cert, key", port, cert, key)

	e.Logger.Info(e.StartTLS(port, cert, key))
	fmt.Println("passed the start")
}
