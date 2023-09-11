package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"github.com/VitoriaXaavier/Rest-GO-API/handlers"
	"github.com/VitoriaXaavier/Rest-GO-API/objects"
	"github.com/VitoriaXaavier/Rest-GO-API/store"
	"github.com/VitoriaXaavier/Rest-GO-API/errors"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router mux.Router
	flushAll func (t *testing.T)
	createOn func (t *testing.T, nome string) *objects.Event
	getOn func (t *testing.T, id string, wantErr bool) *objects.Event
)


func TestMain(t *testing.M) {
	log.Println("Registrando")

	conn := "postgres://user:password@localhost:5432/db?sslmode=disable"
	if c := os.Getenv("DB_CONN"); c != "" {
		conn = c
	}

	router = *mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	st := store.NewPostgresEventStore(conn)
	hnd := handlers.NewEventHandler(st)
	RegisterAllRouter(&router,hnd)
	
	flushAll = func(t *testing.T) {
		db, err := gorm.Open(postgres.Open(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		db.Delete(&objects.Event{}, "1=1")
	}

	createOn = func(t *testing.T, nome string) *objects.Event {
		evt := &objects.Event{
			Nome: nome,
			Descricao: "Descrição de " + nome,
			Website: "https://" + nome + ".com",
			Slot: &objects.TimeSlot{
				StarTime: time.Now().UTC(),
				EndTime: time.Now().UTC().Add(time.Hour),
			},

		}
		err := st.Create(context.TODO(), &objects.CreateRequest{Event: evt})
		if err != nil {
			t.Fatal(err)
		}
		return evt
	}

	getOn = func(t *testing.T, id string, wantErr bool) *objects.Event {
		evt, err := st.Get(context.TODO(), &objects.GetRequest{ID: id})
		if err != nil && wantErr{
			t.Fatal(err)
		}
		return evt
	}
	log.Panicln("começando")
	os.Exit(t.Run())
}

func Do (req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestUnknowEndpoints(t *testing.T) {
	tests := [] struct {
		nome string
		setup func (t *testing.T) *http.Request 
		
	}{ 
		{
		nome: "root",
		setup: 	func (t *testing.T)*http.Request  {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if  err != nil {
				t.Fatal(err)
			}
			return req
		},
	},

	{
		nome: "api-root",
		setup: func (t *testing.T) *http.Request  {
			req, err := http.NewRequest(http.MethodGet, "api/v1/", nil) 
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
	},

}

for _, tt := range tests {
	t.Run(tt.nome, func(t *testing.T) {
		w := Do(tt.setup(t))
		_ = assert.Equal(t, http.StatusNotFound, w.Code) && assert.Equal(t, "404 pagina não encontrada\n", string(w.Body.Bytes()))
	})
}
}

func TestGetEndpoints(t *testing.T) {
	flushAll(t)
	tests := [] struct {
		nome string
		code int
		setup func (t *testing.T) *http.Request

	}{

	{
		nome: "OK",
		setup: func(t *testing.T) *http.Request {
			evt := createOn(t, "OK")
			req, err := http.NewRequest(http.MethodGet, "/api/v1/event?id=" + evt.ID, nil) 
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code: http.StatusOK,
	},
	{
		nome: "NotFound",
		setup: func(t *testing.T) *http.Request {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/event?id=32" , nil)
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code: http.StatusNotFound,
	},
	{
		nome: "WithoutParameter",
		setup: func(t *testing.T) *http.Request {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/event", nil)
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code: http.StatusBadRequest,

	},
}
for _ ,tt := range tests {
	t.Run(tt.nome, func(t *testing.T) {
		w := Do(tt.setup(t))
		assert.Equal(t, tt.code, w.Code)
		got := &objects.EventResponseWrapper{}
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
	})
}
}


func TestListEndpoint(t *testing.T) {
	flushAll(t)
	tests := [] struct {
		nome    string
		code    int
		setup   func(t *testing.T) *http.Request
		listLen int
	}{
		{
			nome: "Zero",
			setup: func(t *testing.T) *http.Request {
				flushAll(t)
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events", nil )
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code: http.StatusOK,
			listLen: 0,
		},
	{
		nome: "All",
		setup: func(t *testing.T) *http.Request {
			_ = createOn(t, "One")
			_ = createOn(t, "two")
			req, err := http.NewRequest(http.MethodGet, "/api/v1/events", nil)
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code:  http.StatusOK,
		listLen: 2,
	},
	{
		nome: "Limited",
		setup: func(t *testing.T) *http.Request {
			_ = createOn(t, "three")
			req, err := http.NewRequest(http.MethodGet, "/api/v1/events?Limit=2", nil)
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code: http.StatusOK,
		listLen: 2,
	},
	{
		nome: "After",
		setup: func(t *testing.T) *http.Request {
			evt := createOn(t, "Four")
			_ = createOn(t, "Five")
			req, err := http.NewRequest(http.MethodGet, "/api/v1/events?after=" +evt.ID, nil) 
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code:  http.StatusOK,
		listLen: 1,
	},
	{
		nome: "Nome",
		setup:  func(t *testing.T) *http.Request {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/events?nome=e", nil)
			if err != nil {
				t.Fatal(err)
			}
			return req
		},
		code: http.StatusOK,
		listLen: 3,
	},
	}
	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			w := Do(tt.setup(t))
			got := &objects.EventResponseWrapper{}
			assert.Equal(t, tt.code,w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Equal(t, len(got.Events), tt.listLen)
		})
	} 

}



func TestCreateEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		message string
		code    int
		evt     *objects.Event
	}{
		{
			name:    "Ok",
			message: "",
			code:    http.StatusOK,
			evt: &objects.Event{
				Nome:        "Help Ok",
				Descricao: "...",
				Website:     "https://help.ok",
				Endereco:     "Help Okland Yard",
				Slot: &objects.TimeSlot{
					StarTime: time.Now().UTC(),
					EndTime:   time.Now().UTC().Add(time.Hour),
				},
			},
		},
		{
			name:    "No Timing",
			message: errors.ErrEventTimingIsRequired.Message,
			code:    errors.ErrEventTimingIsRequired.Code,
			evt: &objects.Event{
				Nome:        "Help Ok",
				Descricao: "...",
				Website:     "https://help.ok",
				Endereco:     "Help Okland Yard",
			},
		},
		{
			name:    "No input",
			message: errors.ErrObjectIsRequired.Message,
			code:    errors.ErrObjectIsRequired.Code,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.evt)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, "/api/v1/event", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			w := Do(req)
			got, gotErr := &objects.EventResponseWrapper{}, &errors.Error{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), gotErr))
			assert.Equal(t, tt.message, gotErr.Message)
			if tt.code == http.StatusOK {
				ok := assert.NotNil(t, got.Event) &&
					assert.NotEmpty(t, got.Event.ID) &&
					assert.NotEmpty(t, got.Event.Criado)
				if ok {
					tt.evt.ID = got.Event.ID
					tt.evt.Criado = got.Event.Criado
					tt.evt.Status = objects.Original
					assert.Equal(t, tt.evt, got.Event)
				}
			}
		})
	}
}

func TestUpdateDetailsEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, evt *objects.Event) (*http.Request, *objects.Event) {
		var (
			b   []byte
			err error
		)
		if evt != nil {
			b, err = json.Marshal(&objects.UpDateDetailsRequest{
				ID:          evt.ID,
				Nome:        evt.Nome,
				Descricao:	 evt.Descricao,
				Website:     evt.Website,
				Endereco:    evt.Endereco,
				Celular: 	 evt.Celular,
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		req, err := http.NewRequest(http.MethodPut, "/api/v1/event/details", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		return req, evt
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, *objects.Event)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, *objects.Event) {
				evt := createOn(t, "Ok")
				evt.Endereco = "Fake"
				evt.Celular = "+23765"
				return reqFn(t, evt)
			},
			code: http.StatusOK,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, *objects.Event) {
				evt := createOn(t, "Ok")
				evt.ID = ""
				return reqFn(t, evt)
			},
			message: errors.ErrEventNotFound.Message,
			code:    http.StatusNotFound,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, *objects.Event) {
				return reqFn(t, nil)
			},
			message: errors.ErrObjectIsRequired.Message,
			code:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, exp := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if exp != nil {
				evt := getOn(t, exp.ID, true)
				assert.Equal(t, exp.Nome, evt.Nome)
				assert.Equal(t, exp.Descricao, evt.Descricao)
				assert.Equal(t, exp.Website, evt.Website)
				assert.Equal(t, exp.Endereco, evt.Endereco)
				assert.Equal(t, exp.Celular, evt.Celular)
			}
		})
	}
}

func TestCancelEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, in *objects.CancelRequest) (*http.Request, string) {
		id := ""
		if in != nil {
			id = in.ID
		}
		req, err := http.NewRequest(http.MethodPatch, "/api/v1/event/cancel?id="+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req, id
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, string)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, string) {
				evt := createOn(t, "Ok")
				return reqFn(t, &objects.CancelRequest{ID: evt.ID})
			},
			code: http.StatusOK,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, nil)
			},
			message: errors.ErrValidEventIDIsRequired.Message,
			code:    http.StatusBadRequest,
		},
		{
			name: "No Found",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, &objects.CancelRequest{ID: "fake"})
			},
			message: errors.ErrEventNotFound.Message,
			code:    http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, id := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if id != "" {
				evt := getOn(t, id, true)
				assert.Equal(t, objects.Cancelado, evt.Status)
			}
		})
	}
}

func TestRescheduleEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, in *objects.RemarcaRequest) (*http.Request, *objects.RemarcaRequest) {
		var (
			b   []byte
			err error
		)
		if in != nil {
			b, err = json.Marshal(in)
			if err != nil {
				t.Fatal(err)
			}
		}
		req, err := http.NewRequest(http.MethodPatch, "/api/v1/event/reschedule", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		return req, in
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, *objects.RemarcaRequest)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, *objects.RemarcaRequest) {
				evt := createOn(t, "Ok")
				return reqFn(t, &objects.RemarcaRequest{
					ID:      evt.ID,
					NewSlot: evt.Slot,
				})
			},
			code: http.StatusOK,
		},
		{
			name: "No Slot",
			setup: func(t *testing.T) (*http.Request, *objects.RemarcaRequest) {
				return reqFn(t, &objects.RemarcaRequest{})
			},
	
			message: errors.ErrEventTimingIsRequired.Message,
			code:    http.StatusBadRequest,
		},
		{
			name: "Invalid Time",
			setup: func(t *testing.T) (*http.Request, *objects.RemarcaRequest) {
				evt := createOn(t, "Ok")
				return reqFn(t, &objects.RemarcaRequest{
					ID:      evt.ID,
					NewSlot: &objects.TimeSlot{},
				})
			},
			message: errors.ErrInvalidTimeFormat.Message,
			code:    http.StatusBadRequest,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, *objects.RemarcaRequest) {
				evt := createOn(t, "Ok")
				return reqFn(t, &objects.RemarcaRequest{NewSlot: evt.Slot})
			},
			message: errors.ErrEventNotFound.Message,
			code:    http.StatusNotFound,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, *objects.RemarcaRequest) {
				return reqFn(t, nil)
			},
			message: errors.ErrObjectIsRequired.Message,
			code:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, exp := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if exp != nil {
				evt := getOn(t, exp.ID, true)
				assert.Equal(t, exp.NewSlot.StarTime.Format(time.RFC3339), evt.Slot.StarTime.UTC().Format(time.RFC3339))
				assert.Equal(t, exp.NewSlot.EndTime.Format(time.RFC3339), evt.Slot.EndTime.UTC().Format(time.RFC3339))
			}
		})
	}
}

func TestDeleteEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, in *objects.DeletRequest) (*http.Request, string) {
		id := ""
		if in != nil {
			id = in.ID
		}
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/event?id="+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req, id
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, string)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, string) {
				evt := createOn(t, "Ok")
				return reqFn(t, &objects.DeletRequest{ID: evt.ID})
			},
			code: http.StatusOK,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, nil)
			},
			message: errors.ErrValidEventIDIsRequired.Message,
			code:    http.StatusBadRequest,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, &objects.DeletRequest{ID: "fake"})
			},
			message: errors.ErrEventNotFound.Message,
			code:    http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, id := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if id != "" {
				assert.Nil(t, getOn(t, id, false))
			}
		})
	}
}