Sınav Bilgi Sistemi

sinav-bilgi-sistemi, Go ile geliştirilen, sınav verilerini merkezi olarak toplayan, yöneten ve uygulamaya sunan bir backend servisidir.

Projenin amacı:
	•	sınav verilerini dış kaynaklardan çekmek
	•	verileri veritabanında saklamak
	•	tekrar eden kayıtları engellemek
	•	admin tarafından yönetilebilir hale getirmek
	•	istemcilere API üzerinden sunmak

Kullanılan teknolojiler
	•	Go
	•	Gin
	•	PostgreSQL
	•	Redis
	•	Docker
	•	JWT Authentication

Şu ana kadar tamamlanan yapılar

Altyapı
	•	Go backend projesi kuruldu
	•	Docker ile PostgreSQL ve Redis ayağa kaldırıldı
	•	migration sistemi kuruldu
	•	proje GitHub’a pushlandı

Kimlik doğrulama ve yetkilendirme
	•	JWT tabanlı auth sistemi geliştirildi
	•	register, login, me endpointleri yazıldı
	•	admin-only route koruması eklendi
	•	middleware tabanlı auth yapısı kuruldu

Exam domain
	•	exams tablosu oluşturuldu
	•	exam CRUD endpointleri yazıldı
	•	sadece adminin create/update/delete yapabildiği yetki sistemi kuruldu
	•	listeleme için pagination ve filtering eklendi

Örnek filtreler:
	•	GET /api/v1/exams?page=1&limit=10
	•	GET /api/v1/exams?source=osym
	•	GET /api/v1/exams?status=published

Provider tabanlı veri çekme yapısı
	•	provider mimarisi kuruldu
	•	mock provider eklendi
	•	osym provider eklendi
	•	provider → service → repository → database akışı kuruldu

Otomatik veri alma
	•	scheduler servisi eklendi
	•	belirli aralıklarla otomatik import çalışıyor
	•	dış kaynaktan gelen veriler upsert mantığı ile kaydediliyor
	•	duplicate kayıt oluşması engellendi

Import log sistemi
	•	her import işlemi loglanıyor
	•	provider adı, durum, import edilen kayıt sayısı ve hata mesajı tutuluyor
	•	import logları API üzerinden listelenebiliyor

Redis cache
	•	GET /api/v1/exams için Redis cache eklendi
	•	create/update/delete/import sonrası cache invalidation yapılıyor
	•	listeleme performansı iyileştirildi

Mimari yapı

Proje layered architecture mantığıyla organize edilmiştir:
	•	domain → veri modelleri
	•	repository → veritabanı erişimi
	•	service → iş mantığı
	•	handler → HTTP katmanı
	•	provider → dış kaynaklardan veri çekme katmanı
	•	middleware → auth ve authorization kontrolleri

Şu an çalışan ana akış
Provider
↓
Import Service
↓
Repository (Upsert)
↓
PostgreSQL
↓
Public API
Scheduler ile bu akış otomatik olarak belirli aralıklarla tetiklenmektedir.

Mevcut durumda çalışan özellikler
	•	JWT auth
	•	role-based authorization
	•	exam CRUD
	•	pagination/filtering
	•	provider-based import
	•	automatic scheduler import
	•	import logging
	•	redis cache
	•	duplicate prevention

Devam eden geliştirme

Şu anda sistem, gerçek ÖSYM kaynağından sınav başlıklarını başarıyla çekmektedir.
Bir sonraki hedef, sınavlara ait:
	•	sınav tarihi
	•	başvuru başlangıç tarihi
	•	başvuru bitiş tarihi
	•	sonuç tarihi

alanlarını güvenilir biçimde parse ederek sisteme eklemektir.
