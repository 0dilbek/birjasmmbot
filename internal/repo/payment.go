package repo

import (
	"database/sql"
	"time"

	"github.com/birjasmm/bot/internal/models"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(userID int64, subType string, amount int, receiptType, receiptFileID, receiptText string) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO payments (user_id, sub_type, amount, status, receipt_type, receipt_file_id, receipt_text)
		VALUES ($1, $2, $3, 'pending', $4, $5, $6)
		RETURNING id`,
		userID, subType, amount, receiptType, receiptFileID, receiptText).Scan(&id)
	return id, err
}

func (r *PaymentRepo) GetByID(id int64) (*models.Payment, error) {
	p := &models.Payment{}
	var reviewedByTgID sql.NullInt64
	var reviewedAt sql.NullTime
	var receiptFileID, receiptText sql.NullString
	err := r.db.QueryRow(`
		SELECT id, user_id, sub_type, amount, status,
		       COALESCE(receipt_type,''), receipt_file_id, receipt_text,
		       reviewed_by_tg_id, reviewed_at, created_at
		FROM payments WHERE id=$1`, id).
		Scan(&p.ID, &p.UserID, &p.SubType, &p.Amount, &p.Status,
			&p.ReceiptType, &receiptFileID, &receiptText,
			&reviewedByTgID, &reviewedAt, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if receiptFileID.Valid {
		p.ReceiptFileID = receiptFileID.String
	}
	if receiptText.Valid {
		p.ReceiptText = receiptText.String
	}
	if reviewedByTgID.Valid {
		v := reviewedByTgID.Int64
		p.ReviewedByTgID = &v
	}
	if reviewedAt.Valid {
		v := reviewedAt.Time
		p.ReviewedAt = &v
	}
	return p, nil
}

func (r *PaymentRepo) Review(id int64, status string, reviewerTgID int64) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE payments SET status=$1, reviewed_by_tg_id=$2, reviewed_at=$3
		WHERE id=$4 AND status='pending'`, status, reviewerTgID, now, id)
	return err
}

func (r *PaymentRepo) SaveAdminMsg(paymentID, adminTgID, chatID int64, messageID int) error {
	_, err := r.db.Exec(`
		INSERT INTO payment_admin_msgs (payment_id, admin_tg_id, chat_id, message_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (payment_id, admin_tg_id) DO UPDATE SET message_id=$4`,
		paymentID, adminTgID, chatID, messageID)
	return err
}

func (r *PaymentRepo) GetAdminMsgs(paymentID int64) ([]*models.PaymentAdminMsg, error) {
	rows, err := r.db.Query(`
		SELECT id, payment_id, admin_tg_id, chat_id, message_id, created_at
		FROM payment_admin_msgs WHERE payment_id=$1`, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.PaymentAdminMsg
	for rows.Next() {
		m := &models.PaymentAdminMsg{}
		if err := rows.Scan(&m.ID, &m.PaymentID, &m.AdminTgID, &m.ChatID, &m.MessageID, &m.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}
