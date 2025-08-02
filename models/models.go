package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Wallet represents a user's wallet
type Wallet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Sport represents a sports category
type Sport struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Event represents a sports event
type Event struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SportID     uuid.UUID `json:"sport_id" db:"sport_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Country     string    `json:"country" db:"country"` // Country where the event takes place
	League      string    `json:"league" db:"league"`   // League or competition name
	HomeTeam    string    `json:"home_team" db:"home_team"`
	AwayTeam    string    `json:"away_team" db:"away_team"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Status      string    `json:"status" db:"status"` // "upcoming", "live", "finished", "cancelled"
	Result      string    `json:"result" db:"result"` // "home_win", "away_win", "draw", "cancelled"
	Version     int       `json:"version" db:"version"` // Version number for audit trail
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Market represents a betting market for an event
type Market struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	Name        string    `json:"name" db:"name"` // "match_winner", "total_goals", "first_scorer", etc.
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`     // "1x2", "over_under", "handicap", "exact_score"
	Status      string    `json:"status" db:"status"` // "open", "suspended", "closed", "settled"
	Version     int       `json:"version" db:"version"` // Version number for audit trail
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Selection represents a betting option within a market
type Selection struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MarketID  uuid.UUID `json:"market_id" db:"market_id"`
	Name      string    `json:"name" db:"name"` // "Home Win", "Away Win", "Draw", "Over 2.5", etc.
	Status    string    `json:"status" db:"status"` // "active", "suspended", "won", "lost"
	Version   int       `json:"version" db:"version"` // Version number for audit trail
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Odds represents odds for a selection
type Odds struct {
	ID         uuid.UUID `json:"id" db:"id"`
	SelectionID uuid.UUID `json:"selection_id" db:"selection_id"`
	Odds       float64   `json:"odds" db:"odds"`
	Version    int       `json:"version" db:"version"` // Version number for audit trail
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Result represents a result for a selection
type Result struct {
	ID         uuid.UUID `json:"id" db:"id"`
	SelectionID uuid.UUID `json:"selection_id" db:"selection_id"`
	Result     string    `json:"result" db:"result"` // "won", "lost", "void", "pending"
	Version    int       `json:"version" db:"version"` // Version number for audit trail
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Bet represents a user's bet
type Bet struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	BetType      string    `json:"bet_type" db:"bet_type"` // "single", "double", "treble", "accumulator"
	Amount       float64   `json:"amount" db:"amount"`
	TotalOdds    float64   `json:"total_odds" db:"total_odds"`
	Status       string    `json:"status" db:"status"` // "pending", "won", "lost", "void", "cancelled"
	PotentialWin float64   `json:"potential_win" db:"potential_win"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BetSelection represents a selection within a bet (for accumulator bets)
type BetSelection struct {
	ID          uuid.UUID `json:"id" db:"id"`
	BetID       uuid.UUID `json:"bet_id" db:"bet_id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	MarketID    uuid.UUID `json:"market_id" db:"market_id"`
	SelectionID uuid.UUID `json:"selection_id" db:"selection_id"`
	Odds        float64   `json:"odds" db:"odds"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Transaction represents a wallet transaction
type Transaction struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Type        string    `json:"type" db:"type"` // "deposit", "withdrawal", "bet", "win", "loss"
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Settlement represents a bet settlement
type Settlement struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	EventID     string             `json:"event_id" db:"event_id"`
	TotalBets   int                `json:"total_bets" db:"total_bets"`
	WinningBets int                `json:"winning_bets" db:"winning_bets"`
	TotalPayout float64            `json:"total_payout" db:"total_payout"`
	Status      string             `json:"status" db:"status"` // "pending", "processing", "completed", "failed"
	Results     []SettlementResult `json:"results" db:"results"`
	ProcessedAt time.Time          `json:"processed_at" db:"processed_at"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// Cashout represents a cashout request
type Cashout struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	BetID         uuid.UUID  `json:"bet_id" db:"bet_id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	EventID       uuid.UUID  `json:"event_id" db:"event_id"`
	CashoutType   string     `json:"cashout_type" db:"cashout_type"` // "full", "partial"
	Status        string     `json:"status" db:"status"`             // "pending", "approved", "rejected", "processed", "cancelled"
	OriginalStake float64    `json:"original_stake" db:"original_stake"`
	CashoutAmount float64    `json:"cashout_amount" db:"cashout_amount"`
	CashoutFee    float64    `json:"cashout_fee" db:"cashout_fee"`
	NetCashout    float64    `json:"net_cashout" db:"net_cashout"`
	CashoutOdds   float64    `json:"cashout_odds" db:"cashout_odds"`
	OriginalOdds  float64    `json:"original_odds" db:"original_odds"`
	Currency      string     `json:"currency" db:"currency"`
	CashoutReason string     `json:"cashout_reason" db:"cashout_reason"`
	ProcessedAt   *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response Models

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// PlaceBetRequest represents a bet placement request
type PlaceBetRequest struct {
	EventID     uuid.UUID `json:"event_id" binding:"required"`
	MarketID    uuid.UUID `json:"market_id" binding:"required"`
	SelectionID uuid.UUID `json:"selection_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
}

// PlaceAccumulatorBetRequest represents a request to place an accumulator bet
type PlaceAccumulatorBetRequest struct {
	Selections []PlaceBetRequest `json:"selections" binding:"required,min=1"`
	Amount     float64           `json:"amount" binding:"required,gt=0"`
}

// DepositRequest represents a deposit request
type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// WithdrawRequest represents a withdrawal request
type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// CashoutRequest represents a cashout request
type CashoutRequest struct {
	BetID         uuid.UUID `json:"bet_id" binding:"required"`
	CashoutType   string    `json:"cashout_type" binding:"required"` // "full", "partial"
	CashoutAmount float64   `json:"cashout_amount"`                  // For partial cashouts
	Reason        string    `json:"reason"`
}

// SettlementRequest represents a settlement request
type SettlementRequest struct {
	EventID          string `json:"event_id" binding:"required"`
	SettlementType   string `json:"settlement_type" binding:"required"` // "automatic", "manual", "cashout"
	SettlementReason string `json:"settlement_reason"`
}

// SettlementResult represents the result of a bet settlement
type SettlementResult struct {
	BetID       uuid.UUID `json:"bet_id" db:"bet_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Odds        float64   `json:"odds" db:"odds"`
	Status      string    `json:"status" db:"status"` // "won", "lost", "void"
	Payout      float64   `json:"payout" db:"payout"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

// Dispute represents a settlement dispute
type Dispute struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	SettlementID uuid.UUID  `json:"settlement_id" db:"settlement_id"`
	BetID        uuid.UUID  `json:"bet_id" db:"bet_id"`
	Reason       string     `json:"reason" db:"reason"`
	Evidence     string     `json:"evidence" db:"evidence"`
	Status       string     `json:"status" db:"status"` // "pending", "upheld", "rejected"
	Notes        string     `json:"notes" db:"notes"`
	ResolvedAt   *time.Time `json:"resolved_at" db:"resolved_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// EventResponse represents an event with additional info
type EventResponse struct {
	Event
	SportName string   `json:"sport_name"`
	Markets   []Market `json:"markets,omitempty"`
}

// MarketResponse represents a market with selections
type MarketResponse struct {
	Market
	Selections []Selection `json:"selections"`
}

// BetResponse represents a bet with additional info
type BetResponse struct {
	Bet
	EventTitle    string `json:"event_title"`
	MarketName    string `json:"market_name"`
	SelectionName string `json:"selection_name"`
	SportName     string `json:"sport_name"`
}

type BetWithSelectionsResponse struct {
	Bet
	Selections []BetSelectionWithDetails `json:"selections"`
}

type BetSelectionWithDetails struct {
	BetSelection
	EventTitle    string `json:"event_title"`
	MarketName    string `json:"market_name"`
	SelectionName string `json:"selection_name"`
	SportName     string `json:"sport_name"`
}

// CashoutCalculation represents a cashout calculation
type CashoutCalculation struct {
	BetID          uuid.UUID `json:"bet_id"`
	OriginalStake  float64   `json:"original_stake"`
	OriginalOdds   float64   `json:"original_odds"`
	CashoutAmount  float64   `json:"cashout_amount"`
	CashoutFee     float64   `json:"cashout_fee"`
	NetCashout     float64   `json:"net_cashout"`
	CashoutOdds    float64   `json:"cashout_odds"`
	RemainingStake float64   `json:"remaining_stake,omitempty"`
}

// CashoutEligibility represents cashout eligibility
type CashoutEligibility struct {
	BetID        uuid.UUID `json:"bet_id"`
	IsEligible   bool      `json:"is_eligible"`
	Reason       string    `json:"reason,omitempty"`
	CashoutValue float64   `json:"cashout_value,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthStatus represents service health status
type HealthStatus struct {
	Service   string                 `json:"service"`
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
}
