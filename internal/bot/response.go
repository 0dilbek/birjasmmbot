package bot

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/models"
)

func (b *Bot) cbRespond(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	taskID := parseInt64(param)

	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil || user.Role != models.RoleExecutor {
		b.send(cq.Message.Chat.ID, loc.T("only_executors"))
		return
	}

	task, _ := b.tasks.GetByID(taskID)
	if task == nil || task.Status != models.TaskOpen {
		b.send(cq.Message.Chat.ID, loc.T("task_unavailable"))
		return
	}
	if task.ResponseCount >= task.MaxResponses {
		b.send(cq.Message.Chat.ID, loc.T("resp_limit_exceeded"))
		return
	}

	exists, _ := b.resps.ExistsForExecutor(taskID, user.ID)
	if exists {
		b.send(cq.Message.Chat.ID, loc.T("already_responded"))
		return
	}

	limits, _ := b.users.GetOrCreateLimits(user.ID, b.cfg.FreeResponses)
	sub, _ := b.users.GetSubscription(user.ID)
	if (limits == nil || limits.FreeResponsesLeft <= 0) && sub == nil {
		b.send(cq.Message.Chat.ID, loc.T("free_limit_exceeded"))
		return
	}

	fsm.SetVal(tgID, "respond_task_id", taskID)
	fsm.Set(tgID, fsm.StateRespondMsg)

	hint := ""
	if sub == nil && limits != nil {
		hint = loc.T("free_resps_hint", limits.FreeResponsesLeft)
	}
	b.send(cq.Message.Chat.ID, loc.T("write_response")+hint)
}

func (b *Bot) onRespondMsg(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	msg := strings.TrimSpace(m.Text)
	if msg == "" {
		b.send(m.Chat.ID, loc.T("msg_empty"))
		return
	}
	fsm.SetVal(tgID, "respond_msg", msg)
	fsm.Set(tgID, fsm.StateRespondPrice)
	b.sendKb(m.Chat.ID, loc.T("enter_price"), loc.SkipPrice())
}

func (b *Bot) onRespondPrice(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	price, err := strconv.ParseInt(strings.TrimSpace(m.Text), 10, 64)
	if err != nil || price < 0 {
		b.send(m.Chat.ID, loc.T("invalid_price"))
		return
	}
	b.submitResponse(m.Chat.ID, m.From.ID, &price)
}

func (b *Bot) submitResponse(chatID, tgID int64, price *int64) {
	loc := b.loc(tgID)
	taskID := fsm.GetInt64(tgID, "respond_task_id")
	msg := fsm.GetStr(tgID, "respond_msg")

	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		fsm.Clear(tgID)
		return
	}

	if _, err := b.resps.Create(taskID, user.ID, msg, price); err != nil {
		b.send(chatID, loc.T("resp_error"))
		log.Printf("Create response error: %v", err)
		return
	}

	sub, _ := b.users.GetSubscription(user.ID)
	if sub == nil {
		b.users.DecrementLimit(user.ID)
	}

	b.notifyClientNewResponse(taskID, user.ID)
	fsm.Clear(tgID)
	b.sendKb(chatID, loc.T("resp_sent"), loc.ExecutorMenu())
}

func (b *Bot) notifyClientNewResponse(taskID, executorUserID int64) {
	task, _ := b.tasks.GetByID(taskID)
	if task == nil {
		return
	}
	client, _ := b.users.GetByID(task.ClientID)
	if client == nil {
		return
	}
	loc := b.loc(client.TelegramID)
	ep, _ := b.users.GetExecutorProfile(executorUserID)
	name := ""
	if ep != nil {
		name = ep.Name
	}
	b.send(client.TelegramID, loc.T("new_resp_notify", task.Title, name))
}
