# Dockerfile for DSL Bypass Ultra System
#
# Bu Dockerfile, uygulamanın tüm bağımlılıklarıyla birlikte
# bir konteyner içinde çalıştırılmasını sağlar. Bu, farklı
# lab ortamlarında tutarlı bir test ortamı sunar.
#
# GÖREV 5 (Deployment) için temel oluşturur.

# Adım 1: Temel imaj olarak resmi Python 3.9 imajını kullan
FROM python:3.9-slim

# Adım 2: Çalışma dizinini ayarla
WORKDIR /app

# Adım 3: Bağımlılıkları kopyala ve kur
# Sadece requirements.txt'yi kopyalayarak Docker katman önbelleğinden
# faydalanıyoruz. Bağımlılıklar değişmediği sürece bu katman yeniden
# oluşturulmaz.
COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Adım 4: Proje dosyalarının tamamını kopyala
COPY . .

# Adım 5: Web arayüzü için portu dışarıya aç
EXPOSE 8080

# Adım 6: Konteyner başlatıldığında çalıştırılacak varsayılan komut
# GÖREV 6'da oluşturulacak ana CLI arayüzü veya web sunucusu burada başlatılabilir.
# Şimdilik web dashboard'ı çalıştıracak şekilde ayarlandı.
CMD ["python", "web/dashboard.py"]