package fsm

import "sync"

type State string

const (
	StateIdle State = ""

	StateChooseRole State = "choose_role"

	StateClientName     State = "client_name"
	StateClientPhone    State = "client_phone"
	StateClientBusiness State = "client_business"
	StateClientCity     State = "client_city"

	StateExecName        State = "exec_name"
	StateExecPhone       State = "exec_phone"
	StateExecCity        State = "exec_city"
	StateExecCategory    State = "exec_category"
	StateExecExperience  State = "exec_experience"
	StateExecPortfolio   State = "exec_portfolio"
	StateExecDescription State = "exec_description"

	StateVerifVideo State = "verif_video"

	StateTaskTitle       State = "task_title"
	StateTaskDescription State = "task_description"
	StateTaskCategory    State = "task_category"
	StateTaskBudgetType  State = "task_budget_type"
	StateTaskBudgetFrom  State = "task_budget_from"
	StateTaskBudgetTo    State = "task_budget_to"
	StateTaskDeadline    State = "task_deadline"
	StateTaskRefs        State = "task_refs"
	StateTaskUrgent      State = "task_urgent"

	StateRespondMsg   State = "respond_msg"
	StateRespondPrice State = "respond_price"

	StateReviewRating  State = "review_rating"
	StateReviewComment State = "review_comment"

	StatePaymentReceipt State = "payment_receipt"
	StateFilterCity     State = "filter_city"

	StateAdminEditUser  State = "admin_edit_user"
	StateClientVerif    State = "client_verif"
	StateFilterTaskCity State = "filter_task_city"
)

var (
	states sync.Map
	data   sync.Map
)

func Get(userID int64) State {
	if s, ok := states.Load(userID); ok {
		return s.(State)
	}
	return StateIdle
}

func Set(userID int64, s State) {
	states.Store(userID, s)
}

func Clear(userID int64) {
	states.Delete(userID)
	data.Delete(userID)
}

func getData(userID int64) map[string]interface{} {
	if d, ok := data.Load(userID); ok {
		return d.(map[string]interface{})
	}
	return make(map[string]interface{})
}

func SetVal(userID int64, key string, val interface{}) {
	d := getData(userID)
	d[key] = val
	data.Store(userID, d)
}

func GetVal(userID int64, key string) (interface{}, bool) {
	d := getData(userID)
	v, ok := d[key]
	return v, ok
}

func GetStr(userID int64, key string) string {
	v, ok := GetVal(userID, key)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func GetInt(userID int64, key string) int {
	v, ok := GetVal(userID, key)
	if !ok {
		return 0
	}
	i, _ := v.(int)
	return i
}

func GetInt64(userID int64, key string) int64 {
	v, ok := GetVal(userID, key)
	if !ok {
		return 0
	}
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	}
	return 0
}

func GetBool(userID int64, key string) bool {
	v, ok := GetVal(userID, key)
	if !ok {
		return false
	}
	b, _ := v.(bool)
	return b
}

func ClearData(userID int64) {
	data.Delete(userID)
}
