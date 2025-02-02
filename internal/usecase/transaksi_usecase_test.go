package usecase_test

import (
	"bank-service-app/internal/domain"
	"bank-service-app/internal/domain/mocks"
	"bank-service-app/internal/usecase"
	"bank-service-app/pkg/logger"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransaksiUsecase_Tabung(t *testing.T) {
	mockNasabahRepo := new(mocks.NasabahRepository)
	mockTransaksiRepo := new(mocks.TransaksiRepository)
	log := logger.NewLogger()

	transaksiUC := usecase.NewTransaksiUsecase(mockNasabahRepo, mockTransaksiRepo, log)

	t.Run("Success Deposit", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(1000000), nil).Once()
		mockTransaksiRepo.On("Create", mock.MatchedBy(func(transaksi *domain.Transaksi) bool {
			return transaksi.NasabahID == int64(1) &&
				transaksi.JenisTransaksi == "CREDIT" &&
				transaksi.Nominal == float64(500000) &&
				transaksi.SaldoAkhir == float64(1500000)
		})).Return(nil).Once()

		newSaldo, err := transaksiUC.Tabung("1234567890", 500000)

		assert.NoError(t, err)
		assert.Equal(t, float64(1500000), newSaldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})

	t.Run("Account Not Found", func(t *testing.T) {
		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nil, errors.New("not found")).Once()

		newSaldo, err := transaksiUC.Tabung("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		assert.Contains(t, err.Error(), "nomor rekening tidak ditemukan")
		mockNasabahRepo.AssertExpectations(t)
	})

	t.Run("Get Saldo Failed", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(0), errors.New("db error")).Once()

		newSaldo, err := transaksiUC.Tabung("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})

	t.Run("Create Transaction Failed", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(1000000), nil).Once()
		mockTransaksiRepo.On("Create", mock.AnythingOfType("*domain.Transaksi")).Return(errors.New("db error")).Once()

		newSaldo, err := transaksiUC.Tabung("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})
}

func TestTransaksiUsecase_Tarik(t *testing.T) {
	mockNasabahRepo := new(mocks.NasabahRepository)
	mockTransaksiRepo := new(mocks.TransaksiRepository)
	log := logger.NewLogger()

	transaksiUC := usecase.NewTransaksiUsecase(mockNasabahRepo, mockTransaksiRepo, log)

	t.Run("Success Withdraw", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(1000000), nil).Once()
		mockTransaksiRepo.On("Create", mock.MatchedBy(func(transaksi *domain.Transaksi) bool {
			return transaksi.NasabahID == int64(1) &&
				transaksi.JenisTransaksi == "DEBIT" &&
				transaksi.Nominal == float64(500000) &&
				transaksi.SaldoAkhir == float64(500000)
		})).Return(nil).Once()

		newSaldo, err := transaksiUC.Tarik("1234567890", 500000)

		assert.NoError(t, err)
		assert.Equal(t, float64(500000), newSaldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})

	t.Run("Account Not Found", func(t *testing.T) {
		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nil, errors.New("not found")).Once()

		newSaldo, err := transaksiUC.Tarik("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		assert.Contains(t, err.Error(), "nomor rekening tidak ditemukan")
		mockNasabahRepo.AssertExpectations(t)
	})

	t.Run("Insufficient Balance", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(400000), nil).Once()

		newSaldo, err := transaksiUC.Tarik("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		assert.Contains(t, err.Error(), "saldo tidak mencukupi")
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})

	t.Run("Get Saldo Failed", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(0), errors.New("db error")).Once()

		newSaldo, err := transaksiUC.Tarik("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})

	t.Run("Create Transaction Failed", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(1000000), nil).Once()
		mockTransaksiRepo.On("Create", mock.AnythingOfType("*domain.Transaksi")).Return(errors.New("db error")).Once()

		newSaldo, err := transaksiUC.Tarik("1234567890", 500000)

		assert.Error(t, err)
		assert.Equal(t, float64(0), newSaldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})
}
