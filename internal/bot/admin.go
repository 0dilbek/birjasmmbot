package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/bot/keyboards"
	"github.com/birjasmm/bot/internal/models"
	"github.com/birjasmm/bot/internal/repo"
)

func (b *Bot) isAdmin(tgID int64) bool {
	return b.adminIDs[tgID]
}

func (b *Bot) cmdAdmin(m *tgbotapi.Message) {
	if !b.isAdmin(m.From.ID) {
		b.send(m.Chat.ID, "Нет доступа.")
		return
	}
	b.sendKb(m.Chat.ID, "🛠 <b>Панель администратора</b>", keyboards.AdminMenu())
}

func (b *Bot) cbAdmin(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")

	parts := strings.Split(param, ":")
	action := parts[0]
	page := 0
	if len(parts) > 1 {
		page = int(parseInt64(parts[len(parts)-1]))
	}

	switch action {
	case "menu":
		edit := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, "🛠 <b>Панель администратора</b>")
		edit.ParseMode = "HTML"
		kb := keyboards.AdminMenu()
		edit.ReplyMarkup = &kb
		b.api.Send(edit)
	case "users":
		b.adminUsers(cq.Message.Chat.ID, cq.Message.MessageID, page)
	case "tasks":
		b.adminTasks(cq.Message.Chat.ID, cq.Message.MessageID, page)
	case "view_task":
		b.adminViewTask(cq.Message.Chat.ID, parseInt64(parts[1]))
	case "verifs":
		b.adminVerifs(cq.Message.Chat.ID, cq.Message.MessageID, page)
	case "view_verif":
		b.adminViewVerif(cq.Message.Chat.ID, parseInt64(parts[1]))
	case "stats":
		b.adminStats(cq.Message.Chat.ID)
	}
}

func (b *Bot) adminViewVerif(chatID int64, verifID int64) {
	v, _ := b.users.GetVerificationByID(verifID)
	if v == nil {
		b.send(chatID, "Заявка не найдена.")
		return
	}
	u, _ := b.users.GetByID(v.UserID)
	uname := "???"
	if u != nil {
		uname = "@" + u.Username
	}

	b.send(chatID, fmt.Sprintf("🔐 Верификация от %s\nUserID: %d\nFile ID: <code>%s</code>", uname, v.UserID, v.VideoFileID))
	b.sendKb(chatID, "Действие:", keyboards.VerifReview(v.ID))
}


func (b *Bot) adminViewTask(chatID int64, taskID int64) {
	t, _ := b.tasks.GetByID(taskID)
	if t == nil {
		b.send(chatID, "Задача не найдена.")
		return
	}
	text := fmt.Sprintf("📋 <b>%s</b>\n\n%s\n\nБюджет: %d-%d\nСтатус: %s", t.Title, t.Description, *t.BudgetFrom, *t.BudgetTo, t.Status)
	b.sendKb(chatID, text, keyboards.TaskAdminActions(t.ID))
}

func (b *Bot) adminUsers(chatID int64, msgID int, page int) {
	pageSize := 10
	offset := page * pageSize
	users, _ := b.users.ListAll(pageSize, offset)
	total, _ := b.users.Count()

	if len(users) == 0 {
		b.send(chatID, "Пользователей нет.")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👥 <b>Список пользователей (Стр %d)</b>\n\n", page+1))
	var ids []int64
	for i, u := range users {
		sb.WriteString(fmt.Sprintf("%d. @%s | %s | %s\n", i+1, u.Username, u.Role, u.Status))
		ids = append(ids, u.ID)
	}

	kbRows := keyboards.NumberedButtons("admin_edit", ids, 5)
	if pgRow := keyboards.AdminPagination("admin:users", page, (total+pageSize-1)/pageSize); pgRow != nil {
		kbRows = append(kbRows, pgRow)
	}
	kbRows = append(kbRows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin:menu")))
	kb := tgbotapi.NewInlineKeyboardMarkup(kbRows...)

	if msgID == 0 {
		b.sendKb(chatID, sb.String(), kb)
	} else {
		edit := tgbotapi.NewEditMessageText(chatID, msgID, sb.String())
		edit.ParseMode = "HTML"
		edit.ReplyMarkup = &kb
		b.api.Send(edit)
	}
}

func (b *Bot) adminTasks(chatID int64, msgID int, page int) {
	pageSize := 10
	tasks, _ := b.tasks.ListOpen(repo.TaskFilter{PageSize: pageSize, Page: page})
	total, _ := b.tasks.CountOpen()

	if len(tasks) == 0 {
		b.send(chatID, "Задач нет.")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 <b>Список задач (Стр %d)</b>\n\n", page+1))
	var ids []int64
	for i, t := range tasks {
		sb.WriteString(fmt.Sprintf("%d. <b>%s</b> | %s\n", i+1, t.Title, t.Status))
		ids = append(ids, t.ID)
	}

	kbRows := keyboards.NumberedButtons("admin:view_task", ids, 5)
	if pgRow := keyboards.AdminPagination("admin:tasks", page, (total+pageSize-1)/pageSize); pgRow != nil {
		kbRows = append(kbRows, pgRow)
	}
	kbRows = append(kbRows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin:menu")))
	kb := tgbotapi.NewInlineKeyboardMarkup(kbRows...)

	if msgID == 0 {
		b.sendKb(chatID, sb.String(), kb)
	} else {
		edit := tgbotapi.NewEditMessageText(chatID, msgID, sb.String())
		edit.ParseMode = "HTML"
		edit.ReplyMarkup = &kb
		b.api.Send(edit)
	}
}

func (b *Bot) adminVerifs(chatID int64, msgID int, page int) {
	// Pending verifications usually don't need pagination if they are few, but let's follow the request
	verifs, _ := b.users.GetPendingVerifications() // This currently returns all pending
	if len(verifs) == 0 {
		b.send(chatID, "Заявок на верификацию нет.")
		return
	}

	pageSize := 10
	total := len(verifs)
	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	if start >= total {
		page = 0
		start = 0
		end = pageSize
		if end > total {
			end = total
		}
	}
	
	pageVerifs := verifs[start:end]

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔐 <b>Верификации (Стр %d)</b>\n\n", page+1))
	var ids []int64
	for i, v := range pageVerifs {
		u, _ := b.users.GetByID(v.UserID)
		uname := "???"
		if u != nil {
			uname = "@" + u.Username
		}
		sb.WriteString(fmt.Sprintf("%d. %s | UserID:%d\n", i+1, uname, v.UserID))
		ids = append(ids, v.ID)
	}

	// For verifications, clicking a number should probably show the video and approve/reject buttons
	kbRows := keyboards.NumberedButtons("admin:view_verif", ids, 5)
	if pgRow := keyboards.AdminPagination("admin:verifs", page, (total+pageSize-1)/pageSize); pgRow != nil {
		kbRows = append(kbRows, pgRow)
	}
	kbRows = append(kbRows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin:menu")))
	kb := tgbotapi.NewInlineKeyboardMarkup(kbRows...)

	if msgID == 0 {
		b.sendKb(chatID, sb.String(), kb)
	} else {
		edit := tgbotapi.NewEditMessageText(chatID, msgID, sb.String())
		edit.ParseMode = "HTML"
		edit.ReplyMarkup = &kb
		b.api.Send(edit)
	}
}


func (b *Bot) adminStats(chatID int64) {
	userCount, _ := b.users.Count()
	taskCount, _ := b.tasks.CountAll()
	openCount, _ := b.tasks.CountOpen()
	respCount, _ := b.resps.CountAll()
	b.send(chatID, fmt.Sprintf(
		"📊 <b>Статистика</b>\n\n👥 Пользователей: %d\n📋 Всего задач: %d\n🟢 Открытых: %d\n💬 Откликов: %d",
		userCount, taskCount, openCount, respCount,
	))
}

func (b *Bot) cbBlockUser(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")
	b.users.Block(parseInt64(param))
	b.send(cq.Message.Chat.ID, "🚫 Пользователь заблокирован.")
}

func (b *Bot) cbUnblockUser(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")
	b.users.Unblock(parseInt64(param))
	b.send(cq.Message.Chat.ID, "✅ Пользователь разблокирован.")
}

func (b *Bot) cbDelTask(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")
	b.tasks.Delete(parseInt64(param))
	b.send(cq.Message.Chat.ID, "🗑 Задача удалена.")
}

// cbAdminEdit shows editable fields for a user
func (b *Bot) cbAdminEdit(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")
	userID := parseInt64(param)
	u, _ := b.users.GetByID(userID)
	if u == nil {
		b.send(cq.Message.Chat.ID, "Пользователь не найден.")
		return
	}
	b.sendKb(cq.Message.Chat.ID,
		fmt.Sprintf("✏️ Редактирование @%s (%s)\nВыберите поле:", u.Username, u.Role),
		keyboards.UserEditFields(userID, string(u.Role)),
	)
}

// cbAdminEditField stores userID+field in FSM and asks for new value
func (b *Bot) cbAdminEditField(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")
	// param = "userID:field"
	parts := strings.SplitN(param, ":", 2)
	if len(parts) != 2 {
		return
	}
	fsm.SetVal(cq.From.ID, "edit_user_id", parseInt64(parts[0]))
	fsm.SetVal(cq.From.ID, "edit_field", parts[1])
	fsm.Set(cq.From.ID, fsm.StateAdminEditUser)

	fieldNames := map[string]string{
		"name": "имя", "city": "город", "business_name": "название бизнеса",
		"description": "описание", "portfolio_links": "ссылки портфолио",
	}
	label := fieldNames[parts[1]]
	if label == "" {
		label = parts[1]
	}
	b.send(cq.Message.Chat.ID, fmt.Sprintf("✏️ Введите новое значение для поля <b>%s</b>:", label))
}

// onAdminEditValue applies the field update
func (b *Bot) onAdminEditValue(m *tgbotapi.Message) {
	if !b.isAdmin(m.From.ID) {
		fsm.Clear(m.From.ID)
		return
	}
	tgID := m.From.ID
	userID := fsm.GetInt64(tgID, "edit_user_id")
	field := fsm.GetStr(tgID, "edit_field")
	value := strings.TrimSpace(m.Text)
	fsm.Clear(tgID)

	if userID == 0 || field == "" || value == "" {
		b.send(m.Chat.ID, "Ошибка: данные не найдены.")
		return
	}

	u, _ := b.users.GetByID(userID)
	if u == nil {
		b.send(m.Chat.ID, "Пользователь не найден.")
		return
	}

	var err error
	if u.Role == models.RoleClient {
		err = b.users.UpdateClientField(userID, field, value)
	} else {
		err = b.users.UpdateExecutorField(userID, field, value)
	}

	if err != nil {
		b.send(m.Chat.ID, "Ошибка при обновлении.")
		return
	}
	b.send(m.Chat.ID, fmt.Sprintf("✅ Поле <b>%s</b> обновлено для @%s.", field, u.Username))
	b.send(u.TelegramID, "ℹ️ Администратор обновил ваш профиль.")
}

// parseInt64 is a shared helper used across admin.go and other files
func parseInt64(s string) int64 {
	s = strings.TrimSpace(s)
	var v int64
	fmt.Sscan(s, &v)
	return v
}
