# Bank Service API Documentation

## Base URL
```
http://localhost:8080
```

## Endpoints

### 1. Register Nasabah
Register a new customer account.

**Endpoint:** `POST /daftar`

**Request Body:**
```json
{
    "nama": "string",
    "nik": "string (16 characters)",
    "no_hp": "string"
}
```

**Success Response (200):**
```json
{
    "no_rekening": "string (10 digits)"
}
```

**Error Response (400):**
```json
{
    "remark": "string (error description)"
}
```

### 2. Menabung (Deposit)
Deposit money into an account.

**Endpoint:** `POST /tabung`

**Request Body:**
```json
{
    "no_rekening": "string",
    "nominal": number
}
```

**Success Response (200):**
```json
{
    "saldo": number
}
```

**Error Response (400):**
```json
{
    "remark": "string (error description)"
}
```

### 3. Tarik Dana (Withdrawal)
Withdraw money from an account.

**Endpoint:** `POST /tarik`

**Request Body:**
```json
{
    "no_rekening": "string",
    "nominal": number
}
```

**Success Response (200):**
```json
{
    "saldo": number
}
```

**Error Response (400):**
```json
{
    "remark": "string (error description)"
}
```

### 4. Cek Saldo (Balance Inquiry)
Check account balance.

**Endpoint:** `GET /saldo/{no_rekening}`

**Parameters:**
- `no_rekening` (path parameter): Account number

**Success Response (200):**
```json
{
    "saldo": number
}
```

**Error Response (400):**
```json
{
    "remark": "string (error description)"
}
```

## Error Codes and Messages

1. Account Registration:
   - "NIK sudah terdaftar"
   - "Nomor HP sudah terdaftar"

2. Transactions:
   - "Nomor rekening tidak ditemukan"
   - "Saldo tidak mencukupi"

## Notes
- All monetary values are in IDR (Indonesian Rupiah)
- Account numbers are 10 digits long
- NIK must be 16 characters