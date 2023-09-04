// onde criamos o server e rotinas das rotas
package main

import (
	"log"
	"net/http"
	"github.com/VitoriaXaavier/Rest-GO-API/handlers"
	"github.com/VitoriaXaavier/Rest-GO-API/store"
	"github.com/gorilla/mux"
)

// Rodar o servidor
type Args struct {
	conn string
	port string
}

// Rodar o servidor com base nos argumentos de Args
func Run(args Args) error {
	router := mux.NewRouter(). 
		PathPrefix("/api/v1/").
		Subrouter()

	st := store.NewPostgresEventStore(args.conn)
	hnd := handlers.NewEventHandler(st)
	RegisterAllRouter(router,hnd)

//Inicio do servidor
	log.Println("Servidor funcionando na porta: ", args.port)
	return http.ListenAndServe(args.port,router)
}

// Registra todas as rotas
func RegisterAllRouter(router *mux.Router, hnd handlers.IEventHandler) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w,r)
		}) 
	})

	// Busca eventos
	router.HandleFunc("/event", hnd.Get).Methods(http.MethodGet)
	// Cria eventos
	router.HandleFunc("/event", hnd.Create).Methods(http.MethodPost)
	// Deleta eventos
	router.HandleFunc("/event", hnd.Delete).Methods(http.MethodDelete)

	// Cancela eventos
	router.HandleFunc("/event/cancel", hnd.Cancel).Methods(http.MethodPatch)
	// Atualiza eventos
	router.HandleFunc("event/details", hnd.UpDateDetails).Methods(http.MethodPut)
	// Remarca eventos
	router.HandleFunc("event/remarca", hnd.Reschedule).Methods(http.MethodPatch)
	// Lista eventos 
	router.HandleFunc("events", hnd.List).Methods(http.MethodGet)
}