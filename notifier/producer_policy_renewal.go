package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gloveboxhq/glovebox-go-code-challenge/notifier/channels/email"
)

const TopicPolicyRenewal = "policy-renewal"

type PolicyRenewalInput struct {
	Recipient   string
	Policy      string
	RenewalDate time.Time
}

type policyRenewalTopicBuilder struct{}

func (policyRenewalTopicBuilder) Topic() string { return TopicPolicyRenewal }

func (policyRenewalTopicBuilder) BuildRequest(_ context.Context, input any) (Request, error) {
	typedInput, ok := input.(PolicyRenewalInput)
	if !ok {
		return Request{}, fmt.Errorf("invalid policy renewal input type")
	}

	if typedInput.Recipient == "" {
		return Request{}, fmt.Errorf("policy renewal requires recipient")
	}

	if typedInput.Policy == "" {
		return Request{}, fmt.Errorf("policy renewal requires policy")
	}

	if typedInput.RenewalDate.IsZero() {
		return Request{}, fmt.Errorf("policy renewal requires renewal date")
	}
	return Request{
		Topic:      TopicPolicyRenewal,
		Recipients: []string{typedInput.Recipient},
		Template:   email.TplPolicyRenewal,
		Vars: map[string]any{
			"policy":      typedInput.Policy,
			"renewalDate": typedInput.RenewalDate.Format(time.DateOnly),
		},
	}, nil
}
