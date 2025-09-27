# k2ray

[![Build Durumu](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/halibiram/k2ray)
[![Lisans: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Sürüm](https://img.shields.io/badge/version-v1.0.0-blue)](https://github.com/halibiram/k2ray/releases)

k2ray, V2Ray proxy aracı için modern, web tabanlı bir yönetim panelidir. V2Ray yapılandırmalarınızı yönetmek, sistem durumunu izlemek ve trafik metriklerini görüntülemek için kullanıcı dostu bir arayüz sağlar; tümünü web tarayıcınızdan yapabilirsiniz.

---

## 📸 Ekran Görüntüleri

*(Uygulamanın görsel bir önizlemesini sunmak için buraya ekran görüntüleri veya GIF'ler ekleyin.)*

![K2Ray Ekran Görüntüsü](https://place-hold.it/800x450/663399/ffffff?text=k2ray%20UI%20Ekran%20Görüntüsü)

---

## ✨ Özellikler

*   **Gelişmiş Yapılandırma Yönetimi**: Birden çok V2Ray yapılandırmasını kolayca oluşturun, düzenleyin, silin ve aralarında geçiş yapın.
*   **QR Kod Entegrasyonu**: QR kodları kullanarak yapılandırmaları sorunsuz bir şekilde dışa aktarın ve içe aktarın. Mobil cihazınızın kamerasıyla tarayın veya bir resim dosyası yükleyin.
*   **Gerçek Zamanlı Sistem İzleme**: CPU, bellek ve disk kullanımı dahil olmak üzere sistem durumunu anlık olarak takip edin.
*   **Canlı Günlük Görüntüleyici**: V2Ray günlüklerini doğrudan web arayüzünden görüntüleyin ve arayın.
*   **Trafik Metrikleri ve Analizi**: Etkileşimli grafiklerle gerçek zamanlı ağ trafiğini, bağlantı istatistiklerini ve veri kullanımını izleyin.
*   **Çok Dilli Destek**: Hem İngilizce hem de Türkçe dillerini destekleyen tam yerelleştirilmiş arayüz.
*   **Özelleştirilebilir Temalar**: Rahat bir kullanıcı deneyimi için aydınlık ve karanlık temalar arasında geçiş yapın.
*   **Güvenli Erişim**: Panelinizi JWT tabanlı kimlik doğrulama ve isteğe bağlı İki Faktörlü Kimlik Doğrulama (2FA) ile koruyun.
*   **RESTful API**: Programatik erişim ve entegrasyon için iyi belgelenmiş bir API.

---

## 🛠️ Başlarken

Bu bölüm, projenin nasıl çalıştırılacağına dair kısa bir genel bakış sunmaktadır. Daha ayrıntılı talimatlar için lütfen `docs/` dizinindeki belgelere bakın.

### ✅ Sistem Gereksinimleri

*   **Go**: Sürüm 1.24 veya üstü
*   **Node.js**: Sürüm 18 veya üstü (npm ile birlikte)
*   **V2Ray**: k2ray'in yönetebileceği çalışan bir V2Ray örneği.

### ⚙️ Kurulum

1.  **Depoyu klonlayın:**
    ```bash
    git clone https://github.com/halibiram/k2ray.git
    cd k2ray
    ```

2.  **Arka Ucu (API Sunucusu) Kurun:**
    ```bash
    # Geliştirme sunucusunu çalıştırın
    go run ./cmd/k2ray
    ```
    API sunucusu genellikle `8080` portunda çalışmaya başlayacaktır.

3.  **Ön Ucu (Web Arayüzü) Kurun:**
    ```bash
    # Web dizinine gidin
    cd web

    # Bağımlılıkları yükleyin
    npm install

    # Geliştirme sunucusunu çalıştırın
    npm run dev
    ```
    Ön uç geliştirme sunucusuna `http://localhost:5173` adresinden erişilebilir.

---

## 🤝 Katkıda Bulunma Yönergeleri

Katkılar, açık kaynak topluluğunu öğrenmek, ilham vermek ve yaratmak için harika bir yer yapan şeydir. Yaptığınız her türlü katkı **büyük takdirle karşılanır**.

Bu projeyi daha iyi hale getirecek bir öneriniz varsa, lütfen depoyu çatallayın (fork) ve bir pull request oluşturun. Ayrıca "enhancement" etiketiyle bir issue de açabilirsiniz.

1.  Projeyi Çatallayın (Fork the Project)
2.  Özellik Dalınızı Oluşturun (`git checkout -b feature/AmazingFeature`)
3.  Değişikliklerinizi Commit'leyin (`git commit -m 'Add some AmazingFeature'`)
4.  Dala Push'layın (`git push origin feature/AmazingFeature`)
5.  Bir Pull Request Açın

---

## 📄 Lisans

Bu proje MIT Lisansı altında dağıtılmaktadır. Daha fazla bilgi için `LICENSE` dosyasına bakın.

---

## 📜 Değişiklik Günlüğü

Değişikliklerin ayrıntılı bir listesi için lütfen [CHANGELOG.md](./CHANGELOG.md) dosyasına bakın.