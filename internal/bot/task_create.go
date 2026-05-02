package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/models"
)

func (b *Bot) startCreateTask(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	user, _ := b.users.GetByTelegramID(m.From.ID)
	if user == nil || user.Role != models.RoleClient {
		b.send(m.Chat.ID, loc.T("only_clients"))
		return
	}
	fsm.ClearData(m.From.ID)
	fsm.Set(m.From.ID, fsm.StateTaskTitle)
	b.send(m.Chat.ID, loc.T("enter_task_title"))
}

func (b *Bot) onTaskTitle(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	title := strings.TrimSpace(m.Text)
	if title == "" || len([]rune(title)) > 50 {
		b.send(m.Chat.ID, loc.T("title_length"))
		return
	}
	fsm.SetVal(m.From.ID, "title", title)
	fsm.Set(m.From.ID, fsm.StateTaskDescription)
	b.send(m.Chat.ID, loc.T("enter_task_desc"))
}

func (b *Bot) onTaskDescription(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	desc := strings.TrimSpace(m.Text)
	if len([]rune(desc)) > 300 {
		b.send(m.Chat.ID, loc.T("desc_too_long"))
		return
	}
	fsm.SetVal(m.From.ID, "description", desc)
	fsm.Set(m.From.ID, fsm.StateTaskCategory)
	b.sendKb(m.Chat.ID, loc.T("choose_task_cat"), loc.TaskCategory())
}

func (b *Bot) cbTaskCategory(cq *tgbotapi.CallbackQuery, cat string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	if fsm.Get(tgID) != fsm.StateTaskCategory {
		return
	}
	fsm.SetVal(tgID, "category", cat)
	fsm.Set(tgID, fsm.StateTaskBudgetType)
	loc := b.loc(tgID)
	b.sendKb(cq.Message.Chat.ID, loc.T("choose_budget"), loc.BudgetType())
}

func (b *Bot) cbBudget(cq *tgbotapi.CallbackQuery, budgetType string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	if fsm.Get(tgID) != fsm.StateTaskBudgetType {
		return
	}
	fsm.SetVal(tgID, "budget_type", budgetType)
	loc := b.loc(tgID)

	switch budgetType {
	case "fixed":
		fsm.Set(tgID, fsm.StateTaskBudgetFrom)
		b.send(cq.Message.Chat.ID, loc.T("enter_amount"))
	case "range":
		fsm.Set(tgID, fsm.StateTaskBudgetFrom)
		b.send(cq.Message.Chat.ID, loc.T("enter_min_amount"))
	case "negotiable":
		fsm.Set(tgID, fsm.StateTaskDeadline)
		b.send(cq.Message.Chat.ID, loc.T("enter_deadline"))
	}
}

func (b *Bot) onTaskBudgetFrom(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	amount, err := strconv.ParseInt(strings.TrimSpace(m.Text), 10, 64)
	if err != nil || amount < 0 {
		b.send(m.Chat.ID, loc.T("invalid_amount"))
		return
	}
	tgID := m.From.ID
	fsm.SetVal(tgID, "budget_from", amount)

	if fsm.GetStr(tgID, "budget_type") == "range" {
		fsm.Set(tgID, fsm.StateTaskBudgetTo)
		b.send(m.Chat.ID, loc.T("enter_max_amount"))
	} else {
		fsm.Set(tgID, fsm.StateTaskDeadline)
		b.send(m.Chat.ID, loc.T("enter_deadline"))
	}
}

func (b *Bot) onTaskBudgetTo(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	amount, err := strconv.ParseInt(strings.TrimSpace(m.Text), 10, 64)
	if err != nil || amount < 0 {
		b.send(m.Chat.ID, loc.T("invalid_amount"))
		return
	}
	fsm.SetVal(m.From.ID, "budget_to", amount)
	fsm.Set(m.From.ID, fsm.StateTaskDeadline)
	b.send(m.Chat.ID, loc.T("enter_deadline"))
}

func (b *Bot) onTaskDeadline(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	deadline, err := time.Parse("02.01.2006", strings.TrimSpace(m.Text))
	if err != nil {
		b.send(m.Chat.ID, loc.T("invalid_date"))
		return
	}
	fsm.SetVal(m.From.ID, "deadline", deadline.Format(time.RFC3339))
	fsm.Set(m.From.ID, fsm.StateTaskRefs)
	b.sendKb(m.Chat.ID, loc.T("enter_refs"), loc.SkipRefs())
}

func (b *Bot) onTaskRefs(m *tgbotapi.Message) {
	fsm.SetVal(m.From.ID, "refs", strings.TrimSpace(m.Text))
	fsm.Set(m.From.ID, fsm.StateTaskUrgent)
	loc := b.loc(m.From.ID)
	b.sendKb(m.Chat.ID, loc.T("is_urgent"), loc.UrgentChoice())
}

func (b *Bot) cbUrgent(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	if fsm.Get(tgID) != fsm.StateTaskUrgent {
		return
	}
	fsm.SetVal(tgID, "urgent", param == "yes")
	b.showTaskConfirm(cq.Message.Chat.ID, tgID)
}

func (b *Bot) cbSkip(cq *tgbotapi.CallbackQuery, what string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)

	switch what {
	case "refs":
		if fsm.Get(tgID) != fsm.StateTaskRefs {
			return
		}
		fsm.SetVal(tgID, "refs", "")
		fsm.Set(tgID, fsm.StateTaskUrgent)
		b.sendKb(cq.Message.Chat.ID, loc.T("is_urgent"), loc.UrgentChoice())
	case "price":
		if fsm.Get(tgID) != fsm.StateRespondPrice {
			return
		}
		b.submitResponse(cq.Message.Chat.ID, tgID, nil)
	case "comment":
		if fsm.Get(tgID) != fsm.StateReviewComment {
			return
		}
		b.submitReview(cq.Message.Chat.ID, tgID, "")
	}
}

func (b *Bot) showTaskConfirm(chatID, tgID int64) {
	loc := b.loc(tgID)

	deadlineStr := fsm.GetStr(tgID, "deadline")
	deadline := ""
	if t, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
		deadline = t.Format("02.01.2006")
	}

	budgetDisplay := fsm.GetStr(tgID, "budget_type")
	switch budgetDisplay {
	case "fixed":
		budgetDisplay = loc.T("btn_budget_fixed")
	case "range":
		budgetDisplay = loc.T("btn_budget_range")
	case "negotiable":
		budgetDisplay = loc.T("budget_negotiable")
	}
	if bf := fsm.GetInt64(tgID, "budget_from"); bf > 0 {
		sum := loc.T("sum_label")
		budgetDisplay = fmt.Sprintf("%d %s", bf, sum)
		if bt := fsm.GetInt64(tgID, "budget_to"); bt > 0 {
			budgetDisplay = fmt.Sprintf("%d – %d %s", bf, bt, sum)
		}
	}

	urgent := loc.T("urgent_no_label")
	if fsm.GetBool(tgID, "urgent") {
		urgent = loc.T("urgent_yes")
	}

	text := loc.T("task_confirm_header") + "\n\n" +
		loc.T("confirm_name") + " " + fsm.GetStr(tgID, "title") + "\n" +
		loc.T("confirm_desc") + " " + fsm.GetStr(tgID, "description") + "\n" +
		loc.T("confirm_cat") + " " + loc.CatLabel(fsm.GetStr(tgID, "category")) + "\n" +
		loc.T("confirm_budget") + " " + budgetDisplay + "\n" +
		loc.T("confirm_deadline") + " " + deadline + "\n" +
		loc.T("confirm_refs") + " " + ifEmpty(fsm.GetStr(tgID, "refs"), "—") + "\n" +
		loc.T("confirm_urgent") + " " + urgent

	b.sendKb(chatID, text, loc.TaskConfirm())
}

func (b *Bot) cbTaskConfirm(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)

	switch param {
	case "cancel":
		fsm.Clear(tgID)
		if user, _ := b.users.GetByTelegramID(tgID); user != nil {
			b.showMainMenu(cq.Message.Chat.ID, user)
		}
	case "edit":
		fsm.ClearData(tgID)
		fsm.Set(tgID, fsm.StateTaskTitle)
		b.send(cq.Message.Chat.ID, loc.T("enter_task_title"))
	case "yes":
		b.publishTask(cq.Message.Chat.ID, tgID)
	}
}

func (b *Bot) publishTask(chatID, tgID int64) {
	loc := b.loc(tgID)
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		return
	}

	deadline, _ := time.Parse(time.RFC3339, fsm.GetStr(tgID, "deadline"))

	var budgetFrom, budgetTo *int64
	if bf := fsm.GetInt64(tgID, "budget_from"); bf > 0 {
		budgetFrom = &bf
	}
	if bt := fsm.GetInt64(tgID, "budget_to"); bt > 0 {
		budgetTo = &bt
	}

	taskID, err := b.tasks.Create(&models.Task{
		ClientID:    user.ID,
		Title:       fsm.GetStr(tgID, "title"),
		Description: fsm.GetStr(tgID, "description"),
		Category:    models.Category(fsm.GetStr(tgID, "category")),
		BudgetType:  models.BudgetType(fsm.GetStr(tgID, "budget_type")),
		BudgetFrom:  budgetFrom,
		BudgetTo:    budgetTo,
		Deadline:    deadline,
		Refs:        fsm.GetStr(tgID, "refs"),
		IsUrgent:    fsm.GetBool(tgID, "urgent"),
	})
	if err != nil {
		b.send(chatID, loc.T("task_error"))
		return
	}

	fsm.Clear(tgID)
	b.sendKb(chatID, loc.T("task_published", taskID), loc.ClientMenu())

	go b.notifyExecutorsNewTask(taskID)
}

func (b *Bot) notifyExecutorsNewTask(taskID int64) {
	task, _ := b.tasks.GetByID(taskID)
	if task == nil {
		return
	}
	executors, _ := b.users.GetExecutorsByCategory(string(task.Category), 50)
	for _, u := range executors {
		if task.IsUrgent {
			sub, _ := b.users.GetSubscription(u.ID)
			if sub == nil || sub.Type != models.SubPro {
				continue
			}
		}
		loc := b.loc(u.TelegramID)
		urgentSuffix := ""
		if task.IsUrgent {
			urgentSuffix = " ⚡"
		}
		b.send(u.TelegramID, loc.T("new_task_notify",
			urgentSuffix, loc.CatLabel(string(task.Category)), task.Title, loc.BudgetText(task),
		))
	}
}
