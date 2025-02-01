package http

import (
	"bank-service-app/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NasabahHandler struct {
	nasabahUsecase   domain.NasabahUsecase
	transaksiUsecase domain.TransaksiUsecase
}

func NewNasabahHandler(e *echo.Echo, nu domain.NasabahUsecase, tu domain.TransaksiUsecase) {
	handler := &NasabahHandler{
		nasabahUsecase:   nu,
		transaksiUsecase: tu,
	}

	// Route List
	e.POST("/daftar", handler.Register)
	e.POST("/tabung", handler.Tabung)
	e.POST("/tarik", handler.Tarik)
	e.GET("/saldo/:no_rekening", handler.GetSaldo)
}

// Register handles nasabah registration
func (h *NasabahHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: "Invalid request format",
		})
	}

	nasabah, err := h.nasabahUsecase.Register(req.Nama, req.NIK, req.NoHP)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, RegisterResponse{
		NoRekening: nasabah.NoRekening,
	})
}

// Tabung handles deposit
func (h *NasabahHandler) Tabung(c echo.Context) error {
	var req TransaksiRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: "Invalid request format",
		})
	}

	saldo, err := h.transaksiUsecase.Tabung(req.NoRekening, req.Nominal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SaldoResponse{
		Saldo: saldo,
	})
}

// Tarik handles withdrawal
func (h *NasabahHandler) Tarik(c echo.Context) error {
	var req TransaksiRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: "Invalid request format",
		})
	}

	saldo, err := h.transaksiUsecase.Tarik(req.NoRekening, req.Nominal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SaldoResponse{
		Saldo: saldo,
	})
}

// GetSaldo handles balance inquiry
func (h *NasabahHandler) GetSaldo(c echo.Context) error {
	noRekening := c.Param("no_rekening")
	if noRekening == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: "Nomor rekening tidak boleh kosong",
		})
	}

	saldo, err := h.nasabahUsecase.GetSaldo(noRekening)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SaldoResponse{
		Saldo: saldo,
	})
}
