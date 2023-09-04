// define toda a base de dados. implementação de armazenamento para eventos 
package store

import (
	"context"
	"log"
	"os"
	"github.com/VitoriaXaavier/Rest-GO-API/errors"
	"github.com/VitoriaXaavier/Rest-GO-API/objects"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Representa o armazenamento de eventos baseado no PostgreSQL
type pg struct {
	db *gorm.DB
}

// Retorna a postgre implementaçao do evento store
func NewPostgresEventStore(conn string) IEventStore {
	// criando a conexao com o banco
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				},

			),
		},
)

	if err != nil {
		panic("Habilitar conexão com database: " + err.Error())
	}

	if err := db.AutoMigrate(&objects.Event{}); err != nil {
		panic(" Habilita a migração para o database: " + err.Error())
	}
	return &pg{db: db}
}

// Implementação dos metodos da interface IEventStore
func (p *pg) Get(ctx context.Context, in *objects.GetRequest) (*objects.Event, error) {
	evt := &objects.Event{}

	// pega o evento com o id igual ao uid do database
	err := p.db.WithContext(ctx).Take(evt, "id = ?", in.ID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.ErrEventNotFound
	}
	return evt, err
}

func (p *pg) List(ctx context.Context, in *objects.ListRequest) ([]*objects.Event, error) {
	if in.Limt == 0 || in.Limt > objects.MaxListLimit {
		in.Limt = objects.MaxListLimit
	}
	query := p.db.WithContext(ctx).Limit(in.Limt)
	if in.After != "" {
		query = query.Where("id > ?", in.After)
	}
	if in.Nome != "" {
		query = query.Where("nome ilike ?" , "%" + in.Nome + "%")
	}
	list := make([]*objects.Event, 0, in.Limt)
	err := query.Order("id").Find(&list).Error

	return list, err
}

func (p *pg) Create(ctx context.Context, in *objects.CreateRequest) error {
	if in.Event == nil {
		return errors.ErrObjectIsRequired
	}
	in.Event.ID = GenerateUniqueId()
	in.Event.Status = objects.Original
	in.Event.Criado = p.db.NowFunc()

	return p.db.WithContext(ctx).
		Create(in.Event).
		Error
}

func (p *pg) UpDateDetails(ctx context.Context, in *objects.UpDateDetailsRequest) error {
	evt := &objects.Event{
		ID: in.ID,
		Nome: in.Nome, 
		Descricao: in.Descricao,
		Website: in.Website,
		Endereco: in.Endereco,
		Celular: in.Celular,
		Atualizado: p.db.NowFunc(),
	}
	return p.db.WithContext(ctx).Model(evt).
	Select("nome","descricao","website","endereco","celular","atualizacao").
	Updates(evt).
	Error
}

func (p *pg) Cancel(ctx context.Context, in *objects.CancelRequest) error {
	evt := &objects.Event{
		ID: in.ID,
		Status: objects.Cancelado,
		Cancelado: p.db.NowFunc(),
	}
	return p.db.WithContext(ctx).Model(evt).
		Select("status", "cancelado").
		Updates(evt).
		Error
}

func (p *pg) Remarca(ctx context.Context, in *objects.RemarcaRequest) error {
	evt := &objects.Event{
		ID: in.ID,
		Slot: in.NewSlot,
		Status: objects.Remarcado,
		Remarcado: p.db.NowFunc(),
	}

	return p.db.WithContext(ctx).Model(evt).
		Select("status", "start_time", "end_time", "remarcado").
		Updates(evt).
		Error
}

func (p *pg) Delete(ctx context.Context, in *objects.DeletRequest) error {
	evt := &objects.Event{
		ID: in.ID,
	}

	return p.db.WithContext(ctx).Model(evt).
		Delete(evt).
		Error
}