package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type recordHandler struct {
	service *recordService
}

func NewRecordHandler(s *recordService) *recordHandler {
	return &recordHandler{
		service: s,
	}
}

func (h *recordHandler) GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	recordList, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recordList)
}

func (h *recordHandler) GetRecordHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	record, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func (h *recordHandler) PostRecordHandler(w http.ResponseWriter, r *http.Request) {
	var record Record
	json.NewDecoder(r.Body).Decode(&record)
	defer r.Body.Close()
	record, err := h.service.Create(record)
	if err != nil {
		http.Error(w, "ошибка создания записи", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}

func (h *recordHandler) PutRecordHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	var record Record
	json.NewDecoder(r.Body).Decode(&record)
	defer r.Body.Close()
	err = h.service.Update(id, record)
	if err != nil {
		http.Error(w, "ошибка обновления записи", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
