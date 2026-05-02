package repo

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/birjasmm/bot/internal/models"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(t *models.Task) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO tasks (client_id, title, description, category, budget_type,
		                   budget_from, budget_to, deadline, refs, is_urgent, status, max_responses)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'open',15)
		RETURNING id`,
		t.ClientID, t.Title, t.Description, t.Category, t.BudgetType,
		t.BudgetFrom, t.BudgetTo, t.Deadline, t.Refs, t.IsUrgent).Scan(&id)
	return id, err
}

func (r *TaskRepo) GetByID(id int64) (*models.Task, error) {
	t := &models.Task{}
	var budgetFrom, budgetTo sql.NullInt64
	err := r.db.QueryRow(`
		SELECT t.id, t.client_id, t.title, t.description, t.category,
		       t.budget_type, t.budget_from, t.budget_to, t.deadline,
		       COALESCE(t.refs,''), t.is_urgent, t.status, t.max_responses,
		       t.created_at, t.updated_at,
		       COALESCE(cp.name, ''),
		       (SELECT COUNT(*) FROM responses WHERE task_id=t.id AND status != 'rejected')
		FROM tasks t
		LEFT JOIN client_profiles cp ON cp.user_id = t.client_id
		WHERE t.id=$1`, id).
		Scan(&t.ID, &t.ClientID, &t.Title, &t.Description, &t.Category,
			&t.BudgetType, &budgetFrom, &budgetTo, &t.Deadline,
			&t.Refs, &t.IsUrgent, &t.Status, &t.MaxResponses,
			&t.CreatedAt, &t.UpdatedAt, &t.ClientName, &t.ResponseCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if budgetFrom.Valid {
		v := budgetFrom.Int64
		t.BudgetFrom = &v
	}
	if budgetTo.Valid {
		v := budgetTo.Int64
		t.BudgetTo = &v
	}
	return t, err
}

func (r *TaskRepo) ListByClient(clientID int64) ([]*models.Task, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.client_id, t.title, t.description, t.category,
		       t.budget_type, t.budget_from, t.budget_to, t.deadline,
		       COALESCE(t.refs,''), t.is_urgent, t.status, t.max_responses,
		       t.created_at, t.updated_at, '',
		       (SELECT COUNT(*) FROM responses WHERE task_id=t.id AND status != 'rejected')
		FROM tasks t
		WHERE t.client_id=$1
		ORDER BY t.created_at DESC`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

type TaskFilter struct {
	Category      string
	City          string
	BudgetMin     int64
	ExcludeUrgent bool
	Page          int
	PageSize      int
}

func (r *TaskRepo) ListOpen(f TaskFilter) ([]*models.Task, error) {
	conds := []string{"t.status = 'open'"}
	args := []interface{}{}
	i := 1

	if f.Category != "" {
		conds = append(conds, fmt.Sprintf("t.category = $%d", i))
		args = append(args, f.Category)
		i++
	}
	if f.City != "" {
		conds = append(conds, fmt.Sprintf("LOWER(COALESCE(cp.city,'')) = LOWER($%d)", i))
		args = append(args, f.City)
		i++
	}
	if f.BudgetMin > 0 {
		conds = append(conds, fmt.Sprintf("(t.budget_from >= $%d OR t.budget_type = 'negotiable')", i))
		args = append(args, f.BudgetMin)
		i++
	}
	if f.ExcludeUrgent {
		conds = append(conds, "t.is_urgent = FALSE")
	}

	where := "WHERE " + strings.Join(conds, " AND ")
	if f.PageSize == 0 {
		f.PageSize = 10
	}
	offset := f.Page * f.PageSize

	query := fmt.Sprintf(`
		SELECT t.id, t.client_id, t.title, t.description, t.category,
		       t.budget_type, t.budget_from, t.budget_to, t.deadline,
		       COALESCE(t.refs,''), t.is_urgent, t.status, t.max_responses,
		       t.created_at, t.updated_at, COALESCE(cp.name,''),
		       (SELECT COUNT(*) FROM responses WHERE task_id=t.id AND status != 'rejected')
		FROM tasks t
		LEFT JOIN client_profiles cp ON cp.user_id = t.client_id
		%s
		ORDER BY t.is_urgent DESC, t.created_at DESC
		LIMIT $%d OFFSET $%d`, where, i, i+1)

	args = append(args, f.PageSize, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

func (r *TaskRepo) UpdateStatus(id int64, status string) error {
	_, err := r.db.Exec(`UPDATE tasks SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *TaskRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
	return err
}

func (r *TaskRepo) CountOpen() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status='open'`).Scan(&count)
	return count, err
}

func (r *TaskRepo) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&count)
	return count, err
}

func scanTasks(rows *sql.Rows) ([]*models.Task, error) {
	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		var budgetFrom, budgetTo sql.NullInt64
		if err := rows.Scan(&t.ID, &t.ClientID, &t.Title, &t.Description, &t.Category,
			&t.BudgetType, &budgetFrom, &budgetTo, &t.Deadline,
			&t.Refs, &t.IsUrgent, &t.Status, &t.MaxResponses,
			&t.CreatedAt, &t.UpdatedAt, &t.ClientName, &t.ResponseCount); err != nil {
			return nil, err
		}
		if budgetFrom.Valid {
			v := budgetFrom.Int64
			t.BudgetFrom = &v
		}
		if budgetTo.Valid {
			v := budgetTo.Int64
			t.BudgetTo = &v
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Response repo methods

type ResponseRepo struct {
	db *sql.DB
}

func NewResponseRepo(db *sql.DB) *ResponseRepo {
	return &ResponseRepo{db: db}
}

func (r *ResponseRepo) Create(taskID, executorID int64, message string, price *int64) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO responses (task_id, executor_id, message, proposed_price, status)
		VALUES ($1, $2, $3, $4, 'active')
		RETURNING id`,
		taskID, executorID, message, price).Scan(&id)
	return id, err
}

func (r *ResponseRepo) CountForTask(taskID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM responses WHERE task_id=$1 AND status != 'rejected'`, taskID).Scan(&count)
	return count, err
}

func (r *ResponseRepo) ExistsForExecutor(taskID, executorID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM responses WHERE task_id=$1 AND executor_id=$2`, taskID, executorID).Scan(&count)
	return count > 0, err
}

func (r *ResponseRepo) ListForTask(taskID int64) ([]*models.Response, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.task_id, r.executor_id, r.message,
		       r.proposed_price, r.status, r.created_at,
		       COALESCE(ep.name,''), COALESCE(ep.rating, 0),
		       COALESCE(ep.is_pro, FALSE), COALESCE(ep.is_verified, FALSE)
		FROM responses r
		LEFT JOIN executor_profiles ep ON ep.user_id = r.executor_id
		WHERE r.task_id=$1
		ORDER BY ep.is_pro DESC, ep.rating DESC, r.created_at ASC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanResponses(rows)
}

func (r *ResponseRepo) ListForExecutor(executorID int64) ([]*models.Response, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.task_id, r.executor_id, r.message,
		       r.proposed_price, r.status, r.created_at,
		       COALESCE(ep.name,''), COALESCE(ep.rating, 0),
		       COALESCE(ep.is_pro, FALSE), COALESCE(ep.is_verified, FALSE)
		FROM responses r
		LEFT JOIN executor_profiles ep ON ep.user_id = r.executor_id
		WHERE r.executor_id=$1
		ORDER BY r.created_at DESC`, executorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanResponses(rows)
}

func (r *ResponseRepo) GetByID(id int64) (*models.Response, error) {
	resp := &models.Response{}
	var price sql.NullInt64
	err := r.db.QueryRow(`
		SELECT r.id, r.task_id, r.executor_id, r.message,
		       r.proposed_price, r.status, r.created_at,
		       COALESCE(ep.name,''), COALESCE(ep.rating, 0),
		       COALESCE(ep.is_pro, FALSE), COALESCE(ep.is_verified, FALSE)
		FROM responses r
		LEFT JOIN executor_profiles ep ON ep.user_id = r.executor_id
		WHERE r.id=$1`, id).
		Scan(&resp.ID, &resp.TaskID, &resp.ExecutorID, &resp.Message,
			&price, &resp.Status, &resp.CreatedAt,
			&resp.ExecutorName, &resp.Rating, &resp.IsPro, &resp.IsVerified)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if price.Valid {
		v := price.Int64
		resp.ProposedPrice = &v
	}
	return resp, err
}

func (r *ResponseRepo) Accept(id int64) error {
	_, err := r.db.Exec(`UPDATE responses SET status='accepted' WHERE id=$1`, id)
	return err
}

func (r *ResponseRepo) RejectOthers(taskID, acceptedResponseID int64) error {
	_, err := r.db.Exec(`
		UPDATE responses SET status='rejected'
		WHERE task_id=$1 AND id != $2 AND status='active'`, taskID, acceptedResponseID)
	return err
}

func (r *ResponseRepo) ListRejectedForTask(taskID int64) ([]*models.Response, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.task_id, r.executor_id, r.message,
		       r.proposed_price, r.status, r.created_at,
		       COALESCE(ep.name,''), COALESCE(ep.rating,0),
		       COALESCE(ep.is_pro,FALSE), COALESCE(ep.is_verified,FALSE)
		FROM responses r
		LEFT JOIN executor_profiles ep ON ep.user_id = r.executor_id
		WHERE r.task_id=$1 AND r.status='rejected'`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanResponses(rows)
}

func (r *ResponseRepo) Withdraw(respID, executorID int64) error {
	_, err := r.db.Exec(`
		UPDATE responses SET status='rejected'
		WHERE id=$1 AND executor_id=$2 AND status='active'`, respID, executorID)
	return err
}

func (r *ResponseRepo) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM responses`).Scan(&count)
	return count, err
}

func scanResponses(rows *sql.Rows) ([]*models.Response, error) {
	var list []*models.Response
	for rows.Next() {
		resp := &models.Response{}
		var price sql.NullInt64
		if err := rows.Scan(&resp.ID, &resp.TaskID, &resp.ExecutorID, &resp.Message,
			&price, &resp.Status, &resp.CreatedAt,
			&resp.ExecutorName, &resp.Rating, &resp.IsPro, &resp.IsVerified); err != nil {
			return nil, err
		}
		if price.Valid {
			v := price.Int64
			resp.ProposedPrice = &v
		}
		list = append(list, resp)
	}
	return list, nil
}

// TaskAssignment

func (r *ResponseRepo) CreateAssignment(taskID, executorID int64) error {
	_, err := r.db.Exec(`
		INSERT INTO task_assignments (task_id, executor_id)
		VALUES ($1, $2)
		ON CONFLICT (task_id) DO UPDATE SET executor_id=$2, assigned_at=NOW()`,
		taskID, executorID)
	return err
}

func (r *ResponseRepo) GetAssignment(taskID int64) (*models.TaskAssignment, error) {
	a := &models.TaskAssignment{}
	err := r.db.QueryRow(`
		SELECT id, task_id, executor_id, assigned_at, completed_at
		FROM task_assignments WHERE task_id=$1`, taskID).
		Scan(&a.ID, &a.TaskID, &a.ExecutorID, &a.AssignedAt, &a.CompletedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *ResponseRepo) CompleteAssignment(taskID int64) error {
	_, err := r.db.Exec(`
		UPDATE task_assignments SET completed_at=NOW() WHERE task_id=$1`, taskID)
	return err
}

func (r *ResponseRepo) GetActiveAssignments(executorID int64) ([]*models.TaskAssignment, error) {
	rows, err := r.db.Query(`
		SELECT id, task_id, executor_id, assigned_at, completed_at
		FROM task_assignments
		WHERE executor_id=$1 AND completed_at IS NULL
		ORDER BY assigned_at DESC`, executorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.TaskAssignment
	for rows.Next() {
		a := &models.TaskAssignment{}
		if err := rows.Scan(&a.ID, &a.TaskID, &a.ExecutorID, &a.AssignedAt, &a.CompletedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

// Review

type ReviewRepo struct {
	db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(taskID, clientID, executorID int64, rating int, comment string) error {
	_, err := r.db.Exec(`
		INSERT INTO reviews (task_id, client_id, executor_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (task_id) DO UPDATE SET rating=$4, comment=$5`,
		taskID, clientID, executorID, rating, comment)
	return err
}

func (r *ReviewRepo) GetByTask(taskID int64) (*models.Review, error) {
	rev := &models.Review{}
	err := r.db.QueryRow(`
		SELECT id, task_id, client_id, executor_id, rating, COALESCE(comment,''), created_at
		FROM reviews WHERE task_id=$1`, taskID).
		Scan(&rev.ID, &rev.TaskID, &rev.ClientID, &rev.ExecutorID,
			&rev.Rating, &rev.Comment, &rev.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return rev, err
}
