package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	// External API tracking
	ExternalKey  string    `json:"external_key" db:"external_key"`   // The Odds API sport key
	Group        string    `json:"group" db:"group"`                 // Sport group (e.g., "American Football")
	HasOutrights bool      `json:"has_outrights" db:"has_outrights"` // Whether sport has outright markets
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Event represents a sports event
type Event struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SportID     uuid.UUID `json:"sport_id" gorm:"type:uuid;not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Country     string    `json:"country" gorm:"not null"`
	League      string    `json:"league" gorm:"not null"`
	HomeTeam    string    `json:"home_team" gorm:"not null"`
	AwayTeam    string    `json:"away_team" gorm:"not null"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null;default:'upcoming'"`
	Result      string    `json:"result"`
	Active      bool      `json:"active" gorm:"not null;default:true"`
	Display     bool      `json:"display" gorm:"not null;default:true"`
	Version     int       `json:"version" gorm:"not null;default:1"`
	// External API tracking
	ExternalID string    `json:"external_id" gorm:"uniqueIndex"` // The Odds API event ID
	SportKey   string    `json:"sport_key" gorm:"index"`         // The Odds API sport key
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// GORM relationships
	Markets []Market `json:"markets,omitempty" gorm:"foreignKey:EventID"`
}

// CreateEventRequest represents a request to create a new event (without version)
type CreateEventRequest struct {
	SportID     uuid.UUID `json:"sport_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Country     string    `json:"country" binding:"required"`
	League      string    `json:"league" binding:"required"`
	HomeTeam    string    `json:"home_team" binding:"required"`
	AwayTeam    string    `json:"away_team" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Status      string    `json:"status" binding:"required"` // "upcoming", "live", "finished", "cancelled"
	Result      string    `json:"result"`                    // "home_win", "away_win", "draw", "cancelled"
	Active      bool      `json:"active"`                    // Whether the event is active for betting
	Display     bool      `json:"display"`                   // Whether the event should be displayed
	// External API tracking
	ExternalID string `json:"external_id"` // The Odds API event ID
	SportKey   string `json:"sport_key"`   // The Odds API sport key
}

// UpdateEventRequest represents a request to update an event (without version)
type UpdateEventRequest struct {
	SportID     uuid.UUID `json:"sport_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Country     string    `json:"country"`
	League      string    `json:"league"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`  // "upcoming", "live", "finished", "cancelled"
	Result      string    `json:"result"`  // "home_win", "away_win", "draw", "cancelled"
	Active      bool      `json:"active"`  // Whether the event is active for betting
	Display     bool      `json:"display"` // Whether the event should be displayed
}

// Market represents a betting market for an event
type Market struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID     uuid.UUID `json:"event_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null;default:'open'"`
	Active      bool      `json:"active" gorm:"not null;default:true"`
	Display     bool      `json:"display" gorm:"not null;default:true"`
	Version     int       `json:"version" gorm:"not null;default:1"`
	// External API tracking
	ExternalKey string    `json:"external_key" gorm:"index"` // The Odds API market key
	LastUpdate  time.Time `json:"last_update"`               // When odds were last updated
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Manual selections field for API responses
	Selections []Selection `json:"selections,omitempty" gorm:"-"`
}

// CreateMarketRequest represents a request to create a new market (without version)
type CreateMarketRequest struct {
	EventID     uuid.UUID `json:"event_id" binding:"required"`
	Name        string    `json:"name" binding:"required"` // "match_winner", "total_goals", "first_scorer", etc.
	Description string    `json:"description"`
	Type        string    `json:"type" binding:"required"`   // "1x2", "over_under", "handicap", "exact_score"
	Status      string    `json:"status" binding:"required"` // "open", "suspended", "closed", "settled"
	Active      bool      `json:"active"`                    // Whether the market is active for betting
	Display     bool      `json:"display"`                   // Whether the market should be displayed
	// External API tracking
	ExternalKey string    `json:"external_key"` // The Odds API market key
	LastUpdate  time.Time `json:"last_update"`  // When odds were last updated
}

// UpdateMarketRequest represents a request to update a market (without version)
type UpdateMarketRequest struct {
	Name        string `json:"name"` // "match_winner", "total_goals", "first_scorer", etc.
	Description string `json:"description"`
	Type        string `json:"type"`    // "1x2", "over_under", "handicap", "exact_score"
	Status      string `json:"status"`  // "open", "suspended", "closed", "settled"
	Active      bool   `json:"active"`  // Whether the market is active for betting
	Display     bool   `json:"display"` // Whether the market should be displayed
}

// Selection represents a betting option within a market
type Selection struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MarketID uuid.UUID `json:"market_id" gorm:"type:uuid;not null"`
	Name     string    `json:"name" gorm:"not null"`
	Status   string    `json:"status" gorm:"not null;default:'active'"`
	Active   bool      `json:"active" gorm:"not null;default:true"`
	Display  bool      `json:"display" gorm:"not null;default:true"`
	Version  int       `json:"version" gorm:"not null;default:1"`
	// External API tracking
	ExternalID string    `json:"external_id" gorm:"uniqueIndex"` // The Odds API selection ID
	Point      float64   `json:"point"`                          // For spreads, totals, etc.
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// GORM relationship for odds
	Odds *Odds `json:"odds,omitempty" gorm:"foreignKey:SelectionID"`
}

// CreateSelectionRequest represents a request to create a new selection (without version)
type CreateSelectionRequest struct {
	MarketID uuid.UUID `json:"market_id" binding:"required"`
	Name     string    `json:"name" binding:"required"`   // "Home Win", "Away Win", "Draw", "Over 2.5", etc.
	Status   string    `json:"status" binding:"required"` // "active", "suspended", "won", "lost"
	Active   bool      `json:"active"`                    // Whether the selection is active for betting
	Display  bool      `json:"display"`                   // Whether the selection should be displayed
	// External API tracking
	ExternalID string  `json:"external_id"` // The Odds API selection ID
	Point      float64 `json:"point"`       // For spreads, totals, etc.
}

// UpdateSelectionRequest represents a request to update a selection (without version)
type UpdateSelectionRequest struct {
	Name    string `json:"name"`    // "Home Win", "Away Win", "Draw", "Over 2.5", etc.
	Status  string `json:"status"`  // "active", "suspended", "won", "lost"
	Active  bool   `json:"active"`  // Whether the selection is active for betting
	Display bool   `json:"display"` // Whether the selection should be displayed
}

// Odds represents odds for a selection
type Odds struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SelectionID uuid.UUID `json:"selection_id" gorm:"type:uuid;not null;uniqueIndex"`
	Odds        float64   `json:"odds" gorm:"not null"`
	Version     int       `json:"version" gorm:"not null;default:1"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// GORM relationships
	// Note: Selection relationship removed to avoid circular reference
}

// CreateOddsRequest represents a request to create odds (without version)
type CreateOddsRequest struct {
	SelectionID uuid.UUID `json:"selection_id" binding:"required"`
	Odds        float64   `json:"odds" binding:"required,gt=0"`
}

// UpdateOddsRequest represents a request to update odds (without version)
type UpdateOddsRequest struct {
	Odds float64 `json:"odds" binding:"required,gt=0"`
}

// Result represents a result for a selection
type Result struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SelectionID uuid.UUID `json:"selection_id" db:"selection_id"`
	Result      string    `json:"result" db:"result"`   // "won", "lost", "void", "pending"
	Version     int       `json:"version" db:"version"` // Version number for audit trail
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateResultRequest represents a request to create a result (without version)
type CreateResultRequest struct {
	SelectionID uuid.UUID `json:"selection_id" binding:"required"`
	Result      string    `json:"result" binding:"required"` // "won", "lost", "void", "pending"
}

// UpdateResultRequest represents a request to update a result (without version)
type UpdateResultRequest struct {
	Result string `json:"result" binding:"required"` // "won", "lost", "void", "pending"
}

// Bet represents a user's bet
type Bet struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	BetType      string     `json:"bet_type" db:"bet_type"` // "single", "double", "treble", "accumulator"
	Amount       float64    `json:"amount" db:"amount"`
	TotalOdds    float64    `json:"total_odds" db:"total_odds"`
	Status       string     `json:"status" db:"status"` // "pending", "won", "lost", "void", "cancelled", "unsettle"
	PotentialWin float64    `json:"potential_win" db:"potential_win"`
	ActualWin    *float64   `json:"actual_win" db:"actual_win"` // Actual amount won (NULL if not settled yet)
	SettledAt    *time.Time `json:"settled_at" db:"settled_at"` // Timestamp when the bet was settled (NULL if not settled yet)
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// BetSelection represents a selection within a bet (for accumulator bets)
type BetSelection struct {
	ID          uuid.UUID `json:"id" db:"id"`
	BetID       uuid.UUID `json:"bet_id" db:"bet_id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	MarketID    uuid.UUID `json:"market_id" db:"market_id"`
	SelectionID uuid.UUID `json:"selection_id" db:"selection_id"`
	Odds        float64   `json:"odds" db:"odds"`
	Result      *string   `json:"result" db:"result"` // Result of the selection: WIN, LOSE, or VOID
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
	EventID            uuid.UUID `json:"event_id" binding:"required"`
	MarketID           uuid.UUID `json:"market_id" binding:"required"`
	SelectionID        uuid.UUID `json:"selection_id" binding:"required"`
	Amount             float64   `json:"amount" binding:"required,gt=0"`
	ExpectedOdds       float64   `json:"expected_odds" binding:"required,gt=0"` // Odds from UI for price validation
	AcceptPriceChanges bool      `json:"accept_price_changes"`                  // Allow user to accept price changes
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

// BeforeCreate hooks to set UUID if not provided
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (m *Market) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (s *Selection) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (o *Odds) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (r *Result) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (s *Sport) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
