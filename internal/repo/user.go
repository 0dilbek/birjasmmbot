package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/birjasmm/bot/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByTelegramID(tgID int64) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, telegram_id, COALESCE(username,''), COALESCE(phone,''),
		       role, status, verification_status, created_at, updated_at
		FROM users WHERE telegram_id = $1`, tgID).
		Scan(&u.ID, &u.TelegramID, &u.Username, &u.Phone,
			&u.Role, &u.Status, &u.VerificationStatus, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) GetByID(id int64) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, telegram_id, COALESCE(username,''), COALESCE(phone,''),
		       role, status, verification_status, created_at, updated_at
		FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.TelegramID, &u.Username, &u.Phone,
			&u.Role, &u.Status, &u.VerificationStatus, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) Create(tgID int64, username, role, lang string) (*models.User, error) {
	if lang == "" {
		lang = "ru"
	}
	u := &models.User{}
	err := r.db.QueryRow(`
		INSERT INTO users (telegram_id, username, role, status, verification_status, language)
		VALUES ($1, $2, $3, 'active', 'none', $4)
		RETURNING id, telegram_id, COALESCE(username,''), COALESCE(phone,''),
		          role, status, verification_status, created_at, updated_at`,
		tgID, username, role, lang).
		Scan(&u.ID, &u.TelegramID, &u.Username, &u.Phone,
			&u.Role, &u.Status, &u.VerificationStatus, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *UserRepo) GetLanguage(tgID int64) (string, error) {
	var lang string
	err := r.db.QueryRow(`SELECT COALESCE(language, 'ru') FROM users WHERE telegram_id=$1`, tgID).Scan(&lang)
	if err == sql.ErrNoRows {
		return "ru", nil
	}
	return lang, err
}

func (r *UserRepo) SetLanguage(userID int64, lang string) error {
	_, err := r.db.Exec(`UPDATE users SET language=$1, updated_at=NOW() WHERE id=$2`, lang, userID)
	return err
}

func (r *UserRepo) SetPhone(userID int64, phone string) error {
	_, err := r.db.Exec(`UPDATE users SET phone=$1, updated_at=NOW() WHERE id=$2`, phone, userID)
	return err
}

func (r *UserRepo) SetRole(userID int64, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`, role, userID)
	return err
}

func (r *UserRepo) SetVerifStatus(userID int64, status string) error {
	_, err := r.db.Exec(`UPDATE users SET verification_status=$1, updated_at=NOW() WHERE id=$2`, status, userID)
	return err
}

func (r *UserRepo) Block(userID int64) error {
	_, err := r.db.Exec(`UPDATE users SET status='blocked', updated_at=NOW() WHERE id=$1`, userID)
	return err
}

func (r *UserRepo) Unblock(userID int64) error {
	_, err := r.db.Exec(`UPDATE users SET status='active', updated_at=NOW() WHERE id=$1`, userID)
	return err
}

func (r *UserRepo) ListAll(limit, offset int) ([]*models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, telegram_id, COALESCE(username,''), COALESCE(phone,''),
		       role, status, verification_status, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.TelegramID, &u.Username, &u.Phone,
			&u.Role, &u.Status, &u.VerificationStatus, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

// ClientProfile

func (r *UserRepo) CreateClientProfile(p *models.ClientProfile) error {
	_, err := r.db.Exec(`
		INSERT INTO client_profiles (user_id, name, business_name, city)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET name=$2, business_name=$3, city=$4`,
		p.UserID, p.Name, p.BusinessName, p.City)
	return err
}

func (r *UserRepo) GetClientProfile(userID int64) (*models.ClientProfile, error) {
	p := &models.ClientProfile{}
	err := r.db.QueryRow(`
		SELECT id, user_id, name, COALESCE(business_name,''), COALESCE(city,''),
		       is_verified, COALESCE(description,'')
		FROM client_profiles WHERE user_id=$1`, userID).
		Scan(&p.ID, &p.UserID, &p.Name, &p.BusinessName, &p.City, &p.IsVerified, &p.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *UserRepo) VerifyClient(userID int64) error {
	_, err := r.db.Exec(`UPDATE client_profiles SET is_verified=TRUE WHERE user_id=$1`, userID)
	return err
}

// ExecutorProfile

func (r *UserRepo) CreateExecutorProfile(p *models.ExecutorProfile) error {
	_, err := r.db.Exec(`
		INSERT INTO executor_profiles (user_id, name, city, category, experience_years, description, portfolio_links)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE
		SET name=$2, city=$3, category=$4, experience_years=$5, description=$6, portfolio_links=$7`,
		p.UserID, p.Name, p.City, p.Category, p.ExperienceYears, p.Description, p.PortfolioLinks)
	return err
}

func (r *UserRepo) GetExecutorProfile(userID int64) (*models.ExecutorProfile, error) {
	p := &models.ExecutorProfile{}
	err := r.db.QueryRow(`
		SELECT id, user_id, name, COALESCE(city,''), COALESCE(category,''),
		       experience_years, COALESCE(description,''), COALESCE(portfolio_links,''),
		       rating, total_orders, completed_orders, response_speed,
		       is_verified, is_pro, created_at
		FROM executor_profiles WHERE user_id=$1`, userID).
		Scan(&p.ID, &p.UserID, &p.Name, &p.City, &p.Category,
			&p.ExperienceYears, &p.Description, &p.PortfolioLinks,
			&p.Rating, &p.TotalOrders, &p.CompletedOrders, &p.ResponseSpeed,
			&p.IsVerified, &p.IsPro, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *UserRepo) SetExecutorPro(userID int64, isPro bool) error {
	_, err := r.db.Exec(`UPDATE executor_profiles SET is_pro=$1 WHERE user_id=$2`, isPro, userID)
	return err
}

func (r *UserRepo) SetExecutorVerified(userID int64) error {
	_, err := r.db.Exec(`UPDATE executor_profiles SET is_verified=TRUE WHERE user_id=$1`, userID)
	return err
}

func (r *UserRepo) UpdateRating(userID int64) error {
	_, err := r.db.Exec(`
		UPDATE executor_profiles ep SET
			rating = (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE executor_id = ep.user_id),
			completed_orders = (SELECT COUNT(*) FROM reviews WHERE executor_id = ep.user_id)
		WHERE ep.user_id = $1`, userID)
	return err
}

// Verification

func (r *UserRepo) CreateVerification(userID int64, fileID string) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO verifications (user_id, video_file_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id`, userID, fileID).Scan(&id)
	return id, err
}


func (r *UserRepo) GetPendingVerifications() ([]*models.Verification, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, video_file_id, status, created_at
		FROM verifications WHERE status='pending' ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Verification
	for rows.Next() {
		v := &models.Verification{}
		if err := rows.Scan(&v.ID, &v.UserID, &v.VideoFileID, &v.Status, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func (r *UserRepo) ReviewVerification(id int64, status string) error {
	now := time.Now()
	_, err := r.db.Exec(`UPDATE verifications SET status=$1, reviewed_at=$2 WHERE id=$3`,
		status, now, id)
	return err
}

func (r *UserRepo) GetVerificationByID(id int64) (*models.Verification, error) {
	v := &models.Verification{}
	err := r.db.QueryRow(`
		SELECT id, user_id, video_file_id, status, created_at
		FROM verifications WHERE id=$1`, id).
		Scan(&v.ID, &v.UserID, &v.VideoFileID, &v.Status, &v.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return v, err
}

func (r *UserRepo) SaveVerificationAdminMsg(verifID, adminTgID, chatID int64, msgID int) error {
	_, err := r.db.Exec(`
		INSERT INTO verification_admin_msgs (verification_id, admin_tg_id, chat_id, message_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (verification_id, admin_tg_id) DO UPDATE 
		SET chat_id=$3, message_id=$4`,
		verifID, adminTgID, chatID, msgID)
	return err
}

func (r *UserRepo) GetVerificationAdminMsgs(verifID int64) ([]models.PaymentAdminMsg, error) {
	rows, err := r.db.Query(`
		SELECT id, verification_id, admin_tg_id, chat_id, message_id, created_at
		FROM verification_admin_msgs WHERE verification_id=$1`, verifID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.PaymentAdminMsg
	for rows.Next() {
		var m models.PaymentAdminMsg
		if err := rows.Scan(&m.ID, &m.PaymentID, &m.AdminTgID, &m.ChatID, &m.MessageID, &m.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}


// UsageLimits

func (r *UserRepo) GetOrCreateLimits(userID int64, defaultFree int) (*models.UsageLimits, error) {
	l := &models.UsageLimits{}
	err := r.db.QueryRow(`
		INSERT INTO usage_limits (user_id, free_responses_left)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET user_id=$1
		RETURNING id, user_id, free_responses_left, last_reset_date`,
		userID, defaultFree).
		Scan(&l.ID, &l.UserID, &l.FreeResponsesLeft, &l.LastResetDate)
	return l, err
}

func (r *UserRepo) DecrementLimit(userID int64) error {
	_, err := r.db.Exec(`
		UPDATE usage_limits SET free_responses_left = free_responses_left - 1
		WHERE user_id=$1 AND free_responses_left > 0`, userID)
	return err
}

// Subscription

func (r *UserRepo) GetSubscription(userID int64) (*models.Subscription, error) {
	s := &models.Subscription{}
	err := r.db.QueryRow(`
		SELECT id, user_id, type, start_date, end_date, is_active
		FROM subscriptions WHERE user_id=$1 AND is_active=TRUE AND end_date > NOW()`, userID).
		Scan(&s.ID, &s.UserID, &s.Type, &s.StartDate, &s.EndDate, &s.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *UserRepo) CreateSubscription(userID int64, subType string, days int) error {
	_, err := r.db.Exec(`
		INSERT INTO subscriptions (user_id, type, end_date, is_active)
		VALUES ($1, $2, NOW() + INTERVAL '1 day' * $3, TRUE)
		ON CONFLICT (user_id) DO UPDATE
		SET type = $2,
		    end_date = CASE
		        WHEN subscriptions.end_date > NOW()
		        THEN subscriptions.end_date + INTERVAL '1 day' * $3
		        ELSE NOW() + INTERVAL '1 day' * $3
		    END,
		    is_active = TRUE`,
		userID, subType, days)
	return err
}

var allowedClientFields = map[string]bool{
	"name": true, "business_name": true, "city": true,
}

var allowedExecutorFields = map[string]bool{
	"name": true, "city": true, "description": true, "portfolio_links": true,
}

func (r *UserRepo) UpdateClientField(userID int64, field, value string) error {
	if !allowedClientFields[field] {
		return nil
	}
	_, err := r.db.Exec(
		fmt.Sprintf(`UPDATE client_profiles SET %s=$1 WHERE user_id=$2`, field),
		value, userID)
	return err
}

func (r *UserRepo) UpdateExecutorField(userID int64, field, value string) error {
	if !allowedExecutorFields[field] {
		return nil
	}
	_, err := r.db.Exec(
		fmt.Sprintf(`UPDATE executor_profiles SET %s=$1 WHERE user_id=$2`, field),
		value, userID)
	return err
}

func (r *UserRepo) LogAction(userID int64, action, data string) {
	r.db.Exec(`INSERT INTO action_logs (user_id, action, data) VALUES ($1, $2, $3)`,
		userID, action, data)
}

func (r *UserRepo) GetExecutorsByCategory(cat string, limit int) ([]*models.User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.telegram_id, COALESCE(u.username,''), COALESCE(u.phone,''),
		       u.role, u.status, u.verification_status, u.created_at, u.updated_at
		FROM users u
		JOIN executor_profiles ep ON ep.user_id = u.id
		WHERE ep.category=$1 AND u.status='active'
		ORDER BY ep.is_pro DESC, ep.rating DESC
		LIMIT $2`, cat, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.TelegramID, &u.Username, &u.Phone,
			&u.Role, &u.Status, &u.VerificationStatus, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) GetTopExecutors(cat string, limit int) ([]*models.ExecutorProfile, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, name, COALESCE(city,''), COALESCE(category,''),
		       experience_years, COALESCE(description,''), COALESCE(portfolio_links,''),
		       rating, total_orders, completed_orders, response_speed,
		       is_verified, is_pro, created_at
		FROM executor_profiles
		WHERE ($1 = '' OR category=$1) AND is_verified=TRUE
		ORDER BY is_pro DESC, rating DESC, completed_orders DESC
		LIMIT $2`, cat, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.ExecutorProfile
	for rows.Next() {
		p := &models.ExecutorProfile{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.City, &p.Category,
			&p.ExperienceYears, &p.Description, &p.PortfolioLinks,
			&p.Rating, &p.TotalOrders, &p.CompletedOrders, &p.ResponseSpeed,
			&p.IsVerified, &p.IsPro, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
