package emailer

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(sender EmailSender, recipient string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	deps := HandlerDependencies{Sender: sender}
	router.POST("/contact", CreateContactHandler(deps, recipient))
	return router
}

func TestCreateContactHandler_Success(t *testing.T) {
	assert := assert.New(t)

	mockSender := &mockEmailSender{
		returnedErr: nil,
	}
	testRecipient := "test-recipient@example.com"

	router := setupTestRouter(mockSender, testRecipient)

	formData := ContactFormData{
		FirstName: "Aka",
		LastName:  "Beex",
		Email:     "akaBeex@example.com",
		Phone:     "08222494342",
		Message:   "Test ini Test",
	}

	jsonData, _ := json.Marshal(formData)
	reqBody := bytes.NewBuffer(jsonData)

	req, _ := http.NewRequest(http.MethodPost, "/contact", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Expected status code 200 OK")

	expectedBody := `{"message":"message received succesfully","status":"success"}`

	assert.JSONEq(expectedBody, w.Body.String(), "Expected success JSON response")

	assert.True(mockSender.called, "Expected Sender.Send to be called")
	assert.Equal(testRecipient, mockSender.toArg, "Recipient email mismatch")
	assert.Contains(mockSender.subjectArg, "Aka Beex", "Subject should contain sender name")
	assert.Contains(mockSender.bodyArg, "akaBeex@example.com", "Body should contain sender email")
	assert.Contains(mockSender.bodyArg, "Test ini Test", "Body should contain the message")
	assert.Contains(mockSender.bodyArg, "08222494342", "Body should contain the phone")
}

func TestCreateContactHandler_BindError(t *testing.T) {
	assert := assert.New(t)
	mockSender := &mockEmailSender{}
	testRecipient := "test-recipient@example.com"

	router := setupTestRouter(mockSender, testRecipient)

	invalidData := `{"firstName":"Jane","lastName":"Doe","email": "jane.doe@example.com"}`
	reqBody := strings.NewReader(invalidData)

	req, _ := http.NewRequest(http.MethodPost, "/contact", reqBody)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Contains(w.Body.String(), `"status":"error"`, "Response should indicate error")
	assert.Contains(w.Body.String(), "Invalid form data", "Response should mention invalid data")

	assert.False(mockSender.called, "Sender.Send should NOT be called on binding error")
}

func TestCreateContactHandler_SendError(t *testing.T) {
	assert := assert.New(t)

	simulatedError := errors.New("Smtp server unavailable")
	mockSender := &mockEmailSender{
		returnedErr: simulatedError,
	}
	testRecipient := "test-recipient@example.com"
	router := setupTestRouter(mockSender, testRecipient)
	formData := ContactFormData{
		FirstName: "I",
		LastName:  "Failure",
		Email:     "failure@example.com",
		Message:   "This should fail.",
	}
	jsonData, _ := json.Marshal(formData)
	reqBody := bytes.NewBuffer(jsonData)

	req, _ := http.NewRequest(http.MethodPost, "/contact", reqBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code, "Expected status code 500 Internal Server Error")

	assert.Contains(w.Body.String(), `"status":"error"`, "Response should indicate error")
	assert.Contains(w.Body.String(), "internal server error", "Response should mention internal error")

	assert.True(mockSender.called, "Sender.Send should be called even if it returns an error")
}
