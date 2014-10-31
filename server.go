package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "log"
    "net/http"
    "time"
)

func main() {
  	api := Api{}
    api.InitDB()
    api.InitSchema()

	handler := rest.ResourceHandler{
        EnableRelaxedContentType: true,
    }
    err := handler.SetRoutes(
        &rest.Route{"GET", "/todos", api.GetAllTodos},
        &rest.Route{"POST", "/todos", api.PostTodo},
        &rest.Route{"GET", "/todos/:id", api.GetTodo},
        &rest.Route{"PUT", "/todos/:id", api.PutTodo},
        &rest.Route{"DELETE", "/todos/:id", api.DeleteTodo},
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Fatal(http.ListenAndServe(":8080", &handler))
}

type Todo struct {
    Id        int64     `json:"id"`
    Text   	  string    `sql:"size:1024" json:"text"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt time.Time `json:"-"`
}

type Api struct {
    DB gorm.DB
}

func (api *Api) InitDB() {
    var err error
    api.DB, err = gorm.Open("mysql", "root:xxx@tcp(podkowik.net:3306)/todo?charset=utf8&parseTime=True")
    if err != nil {
        log.Fatalf("Got error when connect database, the error is '%v'", err)
    }
    api.DB.LogMode(true)
}

func (api *Api) InitSchema() {
    api.DB.AutoMigrate(Todo{})
}

func (api *Api) GetAllTodos(w rest.ResponseWriter, r *rest.Request) {
    todos := []Todo{}
    api.DB.Find(&todos)
    w.WriteJson(&todos)
}

func (api *Api) GetTodo(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    todo := Todo{}
    if api.DB.First(&todo, id).Error != nil {
        rest.NotFound(w, r)
        return
    }
    w.WriteJson(&todo)
}

func (api *Api) PostTodo(w rest.ResponseWriter, r *rest.Request) {
    todo := Todo{}
    if err := r.DecodeJsonPayload(&todo); err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := api.DB.Save(&todo).Error; err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteJson(&todo)
}

func (api *Api) PutTodo(w rest.ResponseWriter, r *rest.Request) {

    id := r.PathParam("id")
    todo := Todo{}
    if api.DB.First(&todo, id).Error != nil {
        rest.NotFound(w, r)
        return
    }

    updated := Todo{}
    if err := r.DecodeJsonPayload(&updated); err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    todo.Text = updated.Text

    if err := api.DB.Save(&todo).Error; err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteJson(&todo)
}

func (api *Api) DeleteTodo(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    todo := Todo{}
    if api.DB.First(&todo, id).Error != nil {
        rest.NotFound(w, r)
        return
    }
    if err := api.DB.Delete(&todo).Error; err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}
