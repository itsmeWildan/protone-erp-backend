# Buku Panduan Integrasi API ProtoERP

Dokumen ini berisi panduan teknis dan alur integrasi API ProtoERP. Dokumen ini dirancang untuk memudahkan tim Frontend dalam menghubungkan antarmuka pengguna (UI) dengan sistem Backend.

> **CATATAN PENTING:**
> Hampir seluruh API (kecuali proses Registrasi dan Login) diwajibkan untuk menyertakan header otentikasi berupa:
> `Authorization: Bearer [TOKEN_ANDA]`
> Token ini didapatkan dari respons API Login. Khusus untuk endpoint download PDF, token disisipkan melalui parameter URL (`?token=...`).

## Tahap 1: Inisialisasi & Login (Wajib Pertama Kali)
Sebelum mengakses modul lain, sistem memerlukan kredensial penyewa (Tenant) dan token akses untuk mengenali siapa yang sedang melakukan permintaan.

| Langkah | Endpoint & Fungsi |
| :--- | :--- |
| **1. Daftarkan Perusahaan**<br>*(Registrasi Tenant)* | **POST** `/api/v1/auth/register-tenant`<br>**Fungsi:** Membuat entitas perusahaan baru di dalam sistem. Response akan mengembalikan detail perusahaan beserta akun admin pertama. |
| **2. Login & Dapatkan Token**<br>*(Otentikasi)* | **POST** `/api/v1/auth/login`<br>**Fungsi:** Digunakan untuk masuk menggunakan email dan password. Response API ini akan mengembalikan **AccessToken**. Harap simpan token ini (misalnya di LocalStorage) untuk digunakan pada seluruh panggilan API berikutnya. |

## Tahap 2: Setup Data Master (Karyawan)
Modul absen, cuti, dan penggajian sangat bergantung pada data karyawan. Pastikan karyawan telah terdaftar sebelum menggunakan fitur operasional harian.

| Langkah | Endpoint & Fungsi |
| :--- | :--- |
| **1. Buat Data Karyawan**<br>*(Create Employee)* | **POST** `/api/v1/employees`<br>**Fungsi:** Mendaftarkan staf atau karyawan baru. Response akan mengembalikan **Employee ID**. Catat ID ini jika Anda membutuhkan referensi spesifik untuk pengujian selanjutnya. |
| **2. Lihat Daftar Karyawan**<br>*(List Employees)* | **GET** `/api/v1/employees`<br>**Fungsi:** Menampilkan seluruh karyawan yang terdaftar di perusahaan Anda. API ini berguna jika Anda perlu mencari ID karyawan tertentu untuk proses lain (misal: memberi jatah cuti). |

## Tahap 3: Alur Kehadiran & Cuti (Operasional Harian)
Bagian ini mencakup rutinitas pencatatan jam kerja harian dan permohonan ketidakhadiran (cuti).

| Langkah | Endpoint & Fungsi |
| :--- | :--- |
| **1. Absen Masuk & Pulang**<br>*(Clock-In / Clock-Out)* | **POST** `/api/v1/attendance/clock-in`<br>**POST** `/api/v1/attendance/clock-out`<br>**Fungsi:** Mencatat waktu kehadiran. Sistem otomatis mengenali identitas karyawan dari token login yang disematkan pada header. |
| **2. Berikan Jatah Cuti**<br>*(Init Balance)* | **POST** `/api/v1/leaves/init-balance`<br>**Fungsi & Syarat:** API ini membutuhkan input **Employee ID** (didapat dari `GET /employees`). API ini bertugas mengalokasikan saldo cuti tahunan (misal: 12 hari) kepada karyawan yang bersangkutan agar mereka bisa mengajukan cuti. |
| **3. Ajukan & Setujui Cuti**<br>*(Request & Approve)* | **POST** `/api/v1/leaves/request`<br>**Fungsi:** Mengirim permohonan cuti baru. Response akan mengembalikan **Leave ID**.<br><br>**PATCH** `/api/v1/leaves/approve`<br>**Fungsi:** Manajer (dengan role yang berwenang) menggunakan **Leave ID** tersebut untuk menyetujui (Approve) permohonan cuti. |

## Tahap 4: Alur Penggajian & Cetak Slip PDF (Siklus Bulanan)
Fitur penggajian dilakukan setiap bulan. PDF slip gaji hanya bisa diunduh apabila periode gaji bulan tersebut telah diproses oleh Admin/HR.

| Langkah | Endpoint & Fungsi |
| :--- | :--- |
| **1. Proses Gaji Bulanan**<br>*(Generate Payroll)* | **POST** `/api/v1/payroll/generate`<br>**Fungsi:** Admin memicu kalkulasi gaji seluruh karyawan dengan mengirimkan parameter bulan (month) dan tahun (year). Response akan mengembalikan **Period ID** (ID Periode Penggajian). |
| **2. Lihat Rincian Gaji**<br>*(My Payslip)* | **GET** `/api/v1/payroll/my-slip?month=5&year=2026`<br>**Fungsi:** Mengambil rincian JSON gaji karyawan yang login. Di dalam JSON responsenya, Anda akan menemukan field **PayrollPeriodID**. ID ini adalah kunci untuk mengunduh versi PDF. |
| **3. Download PDF Slip Gaji**<br>*(Cetak Dokumen)* | **GET** `/api/v1/payroll/slips/{PayrollPeriodID}/download?token=TOKEN_ANDA`<br>**Fungsi:** Endpoint stream untuk mengunduh berkas PDF fisik. Ganti `{PayrollPeriodID}` dengan ID dari langkah 2. Karena proses ini biasanya menggunakan elemen `<a href="...">` di browser yang tidak mengirimkan header otentikasi, Anda **wajib** menyematkan parameter `?token=...` pada URL tersebut. |

## Tahap Tambahan: Lembur, Klaim, dan Anggaran
Ketiga modul ini memiliki karakteristik alur yang hampir sama dengan modul Cuti, yakni berbasis siklus Permohonan (Request) lalu Persetujuan (Approve) menggunakan ID referensi.

| Modul | Alur Endpoint |
| :--- | :--- |
| **Reimbursement (Klaim)** | 1. **POST** `/api/v1/reimbursements/submit` (Dapatkan Claim ID)<br>2. **PATCH** `/api/v1/reimbursements/approve` (Gunakan Claim ID) |
| **Overtime (Lembur)** | 1. **POST** `/api/v1/overtime/request` (Dapatkan Overtime ID)<br>2. **PATCH** `/api/v1/overtime/approve` (Gunakan Overtime ID) |
| **Budget (Anggaran)** | 1. **POST** `/api/v1/budgets/set` (Atur pagu anggaran)<br>2. **GET** `/api/v1/budgets/check` (Cek sisa anggaran) |
