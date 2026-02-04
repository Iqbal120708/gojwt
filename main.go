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
    "fmt"
)

func main() {
    config.Load()
    cfg := config.Get()
    
    db, err := database.NewMySQL()
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to database: %v", err))
    }
    
    ur := repository.NewUserRepo(db)
    uc := usecase.NewUserUseCase(ur)
    uh := handler.NewUserHandler(uc)
    
    r := mux.NewRouter()
	r.HandleFunc("/signup", uh.Create).Methods("POST")
	r.HandleFunc("/login", uh.Login).Methods("POST")
	r.HandleFunc("/refresh", uh.RefreshToken).Methods("POST")
	
	api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)
    
    api.HandleFunc("/profile", uh.Profile).Methods("GET")
    
	http.ListenAndServe(":"+cfg.AppPort, r)
}