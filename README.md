> **Catatan:** Install Docker dan Docker Compose terlebih dahulu sebelum mengikuti panduan ini.

# React Native and Go Task Management Assessment

Panduan menjalankan backend Go dan aplikasi Expo React Native secara lokal, termasuk koneksi dari HP fisik.

## Prasyarat

- Go versi yang sesuai dengan `otomedia-test-backend/go.mod`
- Docker dengan Docker Compose
- Node.js dan npm
- Expo Go di HP untuk menjalankan aplikasi saat development

## 1. Jalankan Backend

Buka terminal pertama:

```bash
cd otomedia-test-backend
docker compose up -d mysql redis
docker compose ps
go run .
```

Tunggu sampai MySQL dan Redis berstatus `healthy` pada `docker compose ps` sebelum menjalankan `go run .`. Backend berjalan pada port `8080` secara default. Pengaturan default database dan Redis sudah sesuai dengan `docker-compose.yml` lokal, sehingga tidak perlu mengisi `.env` backend untuk development.

Migrasi dan data contoh pada folder `migrations/` dijalankan otomatis oleh MySQL saat volume database pertama kali dibuat. Volume yang sudah ada tidak menjalankan ulang file inisialisasi tersebut.

Pastikan backend merespons:

```bash
curl http://localhost:8080/healthz
```

Respons yang diharapkan: `{"status":"ok"}`.

## 2. Atur URL API Frontend

Buka terminal kedua:

```bash
cd otomedia-test-frontend
npm install
cp -n .env.example .env
```

Buka `otomedia-test-frontend/.env` dan atur `EXPO_PUBLIC_API_BASE_URL` sesuai tempat aplikasi dijalankan:

```env
EXPO_PUBLIC_API_BASE_URL=http://localhost:8080
```

Untuk Expo web atau simulator yang berjalan di komputer yang sama, `localhost` biasanya sesuai. Untuk HP fisik, ganti `localhost` dengan alamat IP LAN komputer yang menjalankan backend, contohnya:

```env
EXPO_PUBLIC_API_BASE_URL=http://192.168.1.20:8080
```

Cari alamat IPv4 LAN komputer (Linux):

```bash
hostname -I
```

Di Windows, jalankan perintah berikut pada Command Prompt atau PowerShell, lalu cari `IPv4 Address` pada adaptor Wi-Fi atau Ethernet yang sedang digunakan:

```powershell
ipconfig
```

Di macOS, untuk melihat alamat IPv4 Wi-Fi:

```bash
ipconfig getifaddr en0
```

Jika perintah macOS tidak menampilkan alamat, antarmuka Wi-Fi mungkin menggunakan nama lain. Periksa nama antarmuka dengan `networksetup -listallhardwareports`, lalu ganti `en0` sesuai nama perangkat Wi-Fi.

Pilih alamat jaringan Wi-Fi/LAN yang sama dengan HP, bukan `127.0.0.1` atau alamat loopback lainnya. Alamat IP bisa berubah ketika komputer tersambung ulang ke jaringan.

## 3. Jalankan Frontend di HP

Pastikan HP dan komputer terhubung ke jaringan Wi-Fi yang sama. Jalankan Metro/Expo dari folder frontend:

```bash
npm start -- --lan
```

Buka Expo Go di HP dan pindai QR code yang ditampilkan terminal. Jika nilai `.env` diubah setelah Expo berjalan, hentikan lalu jalankan kembali server Expo agar nilai environment baru dimuat.

Sebelum membuka aplikasi, URL berikut juga bisa diuji dari browser HP (ganti IP sesuai komputer):

```text
http://192.168.1.20:8080/healthz
```

Jika tidak dapat dibuka, pastikan backend masih berjalan, kedua perangkat berada di jaringan yang sama, dan firewall komputer mengizinkan koneksi masuk ke TCP port `8080`. Jaringan kantor atau guest Wi-Fi tertentu dapat mengisolasi perangkat satu sama lain.

> `localhost` pada HP menunjuk ke HP itu sendiri, bukan komputer. Karena itu URL API pada `.env` harus memakai IP LAN komputer saat memakai perangkat fisik.
