package validation

type OutcomeValidator struct {
	SuccessCriteria func(string) bool // A function to validate if the outcome meets success criteria
}

func NewOutcomeValidator(criteria func(string) bool) *OutcomeValidator {
	return &OutcomeValidator{SuccessCriteria: criteria}
}

func (ov *OutcomeValidator) ValidateOutcome(state string) bool {
	return ov.SuccessCriteria(state)
}
