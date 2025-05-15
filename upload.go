package main

import (
	"context"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
)

func UploadImageToCloudinary(c context.Context, file multipart.File, filePath string) (string, error) {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinarySecret)
	if err != nil {
		return "", err
	}
	uploadParams := uploader.UploadParams{
		PublicID: filePath,
	}
	result, err := cld.Upload.Upload(c, file, uploadParams)
	if err != nil {
		return "", err
	}
	imageUrl := result.SecureURL
	log.Println("imageUrl", imageUrl)
	return imageUrl, nil
}

func (r *Router) handleUploadImage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	form, err := ctx.MultipartForm()
	log.Println("🚀 ~ func ~ form:", form)
	if err != nil {
		log.Println("err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	images := []string{}
	log.Println("running upload image...")
	t1 := time.Now()
	files := form.File["images"]
	for _, file := range files {
		image, err := file.Open()
		if err != nil {
			log.Println("file open err:", err)
			continue
		}
		imageUrl, err := UploadImageToCloudinary(c, image, strings.Split(file.Filename, ".")[0])
		if err != nil {
			log.Println("upload img err:", err)
			return
		}
		images = append(images, imageUrl)
	}
	log.Println("upload done:", time.Since(t1))
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: images})
}
