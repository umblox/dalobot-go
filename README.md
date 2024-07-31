# dalobot-go
# OpenWRT Telegram Bot

Telegram bot untuk OpenWRT yang memungkinkan admin mengisi saldo pengguna dan pengguna untuk memeriksa saldo mereka.

## Persyaratan

- OpenWRT dengan Go terinstal
- Akses root ke perangkat OpenWRT
- Token bot Telegram dan ID pengguna admin

## Instalasi

1. **Clone repositori ini** ke perangkat OpenWRT:
   ```sh
   git clone https://github.com/umblox/dalobot-go.git
2. Navigasikan ke direktori proyek :
   ```sh
   cd dalobot-go
3. Kompilasi file GO :
   ```sh
   go build telebot.go
4. Buat direktori untuk file yang diperlukan:
   ```sh
   mkdir -p /root/Telebot-Radius/files
5. Salin file yang diperlukan:
   ```sh
   cp telebot /usr/bin/
   cp telebot.sh /etc/init.d/telebot
   cp auth.txt /root/Telebot-Radius/files/
   cp profile.json /root/Telebot-Radius/files/
6. Ubah izin file :
   ```sh
   chmod +x /usr/bin/telebot
   chmod +x /etc/init.d/telebot
7. Aktifkan dan mulai layanan :
   ```sh  
   /etc/init.d/telebot enable
   /etc/init.d/telebot start
#################################################################

Penggunaan
Perintah Admin
/isi <user_id> <jumlah> - Topup saldo pengguna.
/saldo <user_id> - Cek saldo pengguna tertentu.
Perintah Pengguna
/saldo - Cek saldo Anda.
/menu - Tampilkan menu utama.


