package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/locales"
	"github.com/birjasmm/bot/internal/models"
)

// ─── Role selection ────────────────────────────────────────────────────────────

func (b *Bot) cbRole(cq *tgbotapi.CallbackQuery, role string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)

	lang := fsm.GetStr(tgID, "lang")
	if lang == "" {
		lang = "ru"
	}

	_, err := b.users.Create(tgID, cq.From.UserName, role, lang)
	if err != nil {
		b.send(cq.Message.Chat.ID, loc.T("error_start"))
		return
	}

	if role == string(models.RoleClient) {
		fsm.Set(tgID, fsm.StateClientName)
		b.send(cq.Message.Chat.ID, loc.T("enter_name"))
	} else {
		fsm.Set(tgID, fsm.StateExecName)
		b.send(cq.Message.Chat.ID, loc.T("enter_name"))
	}
}

// ─── Client registration ───────────────────────────────────────────────────────

func (b *Bot) onClientName(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	name := strings.TrimSpace(m.Text)
	if name == "" {
		b.send(m.Chat.ID, loc.T("name_empty"))
		return
	}
	fsm.SetVal(m.From.ID, "name", name)
	fsm.Set(m.From.ID, fsm.StateClientPhone)
	b.sendKb(m.Chat.ID, loc.T("share_phone"), loc.RequestPhone())
}

func (b *Bot) onClientPhone(m *tgbotapi.Message) {
	tgID := m.From.ID
	phone := ""
	if m.Contact != nil {
		phone = m.Contact.PhoneNumber
	} else {
		phone = strings.TrimSpace(m.Text)
	}
	fsm.SetVal(tgID, "phone", phone)
	fsm.Set(tgID, fsm.StateClientBusiness)
	b.sendKb(m.Chat.ID, b.loc(tgID).T("enter_business"), locales.RemoveKeyboard())
}

func (b *Bot) onClientBusiness(m *tgbotapi.Message) {
	fsm.SetVal(m.From.ID, "business", strings.TrimSpace(m.Text))
	fsm.Set(m.From.ID, fsm.StateClientCity)
	b.send(m.Chat.ID, b.loc(m.From.ID).T("enter_city"))
}

func (b *Bot) onClientCity(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		b.send(m.Chat.ID, loc.T("error_start"))
		return
	}

	b.users.SetPhone(user.ID, fsm.GetStr(tgID, "phone"))
	b.users.CreateClientProfile(&models.ClientProfile{
		UserID:       user.ID,
		Name:         fsm.GetStr(tgID, "name"),
		BusinessName: fsm.GetStr(tgID, "business"),
		City:         strings.TrimSpace(m.Text),
	})
	fsm.Clear(tgID)

	b.sendKb(m.Chat.ID, loc.T("reg_done_client"), loc.ClientMenu())
}

// ─── Executor registration ─────────────────────────────────────────────────────

func (b *Bot) onExecName(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	name := strings.TrimSpace(m.Text)
	if name == "" {
		b.send(m.Chat.ID, loc.T("name_empty"))
		return
	}
	fsm.SetVal(m.From.ID, "name", name)
	fsm.Set(m.From.ID, fsm.StateExecPhone)
	b.sendKb(m.Chat.ID, loc.T("share_phone"), loc.RequestPhone())
}

func (b *Bot) onExecPhone(m *tgbotapi.Message) {
	tgID := m.From.ID
	phone := ""
	if m.Contact != nil {
		phone = m.Contact.PhoneNumber
	} else {
		phone = strings.TrimSpace(m.Text)
	}
	fsm.SetVal(tgID, "phone", phone)
	fsm.Set(tgID, fsm.StateExecCity)
	b.sendKb(m.Chat.ID, b.loc(tgID).T("enter_city"), locales.RemoveKeyboard())
}

func (b *Bot) onExecCity(m *tgbotapi.Message) {
	fsm.SetVal(m.From.ID, "city", strings.TrimSpace(m.Text))
	fsm.Set(m.From.ID, fsm.StateExecCategory)
	loc := b.loc(m.From.ID)
	b.sendKb(m.Chat.ID, loc.T("choose_category"), loc.ExecutorCategory())
}

func (b *Bot) cbExecCategory(cq *tgbotapi.CallbackQuery, cat string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	if fsm.Get(tgID) != fsm.StateExecCategory {
		return
	}
	fsm.SetVal(tgID, "category", cat)
	fsm.Set(tgID, fsm.StateExecExperience)
	b.send(cq.Message.Chat.ID, b.loc(tgID).T("enter_experience"))
}

func (b *Bot) onExecExperience(m *tgbotapi.Message) {
	loc := b.loc(m.From.ID)
	exp, err := strconv.Atoi(strings.TrimSpace(m.Text))
	if err != nil || exp < 0 {
		b.send(m.Chat.ID, loc.T("invalid_experience"))
		return
	}
	fsm.SetVal(m.From.ID, "experience", exp)
	fsm.Set(m.From.ID, fsm.StateExecPortfolio)
	b.send(m.Chat.ID, loc.T("enter_portfolio"))
}

func (b *Bot) onExecPortfolio(m *tgbotapi.Message) {
	fsm.SetVal(m.From.ID, "portfolio", strings.TrimSpace(m.Text))
	fsm.Set(m.From.ID, fsm.StateExecDescription)
	b.send(m.Chat.ID, b.loc(m.From.ID).T("enter_description"))
}

func (b *Bot) onExecDescription(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		b.send(m.Chat.ID, loc.T("error_start"))
		return
	}

	b.users.SetPhone(user.ID, fsm.GetStr(tgID, "phone"))
	b.users.CreateExecutorProfile(&models.ExecutorProfile{
		UserID:          user.ID,
		Name:            fsm.GetStr(tgID, "name"),
		City:            fsm.GetStr(tgID, "city"),
		Category:        models.Category(fsm.GetStr(tgID, "category")),
		ExperienceYears: fsm.GetInt(tgID, "experience"),
		PortfolioLinks:  fsm.GetStr(tgID, "portfolio"),
		Description:     strings.TrimSpace(m.Text),
	})
	fsm.Clear(tgID)

	b.sendKb(m.Chat.ID, loc.T("reg_done_exec"), loc.VerifOffer())
}
