package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//func UploadSingleIndex(ctx *gin.Context) {
func UploadSingleIndex(apiServer *APIServer, ginContextPointer *gin.Context) {
	file, err := ginContextPointer.FormFile("file")
	if err != nil {
		ginContextPointer.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	ginContextPointer.JSON(http.StatusOK, gin.H{
		"fileName": file.Filename,
		"size":     file.Size,
		"mimeType": file.Header,
	})
}
