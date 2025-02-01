package http

type RegisterResponse struct {
	NoRekening string `json:"no_rekening"`
}

type SaldoResponse struct {
	Saldo float64 `json:"saldo"`
}

type ErrorResponse struct {
	Remark string `json:"remark"`
}
