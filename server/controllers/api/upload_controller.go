package api

import (
	"bbs-go/package/uploader"
	"io/ioutil"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/services"
)

const uploadMaxBytes int64 = 1024 * 1024 * 3 // 1M

type UploadController struct {
	Ctx iris.Context
}

func (c *UploadController) Post() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	file, header, err := c.Ctx.FormFile("image")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	defer file.Close()

	if header.Size > uploadMaxBytes {
		return simple.JsonErrorMsg("图片不能超过3M")
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	logrus.Info("上传文件：", header.Filename, " size:", header.Size)

	url, err := uploader.PutImage(fileBytes)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// vditor上传
func (c *UploadController) PostEditor() {
	errFiles := make([]string, 0)
	succMap := make(map[string]string)

	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		_, _ = c.Ctx.JSON(iris.Map{
			"msg":  err.Message,
			"code": err.Code,
			"data": iris.Map{
				"errFiles": errFiles,
				"succMap":  succMap,
			},
		})
		return
	}

	maxSize := c.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := c.Ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = c.Ctx.WriteString(err.Error())
		return
	}

	form := c.Ctx.Request().MultipartForm
	files := form.File["file[]"]
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}
		url, err := uploader.PutImage(fileBytes)
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}

		succMap[file.Filename] = url
	}

	_, _ = c.Ctx.JSON(iris.Map{
		"msg":  "",
		"code": 0,
		"data": iris.Map{
			"errFiles": errFiles,
			"succMap":  succMap,
		},
	})
	return

}

// vditor 拷贝第三方图片
func (c *UploadController) PostFetch() {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		_, _ = c.Ctx.JSON(iris.Map{
			"msg":  err.Message,
			"code": err.Code,
			"data": iris.Map{},
		})
		return
	}

	var data map[string]string
	data = make(map[string]string)

	err := c.Ctx.ReadJSON(&data)
	if err != nil {
		_, _ = c.Ctx.JSON(iris.Map{
			"msg":  err.Error(),
			"code": 0,
			"data": iris.Map{},
		})
		return
	}

	url := data["url"]
	output, err := uploader.CopyImage(url)
	if err != nil {
		_, _ = c.Ctx.JSON(iris.Map{
			"msg":  err.Error(),
			"code": 0,
		})
	}
	_, _ = c.Ctx.JSON(iris.Map{
		"msg":  "",
		"code": 0,
		"data": iris.Map{
			"originalURL": url,
			"url":         output,
		},
	})
}
