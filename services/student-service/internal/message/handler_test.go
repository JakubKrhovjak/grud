package message_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"student-service/internal/auth"
	"student-service/internal/message"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SaramaProducerAdapter adapts sarama.SyncProducer to message.Producer interface
type SaramaProducerAdapter struct {
	producer sarama.SyncProducer
	topic    string
}

func (a *SaramaProducerAdapter) SendMessage(key string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: a.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueBytes),
	}

	_, _, err = a.producer.SendMessage(msg)
	return err
}

func (a *SaramaProducerAdapter) Close() error {
	return a.producer.Close()
}

// setupTest creates test dependencies
func setupTest(t *testing.T, mockProducer sarama.SyncProducer) chi.Router {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	adapter := &SaramaProducerAdapter{
		producer: mockProducer,
		topic:    "test-topic",
	}

	service := message.NewService(adapter, logger)
	handler := message.NewHandler(service, logger)

	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	return router
}

func TestMessageHandler(t *testing.T) {
	mockProducer := mocks.NewSyncProducer(t, nil)
	router := setupTest(t, mockProducer)

	t.Run("SendMessage_Success", func(t *testing.T) {
		mockProducer.ExpectSendMessageAndSucceed()

		// Create request payload
		payload := message.SendMessageRequest{
			Message: "Hello from test!",
		}
		body, err := json.Marshal(payload)
		require.NoError(t, err)

		// Create request with auth context
		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Add email to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), auth.EmailKey, "test@example.com")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		// Execute request
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "message sent successfully", response["message"])

		// Sarama mock automatically verifies expectation was met
	})

	t.Run("SendMessage_Unauthorized_NoEmail", func(t *testing.T) {

		payload := message.SendMessageRequest{
			Message: "Hello from test!",
		}
		body, err := json.Marshal(payload)
		require.NoError(t, err)

		// Create request WITHOUT email in context
		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// No message sent - Sarama mock verifies automatically
	})

	t.Run("SendMessage_InvalidJSON", func(t *testing.T) {

		// Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), auth.EmailKey, "test@example.com")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SendMessage_EmptyMessage", func(t *testing.T) {
		// Empty message (should fail validation)
		payload := message.SendMessageRequest{
			Message: "",
		}
		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), auth.EmailKey, "test@example.com")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 400 Bad Request (validation failed)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SendMessage_KafkaError", func(t *testing.T) {
		mockProducer.ExpectSendMessageAndFail(sarama.ErrOutOfBrokers)

		payload := message.SendMessageRequest{
			Message: "Hello from test!",
		}
		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), auth.EmailKey, "test@example.com")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("SendMessage_MultipleMessages", func(t *testing.T) {
		// Expect 2 messages
		mockProducer.ExpectSendMessageAndSucceed()
		mockProducer.ExpectSendMessageAndSucceed()

		// Send first message
		payload1 := message.SendMessageRequest{
			Message: "First message",
		}
		body1, _ := json.Marshal(payload1)

		req1 := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body1))
		req1.Header.Set("Content-Type", "application/json")
		ctx1 := context.WithValue(req1.Context(), auth.EmailKey, "user1@example.com")
		req1 = req1.WithContext(ctx1)

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Send second message
		payload2 := message.SendMessageRequest{
			Message: "Second message",
		}
		body2, _ := json.Marshal(payload2)

		req2 := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body2))
		req2.Header.Set("Content-Type", "application/json")
		ctx2 := context.WithValue(req2.Context(), auth.EmailKey, "user2@example.com")
		req2 = req2.WithContext(ctx2)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Sarama mock automatically verifies both messages were sent
	})

	t.Run("SendMessage_NetworkTimeout", func(t *testing.T) {
		mockProducer.ExpectSendMessageAndFail(sarama.ErrRequestTimedOut)

		payload := message.SendMessageRequest{
			Message: "This will timeout",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), auth.EmailKey, "test@example.com")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("SendMessage_PartialFailure", func(t *testing.T) {

		// First succeeds, second fails
		mockProducer.ExpectSendMessageAndSucceed()
		mockProducer.ExpectSendMessageAndFail(sarama.ErrNotLeaderForPartition)

		// First message succeeds
		payload1 := message.SendMessageRequest{Message: "Should succeed"}
		body1, _ := json.Marshal(payload1)
		req1 := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body1))
		req1.Header.Set("Content-Type", "application/json")
		ctx1 := context.WithValue(req1.Context(), auth.EmailKey, "user1@example.com")
		req1 = req1.WithContext(ctx1)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second message fails
		payload2 := message.SendMessageRequest{Message: "Should fail"}
		body2, _ := json.Marshal(payload2)
		req2 := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body2))
		req2.Header.Set("Content-Type", "application/json")
		ctx2 := context.WithValue(req2.Context(), auth.EmailKey, "user2@example.com")
		req2 = req2.WithContext(ctx2)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusInternalServerError, w2.Code)
	})
}
