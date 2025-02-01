package usecase

import (
	"bank-service-app/internal/domain"
	"fmt"
	"math/rand"
	"time"
)

type nasabahUsecase struct {
	nasabahRepo   domain.NasabahRepository
	transaksiRepo domain.TransaksiRepository
}

func NewNasabahUsecase(nr domain.NasabahRepository, tr domain.TransaksiRepository) domain.NasabahUsecase {
	return &nasabahUsecase{
		nasabahRepo:   nr,
		transaksiRepo: tr,
	}
}

func (u *nasabahUsecase) Register(nama, nik, noHP string) (*domain.Nasabah, error) {
	// Check if NIK already exists
	if _, err := u.nasabahRepo.GetByNIK(nik); err == nil {
		return nil, fmt.Errorf("NIK sudah terdaftar")
	}

	// Check if NoHP already exists
	if _, err := u.nasabahRepo.GetByNoHP(noHP); err == nil {
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
		return nil, err
	}

	return nasabah, nil
}

func (u *nasabahUsecase) GetSaldo(noRekening string) (float64, error) {
	nasabah, err := u.nasabahRepo.GetByNoRekening(noRekening)
	if err != nil {
		return 0, fmt.Errorf("nomor rekening tidak ditemukan")
	}

	saldo, err := u.transaksiRepo.GetSaldoByNasabahID(nasabah.ID)
	if err != nil {
		return 0, err
	}

	return saldo, nil
}

func generateNoRekening() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%010d", r.Intn(9999999999))
}
