package simple

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/mlogclub/simple/strcase"
)

type GenerateStruct struct {
	Name   string
	Fields []GenerateField
}

type GenerateField struct {
	CamelName   string
	NativeField reflect.StructField
}

type InputData struct {
	PkgName   string
	Name      string // FuckShit
	CamelName string // fuckShit
	KebabName string // FuckShit -> fuck-shit
	Fields    []GenerateField
}

func Generate(baseDir, pkgName string, models ...GenerateStruct) {
	for _, model := range models {
		if err := generateRepository(baseDir, pkgName, model); err != nil {
			logrus.Error(err)
		}
		if err := generateService(baseDir, pkgName, model); err != nil {
			logrus.Error(err)
		}
		if err := generateController(baseDir, pkgName, model); err != nil {
			logrus.Error(err)
		}
		if err := generateWeb(baseDir, pkgName, model); err != nil {
			logrus.Error(err)
		}
	}
}

func GetGenerateStruct(s interface{}) GenerateStruct {
	structName := StructName(s)
	structFields := StructFields(s)

	var fields []GenerateField
	for _, f := range structFields {
		if f.Anonymous {
			continue
		}
		fields = append(fields, GenerateField{
			CamelName:   strcase.ToLowerCamel(f.Name),
			NativeField: f,
		})
	}

	return GenerateStruct{
		Name:   structName,
		Fields: fields,
	}
}

func generateRepository(baseDir, pkgName string, s GenerateStruct) error {
	var b bytes.Buffer
	err := repositoryTmpl.Execute(&b, &InputData{
		PkgName:   pkgName,
		Name:      s.Name,
		CamelName: strcase.ToLowerCamel(s.Name),
		KebabName: strcase.ToKebab(s.Name),
		Fields:    s.Fields,
	})
	if err != nil {
		return err
	}
	c := b.String()

	p, err := getFilePath(baseDir, "/repositories/"+strcase.ToSnake(s.Name+"_repository.go"))
	if err != nil {
		return err
	}
	return writeFile(p, c)
}

func generateService(baseDir, pkgName string, s GenerateStruct) error {
	var b bytes.Buffer
	err := serviceTmpl.Execute(&b, &InputData{
		PkgName:   pkgName,
		Name:      s.Name,
		CamelName: strcase.ToLowerCamel(s.Name),
		KebabName: strcase.ToKebab(s.Name),
		Fields:    s.Fields,
	})
	if err != nil {
		return err
	}
	c := b.String()

	p, err := getFilePath(baseDir, "/services/"+strcase.ToSnake(s.Name+"_service.go"))
	if err != nil {
		return err
	}
	return writeFile(p, c)
}

func generateController(baseDir, pkgName string, s GenerateStruct) error {
	var b bytes.Buffer
	err := controllerTmpl.Execute(&b, &InputData{
		PkgName:   pkgName,
		Name:      s.Name,
		CamelName: strcase.ToLowerCamel(s.Name),
		KebabName: strcase.ToKebab(s.Name),
		Fields:    s.Fields,
	})
	if err != nil {
		return err
	}
	c := b.String()

	p, err := getFilePath(baseDir, "/controllers/admin/"+strcase.ToSnake(s.Name+"_controller.go"))
	if err != nil {
		return err
	}
	return writeFile(p, c)
}

func generateWeb(baseDir, pkgName string, s GenerateStruct) error {
	var b bytes.Buffer
	err := viewIndexTmpl.Execute(&b, &InputData{
		PkgName:   pkgName,
		Name:      s.Name,
		KebabName: strcase.ToKebab(s.Name),
		Fields:    s.Fields,
	})
	if err != nil {
		return err
	}
	c := b.String()

	sub := path.Join("/web/admin/src/views/", strcase.ToKebab(s.Name), "Index.vue")

	p, err := getFilePath(baseDir, sub)
	if err != nil {
		return err
	}
	return writeFile(p, c)
}

func getFilePath(baseDir, sub string) (filepath string, err error) {
	filepath = path.Join(baseDir, sub)
	base := path.Dir(filepath)
	err = os.MkdirAll(base, os.ModePerm)
	return
}

func writeFile(filepath string, content string) error {
	exists, err := PathExists(filepath)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("文件已经存在...", filepath)
		filepath = filepath + ".temp"
	}
	return WriteString(filepath, content, true)
}
