package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var localRoom *Room

func main() {

	localRoom = NewRoom(50)
	go localRoom.start()

	mux := http.NewServeMux()
	mux.Handle("/", EnsureMethod(http.MethodGet, room))
	mux.Handle("/join", EnsureMethod(http.MethodPost, join))
	mux.HandleFunc("/ws", ws)

	log.Println("Starting server on port 5555")
	log.Fatal(http.ListenAndServe(":5555", mux))

}

var roomTmp = template.Must(template.ParseFiles(filepath.Join("templates", "room.html")))

func room(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	renderTemplate(roomTmp, w, nil)
}

var chatTmp = template.Must(template.New("chat").ParseFiles(filepath.Join("templates", "chat.html")))

func join(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	//Needs more error checking
	username := r.Form.Get("username")

	data := PageData{
		Username: username,
		ChatData: localRoom.History.RetrieveMessages(),
	}

	renderTemplate(chatTmp, w, data)

}

func ws(w http.ResponseWriter, r *http.Request) {
	NewClient(localRoom, w, r)
}

func renderTemplate(tmp *template.Template, w http.ResponseWriter, data any) {
	if err := tmp.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println(err)
	}
}
