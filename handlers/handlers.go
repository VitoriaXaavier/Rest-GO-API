// contém todas as rotas que o processo requer
package handlers

import (
	 "net/http"
	 "github.com/VitoriaXaavier/Rest-GO-API/store"
)
type IEventHandler interface {
	Get (w http.ResponseWriter, r *http.Request)
	List (w http.ResponseWriter, r *http.Request)
	Create (w http.ResponseWriter, r *http.Request)
	UpDateDetails (w http.ResponseWriter, r *http.Request)
	Cancel (w http.ResponseWriter, r *http.Request)
	Reschedule (w http.ResponseWriter, r *http.Request)
	Delete (w http.ResponseWriter, r *http.Request)
}

//usado para implementar a interface
type handler struct {

}

func NewEventHandler(store store.IEventStore) IEventHandler {
	return &handler{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}

func (h *handler) UpDateDetails(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}

func (h *handler) Cancel(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}

func (h *handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("Implemente-me")
}	