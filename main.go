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

	// データベースに接続
	db.Connect()
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()
	db.DB.AutoMigrate(&models.Count{})

	e.GET("/", Index)
	e.POST("/plus", Plus)
	e.POST("/minus", Minus)
	e.POST("/multi", Multi)
	e.POST("/divide", Divide)
	e.POST("/clear", Clear)

	e.Logger.Fatal(e.Start(":1323"))
}

func Index(c echo.Context) error {
	// もしデータベースが空なのであれば、レコードを作成する
	if err := db.DB.First(&models.Count{}).Error; err != nil {
		db.DB.Create(&models.Count{Number: 0})
	}

	// 最初の要素を取得する
	Count := models.Count{}
	db.DB.Find(&Count)

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		// 取得した Count モデルの Number フィールドを取得する
		"number": Count.Number,
	})
}

// プラスボタンを押したら、数字を増やす
func Plus(c echo.Context) error {
	// 最初の要素を取得する
	Count := models.Count{}
	db.DB.Find(&Count)

	// モデルの Number フィールドを 1 増やす
	Count.Number++
	db.DB.Save(&Count)

	return c.Redirect(http.StatusSeeOther, "/")
}

func Minus(c echo.Context) error {
	// 最初の要素を取得する
	Count := models.Count{}
	db.DB.Find(&Count)

	// モデルの Number フィールドを 1 減らす
	Count.Number--
	db.DB.Save(&Count)

	return c.Redirect(http.StatusSeeOther, "/")
}

func Multi(c echo.Context) error {
	// 最初の要素を取得する
	Count := models.Count{}
	db.DB.Find(&Count)

	// モデルの Number フィールド ×2 する
	Count.Number *= 2
	db.DB.Save(&Count)

	return c.Redirect(http.StatusSeeOther, "/")
}

func Divide(c echo.Context) error {
	// 最初の要素を取得する
	Count := models.Count{}
	db.DB.Find(&Count)

	// モデルの Number フィールドを ÷2 する
	Count.Number /= 2
	db.DB.Save(&Count)

	return c.Redirect(http.StatusSeeOther, "/")
}

func Clear(c echo.Context) error {
	// 最初の要素を取得する
	Count := models.Count{}
	db.DB.Find(&Count)

	// モデルの Number フィールドを 0 にする
	Count.Number = 0
	db.DB.Save(&Count)

	return c.Redirect(http.StatusSeeOther, "/")
}