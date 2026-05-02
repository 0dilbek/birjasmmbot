package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/models"
	"github.com/birjasmm/bot/internal/repo"
)

func (b *Bot) showOpenTasks(chatID int64, f repo.TaskFilter, tgID int64) {
	loc := b.loc(tgID)
	if f.PageSize == 0 {
		f.PageSize = 10
	}
	tasks, _ := b.tasks.ListOpen(f)
	if len(tasks) == 0 {
		b.send(chatID, loc.T("no_open_tasks"))
		return
	}
	b.send(chatID, loc.T("open_tasks_title"))
	for _, t := range tasks {
		urgent := ""
		if t.IsUrgent {
			urgent = " ⚡"
		}
		text := fmt.Sprintf("<b>%s</b>%s\n%s | %s\nОткликов: %d/%d",
			t.Title, urgent, loc.CatLabel(string(t.Category)), loc.BudgetText(t),
			t.ResponseCount, t.MaxResponses)
		b.sendKb(chatID, text, tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(loc.T("btn_view"), fmt.Sprintf("view_task:%d", t.ID)),
			),
		))
	}

	b.sendKb(chatID, loc.T("filter_label"), loc.FilterMenu())
	if len(tasks) == f.PageSize {
		b.sendKb(chatID, loc.T("page_label", f.Page+1), loc.TaskNextPage(f.Page+1, f.Category, f.BudgetMin))
	}
}

func (b *Bot) cbTaskPage(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	// param = "page:cat:budgetMin"
	parts := strings.SplitN(param, ":", 3)
	if len(parts) != 3 {
		return
	}
	page, _ := strconv.Atoi(parts[0])
	cat := parts[1]
	if cat == "_" {
		cat = ""
	}
	budgetMin, _ := strconv.ParseInt(parts[2], 10, 64)

	tgID := cq.From.ID
	f := repo.TaskFilter{Category: cat, BudgetMin: budgetMin, Page: page, PageSize: 10}
	user, _ := b.users.GetByTelegramID(tgID)
	if user != nil && user.Role == models.RoleExecutor {
		sub, _ := b.users.GetSubscription(user.ID)
		if sub == nil || sub.Type != models.SubPro {
			f.ExcludeUrgent = true
		}
	}
	b.showOpenTasks(cq.Message.Chat.ID, f, tgID)
}

func (b *Bot) cbFilterCat(cq *tgbotapi.CallbackQuery, cat string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	f := repo.TaskFilter{PageSize: 10}
	if cat != "all" {
		f.Category = cat
	}
	user, _ := b.users.GetByTelegramID(tgID)
	if user != nil && user.Role == models.RoleExecutor {
		sub, _ := b.users.GetSubscription(user.ID)
		if sub == nil || sub.Type != models.SubPro {
			f.ExcludeUrgent = true
		}
	}
	b.showOpenTasks(cq.Message.Chat.ID, f, tgID)
}

func (b *Bot) cbFilterBudget(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	f := repo.TaskFilter{PageSize: 10}

	parts := strings.SplitN(param, "-", 2)
	if len(parts) == 2 {
		min, _ := strconv.ParseInt(parts[0], 10, 64)
		f.BudgetMin = min
	} else {
		min, _ := strconv.ParseInt(param, 10, 64)
		if min >= 500000 {
			f.BudgetMin = min
		}
	}

	user, _ := b.users.GetByTelegramID(tgID)
	if user != nil && user.Role == models.RoleExecutor {
		sub, _ := b.users.GetSubscription(user.ID)
		if sub == nil || sub.Type != models.SubPro {
			f.ExcludeUrgent = true
		}
	}
	b.showOpenTasks(cq.Message.Chat.ID, f, tgID)
}

func (b *Bot) onFilterCity(m *tgbotapi.Message) {
	city := strings.TrimSpace(m.Text)
	tgID := m.From.ID
	fsm.Clear(tgID)
	loc := b.loc(tgID)
	b.send(m.Chat.ID, loc.T("exec_city_title", city))
	b.showTopExecutors(m.Chat.ID, "", city, 0, tgID)
}

func (b *Bot) onFilterTaskCity(m *tgbotapi.Message) {
	city := strings.TrimSpace(m.Text)
	tgID := m.From.ID
	fsm.Clear(tgID)
	loc := b.loc(tgID)

	f := repo.TaskFilter{City: city, PageSize: 10}
	user, _ := b.users.GetByTelegramID(tgID)
	if user != nil && user.Role == models.RoleExecutor {
		sub, _ := b.users.GetSubscription(user.ID)
		if sub == nil || sub.Type != models.SubPro {
			f.ExcludeUrgent = true
		}
	}
	b.send(m.Chat.ID, loc.T("task_city_title", city))
	b.showOpenTasks(m.Chat.ID, f, tgID)
}

func (b *Bot) cbFilterCity(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	if param == "task" {
		fsm.Set(tgID, fsm.StateFilterTaskCity)
		b.send(cq.Message.Chat.ID, loc.T("task_city_prompt"))
	}
}

func (b *Bot) showTopExecutors(chatID int64, cat string, city string, page int, tgID int64) {
	loc := b.loc(tgID)
	const pageSize = 8
	execs, _ := b.users.GetTopExecutors(cat, pageSize*(page+1))

	start := page * pageSize
	shown := 0
	for i, ep := range execs {
		if i < start {
			continue
		}
		if city != "" && !strings.EqualFold(ep.City, city) {
			continue
		}
		pro := ""
		if ep.IsPro {
			pro = " 💎"
		}
		portfolio := ""
		if ep.PortfolioLinks != "" {
			portfolio = "\n🔗 " + ep.PortfolioLinks
		}
		text := fmt.Sprintf(
			"👤 <b>%s</b>%s\n%s | %s | %d лет опыта\n⭐ %.1f (%d заказов)\n%s%s",
			ep.Name, pro,
			loc.CatLabel(string(ep.Category)), ep.City, ep.ExperienceYears,
			ep.Rating, ep.CompletedOrders,
			ifEmpty(ep.Description, "—"), portfolio,
		)
		b.send(chatID, text)
		shown++
	}
	if shown == 0 {
		b.send(chatID, loc.T("no_verified_exec"))
		return
	}

	b.sendKb(chatID, loc.T("filter_cat_menu"), loc.ExecutorFilterMenu())
	if len(execs) == pageSize*(page+1) {
		b.sendKb(chatID, loc.T("page_label", page+1), loc.ExecNextPage(page+1, cat))
	}
}

func (b *Bot) cbExecFilter(cq *tgbotapi.CallbackQuery, cat string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	if cat == "all" {
		cat = ""
	}
	if cat == "city" {
		fsm.Set(tgID, fsm.StateFilterCity)
		b.send(cq.Message.Chat.ID, loc.T("exec_city_prompt"))
		return
	}
	b.showTopExecutors(cq.Message.Chat.ID, cat, "", 0, tgID)
}

func (b *Bot) cbExecPage(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	// param = "page:cat"
	parts := strings.SplitN(param, ":", 2)
	if len(parts) != 2 {
		return
	}
	page, _ := strconv.Atoi(parts[0])
	cat := parts[1]
	if cat == "_" {
		cat = ""
	}
	b.showTopExecutors(cq.Message.Chat.ID, cat, "", page, cq.From.ID)
}

func (b *Bot) cbViewTask(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	taskID := parseInt64(param)

	task, _ := b.tasks.GetByID(taskID)
	if task == nil {
		b.send(cq.Message.Chat.ID, loc.T("task_not_found"))
		return
	}

	user, _ := b.users.GetByTelegramID(tgID)

	if task.IsUrgent && user != nil && user.Role == models.RoleExecutor {
		sub, _ := b.users.GetSubscription(user.ID)
		if sub == nil || sub.Type != models.SubPro {
			b.send(cq.Message.Chat.ID, loc.T("urgent_pro_only"))
			return
		}
	}

	b.send(cq.Message.Chat.ID, loc.FormatTask(task))

	if user != nil && user.Role == models.RoleClient && user.ID == task.ClientID {
		b.showClientTaskDetail(cq.Message.Chat.ID, task, tgID)
	} else {
		b.showExecutorTaskDetail(cq.Message.Chat.ID, task, user, tgID)
	}
}

func (b *Bot) showClientTaskDetail(chatID int64, task *models.Task, tgID int64) {
	loc := b.loc(tgID)
	resps, _ := b.resps.ListForTask(task.ID)
	if len(resps) == 0 {
		b.send(chatID, loc.T("no_responses"))
	} else {
		for _, r := range resps {
			b.sendKb(chatID, loc.FormatResponse(r), tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(loc.T("btn_view_profile"), fmt.Sprintf("view_resp:%d", r.ID)),
				),
			))
		}
	}
	if task.Status == models.TaskInProgress {
		b.sendKb(chatID, loc.T("task_in_progress_msg"), loc.CompleteTask(task.ID))
	}
}

func (b *Bot) showExecutorTaskDetail(chatID int64, task *models.Task, user *models.User, tgID int64) {
	loc := b.loc(tgID)
	canRespond := task.Status == models.TaskOpen && task.ResponseCount < task.MaxResponses
	if canRespond && user != nil {
		exists, _ := b.resps.ExistsForExecutor(task.ID, user.ID)
		if exists {
			canRespond = false
		}
	}
	b.sendKb(chatID, loc.T("choose_action"), loc.TaskActions(task.ID, canRespond))
}

func (b *Bot) cbViewResponse(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	respID := parseInt64(param)

	resp, _ := b.resps.GetByID(respID)
	if resp == nil {
		b.send(cq.Message.Chat.ID, loc.T("resp_not_found"))
		return
	}

	task, _ := b.tasks.GetByID(resp.TaskID)
	user, _ := b.users.GetByTelegramID(tgID)
	if task == nil || user == nil || user.ID != task.ClientID {
		b.send(cq.Message.Chat.ID, loc.T("no_access"))
		return
	}

	ep, _ := b.users.GetExecutorProfile(resp.ExecutorID)
	sum := loc.T("sum_label")
	text := fmt.Sprintf("👤 <b>%s</b>\n⭐ Рейтинг: %.1f", resp.ExecutorName, resp.Rating)
	if ep != nil {
		text += fmt.Sprintf("\n🎯 %s | %d лет опыта\n🏙 %s\n\n%s",
			loc.CatLabel(string(ep.Category)), ep.ExperienceYears, ep.City, ep.Description)
		if ep.PortfolioLinks != "" {
			text += "\n\n🔗 " + ep.PortfolioLinks
		}
	}
	text += fmt.Sprintf("\n\n💬 <b>Отклик:</b> %s", resp.Message)
	if resp.ProposedPrice != nil {
		text += fmt.Sprintf("\n💰 %d %s", *resp.ProposedPrice, sum)
	}

	b.sendKb(cq.Message.Chat.ID, text, loc.ResponseActions(respID))
}

func (b *Bot) cbWithdrawResp(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	respID := parseInt64(param)

	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		return
	}
	resp, _ := b.resps.GetByID(respID)
	if resp == nil || resp.ExecutorID != user.ID {
		b.send(cq.Message.Chat.ID, loc.T("resp_not_found"))
		return
	}
	if resp.Status != models.RespActive {
		b.send(cq.Message.Chat.ID, loc.T("resp_cant_withdraw"))
		return
	}

	b.resps.Withdraw(respID, user.ID)
	b.send(cq.Message.Chat.ID, loc.T("resp_withdrawn_msg"))
}

func (b *Bot) cbCancelTask(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	taskID := parseInt64(param)

	user, _ := b.users.GetByTelegramID(tgID)
	task, _ := b.tasks.GetByID(taskID)
	if user == nil || task == nil || task.ClientID != user.ID {
		b.send(cq.Message.Chat.ID, loc.T("no_access"))
		return
	}
	if task.Status != models.TaskOpen {
		b.send(cq.Message.Chat.ID, loc.T("task_cancel_forbidden"))
		return
	}

	b.tasks.UpdateStatus(taskID, string(models.TaskCancelled))
	b.send(cq.Message.Chat.ID, loc.T("task_cancelled_msg"))
}
