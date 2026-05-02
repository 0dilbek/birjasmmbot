package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/birjasmm/bot/internal/bot/fsm"
	"github.com/birjasmm/bot/internal/bot/keyboards"
	"github.com/birjasmm/bot/internal/models"
)

const (
	paymentCard = "8600 1234 5678 9000"
	priceBasic  = 50_000
	pricePro    = 100_000
)

func subPrice(subType string) int {
	if subType == string(models.SubPro) {
		return pricePro
	}
	return priceBasic
}

func (b *Bot) startPaymentFlow(chatID, tgID int64, subType string) {
	loc := b.loc(tgID)
	price := subPrice(subType)
	fsm.SetVal(tgID, "pay_type", subType)
	fsm.SetVal(tgID, "pay_amount", price)
	fsm.Set(tgID, fsm.StatePaymentReceipt)

	b.send(chatID, loc.T("payment_prompt", strings.ToUpper(subType), price, paymentCard))
}

func (b *Bot) onPaymentReceipt(m *tgbotapi.Message) {
	tgID := m.From.ID
	loc := b.loc(tgID)
	user, _ := b.users.GetByTelegramID(tgID)
	if user == nil {
		fsm.Clear(tgID)
		return
	}

	subType := fsm.GetStr(tgID, "pay_type")
	amount := fsm.GetInt(tgID, "pay_amount")
	if subType == "" {
		fsm.Clear(tgID)
		return
	}

	var receiptType, fileID, receiptText string
	if m.Photo != nil && len(m.Photo) > 0 {
		receiptType = "photo"
		fileID = m.Photo[len(m.Photo)-1].FileID
	} else if m.Document != nil {
		receiptType = "document"
		fileID = m.Document.FileID
	} else if m.Text != "" {
		receiptType = "text"
		receiptText = m.Text
	} else {
		b.send(m.Chat.ID, loc.T("payment_send_receipt"))
		return
	}

	paymentID, err := b.payments.Create(user.ID, subType, amount, receiptType, fileID, receiptText)
	if err != nil {
		b.send(m.Chat.ID, loc.T("payment_error"))
		log.Printf("payment create error: %v", err)
		fsm.Clear(tgID)
		return
	}

	fsm.Clear(tgID)
	b.notifyAdminsPayment(paymentID, user, subType, amount, receiptType, fileID, receiptText)
	b.sendKb(m.Chat.ID, loc.T("payment_saved"), loc.ExecutorMenu())
}

func (b *Bot) notifyAdminsPayment(paymentID int64, user *models.User, subType string, amount int, receiptType, fileID, receiptText string) {
	summary := fmt.Sprintf(
		"💳 <b>Платёж #%d</b>\n👤 @%s (TG: %d)\n💡 Тариф: <b>%s</b>\n💰 Сумма: %d сум",
		paymentID, user.Username, user.TelegramID, strings.ToUpper(subType), amount,
	)
	kb := keyboards.PaymentApproval(paymentID)

	for adminTgID := range b.adminIDs {
		msgID := b.sendGetID(adminTgID, summary, &kb)
		if msgID > 0 {
			b.payments.SaveAdminMsg(paymentID, adminTgID, adminTgID, msgID)
		}

		switch receiptType {
		case "photo":
			photo := tgbotapi.NewPhoto(adminTgID, tgbotapi.FileID(fileID))
			photo.Caption = "📎 Чек (фото)"
			b.api.Send(photo)
		case "document":
			doc := tgbotapi.NewDocument(adminTgID, tgbotapi.FileID(fileID))
			doc.Caption = "📎 Чек (файл)"
			b.api.Send(doc)
		case "text":
			b.send(adminTgID, "📎 Текст чека:\n"+receiptText)
		}
	}
}

func (b *Bot) cbPaymentApprove(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	if !b.isAdmin(cq.From.ID) {
		return
	}
	b.processPaymentReview(cq, parseInt64(param), "approved")
}

func (b *Bot) cbPaymentReject(cq *tgbotapi.CallbackQuery, param string) {
	b.answer(cq.ID, "")
	if !b.isAdmin(cq.From.ID) {
		return
	}
	b.processPaymentReview(cq, parseInt64(param), "rejected")
}

func (b *Bot) processPaymentReview(cq *tgbotapi.CallbackQuery, paymentID int64, status string) {
	payment, _ := b.payments.GetByID(paymentID)
	if payment == nil {
		b.send(cq.Message.Chat.ID, "Платёж не найден.")
		return
	}
	if payment.Status != models.PaymentPending {
		b.answer(cq.ID, "Этот платёж уже обработан.")
		return
	}

	b.payments.Review(paymentID, status, cq.From.ID)

	reviewer := "@" + cq.From.UserName
	var statusLine string
	if status == "approved" {
		statusLine = "✅ Подтверждён by " + reviewer
		payer, _ := b.users.GetByID(payment.UserID)
		if payer != nil {
			b.users.CreateSubscription(payment.UserID, payment.SubType, 30)
			if payment.SubType == string(models.SubPro) {
				b.users.SetExecutorPro(payment.UserID, true)
			}
			payerLoc := b.loc(payer.TelegramID)
			b.sendKb(payer.TelegramID,
				payerLoc.T("payment_approved", strings.ToUpper(payment.SubType)),
				payerLoc.ExecutorMenu())
		}
	} else {
		statusLine = "❌ Отклонён by " + reviewer
		payer, _ := b.users.GetByID(payment.UserID)
		if payer != nil {
			payerLoc := b.loc(payer.TelegramID)
			b.send(payer.TelegramID, payerLoc.T("payment_rejected"))
		}
	}

	payer, _ := b.users.GetByID(payment.UserID)
	username := fmt.Sprintf("ID:%d", payment.UserID)
	if payer != nil {
		username = "@" + payer.Username
	}
	updatedText := fmt.Sprintf(
		"💳 <b>Платёж #%d</b>\n👤 %s (TG: %d)\n💡 Тариф: <b>%s</b>\n💰 Сумма: %d сум\n\n%s",
		paymentID, username, payment.UserID, strings.ToUpper(payment.SubType), payment.Amount, statusLine,
	)

	adminMsgs, _ := b.payments.GetAdminMsgs(paymentID)
	for _, am := range adminMsgs {
		edit := tgbotapi.NewEditMessageText(am.ChatID, am.MessageID, updatedText)
		edit.ParseMode = "HTML"
		b.api.Send(edit)
	}
}

func (b *Bot) sendGetID(chatID int64, text string, kb interface{}) int {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if kb != nil {
		msg.ReplyMarkup = kb
	}
	sent, err := b.api.Send(msg)
	if err != nil {
		log.Printf("sendGetID error: %v", err)
		return 0
	}
	return sent.MessageID
}
