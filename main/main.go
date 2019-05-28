package main

import (
	
	"log"
	"net/http"
	
)
func main() {
	// inisialisasi file server
	fs := http.FileServer(http.Dir("../webclient"))
	http.Handle("/", fs)

	// konfigurasi websocket
	http.HandleFunc("/ws", handleConnections)

	// listen pesan
	go handlePesans()

	// Menjalankan server di 8080
	log.Println("Aplikasi berjalan di port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


// Konfigurasi upgrader untuk Websocket Connection, upgrade method dari request untuk get connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hubungkan Client
var clients = make(map[*websocket.Conn]bool)
 // BroadCast Pesan dari client agar client lain dapat membaca 
var broadcast = make(chan Pesan)          


// Membangun objek pesan trdiri dari Username dan pesan
type Pesan struct {
	Username string `json:"username"`
	Pesan  string `json:"message"`
}



func handleConnections(w http.ResponseWriter, r *http.Request) {
	// implementasi upgrader
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Pesan
		// Menerima pesan baru dari client dan ubah menjadi Objek pesan
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Setelah menerima, maka pesan akan di broadcast
		broadcast <- msg
	}
}

func handlePesans() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
