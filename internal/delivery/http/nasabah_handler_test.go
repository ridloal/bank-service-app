package http_test

import (
	handler "bank-service-app/internal/delivery/http"
	"bank-service-app/internal/domain"
	"bank-service-app/internal/domain/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupHandler() (*echo.Echo, *mocks.NasabahUsecase, *mocks.TransaksiUsecase) {
	e := echo.New()
	mockNasabahUC := new(mocks.NasabahUsecase)
	mockTransaksiUC := new(mocks.TransaksiUsecase)
	handler.NewNasabahHandler(e, mockNasabahUC, mockTransaksiUC)
	return e, mockNasabahUC, mockTransaksiUC
}

func TestNasabahHandler_Register(t *testing.T) {
	e, mockNasabahUC, _ := setupHandler()

	t.Run("Success Register", func(t *testing.T) {
		reqBody := `{
			"nama": "John Doe",
			"nik": "1234567890",
			"no_hp": "08123456789"
		}`

		expectedNasabah := &domain.Nasabah{
			Nama:       "John Doe",
			NIK:        "1234567890",
			NoHP:       "08123456789",
			NoRekening: "1234567890",
		}

		mockNasabahUC.On("Register",
			"John Doe",
			"1234567890",
			"08123456789",
		).Return(expectedNasabah, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/daftar", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "1234567890", response["no_rekening"])

		mockNasabahUC.AssertExpectations(t)
	})

	t.Run("Invalid Request Format", func(t *testing.T) {
		reqBody := `invalid json`

		req := httptest.NewRequest(http.MethodPost, "/daftar", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request format", response["remark"])
	})

	t.Run("Registration Failed", func(t *testing.T) {
		reqBody := `{
			"nama": "John Doe",
			"nik": "1234567890",
			"no_hp": "08123456789"
		}`

		mockNasabahUC.On("Register",
			"John Doe",
			"1234567890",
			"08123456789",
		).Return(nil, errors.New("NIK sudah terdaftar")).Once()

		req := httptest.NewRequest(http.MethodPost, "/daftar", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "NIK sudah terdaftar", response["remark"])

		mockNasabahUC.AssertExpectations(t)
	})
}

func TestNasabahHandler_Tabung(t *testing.T) {
	e, _, mockTransaksiUC := setupHandler()

	t.Run("Success Deposit", func(t *testing.T) {
		reqBody := `{
			"no_rekening": "1234567890",
			"nominal": 500000
		}`

		mockTransaksiUC.On("Tabung",
			"1234567890",
			float64(500000),
		).Return(float64(1500000), nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/tabung", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1500000), response["saldo"])

		mockTransaksiUC.AssertExpectations(t)
	})

	t.Run("Invalid Request Format", func(t *testing.T) {
		reqBody := `invalid json`

		req := httptest.NewRequest(http.MethodPost, "/tabung", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request format", response["remark"])
	})

	t.Run("Deposit Failed", func(t *testing.T) {
		reqBody := `{
			"no_rekening": "1234567890",
			"nominal": 500000
		}`

		mockTransaksiUC.On("Tabung",
			"1234567890",
			float64(500000),
		).Return(float64(0), errors.New("nomor rekening tidak ditemukan")).Once()

		req := httptest.NewRequest(http.MethodPost, "/tabung", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "nomor rekening tidak ditemukan", response["remark"])

		mockTransaksiUC.AssertExpectations(t)
	})
}

func TestNasabahHandler_Tarik(t *testing.T) {
	e, _, mockTransaksiUC := setupHandler()

	t.Run("Success Withdraw", func(t *testing.T) {
		reqBody := `{
			"no_rekening": "1234567890",
			"nominal": 500000
		}`

		mockTransaksiUC.On("Tarik",
			"1234567890",
			float64(500000),
		).Return(float64(500000), nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/tarik", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(500000), response["saldo"])

		mockTransaksiUC.AssertExpectations(t)
	})

	t.Run("Invalid Request Format", func(t *testing.T) {
		reqBody := `invalid json`

		req := httptest.NewRequest(http.MethodPost, "/tarik", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request format", response["remark"])
	})

	t.Run("Withdraw Failed - Insufficient Balance", func(t *testing.T) {
		reqBody := `{
			"no_rekening": "1234567890",
			"nominal": 500000
		}`

		mockTransaksiUC.On("Tarik",
			"1234567890",
			float64(500000),
		).Return(float64(0), errors.New("saldo tidak mencukupi")).Once()

		req := httptest.NewRequest(http.MethodPost, "/tarik", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "saldo tidak mencukupi", response["remark"])

		mockTransaksiUC.AssertExpectations(t)
	})
}

func TestNasabahHandler_GetSaldo(t *testing.T) {
	e, mockNasabahUC, _ := setupHandler()

	t.Run("Success Get Balance", func(t *testing.T) {
		mockNasabahUC.On("GetSaldo", "1234567890").Return(float64(1000000), nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/saldo/1234567890", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1000000), response["saldo"])

		mockNasabahUC.AssertExpectations(t)
	})

	t.Run("Account Not Found", func(t *testing.T) {
		mockNasabahUC.On("GetSaldo", "1234567890").Return(float64(0), errors.New("nomor rekening tidak ditemukan")).Once()

		req := httptest.NewRequest(http.MethodGet, "/saldo/1234567890", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "nomor rekening tidak ditemukan", response["remark"])

		mockNasabahUC.AssertExpectations(t)
	})
}
