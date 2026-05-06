# ProtoERP Backend

ProtoERP adalah sistem Enterprise Resource Planning (ERP) berbasis **Go (Golang)** yang dirancang dengan **Clean Architecture** dan prinsip **Domain-Driven Design (DDD)**. Backend ini menangani seluruh logika bisnis inti untuk manajemen sumber daya manusia (HR) dan keuangan.

## 🚀 Fitur Utama

- **Multi-Tenancy**: Mendukung banyak perusahaan (tenant) dalam satu instance database.
- **Manajemen Karyawan**: Pengelolaan data master karyawan, jabatan, dan departemen.
- **Absensi (Attendance)**: Fitur Clock-In/Out dengan pencatatan lokasi dan catatan.
- **Manajemen Cuti (Leave)**: Pengajuan cuti, persetujuan manajer, dan pelacakan kuota cuti tahunan.
- **Lembur (Overtime)**: Pengajuan lembur, perhitungan otomatis uang lembur ke dalam gaji.
- **Reimbursement**: Klaim biaya operasional dengan sistem persetujuan.
- **Penggajian (Payroll)**: 
  - Generate slip gaji bulanan otomatis.
  - Perhitungan gaji pokok, tunjangan, potongan, dan lembur.
  - Export slip gaji ke format **PDF**.
  - Integrasi otomatis ke Jurnal Keuangan (Accounting).
- **Keuangan & Anggaran**: 
  - Manajemen Budget per Departemen.
  - Pencatatan Jurnal Umum (General Ledger) otomatis dari transaksi payroll.

## 🛠️ Stack Teknologi

- **Bahasa**: Go (Golang) 1.21+
- **Database**: PostgreSQL (pgx driver)
- **Framework HTTP**: Chi Router (Ringan & Cepat)
- **Library**: 
  - `google/uuid` (Identity)
  - `maroto` (PDF Generation)
  - `golang-jwt` (Authentication)
  - `migrate` (Database Migrations)

## 🏗️ Struktur Proyek (Clean Architecture)

- `internal/domain`: Entitas bisnis dan logika inti (Tanpa dependensi luar).
- `internal/usecase`: Alur kerja aplikasi (Orchestration logic).
- `internal/infrastructure`: Implementasi database (Postgres), Migrasi, dan External Tools.
- `internal/delivery/http`: Handler API, Router, dan Middleware.
- `pkg/`: Library pembantu yang bisa digunakan kembali (Response wrapper, JWT, PDF).

## 🏃 Cara Menjalankan

### 1. Prasyarat
- Instal [Go](https://go.dev/dl/)
- Instal [PostgreSQL](https://www.postgresql.org/download/)

### 2. Konfigurasi Database
Buat database bernama `protone_erp` di PostgreSQL Anda.

### 3. Setup Environment
Buat file `.env` di folder `backend/` (contoh isi ada di bawah):
```env
DB_URL=postgres://user:password@localhost:5432/protone_erp?sslmode=disable
JWT_SECRET=your_secret_key
PORT=3000
```

### 4. Menjalankan Migrasi & Seeder
```bash
# Jalankan migrasi untuk membuat tabel
# (Atau gunakan tool migrasi jika tersedia)

# Jalankan seeder untuk data awal (Tenant BCDE)
go run seed_payroll.go
```

### 5. Jalankan Aplikasi
```bash
cd backend
go run cmd/api/main.go
```

## 🧪 Pengujian API
Anda dapat menguji API menggunakan Postman atau HTTPie. Koleksi API tersedia di folder `docs/` (opsional).

---
Developed with ❤️ by itsmeWildan
