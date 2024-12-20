package server

import (
	"calc-service/internal/calculator"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExpressionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendServerErrorResponse(w, "server error") //500
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expression := req.Expression
	if expression == "" {
		sendErrorResponse(w, "Invalid expression: empty string") //422
		return
	}

	result, err := calculator.Calc(expression)
	if err != nil {
		sendErrorResponse(w, err.Error()) //422
		return
	}

	sendSuccessResponse(w, fmt.Sprintf("%v", result)) // 200
}

// succes answer
func sendSuccessResponse(w http.ResponseWriter, result string) {
	response := Response{Result: result}
	sendJSONResponse(w, response, http.StatusOK)
}

// fail answer
func sendErrorResponse(w http.ResponseWriter, errorMessage string) {
	response := Response{Error: errorMessage}
	sendJSONResponse(w, response, http.StatusUnprocessableEntity)
}

// server error answer
func sendServerErrorResponse(w http.ResponseWriter, errorMessage string) {
	response := Response{Error: errorMessage}
	sendJSONResponse(w, response, http.StatusInternalServerError)
}

// func for sending json answer
func sendJSONResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
