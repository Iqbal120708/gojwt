package main

import (
    "gojwt/pkg/database"
    "gojwt/internal/repository"
    "gojwt/pkg/config"
    "gojwt/internal/middleware"
    "gojwt/internal/usecase"
    "gojwt/internal/handler"
    "github.com/gorilla/mux"
    "net/http"
)

func main() {
    config.Load()
    cfg := config.Get()
    
    db, err := database.NewMySQL()
    if err != nil {
        panic(err)
    }
    
    ur := repository.NewUserRepo(db)
    uc := usecase.NewUserUseCase(ur)
    uh := handler.NewUserHandler(uc)
    
    r := mux.NewRouter()
	r.HandleFunc("/signup", uh.Create).Methods("POST")
	r.HandleFunc("/login", uh.Login).Methods("POST")
	
	api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)
    
    api.HandleFunc("/profile", uh.Profile).Methods("GET")
    
	http.ListenAndServe(":"+cfg.AppPort, r)
}