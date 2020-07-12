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
	path    = os.Getenv("TLS_FILE_PATH")
)

type certFile struct {
	FileName string `json:"FileName"`
}

// POST /api/cert
func postCert(c echo.Context) error {
	var cFile certFile

	fmt.Println("data received")
	if err := c.Bind(&cFile); err != nil {
		fmt.Println("Error binding received data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process content"}
	}

	reg, err := regexp.Compile("([^a-zA-Z0-9\\.\\-]+)")
	if err != nil {
		fmt.Println(err)
	}
	cFile.FileName = reg.ReplaceAllString(cFile.FileName, "")

	fmt.Printf("searching %v for %v\n", path, cFile.FileName)
	// check if file exists in list of files found in shared EmptyDir vol
	if _, err := os.Stat(path + cFile.FileName); err != nil {
		fmt.Println("Error finding file:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to find file"}
	}
	return c.Attachment(path+cFile.FileName, path+cFile.FileName)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/cert", postCert)
	e.POST("/api/cert/", postCert)

	fmt.Println("printing port, cert, key", port, cert, key)

	e.Logger.Info(e.StartTLS(port, cert, key))
}
