# Telegram CRUD Bot Integration — Design Spec

## Overview

Integrate all existing CRUD features (Suhu, Informasi, Penggajian, Pesan, User, Dokumen) into a Telegram bot using inline keyboard UI. The bot will access GORM directly (same pattern as existing controllers) without HTTP overhead.

## Goals

- All 5 CRUD modules accessible via Telegram inline keyboard
- Dokumen module: view-only (list documents with download links)
- No authentication required (internal/academic project)
- Bot runs alongside existing WhatsApp bot and Gin API server
- Consistent UX: menu → submenu → action → confirmation flow

## Architecture

### Approach: Direct GORM Access

The bot calls GORM directly using the shared `*gorm.DB` connection, bypassing HTTP. This mirrors how the existing controllers work but without `*gin.Context`.

### File Structure

```
pcc/
├── telegram/
│   ├── telegram.go      # Bot init, long polling, main handler routing
│   ├── handler.go       # Callback query dispatcher, state machine
│   ├── suhu.go          # CRUD handlers for Suhu module
│   ├── informasi.go     # CRUD handlers for Informasi module
│   ├── penggajian.go    # CRUD handlers for Penggajian module
│   ├── pesan.go         # CRUD handlers for Pesan module
│   ├── user.go          # CRUD handlers for User module
│   └── dokumen.go       # Read-only handler for Dokumen module
├── controllers/          # UNCHANGED
├── models/               # UNCHANGED (used directly by bot)
├── wa/                   # UNCHANGED
├── main.go               # ADD: go telegram.StartBot(db)
├── koneksi.go            # UNCHANGED
└── .env                  # ADD: TELEGRAM_BOT_TOKEN
```

### Integration Point

In `main.go`, add one line before `ai.InitAi()`:

```go
go telegram.StartBot(db)
```

This runs the Telegram bot as a goroutine, same pattern as `wa.KonekWa(db)`.

## Dependencies

- `github.com/go-telegram-bot-api/telegram-bot-api/v5` — Telegram Bot API wrapper for Go

## Environment Variables

Add to `.env`:

```
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

## UI/UX Design

### Main Menu (/start command)

```
🤖 PCC Bot - Menu Utama
Silakan pilih fitur:

[🌡️ Suhu] [📰 Informasi]
[💰 Penggajian] [💬 Pesan]
[👤 User] [📄 Dokumen]
```

### Submenu (example: Suhu)

```
🌡️ Menu Suhu

[📋 Lihat Semua] [➕ Tambah]
[✏️ Ubah] [🗑️ Hapus]
[🔙 Kembali]
```

### CRUD Flows

#### View (Lihat)

Display data as formatted text list:

```
📋 Data Suhu:

ID: 1 | Lokasi: Gedung A | Suhu: 36.5°C
ID: 2 | Lokasi: Gedung B | Suhu: 37.0°C

[🔙 Kembali]
```

#### Add (Tambah) — Step-by-step form

```
Bot: "Masukkan lokasi:"
User: "Gedung A"
Bot: "Masukkan suhu (°C):"
User: "36.5"
Bot: "Konfirmasi simpan?\nLokasi: Gedung A\nSuhu: 36.5°C"
     [✅ Simpan] [❌ Batal]
Bot: "✅ Data suhu berhasil ditambahkan!"
     [🔙 Kembali ke Menu Suhu]
```

#### Edit (Ubah) — Select then edit

```
Bot: "Pilih data yang akan diubah:"
[ID:1] Gedung A - 36.5°C
[ID:2] Gedung B - 37.0°C
[🔙 Batal]
→ User clicks ID:1
Bot: "Masukkan lokasi baru (ketik '-' untuk skip):"
User: "Gedung A Lantai 2"
Bot: "Masukkan suhu baru (ketik '-' untuk skip):"
User: "37.2"
Bot: "Konfirmasi ubah?\nLokasi: Gedung A Lantai 2\nSuhu: 37.2°C"
     [✅ Simpan] [❌ Batal]
```

#### Delete (Hapus) — Select then confirm

```
Bot: "Pilih data yang akan dihapus:"
[ID:1] Gedung A - 36.5°C
[ID:2] Gedung B - 37.0°C
[🔙 Batal]
→ User clicks ID:1
Bot: "Yakin hapus data Gedung A - 36.5°C?"
     [✅ Ya, Hapus] [❌ Batal]
Bot: "✅ Data berhasil dihapus!"
```

## Callback Data Format

Inline keyboard buttons use structured callback data:

| Action | Format | Example |
|---|---|---|
| Navigate to module | `nav:{module}` | `nav:suhu` |
| Navigate to main menu | `nav:main` | `nav:main` |
| Perform action | `act:{module}:{action}` | `act:suhu:lihat` |
| Select item for edit/delete | `pick:{module}:{action}:{id}` | `pick:suhu:ubah:3` |
| Confirm action | `confirm:{module}:{action}[:{id}]` | `confirm:suhu:tambah` |
| Cancel | `cancel` | `cancel` |

## State Management

In-memory map: `map[int64]*UserState` (key = Telegram User ID)

```go
type UserState struct {
    Module    string            // "suhu", "informasi", "penggajian", "pesan", "user"
    Action    string            // "tambah", "ubah"
    Step      int               // Current form step (0-indexed)
    TempData  map[string]string // Temporary form data before saving
    EditID    uint              // ID of record being edited/deleted
}
```

- State is created when user starts a form action (tambah/ubah)
- State is cleared after save, cancel, or navigation
- State persists only in memory (lost on restart)

## Module Specifications

### Suhu
| Operation | Fields | Validation |
|---|---|---|
| View | ID, Lokasi, Suhu, CreatedAt | - |
| Add | Lokasi (string), Suhu (float) | Both required, Suhu must be numeric |
| Edit | Select ID, then edit Lokasi and/or Suhu | '-' to skip field |
| Delete | Select ID, confirm | - |

### Informasi
| Operation | Fields | Validation |
|---|---|---|
| View | ID, Judul, Konten, UrlDokumen | - |
| Add | Judul, Konten, UrlDokumen | All required |
| Edit | Select ID, then edit fields | '-' to skip field |
| Delete | Select ID, confirm | - |

### Penggajian
| Operation | Fields | Validation |
|---|---|---|
| View | ID, NamaPegawai, GajiPokok, JamLembur, GajiKotor, Pajak, GajiBersih | - |
| Add | NamaPegawai, GajiPokok, JamLembur | Auto-calculate: GajiKotor, Pajak (5% if >5jt), GajiBersih |
| Edit | Select ID, then edit fields | Auto-recalculate on change |
| Delete | Select ID, confirm | - |

Penggajian calculation (from existing controller):
- UangLembur = JamLembur × 50000
- GajiKotor = GajiPokok + UangLembur
- Pajak = GajiKotor × 0.05 (if GajiKotor > 5,000,000)
- GajiBersih = GajiKotor - Pajak

### Pesan
| Operation | Fields | Validation |
|---|---|---|
| View | Kode, Balasan, CreatedAt | - |
| Add | Kode, Balasan | Both required |
| Edit | Select by Kode, then edit Balasan | - |
| Delete | Select by Kode, confirm | - |

Note: Pesan uses Kode (string) as primary identifier, not numeric ID.

### User
| Operation | Fields | Validation |
|---|---|---|
| View | ID, Nama, Username | Password NOT displayed |
| Add | Nama, Username, Password | All required, Password SHA1-hashed |
| Edit | Select ID, then edit fields | Password re-hashed if changed |
| Delete | Select ID, confirm | - |

### Dokumen (Read-only)
| Operation | Fields |
|---|---|
| View | ID, NamaDokumen, FileUrl |

Display as clickable links. No add/edit/delete via bot.

## Error Handling

- **Invalid input** (e.g., non-numeric suhu): Bot sends error message and repeats the current step
- **Database error**: Bot sends "Gagal menyimpan data" and returns to submenu
- **Invalid callback data**: Bot ignores and logs error
- **No data found**: Bot shows "Tidak ada data" with back button

## What Does NOT Change

- `controllers/` — All existing HTTP API controllers remain untouched
- `models/` — All models used directly by bot
- `wa/` — WhatsApp bot runs independently alongside
- `koneksi.go` — Shared database connection
- `ai/` — AI module unaffected

## What Changes

- `main.go` — Add import + `go telegram.StartBot(db)`
- `go.mod` — Add `github.com/go-telegram-bot-api/telegram-bot-api/v5`
- `.env` — Add `TELEGRAM_BOT_TOKEN`
- New: `telegram/` directory with 8 Go files
