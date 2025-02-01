package repository

import (
	"bank-service-app/internal/domain"
	"database/sql"
	"time"
)

type postgresNasabahRepository struct {
	db *sql.DB
}

type postgresTransaksiRepository struct {
	db *sql.DB
}

func NewPostgresNasabahRepository(db *sql.DB) domain.NasabahRepository {
	return &postgresNasabahRepository{db: db}
}

func NewPostgresTransaksiRepository(db *sql.DB) domain.TransaksiRepository {
	return &postgresTransaksiRepository{db: db}
}

func (r *postgresNasabahRepository) Create(n *domain.Nasabah) error {
	query := `
        INSERT INTO nasabah (nama, nik, no_hp, no_rekening, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	return r.db.QueryRow(
		query,
		n.Nama,
		n.NIK,
		n.NoHP,
		n.NoRekening,
		time.Now(),
		time.Now(),
	).Scan(&n.ID)
}

func (r *postgresNasabahRepository) GetByNoRekening(noRekening string) (*domain.Nasabah, error) {
	n := &domain.Nasabah{}
	query := `SELECT id, nama, nik, no_hp, no_rekening, created_at, updated_at FROM nasabah WHERE no_rekening = $1`

	err := r.db.QueryRow(query, noRekening).Scan(
		&n.ID,
		&n.Nama,
		&n.NIK,
		&n.NoHP,
		&n.NoRekening,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *postgresNasabahRepository) GetByNoHP(noHP string) (*domain.Nasabah, error) {
	n := &domain.Nasabah{}
	query := `SELECT id, nama, nik, no_hp, no_rekening, created_at, updated_at FROM nasabah WHERE no_hp = $1`

	err := r.db.QueryRow(query, noHP).Scan(
		&n.ID,
		&n.Nama,
		&n.NIK,
		&n.NoHP,
		&n.NoRekening,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *postgresNasabahRepository) GetByNIK(nik string) (*domain.Nasabah, error) {
	n := &domain.Nasabah{}
	query := `SELECT id, nama, nik, no_hp, no_rekening, created_at, updated_at FROM nasabah WHERE nik = $1`

	err := r.db.QueryRow(query, nik).Scan(
		&n.ID,
		&n.Nama,
		&n.NIK,
		&n.NoHP,
		&n.NoRekening,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *postgresTransaksiRepository) Create(t *domain.Transaksi) error {
	query := `
        INSERT INTO transaksi (nasabah_id, jenis_transaksi, nominal, saldo_akhir, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	return r.db.QueryRow(
		query,
		t.NasabahID,
		t.JenisTransaksi,
		t.Nominal,
		t.SaldoAkhir,
		time.Now(),
	).Scan(&t.ID)
}

func (r *postgresTransaksiRepository) GetSaldoByNasabahID(nasabahID int64) (float64, error) {
	var saldo float64
	query := `SELECT COALESCE(saldo_akhir, 0) FROM transaksi WHERE nasabah_id = $1 ORDER BY created_at DESC LIMIT 1`

	err := r.db.QueryRow(query, nasabahID).Scan(&saldo)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return saldo, nil
}
