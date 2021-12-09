package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"count/db"
	"count/models"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = renderer

	db.Connect()
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()
	db.DB.AutoMigrate(&models.Count{})

	e.GET("/", Index)
	e.POST("/plus", Plus)

	e.Logger.Fatal(e.Start(":1323"))
}

func Index(c echo.Context) error {
	if err := db.DB.First(&models.Count{}).Error; err != nil {
		db.DB.Create(&models.Count{Number: 0})
	}

	Count := models.Count{}
	db.DB.Find(&Count)

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"number": Count.Number,
	})
}

func Plus(c echo.Context) error {
	Count := models.Count{}
	db.DB.Find(&Count)

	Count.Number++
	db.DB.Save(&Count)

	return c.Redirect(http.StatusSeeOther, "/")
}