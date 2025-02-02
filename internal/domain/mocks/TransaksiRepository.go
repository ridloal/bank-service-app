// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	domain "bank-service-app/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// TransaksiRepository is an autogenerated mock type for the TransaksiRepository type
type TransaksiRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: transaksi
func (_m *TransaksiRepository) Create(transaksi *domain.Transaksi) error {
	ret := _m.Called(transaksi)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Transaksi) error); ok {
		r0 = rf(transaksi)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSaldoByNasabahID provides a mock function with given fields: nasabahID
func (_m *TransaksiRepository) GetSaldoByNasabahID(nasabahID int64) (float64, error) {
	ret := _m.Called(nasabahID)

	if len(ret) == 0 {
		panic("no return value specified for GetSaldoByNasabahID")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (float64, error)); ok {
		return rf(nasabahID)
	}
	if rf, ok := ret.Get(0).(func(int64) float64); ok {
		r0 = rf(nasabahID)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(nasabahID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransaksiRepository creates a new instance of TransaksiRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransaksiRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransaksiRepository {
	mock := &TransaksiRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
