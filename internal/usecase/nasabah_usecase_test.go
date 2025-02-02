package usecase_test

import (
	"bank-service-app/internal/domain"
	"bank-service-app/internal/domain/mocks"
	"bank-service-app/internal/usecase"
	"bank-service-app/pkg/logger"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNasabahUsecase_Register(t *testing.T) {
	mockNasabahRepo := new(mocks.NasabahRepository)
	mockTransaksiRepo := new(mocks.TransaksiRepository)
	log := logger.NewLogger()

	nasabahUC := usecase.NewNasabahUsecase(mockNasabahRepo, mockTransaksiRepo, log)

	t.Run("Success Register", func(t *testing.T) {
		mockNasabahRepo.On("GetByNIK", "1234567890").Return(nil, errors.New("not found")).Once()
		mockNasabahRepo.On("GetByNoHP", "08123456789").Return(nil, errors.New("not found")).Once()
		mockNasabahRepo.On("Create", mock.AnythingOfType("*domain.Nasabah")).Return(nil).Once()

		nasabah, err := nasabahUC.Register("John Doe", "1234567890", "08123456789")

		assert.NoError(t, err)
		assert.NotNil(t, nasabah)
		assert.Equal(t, "John Doe", nasabah.Nama)
		assert.Equal(t, "1234567890", nasabah.NIK)
		assert.Equal(t, "08123456789", nasabah.NoHP)
		assert.Len(t, nasabah.NoRekening, 10)
		mockNasabahRepo.AssertExpectations(t)
	})

	t.Run("NIK Already Exists", func(t *testing.T) {
		existingNasabah := &domain.Nasabah{
			ID:         1,
			Nama:       "Existing User",
			NIK:        "1234567890",
			NoHP:       "08111111111",
			NoRekening: "1234567890",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockNasabahRepo.On("GetByNIK", "1234567890").Return(existingNasabah, nil).Once()

		nasabah, err := nasabahUC.Register("John Doe", "1234567890", "08123456789")

		assert.Error(t, err)
		assert.Nil(t, nasabah)
		assert.Contains(t, err.Error(), "NIK sudah terdaftar")
		mockNasabahRepo.AssertExpectations(t)
	})

	t.Run("NoHP Already Exists", func(t *testing.T) {
		mockNasabahRepo.On("GetByNIK", "1234567890").Return(nil, errors.New("not found")).Once()

		existingNasabah := &domain.Nasabah{
			ID:         1,
			Nama:       "Existing User",
			NIK:        "9876543210",
			NoHP:       "08123456789",
			NoRekening: "1234567890",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockNasabahRepo.On("GetByNoHP", "08123456789").Return(existingNasabah, nil).Once()

		nasabah, err := nasabahUC.Register("John Doe", "1234567890", "08123456789")

		assert.Error(t, err)
		assert.Nil(t, nasabah)
		assert.Contains(t, err.Error(), "nomor hp sudah terdaftar")
		mockNasabahRepo.AssertExpectations(t)
	})

	t.Run("Create Failed", func(t *testing.T) {
		mockNasabahRepo.On("GetByNIK", "1234567890").Return(nil, errors.New("not found")).Once()
		mockNasabahRepo.On("GetByNoHP", "08123456789").Return(nil, errors.New("not found")).Once()
		mockNasabahRepo.On("Create", mock.AnythingOfType("*domain.Nasabah")).Return(errors.New("db error")).Once()

		nasabah, err := nasabahUC.Register("John Doe", "1234567890", "08123456789")

		assert.Error(t, err)
		assert.Nil(t, nasabah)
		mockNasabahRepo.AssertExpectations(t)
	})
}

func TestNasabahUsecase_GetSaldo(t *testing.T) {
	mockNasabahRepo := new(mocks.NasabahRepository)
	mockTransaksiRepo := new(mocks.TransaksiRepository)
	log := logger.NewLogger()

	nasabahUC := usecase.NewNasabahUsecase(mockNasabahRepo, mockTransaksiRepo, log)

	t.Run("Success Get Saldo", func(t *testing.T) {
		nasabah := &domain.Nasabah{
			ID:         1,
			NoRekening: "1234567890",
		}

		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nasabah, nil).Once()
		mockTransaksiRepo.On("GetSaldoByNasabahID", int64(1)).Return(float64(1000000), nil).Once()

		saldo, err := nasabahUC.GetSaldo("1234567890")

		assert.NoError(t, err)
		assert.Equal(t, float64(1000000), saldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})

	t.Run("Account Not Found", func(t *testing.T) {
		mockNasabahRepo.On("GetByNoRekening", "1234567890").Return(nil, errors.New("not found")).Once()

		saldo, err := nasabahUC.GetSaldo("1234567890")

		assert.Error(t, err)
		assert.Equal(t, float64(0), saldo)
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

		saldo, err := nasabahUC.GetSaldo("1234567890")

		assert.Error(t, err)
		assert.Equal(t, float64(0), saldo)
		mockNasabahRepo.AssertExpectations(t)
		mockTransaksiRepo.AssertExpectations(t)
	})
}
