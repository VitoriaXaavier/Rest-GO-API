// contém todas as rotas que o processo requer
package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/VitoriaXaavier/Rest-GO-API/errors"
	"github.com/VitoriaXaavier/Rest-GO-API/objects"
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
	store store.IEventStore
}

func NewEventHandler(store store.IEventStore) IEventHandler {
	return &handler{store: store}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	 if id == ""{
		WriterError(w, errors.ErrValidEventIDIsRequired)
		return
	 }
	evt, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id})
	if err != nil {
		WriterError(w, err)
		return
	}
	WriterResponse(w, &objects.EventResponseWrapper{Event: evt})
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	after := values.Get("after")
	nome := values.Get("nome")
	limit, err := IntFromString(w, values.Get("limit"))
	if err != nil {
		return
	}

	// lista eventos
	list, err := h.store.List(r.Context(), &objects.ListRequest{
		Limt: limit,
		After: after,
		Nome: nome,
	})
	if err != nil {
		WriterError(w, err)
		return
	}
	WriterResponse(w, &objects.EventResponseWrapper{Events: list})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriterError(w, errors.ErrUnprocessableEntity)
		return
	}
	evt := &objects.Event{}
	if Unmarshal(w, data, evt) != nil {
		return
	}
	if err != ChekSlot(evt.Slot) {
		WriterError(w, err)
		return
	}
	if err = h.store.Create(r.Context(), &objects.CreateRequest{Event: evt}); err != nil {
		WriterError(w, err)
		return
	}
	WriterResponse(w, &objects.EventResponseWrapper{Event: evt})

}

func (h *handler) UpDateDetails(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriterError(w, errors.ErrUnprocessableEntity)
		return
	}
	req := &objects.UpDateDetailsRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}

	// checa se evento existe
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriterError(w, err)
		return
	}

	if err = h.store.UpDateDetails(r.Context(), req); err != nil {
		WriterError(w, err)
		return
	}
	WriterResponse(w, &objects.EventResponseWrapper{})

}

func (h *handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriterError(w, errors.ErrValidEventIDIsRequired)
		return
	}

	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id}); err != nil {
		WriterError(w, err)
		return
	}

	if err := h.store.Cancel(r.Context(), &objects.CancelRequest{ID: id}); err != nil {
		WriterError(w, err)
		return
	}

	WriterResponse(w, &objects.EventResponseWrapper{})
}

func (h *handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriterError(w, errors.ErrUnprocessableEntity)
		return
	}

	req := &objects.RemarcaRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}

	if err := ChekSlot(req.NewSlot); err != nil {
		WriterError(w, err)
		return
	}

	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriterError(w, err)
		return
	}
	if err = h.store.Remarca(r.Context(), req); err != nil {
		WriterError(w, err)
		return
	}

	WriterResponse(w, &objects.EventResponseWrapper{})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriterError(w, errors.ErrValidEventIDIsRequired)
		return
	}
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id}); err != nil {
		WriterError(w, err)
		return
	}
	if err := h.store.Delete(r.Context(), &objects.DeletRequest{ID: id}); err != nil {
		WriterError(w, err)
		return
	}
	WriterResponse(w,&objects.EventResponseWrapper{})
}	