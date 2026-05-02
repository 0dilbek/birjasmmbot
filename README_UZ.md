# BirjaSMM — SMM xizmatlari birjasi uchun Telegram-bot

SMM, videoprodakshn va bloggerlik sohasidagi **buyurtmachilar** (mijozlar) va **ijrochilarni** bog'lash uchun Telegram-bot. Ikki xil interfeys tilini qo'llab-quvvatlaydi: **ruscha** va **o'zbekcha**.

---

## Mundarija

1. [Texnologiyalar](#texnologiyalar)
2. [Loyiha tuzilishi](#loyiha-tuzilishi)
3. [Konfiguratsiya va ishga tushirish](#konfiguratsiya-va-ishga-tushirish)
4. [Ma'lumotlar bazasi](#malumotlar-bazasi)
5. [Foydalanuvchi rollari](#foydalanuvchi-rollari)
6. [Ssenariylar — qadam-baqadam](#ssenariylar)
   - [Yangi foydalanuvchi: til va rol tanlash](#1-yangi-foydalanuvchi)
   - [Buyurtmachini ro'yxatdan o'tkazish](#2-buyurtmachini-royxatdan-otkazish)
   - [Ijrochini ro'yxatdan o'tkazish](#3-ijrochini-royxatdan-otkazish)
   - [Vazifa yaratish (buyurtmachi)](#4-vazifa-yaratish)
   - [Vazifalarni ko'rish va javob berish (ijrochi)](#5-vazifalarni-korish-va-javob-berish)
   - [Ijrochini tanlash (buyurtmachi)](#6-ijrochini-tanlash)
   - [Vazifani yakunlash va sharh qoldirish](#7-vazifani-yakunlash-va-sharh)
   - [Ijrochini verifikatsiyadan o'tkazish](#8-ijrochini-verifikatsiyadan-otkazish)
   - [Buyurtmachini verifikatsiyadan o'tkazish](#9-buyurtmachini-verifikatsiyadan-otkazish)
   - [Obuna va to'lov](#10-obuna-va-tolov)
   - [Tilni o'zgartirish](#11-tilni-ozgartirish)
   - [Administrator paneli](#12-administrator-paneli)
7. [FSM — holatlar mashinasi](#fsm)
8. [Lokalizatsiya](#lokalizatsiya)
9. [Kod tuzilishi](#kod-tuzilishi)

---

## Texnologiyalar

| Komponent | Texnologiya |
|-----------|-----------|
| Til | Go 1.21+ |
| Telegram API | `github.com/go-telegram-bot-api/telegram-bot-api/v5` |
| Ma'lumotlar bazasi | PostgreSQL |
| DB Drayveri | `lib/pq` |
| Konfiguratsiya | `.env` (`github.com/joho/godotenv` orqali) |
| Holatlar (FSM) | In-memory `sync.Map` |

---

## Loyiha tuzilishi

```
BirjaSMM-bot/
├── cmd/bot/
│   └── main.go                  # Kirish nuqtasi, initsializatsiya
├── internal/
│   ├── bot/
│   │   ├── bot.go               # Bot yadrosi, yangilanishlar marshrutizatsiyasi
│   │   ├── menu.go              # Asosiy menyu, profil, mening vazifalarim
│   │   ├── register.go          # Foydalanuvchilarni ro'yxatdan o'tkazish
│   │   ├── task_create.go       # Vazifa yaratish (FSM qadamlari)
│   │   ├── task_view.go         # Vazifalar va ijrochilarni ko'rish
│   │   ├── response.go          # Ijrochining vazifaga javob berishi
│   │   ├── review.go            # Ijrochini tanlash, sharhlar
│   │   ├── payment.go           # Obunalar va to'lov
│   │   ├── verification.go      # Ijrochilar va buyurtmachilarni verifikatsiya qilish
│   │   ├── admin.go             # Administrator paneli
│   │   ├── helpers.go           # Yordamchi funksiyalar
│   │   ├── fsm/
│   │   │   └── fsm.go           # Holatlar mashinasi (in-memory)
│   │   └── keyboards/
│   │       └── keyboards.go     # Faqat administrator uchun klaviaturalar
│   ├── config/
│   │   └── config.go            # .env yuklash
│   ├── db/
│   │   └── db.go                # PostgreSQL ga ulanish, migratsiyalarni ishga tushirish
│   ├── locales/
│   │   └── locales.go           # Barcha RU/UZ tarjimalari va foydalanuvchi klaviaturalari
│   ├── models/
│   │   └── models.go            # Barcha ma'lumotlar tuzilmalari (structs)
│   └── repo/
│       ├── user.go              # Foydalanuvchilar uchun SQL-so'rovlar
│       ├── task.go              # Vazifalar va javoblar uchun SQL-so'rovlar
│       └── payment.go           # To'lovlar uchun SQL-so'rovlar
├── migrations/
│   ├── 001_init.sql             # Asosiy jadvallar
│   ├── 002_payments.sql         # To'lovlar jadvallari
│   ├── 003_language.sql         # users jadvalidagi language ustuni
│   └── 004_verif_admin_msgs.sql # Admin xabarlarini kuzatish jadvali
├── .env.example
├── Makefile
└── go.mod
```

---

## Konfiguratsiya va ishga tushirish

### 1. `.env` faylini nusxalash va to'ldirish

```bash
cp .env.example .env
```

```env
BOT_TOKEN=sizning_tokeningiz_@BotFather_dan
DATABASE_URL=postgres://user:password@localhost:5432/birjasmm?sslmode=disable
ADMIN_IDS=123456789,987654321   # Administratorlar Telegram ID lari vergul bilan
FREE_RESPONSES=5                # Ijrochi uchun bepul javoblar soni
```

### 2. Ma'lumotlar bazasini yaratish

```bash
createdb birjasmm
```

Migratsiyalar bot ishga tushganda **avtomatik** ravishda qo'llaniladi (`migrations/` papkasidagi barcha `.sql` fayllar ketma-ketlikda).

### 3. Botni ishga tushirish

```bash
go run ./cmd/bot/
# yoki
go build -o birjasmm ./cmd/bot/ && ./birjasmm
# yoki Makefile orqali
make run
```

---

## Ma'lumotlar bazasi

### Jadvallar

| Jadval | Vazifasi |
|---------|-----------|
| `users` | Barcha foydalanuvchilar (rol, status, til, verifikatsiya) |
| `client_profiles` | Buyurtmachi profili (ism, biznes, shahar) |
| `executor_profiles` | Ijrochi profili (kategoriya, reyting, portfolio) |
| `tasks` | Buyurtmachilar vazifalari |
| `responses` | Ijrochilarning vazifalarga bergan javoblari |
| `task_assignments` | Vazifa uchun tayinlangan ijrochi |
| `reviews` | Ijrochilar haqida buyurtmachilar sharhlari |
| `subscriptions` | Ijrochilarning faol obunalari |
| `usage_limits` | Ijrochining bepul javoblari limiti |
| `verifications` | Verifikatsiya so'rovlari (video/foto) |
| `payments` | To'lovlar tarixi |
| `payment_admin_msgs` | Administratorlardagi xabarlar ID si (tahrirlash uchun) |
| `verification_admin_msgs`| Administratorlardagi verifikatsiya xabarlari ID si |

---

## Foydalanuvchi rollari

| Rol | Nima qila oladi |
|------|-------|
| **Buyurtmachi** (`client`) | Vazifa yaratish, ijrochini tanlash, vazifani yakunlash, sharh qoldirish, verifikatsiyadan o'tish |
| **Ijrochi** (`executor`) | Vazifalarga javob berish, obuna sotib olish, verifikatsiyadan o'tish, faol buyurtmalarni ko'rish |
| **Administrator** | Foydalanuvchilar, vazifalar, verifikatsiyalar, statistika bilan ishlash, to'lovlarni tasdiqlash |

> Administratorlar `.env` dagi `ADMIN_IDS` orqali aniqlanadi. Ularning interfeysi har doim **rus tilida** bo'ladi.

---

## Ssenariylar

### 1. Yangi foydalanuvchi

**Komanda:** `/start`

1. Bot foydalanuvchi MB (DB) da borligini tekshiradi.
2. Agar bo'lmasa — **til tanlashni** ko'rsatadi: 🇷🇺 Русский / 🇺🇿 O'zbek.
3. Foydalanuvchi tugmani bosadi → til FSM da saqlanadi (ro'yxatdan o'tguncha) va keyin MB ga yoziladi.
4. **Rol tanlash** ko'rsatiladi: Buyurtmachi / Ijrochi.

---

### 2. Buyurtmachini ro'yxatdan o'tkazish

**Trigger:** "Buyurtmachi" rolini tanlash

FSM qadamlari: Ism -> Telefon -> Shahar -> Biznes nomi.
Yakunlangach, buyurtmachi asosiy menyusi ochiladi.

---

### 3. Ijrochini ro'yxatdan o'tkazish

**Trigger:** "Ijrochi" rolini tanlash

FSM qadamlari: Ism -> Telefon -> Shahar -> Kategoriya (SMM/Video/Blogger) -> Tajriba -> Tavsif -> Portfolio.
Yakunlangach, ijrochi asosiy menyusi ochiladi.

---

### 4. Vazifa yaratish

Faqat buyurtmachilar uchun. FSM orqali nomi, tavsifi, kategoriyasi, budjeti, muddati va boshqa ma'lumotlar yig'iladi.
Nashr qilingandan so'ng mos keladigan barcha ijrochilarga xabar boradi.

---

### 5. Vazifalarni ko'rish va javob berish

Ijrochi "Vazifa qidirish" orqali o'z kategoriyasidagi ochiq vazifalarni ko'radi. Vazifaga javob berishda xabar va taklif qilingan narx kiritiladi.
Obunasiz ijrochilar uchun bepul javoblar limiti mavjud.

---

### 6. Ijrochini tanlash

Buyurtmachi o'z vazifasiga kelgan javoblarni ko'radi va ijrochini tanlaydi. Tanlangan ijrochiga buyurtmachi kontaktlari boradi, boshqa ijrochilar esa rad etiladi.

---

### 7. Vazifani yakunlash va sharh

Vazifa yakunlangach, buyurtmachi ijrochiga 1 dan 5 gacha baho beradi va sharh qoldiradi. Bu ijrochining reytingiga ta'sir qiladi.

---

### 8. Verifikatsiya

Ijrochilar videoxabar yuborish orqali, buyurtmachilar esa hujjat rasmini yuborish orqali verifikatsiyadan o'tishlari mumkin. Adminlar tasdiqlagandan so'ng profilda ✅ belgisi paydo bo'ladi.

---

### 9. Obuna va to'lov

Ijrochilar BASIC yoki PRO obunalarini sotib olishlari mumkin. To'lov cheki yuborilgach, adminlar uni tekshirib tasdiqlaydilar.

---

## Lokalizatsiya

Barcha tarjimalar `internal/locales/locales.go` faylida joylashgan. Bot foydalanuvchining bazada saqlangan tiliga qarab javob qaytaradi. Bildirishnomalar ham har bir foydalanuvchining o'z tilida yuboriladi.

---

## Administrator paneli

`/admin` komandasi orqali kiriladi. Unda foydalanuvchilarni boshqarish (bloklash/tahrirlash), vazifalarni o'chirish, verifikatsiya va to'lovlarni tasdiqlash hamda statistika bo'limlari mavjud.
Endilikda admin paneli paginatsiya va qulay interfeys bilan yangilangan.
