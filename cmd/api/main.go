package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net/http"

	//_ "github.com/go-sql-driver/mysql"
	"api/pkg/handlers"
	"api/pkg/products"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	var name, password, dbname string
	flag.StringVar(&name, "user", "", "The name of user")
	flag.StringVar(&password, "p", "", "password")
	flag.StringVar(&dbname, "db", "", "The name of database")
	flag.Parse()
	connStr := "user=" + name + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("no open bd:", err)
		return
	}

	//db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	templates := template.Must(template.ParseGlob("./templates/*"))

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	itemsRepo := products.NewRepository(db)

	handlers := &handlers.ItemsHandler{
		Tmpl:      templates,
		Logger:    logger,
		ItemsRepo: itemsRepo,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Index).Methods("GET")
	r.HandleFunc("/index", handlers.Index).Methods("GET")
	r.HandleFunc("/table", handlers.ListAll).Methods("GET")
	r.HandleFunc("/statistic", handlers.Statistic).Methods("POST")
	r.HandleFunc("/show", handlers.Show).Methods("GET")
	r.HandleFunc("/product", handlers.ListProduct).Methods("POST")
	r.HandleFunc("/upload", handlers.Upload).Methods("GET")

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, r)
}
