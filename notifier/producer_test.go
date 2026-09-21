package main

import (
	"context"
	"testing"
	"time"

	"github.com/gloveboxhq/glovebox-go-code-challenge/notifier/channels/email"
	"github.com/gloveboxhq/glovebox-go-code-challenge/notifier/channels/email/mockemail"
)

func TestNotifyDocumentUpload(t *testing.T) {
	mail := mockemail.NewClient()
	producer := NewProducer(mail)

	err := producer.NotifyTopic(context.Background(), TopicDocumentUpload, DocumentUploadInput{
		Recipient: "user@example.com",
		Document:  "policy.pdf",
	})
	if err != nil {
		t.Fatalf("notify document upload: %v", err)
	}

	logs := mail.SendLogs()
	if logs.IsEmpty() {
		t.Fatal("expected send log")
	}
	last := logs.Last()
	if last.To != "user@example.com" {
		t.Fatalf("expected recipient %q, got %q", "user@example.com", last.To)
	}
	if last.Tpl != email.TplDocumentUpload {
		t.Fatalf("expected template %q, got %q", email.TplDocumentUpload, last.Tpl)
	}
}

func TestNotifyUnknownTopic(t *testing.T) {
	err := NewProducer(mockemail.NewClient()).NotifyTopic(context.Background(), "unknown", nil)
	if err == nil {
		t.Fatal("expected unknown topic error")
	}
}

func TestNotifyPolicyRenewal(t *testing.T) {
	mail := mockemail.NewClient()
	producer := NewProducer(mail)
	renewalDate := time.Now().Add(30 * 24 * time.Hour)

	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "invalid input type",
			input: "invalid",
		},
		{
			name: "missing recipient",
			input: PolicyRenewalInput{
				Policy:      "HO-4729183",
				RenewalDate: renewalDate,
			},
		},
		{
			name: "missing policy number",
			input: PolicyRenewalInput{
				Recipient:   "user@example.com",
				RenewalDate: renewalDate,
			},
		},
		{
			name: "missing renewal date",
			input: PolicyRenewalInput{
				Recipient: "user@example.com",
				Policy:    "HO-4729183",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := producer.NotifyTopic(context.Background(), TopicPolicyRenewal, tt.input)
			if err == nil {
				t.Fatalf("expected error for test case %q", tt.name)
			}
		})
	}

	t.Run("valid input", func(t *testing.T) {
		mail.FlushSendLogs()

		err := producer.NotifyTopic(context.Background(), TopicPolicyRenewal, PolicyRenewalInput{
			Recipient:   "user@example.com",
			Policy:      "HO-4729183",
			RenewalDate: renewalDate,
		})

		if err != nil {
			t.Fatalf("notify policy renewal: %v", err)
		}

		logs := mail.SendLogs()
		if logs.IsEmpty() {
			t.Fatal("expected send log")
		}

		last := logs.Last()
		if last.Tpl != email.TplPolicyRenewal {
			t.Fatalf(
				"expected template %q, got %q",
				email.TplPolicyRenewal,
				last.Tpl,
			)
		}
	})

}
