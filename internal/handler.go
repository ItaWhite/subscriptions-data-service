package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type recordHandler struct {
	service *recordService
}

func NewRecordHandler(s *recordService) *recordHandler {
	return &recordHandler{
		service: s,
	}
}

func parseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func toModel(dto RecordDto) (Record, error) {
	userId, err := uuid.Parse(dto.UserID)
	if err != nil {
		return Record{}, err
	}
	start, err := parseMonthYear(dto.StartDate)
	if err != nil {
		return Record{}, err
	}
	var end *time.Time
	if dto.EndDate != "" {
		t, err := parseMonthYear(dto.EndDate)
		if err != nil {
			return Record{}, err
		}
		end = &t
	}
	return Record{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      userId,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func toDto(record Record) RecordDto {
	start := fmt.Sprintf("%02d-%d", record.StartDate.Month(), record.StartDate.Year())
	var end string
	if record.EndDate != nil {
		end = fmt.Sprintf("%02d-%d", record.EndDate.Month(), record.EndDate.Year())
	}
	return RecordDto{
		Id:          record.Id,
		ServiceName: record.ServiceName,
		Price:       record.Price,
		UserID:      record.UserID.String(),
		StartDate:   start,
		EndDate:     end,
	}
}

func (h *recordHandler) GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	recordList, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	dtoList := make([]RecordDto, 0, len(recordList))
	for _, r := range recordList {
		dtoList = append(dtoList, toDto(r))
	}
	json.NewEncoder(w).Encode(dtoList)
}

func (h *recordHandler) GetRecordHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	record, err := h.service.GetByID(r.Context(), id)
	switch {
	case errors.Is(err, ErrRecordNotFound):
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toDto(record))
}

func (h *recordHandler) PostRecordHandler(w http.ResponseWriter, r *http.Request) {
	var dto RecordDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	record, err := toModel(dto)
	if err != nil {
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	record, err = h.service.Create(r.Context(), record)
	if err != nil {
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toDto(record))
}

func (h *recordHandler) PutRecordHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	var dto RecordDto
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	record, err := toModel(dto)
	if err != nil {
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	err = h.service.Update(r.Context(), id, record)
	switch {
	case errors.Is(err, ErrRecordNotFound):
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *recordHandler) DeleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	err = h.service.Delete(r.Context(), id)
	switch {
	case errors.Is(err, ErrRecordNotFound):
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
