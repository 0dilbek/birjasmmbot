package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/locales"
	"github.com/birjasmm/bot/internal/models"
	"github.com/birjasmm/bot/internal/repo"
)

func (b *Bot) cmdStart(m *tgbotapi.Message) {
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user != nil {
		b.showMainMenu(m.Chat.ID, user)
		return
	}
	fsm.Clear(m.From.ID)
	b.sendKb(m.Chat.ID, locales.New("ru").T("lang_choice"), locales.LangChoice())
}

func (b *Bot) showMainMenu(chatID int64, user *models.User) {
	loc := b.loc(user.TelegramID)
	b.showMainMenuLoc(chatID, user, loc)
}

func (b *Bot) showMainMenuLoc(chatID int64, user *models.User, loc *locales.Loc) {
	if user.Role == models.RoleClient {
		cp, _ := b.users.GetClientProfile(user.ID)
		name := ""
		if cp != nil {
			name = cp.Name
		}
		b.sendKb(chatID, loc.T("client_welcome", name), loc.ClientMenu())
	} else {
		ep, _ := b.users.GetExecutorProfile(user.ID)
		name := ""
		if ep != nil {
			name = ep.Name
		}
		b.sendKb(chatID, loc.T("exec_welcome", name), loc.ExecutorMenu())
	}
}

func (b *Bot) handleMenu(m *tgbotapi.Message) {
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user == nil {
		b.cmdStart(m)
		return
	}
	if user.Status == models.StatusBlocked {
		b.send(m.Chat.ID, b.loc(m.From.ID).T("account_blocked"))
		return
	}

	switch locales.MenuAction(m.Text) {
	case "create_task":
		b.startCreateTask(m)
	case "my_tasks":
		b.showMyTasks(m)
	case "profile":
		b.showProfile(m)
	case "executors":
		b.showTopExecutors(m.Chat.ID, "", "", 0, m.From.ID)
	case "become_verified":
		b.startClientVerif(m)
	case "in_progress":
		b.showExecutorActiveTasks(m)
	case "find_orders":
		b.showTasksForExecutor(m)
	case "my_responses":
		b.showMyResponses(m)
	case "subscription":
		b.showSubscription(m.Chat.ID, m.From.ID)
	case "change_lang":
		b.sendKb(m.Chat.ID, locales.New("ru").T("lang_choice"), locales.LangChoice())
	default:
		b.showMainMenu(m.Chat.ID, user)
	}
}

func (b *Bot) showTasksForExecutor(m *tgbotapi.Message) {
	user, _ := b.users.GetByTelegramID(m.From.ID)
	f := repo.TaskFilter{PageSize: 10}
	if user != nil {
		ep, _ := b.users.GetExecutorProfile(user.ID)
		if ep != nil && string(ep.Category) != "" {
			f.Category = string(ep.Category)
		}
		sub, _ := b.users.GetSubscription(user.ID)
		if sub == nil || sub.Type != models.SubPro {
			f.ExcludeUrgent = true
		}
	}
	b.showOpenTasks(m.Chat.ID, f, m.From.ID)
}

func (b *Bot) showMyTasks(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user == nil {
		return
	}
	tasks, _ := b.tasks.ListByClient(user.ID)
	if len(tasks) == 0 {
		b.send(m.Chat.ID, loc.T("no_tasks"))
		return
	}
	for _, t := range tasks {
		urgent := ""
		if t.IsUrgent {
			urgent = " ⚡"
		}
		text := fmt.Sprintf("%s <b>%s</b>%s\nОткликов: %d",
			taskEmoji(t.Status), t.Title, urgent, t.ResponseCount)

		rows := [][]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(loc.T("btn_view_detail"), fmt.Sprintf("view_task:%d", t.ID)),
			),
		}
		if t.Status == models.TaskOpen {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(loc.T("btn_cancel_inline"), fmt.Sprintf("cancel_task:%d", t.ID)),
			))
		}
		b.sendKb(m.Chat.ID, text, tgbotapi.NewInlineKeyboardMarkup(rows...))
	}
}

func (b *Bot) showProfile(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user == nil {
		return
	}

	if user.Role == models.RoleClient {
		cp, _ := b.users.GetClientProfile(user.ID)
		if cp == nil {
			b.send(m.Chat.ID, loc.T("profile_not_found"))
			return
		}
		verified := loc.T("profile_not_verified")
		if cp.IsVerified {
			verified = loc.T("profile_verified")
		}
		b.send(m.Chat.ID, loc.T("client_profile_text",
			cp.Name, cp.BusinessName, cp.City, verified, user.Phone))
		return
	}

	ep, _ := b.users.GetExecutorProfile(user.ID)
	if ep == nil {
		b.send(m.Chat.ID, loc.T("profile_not_found"))
		return
	}

	sub, _ := b.users.GetSubscription(user.ID)
	limits, _ := b.users.GetOrCreateLimits(user.ID, b.cfg.FreeResponses)

	subText := loc.T("sub_none_label")
	if sub != nil {
		subText = loc.T("sub_active_label", strings.ToUpper(string(sub.Type)), sub.EndDate.Format("02.01.2006"))
	}
	freeLeft := 0
	if limits != nil {
		freeLeft = limits.FreeResponsesLeft
	}
	verifiedLabel := loc.T("exec_not_verified_label")
	if ep.IsVerified {
		verifiedLabel = loc.T("exec_verified_label")
	}

	b.send(m.Chat.ID, loc.T("exec_profile_text",
		ep.Name, loc.CatLabel(string(ep.Category)), ep.City,
		ep.ExperienceYears, ep.Rating, ep.CompletedOrders,
		verifiedLabel, subText, freeLeft,
		ifEmpty(ep.PortfolioLinks, "—"),
	))
}

func (b *Bot) showMyResponses(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user == nil {
		return
	}
	resps, _ := b.resps.ListForExecutor(user.ID)
	if len(resps) == 0 {
		b.send(m.Chat.ID, loc.T("no_my_responses"))
		return
	}
	for _, r := range resps {
		task, _ := b.tasks.GetByID(r.TaskID)
		title := fmt.Sprintf("#%d", r.TaskID)
		if task != nil {
			title = task.Title
		}
		text := fmt.Sprintf("📋 <b>%s</b>\n%s", title, loc.RespLabel(r.Status))
		if r.Status == models.RespActive {
			b.sendKb(m.Chat.ID, text, loc.WithdrawResponse(r.ID))
		} else {
			b.send(m.Chat.ID, text)
		}
	}
}

func (b *Bot) showExecutorActiveTasks(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user == nil {
		return
	}
	assignments, _ := b.resps.GetActiveAssignments(user.ID)
	if len(assignments) == 0 {
		b.send(m.Chat.ID, loc.T("no_active_tasks"))
		return
	}
	b.send(m.Chat.ID, loc.T("active_tasks_title"))
	for _, a := range assignments {
		task, _ := b.tasks.GetByID(a.TaskID)
		if task == nil {
			continue
		}
		client, _ := b.users.GetByID(task.ClientID)
		cp, _ := b.users.GetClientProfile(task.ClientID)

		clientName := fmt.Sprintf("ID:%d", task.ClientID)
		clientContact := ""
		if client != nil {
			clientContact = loc.T("client_tg_line", client.Username)
			if client.Phone != "" {
				clientContact += loc.T("client_phone_line", client.Phone)
			}
		}
		if cp != nil {
			clientName = cp.Name
		}

		b.send(m.Chat.ID, loc.T("active_task_text",
			task.Title, loc.BudgetText(task),
			clientName, clientContact,
			a.AssignedAt.Format("02.01.2006"),
		))
	}
}

func (b *Bot) showSubscription(chatID, tgID int64) {
	loc := b.loc(tgID)
	b.sendKb(chatID, loc.T("sub_menu"), loc.SubscriptionMenu())
}

func (b *Bot) cbSub(cq *tgbotapi.CallbackQuery, subType string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		return
	}
	if user.Role != models.RoleExecutor {
		b.send(cq.Message.Chat.ID, b.loc(tgID).T("sub_only_exec"))
		return
	}
	b.startPaymentFlow(cq.Message.Chat.ID, tgID, subType)
}

func (b *Bot) cbBack(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		return
	}
	switch param {
	case "tasks":
		b.showOpenTasks(cq.Message.Chat.ID, repo.TaskFilter{PageSize: 10}, cq.From.ID)
	default:
		b.showMainMenu(cq.Message.Chat.ID, user)
	}
}
