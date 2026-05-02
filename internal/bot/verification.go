package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/bot/keyboards"
	"github.com/birjasmm/bot/internal/models"
)

// ── Executor verification ────────────────────────────────────────────────────

func (b *Bot) cbVerif(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)

	if param == "skip" {
		user, _ := b.users.GetByTelegramID(tgID)
		if user != nil {
			b.showMainMenu(cq.Message.Chat.ID, user)
		}
		return
	}

	if param == "start" {
		code := fmt.Sprintf("%04d", tgID%10000)
		fsm.Set(tgID, fsm.StateVerifVideo)
		b.sendKb(cq.Message.Chat.ID,
			loc.T("verif_code_prompt", code),
			loc.VerifVideo())
	}
}

func (b *Bot) onVerifVideo(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	if m.Video == nil && m.VideoNote == nil {
		b.send(m.Chat.ID, loc.T("verif_send_video"))
		return
	}

	fileID := ""
	if m.Video != nil {
		fileID = m.Video.FileID
	} else {
		fileID = m.VideoNote.FileID
	}

	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		return
	}

	verifID, _ := b.users.CreateVerification(user.ID, fileID)
	b.users.SetVerifStatus(user.ID, string(models.VerifPending))
	fsm.Clear(tgID)

	b.sendKb(m.Chat.ID, loc.T("verif_sent"), loc.ExecutorMenu())

	adminKb := keyboards.VerifReview(verifID)
	for adminID := range b.adminIDs {
		msg := tgbotapi.NewMessage(adminID, fmt.Sprintf("🔐 <b>Новая верификация (Исполнитель)</b>\nОт: @%s\nUserID: %d", m.From.UserName, user.ID))
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = adminKb
		sent, err := b.api.Send(msg)
		if err == nil {
			b.users.SaveVerificationAdminMsg(verifID, adminID, sent.Chat.ID, sent.MessageID)
		}
	}
}

func (b *Bot) cbVerifApprove(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")

	verifID := parseInt64(param)
	verif, _ := b.users.GetVerificationByID(verifID)
	if verif == nil || verif.Status != "pending" {
		b.send(cq.Message.Chat.ID, "Заявка уже обработана или не найдена.")
		return
	}

	b.users.ReviewVerification(verifID, "approved")
	b.users.SetVerifStatus(verif.UserID, string(models.VerifVerified))

	u, _ := b.users.GetByID(verif.UserID)
	if u != nil {
		userLoc := b.loc(u.TelegramID)
		if u.Role == models.RoleExecutor {
			b.users.SetExecutorVerified(verif.UserID)
			b.send(u.TelegramID, userLoc.T("verif_approved_exec"))
		} else {
			b.users.VerifyClient(verif.UserID)
			b.sendKb(u.TelegramID, userLoc.T("verif_approved_client"), userLoc.ClientMenu())
		}
	}
	
	// Broadcast update
	msgs, _ := b.users.GetVerificationAdminMsgs(verifID)
	for _, m := range msgs {
		edit := tgbotapi.NewEditMessageText(m.ChatID, m.MessageID, 
			fmt.Sprintf("✅ Верификация <b>ID:%d</b> одобрена администратором @%s", verif.UserID, cq.From.UserName))
		edit.ParseMode = "HTML"
		b.api.Send(edit)
	}
}

func (b *Bot) cbVerifReject(cq *tgbotapi.CallbackQuery, param string) {
	if !b.isAdmin(cq.From.ID) {
		b.answer(cq.ID, "Нет доступа")
		return
	}
	b.answer(cq.ID, "")

	verifID := parseInt64(param)
	verif, _ := b.users.GetVerificationByID(verifID)
	if verif == nil || verif.Status != "pending" {
		return
	}

	b.users.ReviewVerification(verifID, "rejected")
	b.users.SetVerifStatus(verif.UserID, string(models.VerifNone))

	if u, _ := b.users.GetByID(verif.UserID); u != nil {
		b.send(u.TelegramID, b.loc(u.TelegramID).T("verif_rejected_user"))
	}
	
	// Broadcast update
	msgs, _ := b.users.GetVerificationAdminMsgs(verifID)
	for _, m := range msgs {
		edit := tgbotapi.NewEditMessageText(m.ChatID, m.MessageID, 
			fmt.Sprintf("❌ Верификация <b>ID:%d</b> отклонена администратором @%s", verif.UserID, cq.From.UserName))
		edit.ParseMode = "HTML"
		b.api.Send(edit)
	}
}

// ── Client verification ──────────────────────────────────────────────────────

func (b *Bot) startClientVerif(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		return
	}
	if user.VerificationStatus == models.VerifVerified {
		b.send(m.Chat.ID, loc.T("already_verified"))
		return
	}
	if user.VerificationStatus == models.VerifPending {
		b.send(m.Chat.ID, loc.T("verif_pending_status"))
		return
	}
	b.sendKb(m.Chat.ID, loc.T("verif_client_offer"), loc.ClientVerifOffer())
}

func (b *Bot) cbClientVerif(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	tgID := cq.From.ID
	loc := b.loc(tgID)
	if param == "cancel" {
		user, _ := b.users.GetByTelegramID(tgID)
		if user != nil {
			b.showMainMenu(cq.Message.Chat.ID, user)
		}
		return
	}
	if param == "start" {
		fsm.Set(tgID, fsm.StateClientVerif)
		b.send(cq.Message.Chat.ID, loc.T("verif_client_doc_prompt"))
	}
}

func (b *Bot) onClientVerifDoc(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		fsm.Clear(tgID)
		return
	}

	var fileID string
	if m.Photo != nil && len(m.Photo) > 0 {
		fileID = m.Photo[len(m.Photo)-1].FileID
	} else if m.Document != nil {
		fileID = m.Document.FileID
	} else {
		b.send(m.Chat.ID, loc.T("verif_client_send_doc"))
		return
	}

	verifID, _ := b.users.CreateVerification(user.ID, fileID)
	b.users.SetVerifStatus(user.ID, string(models.VerifPending))
	fsm.Clear(tgID)

	b.sendKb(m.Chat.ID, loc.T("verif_client_sent"), loc.ClientMenu())

	adminKb := keyboards.VerifReview(verifID)
	for adminID := range b.adminIDs {
		caption := fmt.Sprintf("📄 <b>Верификация заказчика</b>\nОт: @%s\nUserID: %d", m.From.UserName, user.ID)
		
		var sent tgbotapi.Message
		var err error
		if m.Photo != nil {
			photo := tgbotapi.NewPhoto(adminID, tgbotapi.FileID(fileID))
			photo.Caption = caption
			photo.ParseMode = "HTML"
			photo.ReplyMarkup = adminKb
			sent, err = b.api.Send(photo)
		} else {
			doc := tgbotapi.NewDocument(adminID, tgbotapi.FileID(fileID))
			doc.Caption = caption
			doc.ParseMode = "HTML"
			doc.ReplyMarkup = adminKb
			sent, err = b.api.Send(doc)
		}
		
		if err == nil {
			b.users.SaveVerificationAdminMsg(verifID, adminID, sent.Chat.ID, sent.MessageID)
		}
	}
}
