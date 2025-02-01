package usecase

import (
	"bank-service-app/internal/domain"
	"bank-service-app/pkg/logger"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type nasabahUsecase struct {
	nasabahRepo   domain.NasabahRepository
	transaksiRepo domain.TransaksiRepository
	log           *logger.Logger
}

func NewNasabahUsecase(nr domain.NasabahRepository, tr domain.TransaksiRepository, log *logger.Logger) domain.NasabahUsecase {
	return &nasabahUsecase{
		nasabahRepo:   nr,
		transaksiRepo: tr,
		log:           log,
	}
}

func (u *nasabahUsecase) Register(nama, nik, noHP string) (*domain.Nasabah, error) {
	u.log.InfoWithContext("Starting nasabah registration",
		zap.String("nama", nama),
		zap.String("nik", nik),
		zap.String("no_hp", noHP),
	)

	// Check if NIK already exists
	if _, err := u.nasabahRepo.GetByNIK(nik); err == nil {
		u.log.WarnWithContext("NIK already registered",
			zap.String("nik", nik),
		)
		return nil, fmt.Errorf("NIK sudah terdaftar")
	}

	// Check if NoHP already exists
	if _, err := u.nasabahRepo.GetByNoHP(noHP); err == nil {
		u.log.WarnWithContext("Phone number already registered",
			zap.String("no_hp", noHP),
		)
		return nil, fmt.Errorf("nomor hp sudah terdaftar")
	}

	// Generate unique account number
	noRekening := generateNoRekening()

	nasabah := &domain.Nasabah{
		Nama:       nama,
		NIK:        nik,
		NoHP:       noHP,
		NoRekening: noRekening,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := u.nasabahRepo.Create(nasabah); err != nil {
		u.log.ErrorWithContext("Failed to create nasabah", err,
			zap.String("nik", nik),
			zap.String("no_hp", noHP),
		)
		return nil, err
	}

	u.log.InfoWithContext("Successfully registered nasabah",
		zap.String("no_rekening", noRekening),
		zap.String("nik", nik),
	)

	return nasabah, nil
}

func (u *nasabahUsecase) GetSaldo(noRekening string) (float64, error) {
	nasabah, err := u.nasabahRepo.GetByNoRekening(noRekening)
	if err != nil {
		u.log.ErrorWithContext("Failed to get nasabah by no rekening", err, zap.String("no_rekening", noRekening))
		return 0, fmt.Errorf("nomor rekening tidak ditemukan")
	}

	saldo, err := u.transaksiRepo.GetSaldoByNasabahID(nasabah.ID)
	if err != nil {
		u.log.ErrorWithContext("Failed to get saldo by nasabah id", err, zap.String("no_rekening", noRekening))
		return 0, err
	}

	return saldo, nil
}

func generateNoRekening() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%010d", r.Intn(9999999999))
}
