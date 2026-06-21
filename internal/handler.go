package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

// GetRecordsHandler
// @Summary Получить список подписок
// @Tags records
// @Produce json
// @Success 200 {array} RecordDto
// @Failure 500 {string} string "внутренняя ошибка"
// @Router /records [get]
func (h *recordHandler) GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"method", r.Method,
		"path", r.URL.Path,
	)
	recordList, err := h.service.GetAll(r.Context())
	if err != nil {
		logger.Error("get records failed", "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	dtoList := make([]RecordDto, 0, len(recordList))
	for _, r := range recordList {
		dtoList = append(dtoList, toDto(r))
	}
	logger.Info("records found", "count", len(dtoList))
	err = json.NewEncoder(w).Encode(dtoList)
	if err != nil {
		logger.Error("error encoding dtoList", "error", err)
	}
}

// GetRecordHandler
// @Summary Получить подписку по ID
// @Tags records
// @Produce json
// @Param id path int true "ID записи"
// @Success 200 {object} RecordDto
// @Failure 400 {string} string "некорректный id"
// @Failure 404 {string} string "запись не найдена"
// @Failure 500 {string} string "внутренняя ошибка"
// @Router /records/{id} [get]
func (h *recordHandler) GetRecordHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"method", r.Method,
		"path", r.URL.Path,
	)
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Warn("invalid id", "strId", strId)
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	record, err := h.service.GetByID(r.Context(), id)
	switch {
	case errors.Is(err, ErrRecordNotFound):
		logger.Warn("record not found", "id", id)
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	case err != nil:
		logger.Error("get record failed", "id", id, "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	logger.Info("record found", "id", id)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(toDto(record))
	if err != nil {
		logger.Error("error encoding dto", "err", err)
	}
}

// PostRecordHandler
// @Summary Создать подписку
// @Tags records
// @Accept json
// @Produce json
// @Param request body RecordDto true "Данные подписки"
// @Success 201 {object} RecordDto "Созданная запись"
// @Failure 400 {string} string "некорректный JSON"
// @Failure 500 {string} string "внутренняя ошибка"
// @Router /records [post]
func (h *recordHandler) PostRecordHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"method", r.Method,
		"path", r.URL.Path,
	)
	var dto RecordDto
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		logger.Warn("invalid json", "error", err)
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	record, err := toModel(dto)
	if err != nil {
		logger.Error("parse dto failed", "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	record, err = h.service.Create(r.Context(), record)
	if err != nil {
		logger.Error("create record failed", "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	logger.Info("record created", "id", record.Id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(toDto(record))
	if err != nil {
		logger.Error("error encoding dto", "err", err)
	}
}

// PutRecordHandler
// @Summary Обновить подписку
// @Tags records
// @Accept json
// @Param id path int true "ID записи"
// @Param request body RecordDto true "Новые данные"
// @Success 204 "No Content"
// @Failure 400 {string} string "некорректный JSON или id"
// @Failure 404 {string} string "запись не найдена"
// @Failure 500 {string} string "внутренняя ошибка"
// @Router /records/{id} [put]
func (h *recordHandler) PutRecordHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"method", r.Method,
		"path", r.URL.Path,
	)
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Warn("invalid id", "strId", strId)
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	var dto RecordDto
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		logger.Warn("invalid json", "error", err)
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	record, err := toModel(dto)
	if err != nil {
		logger.Error("parse dto failed", "id", id, "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	err = h.service.Update(r.Context(), id, record)
	switch {
	case errors.Is(err, ErrRecordNotFound):
		logger.Warn("record not found", "id", id)
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	case err != nil:
		logger.Error("update record failed", "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	logger.Info("record updated", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteRecordHandler
// @Summary Удалить подписку
// @Tags records
// @Param id path int true "ID записи"
// @Success 204 "No Content"
// @Failure 400 {string} string "некорректный id"
// @Failure 404 {string} string "запись не найдена"
// @Failure 500 {string} string "внутренняя ошибка"
// @Router /records/{id} [delete]
func (h *recordHandler) DeleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"method", r.Method,
		"path", r.URL.Path,
	)
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Warn("invalid id", "strId", strId)
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	err = h.service.Delete(r.Context(), id)
	switch {
	case errors.Is(err, ErrRecordNotFound):
		logger.Warn("record not found", "id", id)
		http.Error(w, "запись не найдена", http.StatusNotFound)
		return
	case err != nil:
		logger.Error("delete record failed", "id", id, "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	logger.Info("record deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

// GetTotalPrice
// @Summary Получить сумму подписок за период
// @Tags records
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param from query string true "Начало периода (MM-YYYY)"
// @Param to query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} map[string]int
// @Failure 400 {string} string "некорректный формат даты"
// @Failure 500 {string} string "внутренняя ошибка"
// @Router /records/total [get]
func (h *recordHandler) GetTotalPrice(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"method", r.Method,
		"path", r.URL.Path,
	)
	q := r.URL.Query()
	userIDStr := q.Get("user_id")
	serviceName := q.Get("service_name")
	fromStr := q.Get("from")
	toStr := q.Get("to")
	from, err := parseMonthYear(fromStr)
	if err != nil {
		logger.Warn("invalid date", "fromStr", fromStr)
		http.Error(w, "некорректный формат даты", http.StatusBadRequest)
		return
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		logger.Warn("invalid date", "toStr", toStr)
		http.Error(w, "некорректный формат даты", http.StatusBadRequest)
		return
	}
	total, err := h.service.TotalPrice(r.Context(), userIDStr, serviceName, from, to)
	if err != nil {
		logger.Error("get total price failed", "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	logger.Info("total price calculated", "total", total)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]int{"total": total})
	if err != nil {
		logger.Error("error encoding map of total", "error", err)
	}
}
