package finance

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string
type AccountNormal string
type JournalStatus string

const (
	TypeAsset     AccountType = "asset"
	TypeLiability AccountType = "liability"
	TypeEquity    AccountType = "equity"
	TypeRevenue   AccountType = "revenue"
	TypeExpense   AccountType = "expense"

	NormalDebit  AccountNormal = "debit"
	NormalCredit AccountNormal = "credit"

	StatusDraft    JournalStatus = "draft"
	StatusPosted   JournalStatus = "posted"
	StatusReversed JournalStatus = "reversed"
)

type ChartOfAccount struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Code          string
	Name          string
	Type          AccountType
	NormalBalance AccountNormal
}



type JournalEntry struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	JournalNo   string
	Date        time.Time
	Description string
	Status      JournalStatus
	SourceType  string
	SourceID    *uuid.UUID
	TotalDebit  float64
	TotalCredit float64
	Lines       []JournalLine
}

type JournalLine struct {
	ID             uuid.UUID
	JournalEntryID uuid.UUID
	COAID          uuid.UUID
	Description    string
	Debit          float64
	Credit         float64
}

func (j *JournalEntry) AddLine(coaID uuid.UUID, desc string, debit, credit float64) {
	j.Lines = append(j.Lines, JournalLine{
		ID:             uuid.New(),
		JournalEntryID: j.ID,
		COAID:          coaID,
		Description:    desc,
		Debit:          debit,
		Credit:         credit,
	})
	j.TotalDebit += debit
	j.TotalCredit += credit
}

func (j *JournalEntry) IsBalanced() bool {
	return j.TotalDebit == j.TotalCredit
}
