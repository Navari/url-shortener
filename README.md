# 🔗 URL Shortener

Modern, yüksek performanslı URL kısaltma servisi. Go + PostgreSQL + Redis teknolojileri ile geliştirilmiştir.

## ✨ Özellikler

- ⚡ **Yüksek Performans**: Redis cache ile milisaniye düzeyinde yanıt süreleri
- 🔒 **Güvenli**: Bearer token auth ve güvenli random code üretimi
- 📊 **İstatistikler**: URL istatistikleri ve kullanım bilgileri
- ⏰ **Expiry Desteği**: URL'ler için son kullanma tarihi belirleme
- 🐳 **Docker Ready**: Docker Compose ile tek komutla çalıştırma
- 📚 **Swagger API**: Otomatik API dokümantasyonu
- 🧪 **Test Coverage**: %80+ test kapsamı
- 🏗️ **Clean Architecture**: Katmanlı mimari ve dependency injection

## 🚀 Hızlı Başlangıç

### Docker ile Çalıştırma (Önerilen)

```bash
# Repository'yi klonlayın
git clone <repository-url>
cd shortener

# Tüm servisleri başlatın
make up

# Logları takip edin
make logs
```

Servis `http://localhost:8080` adresinde çalışacaktır.

### Manuel Kurulum

#### Gereksinimler

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

#### Kurulum Adımları

```bash
# Bağımlılıkları yükleyin
make deps

# Ortam değişkenlerini ayarlayın
cp env.example .env
# .env dosyasını düzenleyin (özellikle AUTH_TOKEN'ı değiştirin)

# Uygulamayı çalıştırın
make dev
```

## 📖 API Kullanımı

### URL Kısaltma (Auth Token Gerekli)

```bash
# Basit kullanım (expires_at isteğe bağlı)
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-token-here" \
  -d '{
    "url": "https://www.example.com/very/long/url/path"
  }'

# Son kullanma tarihi ile (isteğe bağlı)
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-token-here" \
  -d '{
    "url": "https://www.example.com/very/long/url/path",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

**Yanıt:**
```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://www.example.com/very/long/url/path",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

### URL Yönlendirme

```bash
curl -L http://localhost:8080/abc123
# Otomatik olarak orijinal URL'e yönlendirir
```

### İstatistik Görüntüleme

```bash
curl http://localhost:8080/api/v1/stats/abc123
```

**Yanıt:**
```json
{
  "short_code": "abc123",
  "original_url": "https://www.example.com/very/long/url/path",
  "click_count": 42,
  "created_at": "2024-01-15T10:30:00Z",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

### Health Check

```bash
curl http://localhost:8080/healthz
```

## 📚 API Dokümantasyonu

Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Swagger dokümantasyonunu güncellemek için:
```bash
make docs
```

## 🛠️ Geliştirme

### Mevcut Komutlar

```bash
# Uygulamayı build et
make build

# Testleri çalıştır
make test

# Test coverage raporu oluştur
make coverage

# Linting işlemi
make lint

# Geliştirme modunda çalıştır
make dev

# Docker servisleri
make up          # Servisleri başlat
make down        # Servisleri durdur
make logs        # Logları görüntüle
make restart     # Servisleri yeniden başlat
make rebuild     # Build ve restart

# Dokümantasyon
make docs        # Swagger dokümantasyonu güncelle

# Yardım
make help        # Tüm komutları listele
```

### Test Çalıştırma

```bash
# Tüm testler
make test

# Coverage raporu
make coverage
# coverage.html dosyası oluşturulacak
```

### Proje Yapısı

```
.
├── cmd/
│   └── server/          # Ana uygulama giriş noktası
├── internal/
│   ├── cache/           # Redis cache implementasyonu
│   ├── config/          # Yapılandırma yönetimi
│   ├── handler/         # HTTP handler'ları
│   ├── logger/          # Loglama utilities
│   ├── model/           # Veri modelleri
│   ├── repository/      # Veritabanı katmanı
│   └── service/         # İş mantığı katmanı
├── docker-compose.yml   # Docker Compose yapılandırması
├── Dockerfile          # Docker build dosyası
├── Makefile           # Build ve operasyon komutları
└── README.md          # Bu dosya
```

## ⚙️ Yapılandırma

### Ortam Değişkenleri

| Değişken | Açıklama | Varsayılan |
|----------|----------|------------|
| `PORT` | Server portu | `8080` |
| `ENV` | Ortam (development/production) | `development` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL kullanıcı | `postgres` |
| `DB_PASSWORD` | PostgreSQL şifre | `postgres` |
| `DB_NAME` | Veritabanı adı | `shortener_db` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `BASE_URL` | Temel URL | `http://localhost:8080` |
| `CACHE_TTL` | Cache süresi (saniye) | `3600` |
| `SHORT_CODE_LENGTH` | Kısa kod uzunluğu | `6` |

## 🚀 Production Dağıtımı

### Docker ile Production

```bash
# Production için build
docker build -t shortener:latest .

# Production ortamında çalıştır
docker-compose -f docker-compose.prod.yml up -d
```

### Performans Optimizasyonları

- Redis connection pooling
- Database connection pooling
- Graceful shutdown
- Health check endpoints
- Structured logging

## 🧪 Test ve Kalite

- Unit testler: Service ve repository katmanları
- Integration testler: API endpoint'leri
- Mock'lar: Testify ile dependency mocking
- Coverage: %80+ hedeflenen kapsam

## 📊 İzleme ve Metrikler

- Structured logging (Zap)
- Health check endpoint
- Swagger API dokümantasyonu
- Redis ve PostgreSQL health check'leri

## 🤝 Katkıda Bulunma

1. Fork edin
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Commit edin (`git commit -m 'Add amazing feature'`)
4. Branch'i push edin (`git push origin feature/amazing-feature`)
5. Pull Request açın

## 📝 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakın.

## 🎯 Roadmap

- [ ] Rate limiting middleware
- [ ] Prometheus metrics
- [ ] Custom domain desteği
- [ ] Bulk URL creation API
- [ ] URL preview endpoint
- [ ] Admin dashboard
- [ ] Database migrations
- [ ] Kubernetes deployment manifests 