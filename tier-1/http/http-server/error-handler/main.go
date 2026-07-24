package main

import (
	"encoding/json"
	"fmt"
	"tier-1/error-handler/domain"
)

// Mapeo de codigo de dominio -> status HTTP.
// Vive en la capa de interfaces (handlers), NO en el dominio --
// el dominio no debe saber que existe HTTP.
var codeToHTTPStatus = map[domain.ErrorCode]int{
	domain.ErrKeyNotFound:        404,
	domain.ErrTenantNotFound:     404,
	domain.ErrKeyInactive:        409, // conflicto de estado, no "no encontrado"
	domain.ErrInvalidInput:       400,
	domain.ErrHSMOperationFailed: 502, // el HSM es un backend externo -- bad gateway, no 500 generico
	domain.ErrK8sSecretNotFound:  500, // esto SI es un fallo de infra propio, 500 es correcto aqui
}

// public response - code + message, Details never
type errorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func handleError(err error) (statusCode int, body errorResponse, logLine string) {
	var de *domain.DomainError
	if domainErr, ok := err.(*domain.DomainError); ok {
		de = domainErr
	}

	if de == nil {
		// error no reconocido -- nunca exponer err.Error() crudo
		return 500, errorResponse{
			Error:   "internal_error",
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occurred",
		}, fmt.Sprintf("unmapped error: %v", err)
	}

	status, ok := codeToHTTPStatus[de.Code]
	if !ok {
		status = 500
	}

	// Log interno: SI incluye Details (secretName, namespace, etc.)
	logJSON, _ := json.Marshal(de.Details)
	logLine = fmt.Sprintf("code=%s message=%s details=%s", de.Code, de.Message, logJSON)

	// Respuesta al cliente: NUNCA Details
	return status, errorResponse{
		Error:   "domain_error",
		Code:    string(de.Code),
		Message: de.Message,
	}, logLine
}

func main() {
	casos := []*domain.DomainError{
		{Code: domain.ErrKeyNotFound, Message: "The requested key was not found"},
		{
			Code:    domain.ErrK8sSecretNotFound,
			Message: "The Kubernetes secret was not found",
			Details: map[string]interface{}{
				"secretName": "hsm-master-keys",
				"namespace":  "production",
			},
		},
		{Code: domain.ErrHSMOperationFailed, Message: "The HSM operation failed"},
	}

	for _, de := range casos {
		status, body, logLine := handleError(de)
		bodyJSON, _ := json.Marshal(body)
		fmt.Printf("--- %s ---\n", de.Code)
		fmt.Println("  HTTP status:", status)
		fmt.Println("  Client response:", string(bodyJSON))
		fmt.Println("  Log (with Details):", logLine)
		fmt.Println()
	}
}
