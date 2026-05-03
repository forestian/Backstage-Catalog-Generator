package normalize

import "testing"

func TestEntityName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"payment-api", "payment-api"},
		{"Payment API", "payment-api"},
		{"worker_service", "worker-service"},
		{"My Service 123", "my-service-123"},
		{"--leading-dashes--", "leading-dashes"},
		{"UPPERCASE", "uppercase"},
		{"special!@#chars", "specialchars"},
	}
	for _, c := range cases {
		got := EntityName(c.in)
		if got != c.want {
			t.Errorf("EntityName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
