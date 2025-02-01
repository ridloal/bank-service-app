package http

type RegisterRequest struct {
	Nama string `json:"nama"`
	NIK  string `json:"nik"`
	NoHP string `json:"no_hp"`
}

type TransaksiRequest struct {
	NoRekening string  `json:"no_rekening"`
	Nominal    float64 `json:"nominal"`
}
