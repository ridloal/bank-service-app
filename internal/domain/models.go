package domain

import "time"

type Nasabah struct {
	ID         int64     `json:"id"`
	Nama       string    `json:"nama"`
	NIK        string    `json:"nik"`
	NoHP       string    `json:"no_hp"`
	NoRekening string    `json:"no_rekening"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Transaksi struct {
	ID             int64     `json:"id"`
	NasabahID      int64     `json:"nasabah_id"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	Nominal        float64   `json:"nominal"`
	SaldoAkhir     float64   `json:"saldo_akhir"`
	CreatedAt      time.Time `json:"created_at"`
}

// Repository interfaces
type NasabahRepository interface {
	Create(nasabah *Nasabah) error
	GetByNoRekening(noRekening string) (*Nasabah, error)
	GetByNIK(nik string) (*Nasabah, error)
	GetByNoHP(noHP string) (*Nasabah, error)
}

type TransaksiRepository interface {
	Create(transaksi *Transaksi) error
	GetSaldoByNasabahID(nasabahID int64) (float64, error)
}

// Usecase interfaces
type NasabahUsecase interface {
	Register(nama, nik, noHP string) (*Nasabah, error)
	GetSaldo(noRekening string) (float64, error)
}

type TransaksiUsecase interface {
	Tabung(noRekening string, nominal float64) (float64, error)
	Tarik(noRekening string, nominal float64) (float64, error)
}
