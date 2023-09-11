package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
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

	if !slot.StarTime.After(time.Time{}) {
		return errors.ErrInvalidTimeFormat
	}

	if !slot.EndTime.After(time.Time{}) {
		return errors.ErrInvalidTimeFormat
	}
	return nil
}

