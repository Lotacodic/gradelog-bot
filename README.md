# 🎓 GradeLog Bot

> A Telegram bot that helps students track, calculate, and manage their academic grades — anytime, anywhere.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Version](https://img.shields.io/badge/version-v1.0.0-blue)](https://github.com/Lotacodic/gradelog-bot/releases/tag/v1.0.0)
[![Telegram Bot](https://img.shields.io/badge/Telegram-Bot-2CA5E0?style=flat&logo=telegram)](https://t.me/testgradelogbot)

---

## 📌 Overview

GradeLog Bot is a fully functional Telegram bot built with Go that allows students to securely register, log their courses, and instantly calculate their GPA and CGPA — all from within Telegram. No spreadsheets, no manual calculations, no friction.

The project was built from the ground up as a real-world backend system, covering authentication, data persistence, business logic, and a conversational bot interface.

---

## 🚀 Live Demo

Try it live on Telegram: [@testgradelogbot](https://t.me/testgradelogbot)

---

## ✨ Features

- 🔐 **Secure Authentication** — Register and login with bcrypt-hashed passwords
- 📚 **Semester Management** — Create and organize courses by semester
- 📝 **Course CRUD** — Add, view, update, and delete courses with full validation
- 📊 **GPA Calculation** — Instant per-semester GPA on a 5.0 scale
- 🏆 **CGPA Calculation** — Cumulative GPA calculated across all semesters
- 💬 **Conversational UI** — Clean inline keyboard interface built for Telegram
- 💾 **Data Persistence** — JSON-based local database that survives restarts
- 🔒 **Session Management** — Per-user session tracking for multi-step conversations

---

## 🏗️ System Architecture

```
gradelog-bot/
├── main.go               # Entry point — loads env and starts the bot
├── bot/
│   └── bot.go            # Telegram bot handlers, session management, UI logic
├── services/
│   └── service.go        # Business logic — auth, CRUD, GPA/CGPA calculation
├── storage/
│   └── storage.go        # Data persistence — JSON load and save
├── models/
│   └── models.go         # Core data structures — Student, Course, Semester, Database
├── go.mod                # Go module definition
└── go.sum                # Dependency checksums
```

**Request Flow:**
```
Telegram User → bot/bot.go → services/service.go → storage/storage.go → database.json
```

The architecture follows a clean layered approach — the bot layer handles user interaction, the service layer owns business logic, and the storage layer handles persistence. Each layer is fully decoupled.

---

## 🛠️ Technical Stack

| Technology | Purpose | Why It Was Chosen |
|---|---|---|
| **Go** | Core language | Statically typed, fast compilation, excellent concurrency support |
| **telebot.v3** | Telegram bot framework | Clean API, inline keyboard support, active maintenance |
| **bcrypt** | Password hashing | Industry-standard one-way hashing for secure credential storage |
| **UUID** | Unique student IDs | Collision-resistant IDs for reliable student identification |
| **encoding/json** | Data persistence | Zero-dependency, human-readable local storage |

---

## ⚙️ Getting Started

### Prerequisites

- [Go 1.21+](https://golang.org/dl/) installed
- A Telegram Bot Token from [@BotFather](https://t.me/BotFather)

### Installation

**1. Clone the repository**
```bash
git clone https://github.com/Lotacodic/gradelog-bot.git
cd gradelog-bot
```

**2. Install dependencies**
```bash
go mod tidy
```

**3. Set your bot token**

Create a `.env` file in the root directory:
```bash
BOT_TOKEN=your_telegram_bot_token_here
```

Or export it directly in your terminal:
```bash
export BOT_TOKEN=your_telegram_bot_token_here
```

**4. Run the bot**
```bash
go run main.go
```

**5. Open Telegram and send `/start` to your bot**

---

## 📖 Usage

| Action | How |
|---|---|
| Register | Tap 📝 Register and follow the prompts |
| Login | Tap 🔐 Login with your email and password |
| Add Semester | Tap ➕ Add Semester from the dashboard |
| Add Course | Tap ➕ Add Course, select semester, enter details |
| View Courses | Tap 📚 View Courses and select a semester |
| Update Course | Tap ✏️ Update Course and follow the prompts |
| Delete Course | Tap 🗑️ Delete Course and confirm |
| Calculate GPA | Tap 📊 Calculate GPA and select a semester |
| Calculate CGPA | Tap 🏆 Calculate CGPA for cumulative result |

---

## 🗺️ Roadmap

GradeLog Bot is designed to grow beyond grade tracking into a full academic assistant. Planned features include:

### 🤖 Academic Intelligence
- AI-powered study tips and recommendations based on weak courses
- Grade trend analysis — detect if a student is improving or declining across semesters
- Predictive alerts when a student is at risk of a GPA drop

### 🔔 Notifications & Reminders
- Scheduled reminders to log grades after exam periods
- Automatic alerts when CGPA drops below a student-defined threshold
- End-of-semester performance summary delivered via Telegram

### 📈 Data & Insights
- Export full academic transcript as a PDF directly from the bot
- Visual grade distribution charts per semester
- Best and worst performing course highlights
- Semester-over-semester performance comparison

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome. Feel free to open an issue or submit a pull request.

---

## 📄 License

This project is licensed under the MIT License.

---

<p align="center">Built with ❤️ and Go</p>
