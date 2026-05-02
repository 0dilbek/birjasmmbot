package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Admin-facing keyboards (always shown in Russian to hardcoded admin IDs).

func AdminMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Пользователи", "admin:users"),
			tgbotapi.NewInlineKeyboardButtonData("📋 Задачи", "admin:tasks"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔐 Верификации", "admin:verifs"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "admin:stats"),
		),
	)
}

func VerifReview(verifID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Одобрить", fmt.Sprintf("verif_ok:%d", verifID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", fmt.Sprintf("verif_no:%d", verifID)),
		),
	)
}

func UserActions(userID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫 Заблокировать", fmt.Sprintf("block_user:%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("✅ Разблокировать", fmt.Sprintf("unblock_user:%d", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать", fmt.Sprintf("admin_edit:%d", userID)),
		),
	)
}

func UserEditFields(userID int64, role string) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📛 Имя", fmt.Sprintf("admin_edit_field:%d:name", userID)),
			tgbotapi.NewInlineKeyboardButtonData("🏙 Город", fmt.Sprintf("admin_edit_field:%d:city", userID)),
		),
	}
	if role == "client" {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏢 Бизнес", fmt.Sprintf("admin_edit_field:%d:business_name", userID)),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Описание", fmt.Sprintf("admin_edit_field:%d:description", userID)),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Портфолио", fmt.Sprintf("admin_edit_field:%d:portfolio_links", userID)),
		))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func TaskAdminActions(taskID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", fmt.Sprintf("del_task:%d", taskID)),
		),
	)
}

func PaymentApproval(paymentID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", fmt.Sprintf("pay_ok:%d", paymentID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", fmt.Sprintf("pay_no:%d", paymentID)),
		),
	)
}

func ClientVerifAdminActions(userID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Верифицировать", fmt.Sprintf("verif_client_ok:%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", fmt.Sprintf("verif_client_no:%d", userID)),
		),
	)
}
func AdminPagination(action string, page, totalPages int) []tgbotapi.InlineKeyboardButton {
	if totalPages <= 1 {
		return nil
	}
	var row []tgbotapi.InlineKeyboardButton
	if page > 0 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("%s:%d", action, page-1)))
	}
	row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page+1, totalPages), "noop"))
	if page < totalPages-1 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("%s:%d", action, page+1)))
	}
	return row
}

func NumberedButtons(actionPrefix string, ids []int64, cols int) [][]tgbotapi.InlineKeyboardButton {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton
	for i, id := range ids {
		currentRow = append(currentRow, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+1), fmt.Sprintf("%s:%d", actionPrefix, id)))
		if len(currentRow) == cols {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}
	return rows
}
