package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
    profilesFilePath = "/root/Telebot-Radius/files/profiles.json"
)

var (
    botToken string
    adminID  int64
)

type UserProfile struct {
    Balance  int    `json:"balance"`
    Username string `json:"username"`
}

type Profiles map[int64]UserProfile

func loadAuth() {
    data, err := ioutil.ReadFile("/root/Telebot-Radius/files/auth")
    if err != nil {
        log.Fatalf("Error reading auth file: %v", err)
    }
    lines := strings.Split(string(data), "\n")
    if len(lines) >= 2 {
        botToken = strings.TrimSpace(lines[0])
        adminID, err = strconv.ParseInt(strings.TrimSpace(lines[1]), 10, 64)
        if err != nil {
            log.Fatalf("Invalid admin ID in auth file: %v", err)
        }
    } else {
        log.Fatal("Auth file must have at least 2 lines (token and admin chat ID).")
    }
}

func readProfiles() Profiles {
    data, err := ioutil.ReadFile(profilesFilePath)
    if err != nil {
        if os.IsNotExist(err) {
            return Profiles{}
        }
        log.Fatalf("Error reading profiles file: %v", err)
    }
    var profiles Profiles
    if err := json.Unmarshal(data, &profiles); err != nil {
        log.Fatalf("Error parsing profiles file: %v", err)
    }
    return profiles
}

func writeProfiles(profiles Profiles) {
    data, err := json.MarshalIndent(profiles, "", "  ")
    if err != nil {
        log.Fatalf("Error marshalling profiles: %v", err)
    }
    if err := ioutil.WriteFile(profilesFilePath, data, 0644); err != nil {
        log.Fatalf("Error writing profiles file: %v", err)
    }
}

func isAdmin(userID int64) bool {
    return userID == adminID
}

func handleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update, profiles Profiles) {
    userID := update.Message.From.ID
    username := update.Message.From.UserName
    profile, exists := profiles[userID]
    if !exists {
        profile = UserProfile{Balance: 0, Username: username}
        profiles[userID] = profile
        writeProfiles(profiles)
    }

    if isAdmin(userID) {
        commands := "ADMIN ACCESS\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nDaftar perintah:\n1. /isi - Topup saldo user\n2. /saldo - Cek saldo user\n"
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, commands)
        bot.Send(msg)
    } else {
        welcomeMessage := fmt.Sprintf("▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nACCESS USER\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nSelamat datang di Arneta.ID DaloBot\nUsername  :  %s\nUser Id        :  %d\nBalance      :   Rp.%d\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬", username, userID, profile.Balance)
        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("DAPATKAN KODE", "get_code"),
                tgbotapi.NewInlineKeyboardButtonData("TOPUP SALDO", "start_topup"),
            ),
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonURL("HUBUNGI ADMIN", "https://t.me/arnetadotid"),
            ),
        )
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    }
}

func handleMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, profiles Profiles) {
    userID := update.Message.From.ID
    profile, exists := profiles[userID]

    if isAdmin(userID) {
        commands := "▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nADMIN ACCESS\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nDaftar perintah:\n1. /isi - Topup saldo user\n2. /saldo - Cek saldo user\n"
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, commands)
        bot.Send(msg)
    } else if exists {
        balance := profile.Balance
        username := profile.Username
        welcomeMessage := fmt.Sprintf("▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nACCESS USER\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬\nSelamat datang di Arneta.ID DaloBot\nUsername  :  %s\nUser Id        :  %d\nBalance      :   Rp.%d\n▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬", username, userID, balance)
        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("DAPATKAN KODE", "get_code"),
                tgbotapi.NewInlineKeyboardButtonData("TOPUP SALDO", "start_topup"),
            ),
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonURL("HUBUNGI ADMIN", "https://t.me/arnetadotid"),
            ),
        )
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    } else {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Data saldo Anda tidak ditemukan.")
        bot.Send(msg)
    }
}

func handleAddBalance(bot *tgbotapi.BotAPI, update tgbotapi.Update, profiles Profiles) {
    userID := update.Message.From.ID
    if !isAdmin(userID) {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Anda tidak memiliki izin untuk menggunakan perintah ini.")
        bot.Send(msg)
        return
    }

    args := strings.Split(update.Message.Text, " ")
    if len(args) != 3 {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gunakan format: /isi <user_id> <jumlah>")
        bot.Send(msg)
        return
    }

    targetUserID, err := strconv.ParseInt(args[1], 10, 64)
    if err != nil {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ID pengguna tidak valid.")
        bot.Send(msg)
        return
    }

    amount, err := strconv.Atoi(args[2])
    if err != nil {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Jumlah tidak valid.")
        bot.Send(msg)
        return
    }

    profile, exists := profiles[targetUserID]
    if !exists {
        profile = UserProfile{Balance: 0, Username: "tidak memiliki username"}
    }
    profile.Balance += amount
    profiles[targetUserID] = profile
    writeProfiles(profiles)

    adminMessage := fmt.Sprintf("💰 Saldo pengguna %d berhasil ditambahkan sebesar Rp.%d.", targetUserID, amount)
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, adminMessage)
    bot.Send(msg)

    userMessage := fmt.Sprintf("Saldo Anda telah ditambahkan sebesar Rp.%d oleh admin.", amount)
    msg = tgbotapi.NewMessage(targetUserID, userMessage)
    bot.Send(msg)
}

func handleCheckBalance(bot *tgbotapi.BotAPI, update tgbotapi.Update, profiles Profiles) {
    userID := update.Message.From.ID

    if isAdmin(userID) {
        args := strings.Split(update.Message.Text, " ")
        if len(args) != 2 {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Gunakan format: /saldo <user_id>")
            bot.Send(msg)
            return
        }

        targetUserID, err := strconv.ParseInt(args[1], 10, 64)
        if err != nil {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ID pengguna tidak valid.")
            bot.Send(msg)
            return
        }

        profile, exists := profiles[targetUserID]
        if !exists {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Data saldo pengguna tidak ditemukan.")
            bot.Send(msg)
            return
        }

        balanceMessage := fmt.Sprintf("💰 Saldo pengguna %d adalah Rp.%d.", targetUserID, profile.Balance)
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, balanceMessage)
        bot.Send(msg)
    } else {
        profile, exists := profiles[userID]
        if !exists {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Data saldo Anda tidak ditemukan.")
            bot.Send(msg)
            return
        }

        balanceMessage := fmt.Sprintf("💰 Saldo Anda adalah Rp.%d.", profile.Balance)
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, balanceMessage)
        bot.Send(msg)
    }
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, profiles Profiles) {
    callback := update.CallbackQuery
    userID := callback.From.ID

    switch callback.Data {
    case "get_code":
        msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Kode voucher Anda: ABC123456")
        bot.Send(msg)
    case "start_topup":
        msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Untuk topup saldo, hubungi admin di @arnetadotid.")
        bot.Send(msg)
    }
}

func main() {
    loadAuth()

    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Fatalf("Error creating bot: %v", err)
    }

    bot.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    profiles := readProfiles()

    for update := range updates {
        if update.Message != nil {
            switch update.Message.Command() {
            case "start":
                handleStart(bot, update, profiles)
            case "menu":
                handleMenu(bot, update, profiles)
            case "isi":
                handleAddBalance(bot, update, profiles)
            case "saldo":
                handleCheckBalance(bot, update, profiles)
            }
        } else if update.CallbackQuery != nil {
            handleCallbackQuery(bot, update, profiles)
        }
    }
}
