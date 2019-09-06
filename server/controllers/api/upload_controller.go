package api

import (
	"io/ioutil"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/oss"
	"github.com/mlogclub/bbs-go/services"
)

const uploadMaxBytes int64 = 1024 * 1024 * 3 // 1M

type UploadController struct {
	Ctx iris.Context
}

func (this *UploadController) Post() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	file, header, err := this.Ctx.FormFile("image")
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

	url, err := oss.PutImage(fileBytes)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// vditor上传
func (this *UploadController) PostEditor() {
	errFiles := make([]string, 0)
	succMap := make(map[string]string)

	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		_, _ = this.Ctx.JSON(iris.Map{
			"msg":  "请先登录",
			"code": 1,
			"data": iris.Map{
				"errFiles": errFiles,
				"succMap":  succMap,
			},
		})
		return
	}

	maxSize := this.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := this.Ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		this.Ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = this.Ctx.WriteString(err.Error())
		return
	}

	form := this.Ctx.Request().MultipartForm
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
		url, err := oss.PutImage(fileBytes)
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}

		succMap[file.Filename] = url
	}

	_, _ = this.Ctx.JSON(iris.Map{
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
func (this *UploadController) PostFetch() {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		_, _ = this.Ctx.JSON(iris.Map{
			"msg":  "请先登录",
			"code": 1,
			"data": iris.Map{

			},
		})
		return
	}

	var data map[string]string
	data = make(map[string]string)

	err := this.Ctx.ReadJSON(&data)
	if err != nil {
		_, _ = this.Ctx.JSON(iris.Map{
			"msg":  err.Error(),
			"code": 0,
			"data": iris.Map{

			},
		})
		return
	}

	url := data["url"]
	output, err := oss.CopyImage(url)
	if err != nil {
		_, _ = this.Ctx.JSON(iris.Map{
			"msg":  err.Error(),
			"code": 0,
		})
	}
	_, _ = this.Ctx.JSON(iris.Map{
		"msg":  "",
		"code": 0,
		"data": iris.Map{
			"originalURL": url,
			"url":         output,
		},
	})
}
