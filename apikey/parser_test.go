package apikey

import "testing"

func TestParseCases(t *testing.T) {
	validSecret := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	tests := []struct {
		name   string
		raw    string
		reason ValidationFailureReason
	}{
		{name: "empty", raw: "", reason: FailureMissing},
		{name: "whitespace", raw: "   ", reason: FailureMissing},
		{name: "wrong prefix", raw: "bk_crm_7F3K9Q2M8N4P6R1T_" + validSecret, reason: FailureMalformed},
		{name: "missing app", raw: "ak__7F3K9Q2M8N4P6R1T_" + validSecret, reason: FailureMalformed},
		{name: "long app", raw: "ak_abcd_7F3K9Q2M8N4P6R1T_" + validSecret, reason: FailureMalformed},
		{name: "uppercase app", raw: "ak_A1_7F3K9Q2M8N4P6R1T_" + validSecret, reason: FailureMalformed},
		{name: "short public key", raw: "ak_crm_ABC_" + validSecret, reason: FailureMalformed},
		{name: "long public key", raw: "ak_crm_7F3K9Q2M8N4P6R1TX_" + validSecret, reason: FailureMalformed},
		{name: "invalid public key letters", raw: "ak_crm_7F3K9Q2M8N4P6R1O_" + validSecret, reason: FailureMalformed},
		{name: "missing secret", raw: "ak_crm_7F3K9Q2M8N4P6R1T_", reason: FailureMalformed},
		{name: "invalid secret characters", raw: "ak_crm_7F3K9Q2M8N4P6R1T_***", reason: FailureMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, reason := Parse(tt.raw)
			if reason != tt.reason {
				t.Fatalf("reason = %q, want %q", reason, tt.reason)
			}
		})
	}
}

func TestParseSecretWithUnderscore(t *testing.T) {
	raw := "ak_crm_7F3K9Q2M8N4P6R1T_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA_"
	parsed, reason := Parse(raw)
	if reason != FailureNone {
		t.Fatalf("reason = %q", reason)
	}
	if parsed.Secret == "" || parsed.App != "crm" || parsed.PublicKey != "7F3K9Q2M8N4P6R1T" {
		t.Fatalf("unexpected parsed key: %#v", parsed)
	}
}
