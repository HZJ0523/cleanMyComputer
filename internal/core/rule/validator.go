package rule

import (
	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateRule(rule *models.CleanRule) error {
	return rule.Validate()
}

func (v *Validator) ValidateRules(rules []*models.CleanRule) []error {
	var errs []error
	for _, rule := range rules {
		if err := v.ValidateRule(rule); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
