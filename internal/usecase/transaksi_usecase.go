package usecase

import (
	"bank-service-app/internal/domain"
	"fmt"
	"time"
)

type transaksiUsecase struct {
	nasabahRepo   domain.NasabahRepository
	transaksiRepo domain.TransaksiRepository
}

func NewTransaksiUsecase(nr domain.NasabahRepository, tr domain.TransaksiRepository) domain.TransaksiUsecase {
	return &transaksiUsecase{
		nasabahRepo:   nr,
		transaksiRepo: tr,
	}
}

func (u *transaksiUsecase) Tabung(noRekening string, nominal float64) (float64, error) {
	nasabah, err := u.nasabahRepo.GetByNoRekening(noRekening)
	if err != nil {
		return 0, fmt.Errorf("nomor rekening tidak ditemukan")
	}

	currentSaldo, err := u.transaksiRepo.GetSaldoByNasabahID(nasabah.ID)
	if err != nil {
		return 0, err
	}

	newSaldo := currentSaldo + nominal
	transaksi := &domain.Transaksi{
		NasabahID:      nasabah.ID,
		JenisTransaksi: "CREDIT",
		Nominal:        nominal,
		SaldoAkhir:     newSaldo,
		CreatedAt:      time.Now(),
	}

	if err := u.transaksiRepo.Create(transaksi); err != nil {
		return 0, err
	}

	return newSaldo, nil
}

func (u *transaksiUsecase) Tarik(noRekening string, nominal float64) (float64, error) {
	nasabah, err := u.nasabahRepo.GetByNoRekening(noRekening)
	if err != nil {
		return 0, fmt.Errorf("nomor rekening tidak ditemukan")
	}

	currentSaldo, err := u.transaksiRepo.GetSaldoByNasabahID(nasabah.ID)
	if err != nil {
		return 0, err
	}

	if currentSaldo < nominal {
		return 0, fmt.Errorf("saldo tidak mencukupi")
	}

	newSaldo := currentSaldo - nominal
	transaksi := &domain.Transaksi{
		NasabahID:      nasabah.ID,
		JenisTransaksi: "DEBIT",
		Nominal:        nominal,
		SaldoAkhir:     newSaldo,
		CreatedAt:      time.Now(),
	}

	if err := u.transaksiRepo.Create(transaksi); err != nil {
		return 0, err
	}

	return newSaldo, nil
}
