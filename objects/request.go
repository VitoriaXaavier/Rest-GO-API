// o define o objeto do nosso evento junto com outros objetos
package objects

import (
	"encoding/json"
	"net/http"
)

// limite
const MaxListLimit = 200

// recuperar um unico evento
type GetRequest struct {
	ID string `json:"id"`
}

// recuperar multiplos eventos
type ListRequest struct {
	Limt int `json: "limit"`
	After string `json: "after"`
	Nome string `json: "nome"`
}

// Criar um novo evento
type CreateRequest struct {
	Event *Event `json: "event"`
}

//Atualizar um evento
type UpDateDetailsRequest struct {
	ID string `json: "id"`
	Nome string `json: "nome"`
	Descricao string `json: "descricao"`
	Website string `json: "website"`
	Endereco string `json: "endereco"`
	Celular string `json: "celular"`
}

// Cancela um evento
type CancelRequest struct {
	ID string `json: "id"`
}

// Remarca um evento
type RemarcaRequest struct {
	ID string `json: "id"`
	NewSlot *TimeSlot `json: "newslot"`
}

// Deleta um evento
type DeletRequest struct {
	ID string `json: "id"`
}

// Resposta de qualquer solicitaçao de evento
type EventResponseWrapper struct {
	Event *Event `json: "event, omitempty"`
	Events []*Event `json: "events, omitempty"`
	Code int `json: "-"`
}

// Converte EventResponseWrapper  em json
func (e *EventResponseWrapper) JSON() []byte {
	if e == nil {
		return []byte( "{}")
	}
	res, _ := json.Marshal(e)
	return res
}

// Retorna o status code
func (e *EventResponseWrapper) StatusCode() int {
	if e == nil || e.Code == 0 {
		return http.StatusOK
	}
	return e.Code
}