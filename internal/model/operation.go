package model

type Mode string

const (
	DEBIT  = Mode("DEBIT")
	CREDIT = Mode("CREDIT")
)

type OperationType struct {
	ID          int
	Description string
	Mode        Mode
}

func (o OperationType) ValidOperation(amount float64) bool {
	var ok bool
	switch o.Mode {
	case DEBIT:
		ok = amount < 0
	case CREDIT:
		ok = amount > 0
	}

	return ok
}
