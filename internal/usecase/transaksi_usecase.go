package usecase

import (
	"bank-service-app/internal/domain"
	"bank-service-app/pkg/logger"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type transaksiUsecase struct {
	nasabahRepo   domain.NasabahRepository
	transaksiRepo domain.TransaksiRepository
	log           *logger.Logger
}

func NewTransaksiUsecase(nr domain.NasabahRepository, tr domain.TransaksiRepository, log *logger.Logger) domain.TransaksiUsecase {
	return &transaksiUsecase{
		nasabahRepo:   nr,
		transaksiRepo: tr,
		log:           log,
	}
}

func (u *transaksiUsecase) Tabung(noRekening string, nominal float64) (float64, error) {
	u.log.InfoWithContext("Starting Deposit Process",
		zap.String("no_rekening", noRekening),
		zap.Float64("nominal", nominal),
	)

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

	u.log.InfoWithContext("Deposit Process Success",
		zap.String("no_rekening", noRekening),
		zap.Int64("transaksi_nasabah_id", transaksi.NasabahID),
		zap.String("transaksi_jenis", transaksi.JenisTransaksi),
		zap.Float64("transaksi_nominal", transaksi.Nominal),
		zap.Float64("transaksi_saldo_akhir", transaksi.SaldoAkhir),
	)

	return newSaldo, nil
}

func (u *transaksiUsecase) Tarik(noRekening string, nominal float64) (float64, error) {
	u.log.InfoWithContext("Starting Withdraw Process",
		zap.String("no_rekening", noRekening),
		zap.Float64("nominal", nominal),
	)

	nasabah, err := u.nasabahRepo.GetByNoRekening(noRekening)
	if err != nil {
		return 0, fmt.Errorf("nomor rekening tidak ditemukan")
	}

	currentSaldo, err := u.transaksiRepo.GetSaldoByNasabahID(nasabah.ID)
	if err != nil {
		return 0, err
	}

	if currentSaldo < nominal {
		u.log.WarnWithContext("Insufficient Balance",
			zap.String("no_rekening", noRekening),
			zap.Float64("saldo", currentSaldo),
			zap.Float64("nominal", nominal),
		)
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

	u.log.InfoWithContext("Withdraw Process Success",
		zap.String("no_rekening", noRekening),
		zap.Int64("transaksi_nasabah_id", transaksi.NasabahID),
		zap.String("transaksi_jenis", transaksi.JenisTransaksi),
		zap.Float64("transaksi_nominal", transaksi.Nominal),
		zap.Float64("transaksi_saldo_akhir", transaksi.SaldoAkhir),
	)

	return newSaldo, nil
}
