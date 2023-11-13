package store

import (
	"context"
	"math/rand"
	"fmt"
	"time"
	"github.com/VitoriaXaavier/Rest-GO-API/objects"
)

// Interface do banco para armazenar eventos
type IEventStore interface {
	Get (ctx context.Context, in *objects.GetRequest) (*objects.Event, error)
	List (ctx context.Context, in *objects.ListRequest) ([]*objects.Event, error)
	Create (ctx context.Context, in *objects.CreateRequest) error
	UpDateDetails (ctx context.Context, in *objects.UpDateDetailsRequest) error
	Cancel 	(ctx context.Context, in *objects.CancelRequest) error
	Remarca (ctx context.Context, in *objects.RemarcaRequest) error
	Delete (ctx context.Context, in *objects.DeletRequest) error
}

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

// Irá retornar um id exclusivel classificado com base no tempo
func GenerateUniqueId() string {
	word := []byte("0987654321")
	rand.Shuffle(len(word), func(i,j  int) {
		word[i], word[j] = word[j], word[i]
	})
	now := time.Now().UTC()
	timeFormatted := now.Format("06-01-02 15:04:05")
	return fmt.Sprintf("%s-%s",timeFormatted, string(word))
}