package repository

import (
	"bank-service-app/pkg/logger"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresNasabahRepository(t *testing.T) {
	db := &sql.DB{}
	log := logger.NewLogger()

	repo := NewPostgresNasabahRepository(db, log)

	assert.NotNil(t, repo)
	assert.IsType(t, &postgresNasabahRepository{}, repo)
}

func TestNewPostgresTransaksiRepository(t *testing.T) {
	db := &sql.DB{}
	log := logger.NewLogger()

	repo := NewPostgresTransaksiRepository(db, log)

	assert.NotNil(t, repo)
	assert.IsType(t, &postgresTransaksiRepository{}, repo)
}
