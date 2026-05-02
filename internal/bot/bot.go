package bot

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/config"
	"github.com/birjasmm/bot/internal/locales"
	"github.com/birjasmm/bot/internal/repo"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cfg      *config.Config
	users    *repo.UserRepo
	tasks    *repo.TaskRepo
	resps    *repo.ResponseRepo
	reviews  *repo.ReviewRepo
	payments *repo.PaymentRepo
	adminIDs map[int64]bool
}

func New(cfg *config.Config, db *sql.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized as @%s", api.Self.UserName)

	adminIDs := make(map[int64]bool)
	for _, id := range cfg.AdminIDs {
		adminIDs[id] = true
	}

	return &Bot{
		api:      api,
		cfg:      cfg,
		users:    repo.NewUserRepo(db),
		tasks:    repo.NewTaskRepo(db),
		resps:    repo.NewResponseRepo(db),
		reviews:  repo.NewReviewRepo(db),
		payments: repo.NewPaymentRepo(db),
		adminIDs: adminIDs,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for update := range b.api.GetUpdatesChan(u) {
		go b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
		}
	}()
	if u.CallbackQuery != nil {
		b.handleCallback(u.CallbackQuery)
		return
	}
	if u.Message != nil {
		b.handleMessage(u.Message)
	}
}

// loc returns a Localizer for the given telegram user ID.
// Checks FSM data first (pre-registration), then DB.
func (b *Bot) loc(tgID int64) *locales.Loc {
	if lang := fsm.GetStr(tgID, "lang"); lang != "" {
		return locales.New(lang)
	}
	lang, _ := b.users.GetLanguage(tgID)
	return locales.New(lang)
}

func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	b.api.Send(msg)
}

func (b *Bot) sendKb(chatID int64, text string, kb interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = kb
	b.api.Send(msg)
}

func (b *Bot) answer(id, text string) {
	b.api.Request(tgbotapi.NewCallback(id, text))
}

func (b *Bot) handleMessage(m *tgbotapi.Message) {
	tgID := m.From.ID
	b.users.LogAction(tgID, "msg", m.Text)

	if m.IsCommand() {
		switch m.Command() {
		case "start":
			b.cmdStart(m)
		case "admin":
			b.cmdAdmin(m)
		}
		return
	}

	switch fsm.Get(tgID) {
	case fsm.StateClientName:
		b.onClientName(m)
	case fsm.StateClientPhone:
		b.onClientPhone(m)
	case fsm.StateClientBusiness:
		b.onClientBusiness(m)
	case fsm.StateClientCity:
		b.onClientCity(m)
	case fsm.StateExecName:
		b.onExecName(m)
	case fsm.StateExecPhone:
		b.onExecPhone(m)
	case fsm.StateExecCity:
		b.onExecCity(m)
	case fsm.StateExecExperience:
		b.onExecExperience(m)
	case fsm.StateExecPortfolio:
		b.onExecPortfolio(m)
	case fsm.StateExecDescription:
		b.onExecDescription(m)
	case fsm.StateVerifVideo:
		b.onVerifVideo(m)
	case fsm.StateClientVerif:
		b.onClientVerifDoc(m)
	case fsm.StateTaskTitle:
		b.onTaskTitle(m)
	case fsm.StateTaskDescription:
		b.onTaskDescription(m)
	case fsm.StateTaskBudgetFrom:
		b.onTaskBudgetFrom(m)
	case fsm.StateTaskBudgetTo:
		b.onTaskBudgetTo(m)
	case fsm.StateTaskDeadline:
		b.onTaskDeadline(m)
	case fsm.StateTaskRefs:
		b.onTaskRefs(m)
	case fsm.StateRespondMsg:
		b.onRespondMsg(m)
	case fsm.StateRespondPrice:
		b.onRespondPrice(m)
	case fsm.StateReviewComment:
		b.onReviewComment(m)
	case fsm.StatePaymentReceipt:
		b.onPaymentReceipt(m)
	case fsm.StateFilterCity:
		b.onFilterCity(m)
	case fsm.StateFilterTaskCity:
		b.onFilterTaskCity(m)
	case fsm.StateAdminEditUser:
		b.onAdminEditValue(m)
	default:
		b.handleMenu(m)
	}
}

func (b *Bot) handleCallback(cq *tgbotapi.CallbackQuery) {
	b.users.LogAction(cq.From.ID, "cb", cq.Data)
	action, param := splitCallback(cq.Data)

	switch action {
	case "lang":
		b.cbLang(cq, param)
	case "role":
		b.cbRole(cq, param)
	case "exec_cat":
		b.cbExecCategory(cq, param)
	case "verif":
		b.cbVerif(cq, param)
	case "client_verif":
		b.cbClientVerif(cq, param)
	case "task_cat":
		b.cbTaskCategory(cq, param)
	case "budget":
		b.cbBudget(cq, param)
	case "urgent":
		b.cbUrgent(cq, param)
	case "task_confirm":
		b.cbTaskConfirm(cq, param)
	case "skip":
		b.cbSkip(cq, param)
	case "filter_cat":
		b.cbFilterCat(cq, param)
	case "filter_budget":
		b.cbFilterBudget(cq, param)
	case "filter_city":
		b.cbFilterCity(cq, param)
	case "exec_filter":
		b.cbExecFilter(cq, param)
	case "task_page":
		b.cbTaskPage(cq, param)
	case "exec_page":
		b.cbExecPage(cq, param)
	case "view_task":
		b.cbViewTask(cq, param)
	case "view_resp":
		b.cbViewResponse(cq, param)
	case "respond":
		b.cbRespond(cq, param)
	case "select_exec":
		b.cbSelectExec(cq, param)
	case "complete":
		b.cbComplete(cq, param)
	case "rate":
		b.cbRate(cq, param)
	case "sub":
		b.cbSub(cq, param)
	case "pay_ok":
		b.cbPaymentApprove(cq, param)
	case "pay_no":
		b.cbPaymentReject(cq, param)
	case "verif_ok":
		b.cbVerifApprove(cq, param)
	case "verif_no":
		b.cbVerifReject(cq, param)
	case "admin":
		b.cbAdmin(cq, param)
	case "admin_edit":
		b.cbAdminEdit(cq, param)
	case "admin_edit_field":
		b.cbAdminEditField(cq, param)
	case "block_user":
		b.cbBlockUser(cq, param)
	case "unblock_user":
		b.cbUnblockUser(cq, param)
	case "del_task":
		b.cbDelTask(cq, param)
	case "withdraw_resp":
		b.cbWithdrawResp(cq, param)
	case "cancel_task":
		b.cbCancelTask(cq, param)
	case "back":
		b.cbBack(cq, param)
	default:
		b.answer(cq.ID, "")
	}
}

// cbLang handles language selection / change.
func (b *Bot) cbLang(cq *tgbotapi.CallbackQuery, lang string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID

	user, _ := b.users.GetByTelegramID(tgID)
	if user != nil {
		b.users.SetLanguage(user.ID, lang)
	} else {
		fsm.SetVal(tgID, "lang", lang)
	}

	loc := locales.New(lang)

	if user == nil {
		b.sendKb(cq.Message.Chat.ID, loc.T("welcome_new"), loc.RoleChoice())
		return
	}

	changedKey := "lang_changed_ru"
	if lang == "uz" {
		changedKey = "lang_changed_uz"
	}
	b.send(cq.Message.Chat.ID, loc.T(changedKey))
	b.showMainMenuLoc(cq.Message.Chat.ID, user, loc)
}
