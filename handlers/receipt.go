package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/suryamp/receipt-processor/models"
    "github.com/suryamp/receipt-processor/processor"
    "github.com/suryamp/receipt-processor/validator"
)

type Handler struct {
    processor processor.ReceiptProcessor
}

func NewHandler(p processor.ReceiptProcessor) *Handler {
    return &Handler{processor: p}
}

func (h *Handler) ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
    var receipt models.Receipt
    
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
        return
    }

    if err := validator.ValidateReceipt(receipt); err != nil {
        http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
        return
    }

    id, err := h.processor.ProcessReceipt(receipt)
    if err != nil {
        http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.ProcessResponse{ID: id})
}

func (h *Handler) GetPointsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    points, err := h.processor.GetPoints(id)
    if err != nil {
        http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.PointsResponse{Points: points})
}