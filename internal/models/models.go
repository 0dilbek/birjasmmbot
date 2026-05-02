package models

import "time"

type Role string

const (
	RoleClient   Role = "client"
	RoleExecutor Role = "executor"
)

type UserStatus string

const (
	StatusActive  UserStatus = "active"
	StatusBlocked UserStatus = "blocked"
)

type VerifStatus string

const (
	VerifNone     VerifStatus = "none"
	VerifPending  VerifStatus = "pending"
	VerifVerified VerifStatus = "verified"
)

type Category string

const (
	CatSMM     Category = "smm"
	CatVideo   Category = "video"
	CatBlogger Category = "blogger"
)

type BudgetType string

const (
	BudgetFixed      BudgetType = "fixed"
	BudgetRange      BudgetType = "range"
	BudgetNegotiable BudgetType = "negotiable"
)

type TaskStatus string

const (
	TaskOpen       TaskStatus = "open"
	TaskInProgress TaskStatus = "in_progress"
	TaskCompleted  TaskStatus = "completed"
	TaskCancelled  TaskStatus = "cancelled"
)

type RespStatus string

const (
	RespActive   RespStatus = "active"
	RespAccepted RespStatus = "accepted"
	RespRejected RespStatus = "rejected"
)

type SubType string

const (
	SubBasic SubType = "basic"
	SubPro   SubType = "pro"
)

type User struct {
	ID                 int64
	TelegramID         int64
	Username           string
	Phone              string
	Role               Role
	Status             UserStatus
	VerificationStatus VerifStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ClientProfile struct {
	ID           int64
	UserID       int64
	Name         string
	BusinessName string
	City         string
	IsVerified   bool
	Description  string
}

type ExecutorProfile struct {
	ID              int64
	UserID          int64
	Name            string
	City            string
	Category        Category
	ExperienceYears int
	Description     string
	PortfolioLinks  string
	Rating          float64
	TotalOrders     int
	CompletedOrders int
	ResponseSpeed   float64
	IsVerified      bool
	IsPro           bool
	CreatedAt       time.Time
}

type Verification struct {
	ID          int64
	UserID      int64
	VideoFileID string
	Status      string
	CreatedAt   time.Time
	ReviewedAt  *time.Time
}

type Task struct {
	ID           int64
	ClientID     int64
	Title        string
	Description  string
	Category     Category
	BudgetType   BudgetType
	BudgetFrom   *int64
	BudgetTo     *int64
	Deadline     time.Time
	Refs         string
	IsUrgent     bool
	Status       TaskStatus
	MaxResponses int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// Joined fields
	ClientName    string
	ResponseCount int
}

type Response struct {
	ID            int64
	TaskID        int64
	ExecutorID    int64
	Message       string
	ProposedPrice *int64
	Status        RespStatus
	CreatedAt     time.Time
	// Joined fields
	ExecutorName string
	Rating       float64
	IsPro        bool
	IsVerified   bool
}

type TaskAssignment struct {
	ID          int64
	TaskID      int64
	ExecutorID  int64
	AssignedAt  time.Time
	CompletedAt *time.Time
}

type Review struct {
	ID         int64
	TaskID     int64
	ClientID   int64
	ExecutorID int64
	Rating     int
	Comment    string
	CreatedAt  time.Time
}

type Subscription struct {
	ID        int64
	UserID    int64
	Type      SubType
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

type UsageLimits struct {
	ID                int64
	UserID            int64
	FreeResponsesLeft int
	LastResetDate     time.Time
}

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentApproved PaymentStatus = "approved"
	PaymentRejected PaymentStatus = "rejected"
)

type Payment struct {
	ID             int64
	UserID         int64
	SubType        string
	Amount         int
	Status         PaymentStatus
	ReceiptType    string
	ReceiptFileID  string
	ReceiptText    string
	ReviewedByTgID *int64
	ReviewedAt     *time.Time
	CreatedAt      time.Time
}

type PaymentAdminMsg struct {
	ID        int64
	PaymentID int64
	AdminTgID int64
	ChatID    int64
	MessageID int
	CreatedAt time.Time
}
