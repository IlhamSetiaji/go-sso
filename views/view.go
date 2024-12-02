package views

import (
	"app/go-sso/utils"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

var LayoutDir = "views/layouts"

func NewView(layout string, files ...string) *View {
	files = append(layoutFiles(), files...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) Render(c *gin.Context, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		dataMap = make(map[string]interface{})
	}

	if _, exists := dataMap["Title"]; !exists {
		dataMap["Title"] = "Default Title"
	}
	if _, exists := dataMap["AssetBase"]; !exists {
		dataMap["AssetBase"] = "/assets"
	}
	if _, exists := dataMap["CsrfToken"]; !exists {
		dataMap["CsrfToken"] = csrf.GetToken(c)
	}

	session := sessions.Default(c)
	dataMap["Errors"] = session.Flashes("errors")
	dataMap["Status"] = session.Get("status")
	dataMap["Success"] = session.Get("success")
	dataMap["Error"] = session.Get("error")
	dataMap["Warning"] = session.Get("warning")
	dataMap["Profile"] = session.Get("profile")
	dataMap["CurrentPath"] = c.Request.URL.Path
	sessionKeys := []string{"success", "error", "status", "errors", "warning"}
	for _, key := range sessionKeys {
		if dataMap[key] != nil {
			session.Delete(key)
		}
	}
	session.Save()

	templateHelper := utils.NewTemplateHelper(c)

	dataMap["HasPermission"] = templateHelper.HasPermission

	err := v.Template.ExecuteTemplate(c.Writer, v.Layout, dataMap)
	if err != nil {
		c.String(500, err.Error())
	}
}

func layoutFiles() []string {
	files, err := ioutil.ReadDir(LayoutDir)
	if err != nil {
		panic(err)
	}

	var layoutFiles []string
	for _, file := range files {
		if !file.IsDir() {
			layoutFiles = append(layoutFiles, filepath.Join(LayoutDir, file.Name()))
		}
	}
	fmt.Printf("files: %v\n", layoutFiles)
	return layoutFiles
}
