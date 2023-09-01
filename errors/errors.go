//Contém todos os erros encontrados enquanto  processa a api 
package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	// Erro interno http 500
	ErrInternal = &Error {
		Code: http.StatusInternalServerError,
		Message: "Algo está errado",
	}

	// Erro não processável http 422
	ErrUnprocessableEntity = &Error {
		Code: http.StatusUnprocessableEntity,
		Message: "Entidade não processável",
	}

	// Requisição ruim http 400
	ErrBadRequest = &Error {
		Code: http.StatusBadRequest,
		Message: "Argumento invalido",
	}

	// Evento não encontrado http 404
	ErrEventNotFound = &Error {
		Code: http.StatusNotFound,
		Message: "Evento não encontrado",
	}

	// Objeto de solicitação deve ser provido http 400
	ErrObjectIsRequired = &Error {
		Code: http.StatusBadRequest,
		Message: "O objeto de solicitação deve ser fornecido",
	}

	// Necessário um ID válido http 400
	ErrValidEventIDIsRequired = &Error {
		Code: http.StatusBadRequest,
		Message: "É necessário um ID de evento válido",
	}

	// Deve fornecer a hora do inicio e fim http 400
	ErrEventTimingIsRequired = &Error { 
		Code: http.StatusBadRequest,
		Message: "A hora de início e de término do evento deve ser fornecida",
	}

	// Limite deve ser um valor inteiro http 400
	ErrInvalidLimit = &Error {
		Code: http.StatusBadRequest,
		Message: "O limite deve ser um valor inteiro",
	}

	//Formato errado para o tempo http 400
	ErrInvalidTimeFormat = &Error {
		Code: http.StatusBadRequest,
		Message: "O tempo deve ser passado no formato RFC3339: " + time.RFC3339 , 
	}
	
)

type Error struct {
	Code    int
	Message string
}

// Implementa a struct de erro
func( err *Error) Error() string {
	return err.String()
}

func(err *Error) String() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("Error: code=%s menssagem=%s ", http.StatusText(err.Code), err.Message)
}

// Converte o erro em JSON
func(err *Error) JSON() []byte {
	if err == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(err)
	return res
}

// Implementa o status code
func(err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}
	return err.Code
}