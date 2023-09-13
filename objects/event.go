// o define o objeto do nosso evento junto com outros objetos
package objects

import "time"

// Define o status do evento
type EventStatus string

const (
	Original  EventStatus = "original"
	Cancelado EventStatus = "cancelado"
	Remarcado EventStatus = "remarcado"
)

type TimeSlot struct {
	StarTime time.Time `json:"star_time,omitempty"`
	EndTime  time.Time `json:"end_time,omitempty"`
}

// Eventos de objetos da api
type Event struct {
	// identificador
	ID string `gorm:"primary_key" json:"id,omitempty"`

	Nome      string `json:"nome,omitempty"`
	Descricao string `json:"descricao,omitempty"`
	Website   string `json:"website,omitempty"`
	Endereco  string `json:"endereco,omitempty"`
	Celular   string `json:"celular,omitempty"`

	Slot   *TimeSlot   `gorm:"embedded" json:"slot,omitempty"`
	Status EventStatus `json:"status,omitempty"`

	Criado     time.Time `json:"criado,omitempty"`
	Atualizado time.Time `json:"atualizado,omitempty"`
	Cancelado  time.Time `json:"cancelado,omitempty"`
	Remarcado  time.Time `json:"remarcado,omitempty"`
}
