package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/locales"
	"github.com/birjasmm/bot/internal/models"
)

func (b *Bot) cbSelectExec(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	respID := parseInt64(param)

	resp, _ := b.resps.GetByID(respID)
	if resp == nil {
		b.send(cq.Message.Chat.ID, loc.T("exec_not_found"))
		return
	}
	task, _ := b.tasks.GetByID(resp.TaskID)
	user, _ := b.users.GetByTelegramID(tgID)
	if task == nil || user == nil || user.ID != task.ClientID {
		b.send(cq.Message.Chat.ID, loc.T("no_access"))
		return
	}

	b.resps.Accept(respID)
	b.resps.RejectOthers(task.ID, respID)
	b.resps.CreateAssignment(task.ID, resp.ExecutorID)
	b.tasks.UpdateStatus(task.ID, string(models.TaskInProgress))

	b.notifyExecutorSelected(task, user, resp.ExecutorID)
	go b.notifyRejectedExecutors(task.ID, resp.ExecutorID)
	b.send(cq.Message.Chat.ID, loc.T("exec_selected", resp.ExecutorName))
}

func (b *Bot) notifyRejectedExecutors(taskID, acceptedExecutorID int64) {
	rejected, _ := b.resps.ListRejectedForTask(taskID)
	task, _ := b.tasks.GetByID(taskID)
	title := fmt.Sprintf("#%d", taskID)
	if task != nil {
		title = task.Title
	}
	for _, r := range rejected {
		if r.ExecutorID == acceptedExecutorID {
			continue
		}
		u, _ := b.users.GetByID(r.ExecutorID)
		if u != nil {
			loc := b.loc(u.TelegramID)
			b.send(u.TelegramID, loc.T("rejected_notify", title))
		}
	}
}

func (b *Bot) notifyExecutorSelected(task *models.Task, client *models.User, executorID int64) {
	executor, _ := b.users.GetByID(executorID)
	if executor == nil {
		return
	}
	loc := b.loc(executor.TelegramID)
	cp, _ := b.users.GetClientProfile(client.ID)
	clientName := ""
	if cp != nil {
		clientName = cp.Name
	}
	b.send(executor.TelegramID, loc.T("selected_notify",
		task.Title, clientName, client.Username, client.Phone))
}

func (b *Bot) cbComplete(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	taskID := parseInt64(param)

	task, _ := b.tasks.GetByID(taskID)
	user, _ := b.users.GetByTelegramID(tgID)
	if task == nil || user == nil || user.ID != task.ClientID {
		b.send(cq.Message.Chat.ID, loc.T("no_access"))
		return
	}

	b.tasks.UpdateStatus(taskID, string(models.TaskCompleted))

	assign, _ := b.resps.GetAssignment(taskID)
	if assign == nil {
		b.send(cq.Message.Chat.ID, loc.T("task_complete_simple"))
		return
	}
	b.resps.CompleteAssignment(taskID)

	ep, _ := b.users.GetExecutorProfile(assign.ExecutorID)
	if ep == nil {
		b.send(cq.Message.Chat.ID, loc.T("task_complete_simple"))
		return
	}

	fsm.SetVal(tgID, "review_task_id", taskID)
	fsm.SetVal(tgID, "review_executor_id", assign.ExecutorID)
	fsm.Set(tgID, fsm.StateReviewRating)
	b.sendKb(cq.Message.Chat.ID,
		loc.T("task_complete_review", ep.Name),
		locales.RatingKeyboard(taskID))
}

func (b *Bot) cbRate(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)

	// param = "rating:taskID"
	parts := strings.SplitN(param, ":", 2)
	if len(parts) != 2 {
		return
	}
	rating, _ := strconv.Atoi(parts[0])
	taskID, _ := strconv.ParseInt(parts[1], 10, 64)

	fsm.SetVal(tgID, "review_rating", rating)
	fsm.SetVal(tgID, "review_task_id", taskID)
	fsm.Set(tgID, fsm.StateReviewComment)
	b.sendKb(cq.Message.Chat.ID, loc.T("write_review"), loc.SkipReviewComment())
}

func (b *Bot) onReviewComment(m *tgbotapi.Message) {
	b.submitReview(m.Chat.ID, m.From.ID, strings.TrimSpace(m.Text))
}

func (b *Bot) submitReview(chatID, tgID int64, comment string) {
	loc := b.loc(tgID)
	taskID := fsm.GetInt64(tgID, "review_task_id")
	executorID := fsm.GetInt64(tgID, "review_executor_id")
	rating := fsm.GetInt(tgID, "review_rating")

	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil || taskID == 0 {
		fsm.Clear(tgID)
		return
	}

	b.reviews.Create(taskID, user.ID, executorID, rating, comment)
	b.users.UpdateRating(executorID)
	fsm.Clear(tgID)

	b.sendKb(chatID, loc.T("thanks_review"), loc.ClientMenu())
}
