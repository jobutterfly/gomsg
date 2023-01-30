package main


import (
	"log"
	"net/http"
	"os"

	"github.com/enzdor/gomsg/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	errEnv := godotenv.Load()
	if errEnv != nil {
	    log.Fatal(errEnv)
	}
	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")

	db := controllers.NewDB(user, pass, name)
	h := controllers.NewHandler(db)

	http.HandleFunc("/", h.ServeIndex)
	http.HandleFunc("/board/", h.ServeBoard)
	http.HandleFunc("/thread/", h.ServeThread)
	http.HandleFunc("/post/", h.ServePost)
	http.HandleFunc("/reply/", h.ServeReply)
	http.HandleFunc("/kill/", h.ServeKill)
	http.HandleFunc("/error/", h.ServeError)
	
	log.Print("Listening on port :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
	    log.Fatal(err)
	}
}









