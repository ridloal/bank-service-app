package repository

import (
	"bank-service-app/internal/domain"
	"bank-service-app/pkg/logger"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type postgresNasabahRepository struct {
	db  *sql.DB
	log *logger.Logger
}

type postgresTransaksiRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewPostgresNasabahRepository(db *sql.DB, log *logger.Logger) domain.NasabahRepository {
	return &postgresNasabahRepository{
		db:  db,
		log: log,
	}
}

func NewPostgresTransaksiRepository(db *sql.DB, log *logger.Logger) domain.TransaksiRepository {
	return &postgresTransaksiRepository{
		db:  db,
		log: log,
	}
}

func (r *postgresNasabahRepository) Create(n *domain.Nasabah) error {
	r.log.InfoWithContext("Creating new nasabah",
		zap.String("nik", n.NIK),
		zap.String("no_hp", n.NoHP),
	)

	query := `
        INSERT INTO nasabah (nama, nik, no_hp, no_rekening, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		n.Nama,
		n.NIK,
		n.NoHP,
		n.NoRekening,
		n.CreatedAt,
		n.UpdatedAt,
	).Scan(&n.ID)

	if err != nil {
		r.log.ErrorWithContext("Failed to create nasabah", err,
			zap.String("nik", n.NIK),
			zap.String("no_hp", n.NoHP),
		)
		return err
	}

	return nil
}

func (r *postgresNasabahRepository) GetByNoRekening(noRekening string) (*domain.Nasabah, error) {
	r.log.InfoWithContext("Get nasabah by no rekening", zap.String("no_rekening", noRekening))

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
		r.log.ErrorWithContext("Failed to get nasabah by no rekening", err, zap.String("no_rekening", noRekening))
		return nil, err
	}
	return n, nil
}

func (r *postgresNasabahRepository) GetByNoHP(noHP string) (*domain.Nasabah, error) {
	r.log.InfoWithContext("Get nasabah by no hp", zap.String("no_hp", noHP))

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
		r.log.ErrorWithContext("Failed to get nasabah by no hp", err, zap.String("no_hp", noHP))
		return nil, err
	}
	return n, nil
}

func (r *postgresNasabahRepository) GetByNIK(nik string) (*domain.Nasabah, error) {
	r.log.InfoWithContext("Get nasabah by nik", zap.String("nik", nik))

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
		r.log.ErrorWithContext("Failed to get nasabah by nik", err, zap.String("nik", nik))
		return nil, err
	}
	return n, nil
}

func (r *postgresTransaksiRepository) Create(t *domain.Transaksi) error {
	r.log.InfoWithContext("Creating new nasabah",
		zap.Int64("nasabah_id", t.NasabahID),
		zap.String("jenis_transaksi", t.JenisTransaksi),
		zap.Float64("nominal", t.Nominal),
		zap.Float64("saldo_akhir", t.SaldoAkhir),
	)

	query := `
        INSERT INTO transaksi (nasabah_id, jenis_transaksi, nominal, saldo_akhir, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		t.NasabahID,
		t.JenisTransaksi,
		t.Nominal,
		t.SaldoAkhir,
		time.Now(),
	).Scan(&t.ID)

	if err != nil {
		r.log.ErrorWithContext("Failed to create transaksi", err,
			zap.Int64("nasabah_id", t.NasabahID),
			zap.String("jenis_transaksi", t.JenisTransaksi),
			zap.Float64("nominal", t.Nominal),
			zap.Float64("saldo_akhir", t.SaldoAkhir),
		)
		return err
	}

	return nil
}

func (r *postgresTransaksiRepository) GetSaldoByNasabahID(nasabahID int64) (float64, error) {
	r.log.InfoWithContext("Get saldo by nasabah ID", zap.Int64("nasabah_id", nasabahID))

	var saldo float64
	query := `SELECT COALESCE(saldo_akhir, 0) FROM transaksi WHERE nasabah_id = $1 ORDER BY created_at DESC LIMIT 1`

	err := r.db.QueryRow(query, nasabahID).Scan(&saldo)
	if err == sql.ErrNoRows {
		r.log.Info("No saldo found for nasabah with ID", zap.Int64("nasabah_id", nasabahID))
		return 0, nil
	}
	if err != nil {
		r.log.ErrorWithContext("Failed to get saldo by nasabah ID", err, zap.Int64("nasabah_id", nasabahID))
		return 0, err
	}
	return saldo, nil
}
