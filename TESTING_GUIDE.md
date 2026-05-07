# 🚀 ProtoERP Full Testing Guide

Panduan ini menjelaskan alur pengujian lengkap dari pendaftaran pengguna hingga pencetakan slip gaji PDF.

## 🔐 Tahap 1: Autentikasi (Pintu Masuk)

### 1. Register Tenant & Admin
Jika belum punya akun, daftar dulu:
- **POST** `/api/v1/auth/register`
- **Body**:
  ```json
  {
    "company_name": "PT Maju Mundur",
    "email": "admin@maju.com",
    "password": "password123",
    "full_name": "Administrator"
  }
  ```

### 2. Login
- **POST** `/api/v1/auth/login`
- **Body**:
  ```json
  {
    "email": "admin@maju.com",
    "password": "password123"
  }
  ```
- **OUTPUT**: Anda akan mendapatkan `token`. 
- **PENTING**: Copy token tersebut dan masukkan ke tab **Auth -> Bearer Token** di HTTPie untuk semua request selanjutnya.

---

## 👥 Tahap 2: Manajemen Karyawan
- **GET** `/api/v1/employees`: Untuk melihat daftar karyawan.
- **Dapatkan ID Karyawan**: Cari field `"id"` (Contoh: `1bf46483...`). Anda butuh ID ini untuk beberapa API.

---

## 💰 Tahap 3: Siklus Payroll (Gaji)

### 1. Generate Payroll
Membuat draf gaji untuk seluruh karyawan di bulan tertentu.
- **POST** `/api/v1/payroll/generate`
- **Body**: `{"month": 5, "year": 2026}`

### 2. Lihat Slip Saya
Mengecek apakah draf sudah ada dan mengambil ID Periode.
- **GET** `/api/v1/payroll/my-slip?month=5&year=2026`
- **ID YANG DICARI**: 
  - **`PayrollPeriodID`**: Digunakan untuk Download PDF.
  - **`ID` (paling atas)**: Adalah ID Slip Gaji Anda.

### 3. Approve Payroll (Persetujuan)
Mengunci draf agar statusnya berubah dari `draft` menjadi `approved`.
- **PATCH** `/api/v1/payroll/approve`
- **Body**: `{"month": 5, "year": 2026}`

### 4. Pay Payroll (Pembayaran)
Tahap akhir! Memotong anggaran departemen dan mencatat ke jurnal keuangan.
- **POST** `/api/v1/payroll/pay`
- **Body**: `{"month": 5, "year": 2026}`

---

## 📄 Tahap 4: Download PDF Slip Gaji

Di tahap ini sering terjadi kebingungan ID. Ingat rumus ini:
**URL: `/api/v1/payroll/slips/{PeriodID}/download`**

1. **Mana ID yang benar?**
   - Buka respon dari `GET /api/v1/payroll/my-slip`.
   - Cari field bernama **`PayrollPeriodID`**.
   - **JANGAN** gunakan field `ID` yang paling atas (itu ID Slip, bukan ID Periode).

2. **Kenapa di HTTPie tidak bisa di-Save?**
   - HTTPie versi desktop kadang mendeteksi PDF sebagai teks biasa sehingga tombol "Save body" mati.
   - **Solusi**: Gunakan browser atau PowerShell jika tombol Save di HTTPie tidak aktif.

3. **Ingat Masa Berlaku Token!**
   - Token JWT memiliki masa berlaku (biasanya 1 jam).
   - Jika muncul error `Invalid or expired token` saat download, silakan **Login ulang** untuk mendapatkan Token baru.

---

## 📊 Tahap 5: Verifikasi Dashboard & Budget

### 1. Statistik Dashboard
- **GET** `/api/v1/dashboard/stats`
- **Tujuan**: Memastikan total pengeluaran gaji muncul di chart.

### 2. Sisa Anggaran (Budget)
- **GET** `/api/v1/budgets/check?department_id={ID}&month=5&year=2026`
- **Dapatkan Department ID**: Bisa dilihat di data karyawan atau menu "GET DEPARTEMEN ID".
- **Hasil**: Nilai `spent` harus bertambah dan `remaining` berkurang.
