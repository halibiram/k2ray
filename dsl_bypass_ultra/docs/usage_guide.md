# DSL Bypass Ultra - Kullanım Kılavuzu (Final Sürüm)

Bu kılavuz, sistemi kurmayı ve yeni komut satırı arayüzünü (CLI) kullanarak yönetmeyi açıklar.

## 1. Kurulum

Başlamadan önce, projenin ana dizinindeki `scripts/install.sh` betiğini çalıştırarak sistemi kurun. Bu betik, gerekli tüm bağımlılıkları yükleyecek ve sistemi çalıştırmak için gereken `config.yaml` ve `start.sh` dosyalarını oluşturacaktır.

```bash
cd /path/to/dsl_bypass_ultra
./scripts/install.sh
```

Kurulum tamamlandıktan sonra, projenin ana dizininde oluşturulan `config.yaml` dosyasını bir metin düzenleyici ile açın ve modeminizin yönetici şifresini `password` alanına girin. **Bu adım zorunludur.**

```yaml
modem:
  host: "192.168.1.1"
  username: "admin"
  password: "YOUR_PASSWORD_HERE" # <--- ŞİFRENİZİ BURAYA GİRİN
```

## 2. Kullanım

Sistem, projenin ana dizininde oluşturulan `start.sh` betiği aracılığıyla yönetilir.

### a. Sistemi Çalıştırma (`run` komutu)

Bu, sistemin ana çalışma komutudur. Bypass işlemini yürütür, performans monitörünü arka planda başlatır ve gerçek zamanlı izleme için web panelini (`http://0.0.0.0:8080`) sunar.

```bash
./start.sh run
```

**Seçenekler:**
*   `--profile <profil_adı>`: `config.yaml` içinde tanımlı farklı bir optimizasyon profili kullanmak için.
    ```bash
    ./start.sh run --profile stable
    ```
*   `--no-bypass`: Bypass işlemini atlayıp sadece izleme modunda başlatmak için.
    ```bash
    ./start.sh run --no-bypass
    ```

### b. Anlık Durumu Kontrol Etme (`status` komutu)

Sistemi tamamen başlatmadan, sadece o anki modem durumunu hızlıca kontrol etmek için `status` alt komutunu kullanın. Betik, durumu yazdırır ve çıkar.

```bash
./start.sh status
```

**Örnek Çıktı:**
```
📊 Fetching current modem status...

--- CURRENT DSL STATUS ---
status              : Up
snr_margin_down     : 25.0
data_rate_down      : 32.54 Mbps
...
--------------------------
```