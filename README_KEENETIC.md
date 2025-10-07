# K2Ray Keenetic Extra DSL KN2112 Kurulum Rehberi

Bu rehber, K2Ray'i Keenetic Extra DSL KN2112 modem/router'ında Entware sistemi üzerinde kurmanız için hazırlanmıştır.

## 📋 Ön Gereksinimler

### Donanım Gereksinimleri
- **Keenetic Extra DSL KN2112** modem/router
- **USB bellek** (en az 4GB, Entware kurulumu için)
- **SSH erişimi** etkinleştirilmiş olmalı

### Yazılım Gereksinimleri
- **Entware** paket sistemi kurulmuş olmalı
- **SSH istemci** (Windows için PuTTY, macOS/Linux için terminal)
- **Root erişimi** Keenetic cihazına

## 🔧 1. Adım: Entware Kurulumu

Eğer Entware henüz kurulmamışsa:

1. Keenetic web arayüzüne gidin (http://192.168.1.1)
2. **Genel sistem ayarları** → **Güncellemeler ve paketler** → **Sistem dosyaları**
3. **Entware** paketini etkinleştirin
4. USB belleği cihaza takın ve sistemi yeniden başlatın

## 🚀 2. Adım: K2Ray Kurulumu

### Otomatik Kurulum (Önerilen)

1. **SSH ile Keenetic'e bağlanın:**
   ```bash
   ssh admin@192.168.1.1
   # Şifrenizi girin ve root olun
   su
   ```

2. **K2Ray paketini indirin:**
   ```bash
   cd /tmp
   wget https://github.com/halibiram/k2ray/releases/latest/download/k2ray-keenetic-latest.tar.gz
   tar -xzf k2ray-keenetic-latest.tar.gz
   cd k2ray-keenetic-*
   ```

3. **Kurulum scriptini çalıştırın:**
   ```bash
   chmod +x deployments/entware/scripts/install_entware.sh
   ./deployments/entware/scripts/install_entware.sh
   ```

### Manuel Kurulum

1. **Gerekli bağımlılıkları yükleyin:**
   ```bash
   opkg update
   opkg install curl wget ca-certificates
   ```

2. **Dizinleri oluşturun:**
   ```bash
   mkdir -p /opt/bin /opt/etc/k2ray /opt/var/lib/k2ray /opt/var/log
   ```

3. **K2Ray binary'sini kopyalayın:**
   ```bash
   cp k2ray /opt/bin/
   chmod +x /opt/bin/k2ray
   ```

4. **Init script'i kurun:**
   ```bash
   cp deployments/entware/init.d/S99k2ray /opt/etc/init.d/
   chmod +x /opt/etc/init.d/S99k2ray
   ```

## ⚙️ 3. Adım: Yapılandırma

1. **Yapılandırma dosyasını düzenleyin:**
   ```bash
   vi /opt/etc/k2ray/config.yaml
   ```

2. **Önemli ayarlar:**
   ```yaml
   modem:
     host: "192.168.1.1"
     username: "admin"
     password: "KEENETIC_ADMIN_SIFRENIZ"  # Değiştirin!
   
   server:
     host: "0.0.0.0"
     port: 8080
   
   security:
     jwt_secret: "GUVENLIBIR_SECRET_KEY_OLUSTURUN"  # Değiştirin!
   ```

## 🎯 4. Adım: Servis Yönetimi

### Servis Komutları
```bash
# Servisi başlat
/opt/etc/init.d/S99k2ray start

# Servis durumunu kontrol et
/opt/etc/init.d/S99k2ray status

# Servisi durdur
/opt/etc/init.d/S99k2ray stop

# Servisi yeniden başlat
/opt/etc/init.d/S99k2ray restart
```

### Log Kontrolü
```bash
# Gerçek zamanlı log takibi
tail -f /opt/var/log/k2ray.log

# Son log kayıtları
tail -100 /opt/var/log/k2ray.log
```

## 🌐 5. Adım: Web Arayüzüne Erişim

1. **Web tarayıcınızı açın**
2. **Şu adrese gidin:** http://192.168.1.1:8080
3. **İlk kurulum sihirbazını takip edin**

## 📊 Özellikler ve Kullanım

### DSL Hat Monitoring
- Gerçek zamanlı SNR (Signal-to-Noise Ratio) takibi
- Hat zayıflatma (Attenuation) ölçümü  
- Bağlantı hızı monitörü
- Hata sayacları

### V2Ray Yönetimi
- Kolay yapılandırma yönetimi
- QR kod desteği
- Çoklu profil desteği
- Gerçek zamanlı istatistikler

### Sistem Monitörü
- CPU kullanımı
- Bellek kullanımı
- Ağ trafiği
- Disk kullanımı

## 🔧 Sorun Giderme

### Genel Sorunlar

1. **Servis başlamıyor:**
   ```bash
   # Log dosyasını kontrol edin
   cat /opt/var/log/k2ray.log
   
   # Binary'nin çalıştırılabilir olup olmadığını kontrol edin
   ls -la /opt/bin/k2ray
   
   # Yapılandırma dosyasının doğruluğunu kontrol edin
   /opt/bin/k2ray --config /opt/etc/k2ray --check-config
   ```

2. **Web arayüzüne erişim yok:**
   ```bash
   # Port'un açık olup olmadığını kontrol edin
   netstat -ln | grep :8080
   
   # Firewall kurallarını kontrol edin (varsa)
   iptables -L | grep 8080
   ```

3. **DSL bilgilerine erişim yok:**
   ```bash
   # Modem şifresinin doğruluğunu kontrol edin
   curl -u admin:SIFRENIZ http://192.168.1.1/
   ```

### Log Seviyeleri
```yaml
logging:
  level: "debug"  # Detaylı hata ayıklama için
```

## 🔄 Güncelleme

1. **Yeni sürümü indirin:**
   ```bash
   cd /tmp
   wget https://github.com/halibiram/k2ray/releases/latest/download/k2ray-keenetic-latest.tar.gz
   tar -xzf k2ray-keenetic-latest.tar.gz
   ```

2. **Servisi durdurun:**
   ```bash
   /opt/etc/init.d/S99k2ray stop
   ```

3. **Binary'i güncelleyin:**
   ```bash
   cp k2ray-keenetic-*/k2ray /opt/bin/
   chmod +x /opt/bin/k2ray
   ```

4. **Servisi başlatın:**
   ```bash
   /opt/etc/init.d/S99k2ray start
   ```

## 🗑️ Kaldırma

```bash
# Servisi durdur
/opt/etc/init.d/S99k2ray stop

# Dosyaları sil
rm -f /opt/bin/k2ray
rm -f /opt/etc/init.d/S99k2ray
rm -rf /opt/etc/k2ray
rm -rf /opt/var/lib/k2ray
rm -f /opt/var/log/k2ray.log
```

## 🏗️ Geliştirici Notları

### Cross-Compilation
Bu proje, MIPS little-endian mimarisi için derlenmiştir (Keenetic KN2112):
```bash
make build-keenetic
```

### Build Ortamı
- **Go 1.20+** gereklidir
- **CGO_ENABLED=0** static linking için
- **GOOS=linux GOARCH=mipsle** hedef platform

## 📞 Destek

- **GitHub Issues:** https://github.com/halibiram/k2ray/issues
- **Dokümantasyon:** https://github.com/halibiram/k2ray/wiki
- **Keenetic Forum:** https://forum.keenetic.com

## ⚖️ Lisans

Bu proje MIT lisansı altında dağıtılmaktadır. Detaylar için `LICENSE` dosyasına bakın.

---

**Not:** Bu kurulum rehberi Keenetic Extra DSL KN2112 için optimize edilmiştir. Diğer Keenetic modelleri için ufak değişiklikler gerekebilir.