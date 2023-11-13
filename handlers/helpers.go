package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	//"time"
	"github.com/VitoriaXaavier/Rest-GO-API/errors"
	"github.com/VitoriaXaavier/Rest-GO-API/objects"
)

type Response interface {
	JSON() []byte
	StatusCode() int
}

func WriterResponse(w http.ResponseWriter, resp Response) {
	w.WriteHeader(resp.StatusCode())
	_, _ = w.Write(resp.JSON())
}

func WriterError(w http.ResponseWriter, err error) {
	res, ok := err.(*errors.Error)
	if !ok {
		log.Println(err)
		res = errors.ErrInternal
	}
	WriterResponse(w, res)
}

func IntFromString(w http.ResponseWriter, v string) (int, error) {
	if v == "" {
		return 0, nil
	}
	res, err := strconv.Atoi(v)
	if err != nil {
		log.Panicln(err)
		WriterError(w, errors.ErrInvalidLimit)
	}
	return res, err
}

func Unmarshal( w http.ResponseWriter, data []byte, v interface{}) (error) {
	if d := string(data); d == "null" || d == "" {
		WriterError(w, errors.ErrObjectIsRequired)
		return errors.ErrObjectIsRequired
	}

	err := json.Unmarshal(data, v) 
	if err != nil {
		log.Panicln(err)
		WriterError(w, errors.ErrBadRequest)
	}
	return err
}

func ChekSlot(slot *objects.TimeSlot) error {
	if slot == nil {
		return errors.ErrEventTimingIsRequired
	}
	layout := "2006-01-02 15:04:05"
	if slot.StarTime.IsZero() || slot.StarTime.Format(layout) != slot.StarTime.Format(layout) {
		return errors.ErrInvalidTimeFormat
	}

	// Verifica se EndTime está no formato desejado
	if slot.EndTime.IsZero() || slot.EndTime.Format(layout) != slot.EndTime.Format(layout) {
		return errors.ErrInvalidTimeFormat
	}
	return nil
}

