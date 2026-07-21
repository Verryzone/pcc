# Telegram CRUD Bot — Implementation Plan

## Spec Reference
`docs/superpowers/specs/2026-07-15-telegram-crud-bot-design.md`

## Prerequisites

### 1. Add dependency
```bash
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
```

### 2. Add environment variable
Add `TELEGRAM_BOT_TOKEN=<token>` to `.env`

---

## Implementation Steps

### Step 1: Create `telegram/telegram.go` — Bot initialization & main handler

**Files:** `telegram/telegram.go` (new), `main.go` (edit), `.env` (edit)

Responsibilities:
- `StartBot(db *gorm.DB)` function — entry point called from main.go
- Read `TELEGRAM_BOT_TOKEN` from env
- Create bot client via `tgbotapi.NewBotAPI()`
- Set up long polling with `bot.GetUpdatesChan()`
- Main update loop: dispatch to callback handler or command handler
- Handle `/start` command → show main menu inline keyboard

**main.go changes:**
- Add import `"main/telegram"`
- Add `go telegram.StartBot(db)` before `ai.InitAi()`

**.env changes:**
- Add `TELEGRAM_BOT_TOKEN=`

---

### Step 2: Create `telegram/handler.go` — State management & callback dispatcher

**Files:** `telegram/handler.go` (new)

Responsibilities:
- Define `UserState` struct (Module, Action, Step, TempData, EditID)
- Define `userStates` map: `map[int64]*UserState`
- `handleCallbackQuery(bot, update, db)` — main dispatcher
- Parse callback data prefix (`nav:`, `act:`, `pick:`, `confirm:`, `cancel`)
- Route to appropriate module handler
- `handleTextMessage(bot, update, db)` — for form input steps
- Helper functions: `showMainMenu()`, `showSubmenu(module)`, `sendOrEdit(chatID, messageID, text, keyboard)`

---

### Step 3: Create `telegram/suhu.go` — Suhu CRUD

**Files:** `telegram/suhu.go` (new)

Functions:
- `handleSuhuCallback(bot, update, db, action)` — route to lihat/tambah/ubah/hapus
- `suhuLihat(chatID, db)` — query all Suhu, format as text list
- `suhuMulaiTambah(chatID, state)` — start tambah flow, set step=0
- `suhuMulaiUbah(chatID, db)` — show item selection keyboard
- `suhuMulaiHapus(chatID, db)` — show item selection keyboard
- `suhuHandleInput(chatID, state, text, db)` — process form input per step
- `suhuKonfirmasi(chatID, state, db)` — save/update/delete to database
- `suhuSubmenu(chatID)` — show suhu submenu keyboard

**Models used:** `models.Suhu` (Id, Lokasi, Suhu, CreatedAt)

---

### Step 4: Create `telegram/informasi.go` — Informasi CRUD

**Files:** `telegram/informasi.go` (new)

Same pattern as suhu.go:
- `handleInformasiCallback(bot, update, db, action)`
- `informasiLihat`, `informasiMulaiTambah`, `informasiMulaiUbah`, `informasiMulaiHapus`
- `informasiHandleInput`, `informasiKonfirmasi`, `informasiSubmenu`

**Models used:** `models.Informasi` (ID via gorm.Model, Judul, Konten, UrlDokumen)

---

### Step 5: Create `telegram/penggajian.go` — Penggajian CRUD

**Files:** `telegram/penggajian.go` (new)

Same pattern, plus auto-calculation:
- Input: NamaPegawai, GajiPokok, JamLembur
- Auto-calculate: GajiKotor, Pajak, GajiBersih (same logic as `controllers.TambahPenggajian`)
- Display shows all fields including calculated ones

**Models used:** `models.Penggajian` (ID, NamaPegawai, GajiPokok, JamLembur, GajiKotor, Pajak, GajiBersih)

---

### Step 6: Create `telegram/pesan.go` — Pesan CRUD

**Files:** `telegram/pesan.go` (new)

Same pattern, but uses Kode (string) as identifier:
- Selection shows Kode instead of numeric ID
- Input: Kode, Balasan

**Models used:** `models.Pesan` (Kode, Balasan, CreatedAt, UpdatedAt)

---

### Step 7: Create `telegram/user.go` — User CRUD

**Files:** `telegram/user.go` (new)

Same pattern, plus:
- Password is SHA1-hashed before saving (same as `controllers.UserTambah`)
- View does NOT display password
- Input: Nama, Username, Password

**Models used:** `models.User` (ID, Nama, Username, Password)

---

### Step 8: Create `telegram/dokumen.go` — Dokumen read-only

**Files:** `telegram/dokumen.go` (new)

Functions:
- `handleDokumenCallback(bot, update, db)` — only lihat action
- `dokumenLihat(chatID, db)` — list all documents with FileUrl as clickable links

**Models used:** `models.Dokumen` (ID, NamaDokumen, FileId, FileUrl)

---

### Step 9: Testing & Verification

- Run `go build` to verify no compilation errors
- Test each module manually via Telegram:
  1. `/start` → main menu appears
  2. Each module → submenu → lihat/tambah/ubah/hapus
  3. Verify data persists in MySQL
  4. Verify error handling (invalid input, empty data)
- Run `go vet ./...` for static analysis

---

## File Summary

| File | Action | Description |
|---|---|---|
| `.env` | EDIT | Add TELEGRAM_BOT_TOKEN |
| `go.mod` | AUTO | go get adds dependency |
| `main.go` | EDIT | Add import + go telegram.StartBot(db) |
| `telegram/telegram.go` | CREATE | Bot init, polling, /start command |
| `telegram/handler.go` | CREATE | State mgmt, callback dispatcher, text handler |
| `telegram/suhu.go` | CREATE | Suhu CRUD handlers |
| `telegram/informasi.go` | CREATE | Informasi CRUD handlers |
| `telegram/penggajian.go` | CREATE | Penggajian CRUD handlers |
| `telegram/pesan.go` | CREATE | Pesan CRUD handlers |
| `telegram/user.go` | CREATE | User CRUD handlers |
| `telegram/dokumen.go` | CREATE | Dokumen read-only handler |
