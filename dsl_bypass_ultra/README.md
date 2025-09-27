# Google Jules Ajanı - Keenetic DSLAM Bypass Ultra System

🔬 **KAPALI LAB ORTAMI UYARISI**

**ÖNEMLİ:** Bu sistem, yalnızca üniversite kapalı laboratuvar ortamında, akademik araştırma amacıyla geliştirilmektedir. Gerçek (production) ağlarda kullanılması kesinlikle hedeflenmemektedir ve yasa dışı olabilir. Tüm testler, bir DSLAM ve bir Keenetic modemden oluşan tamamen izole bir ağda gerçekleştirilecektir.

---

## 🎯 Proje Hedefi

Bu projenin amacı, kapalı bir laboratuvar ortamında Keenetic Extra DSL gibi modemler için, DSLAM (Digital Subscriber Line Access Multiplexer) tarafından uygulanan profil kısıtlamalarını aşan (bypass) ultra gelişmiş bir sistem geliştirmektir. Sistem, hat parametrelerini (SNR, zayıflama vb.) akıllıca manipüle ederek, modem ve DSLAM arasındaki anlaşmayı daha yüksek hız profilleri üzerinden yapmaya zorlamayı hedefler.

**Ana Hedef:** 300 metre mesafedeki bir hattı 5 metre gibi simüle ederek ve diğer parametreleri optimize ederek, standart bir 30 Mbps VDSL profilini **100+ Mbps** hızlarına çıkarmak.

## 🚀 Temel Özellikler (Geliştirme Hedefleri)

- **Tek Komutla Kurulum:** `./scripts/install.sh` ile tüm sistemin 5 dakika içinde kurulabilmesi.
- **Otomatik Keşif:** Ağdaki Keenetic modemin otomatik olarak bulunması.
- **Akıllı Parametre Motoru:** Mevcut hat değerlerine göre en uygun bypass parametrelerini hesaplayan çekirdek motor.
- **DSLAM Spoofing:** DSLAM'a hattın fiziksel durumundan çok daha iyi olduğunu "söyleyerek" daha yüksek hız profillerine erişim sağlama.
- **Gerçek Zamanlı İzleme:** Performans metriklerini (hız, SNR, CRC hataları) anlık olarak takip eden bir web arayüzü.
- **Güvenlik ve Geri Yükleme:** "Modem brick" riskini en aza indirmek için otomatik yapılandırma yedekleme ve tek tuşla geri yükleme mekanizmaları.
- **AI Tabanlı Optimizasyon (Gelecek Hedefi):** Zamanla en iyi ayarları kendi kendine öğrenen bir yapay zeka modülü.

## 🏛️ Proje Mimarisi

Proje, aşağıdaki modüler yapı üzerine inşa edilmiştir:

- `core/`: Sistemin ana mantığını içeren çekirdek bileşenler (modem arayüzü, parametre manipülatörü, bypass orkestratörü).
- `engines/`: VDSL, SNR, hız artırma gibi özel optimizasyon motorları.
- `ai/`: Gelecekteki makine öğrenmesi tabanlı optimizasyon modülleri.
- `utils/`: Ağ tarayıcı, konfigürasyon oluşturucu gibi yardımcı araçlar.
- `web/`: Flask tabanlı gerçek zamanlı izleme paneli ve API sunucusu.
- `scripts/`: Kurulum, dağıtım ve test betikleri.
- `config/`: Optimizasyon kuralları, modem profilleri gibi yapılandırma dosyaları.
- `docs/`: Proje dokümantasyonu.

## 📦 Kurulum ve Kullanım (GÖREV 5 Sonrası)

1.  Projeyi klonlayın: `git clone ...`
2.  Kurulum betiğini çalıştırın: `cd dsl_bypass_ultra/scripts && ./install.sh`
3.  Konfigürasyon dosyasını düzenleyin: `config/default.yaml`
4.  Sistemi başlatın: `python main.py` (veya benzeri bir komut)

---
Bu README dosyası, projenin GÖREV 1 aşamasında oluşturulmuştur ve ilerleyen görevlerde güncellenecektir.