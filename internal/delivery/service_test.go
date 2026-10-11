package delivery

import "testing"

func TestService_Calculate(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []struct {
		name        string
		orderAmount int64
		want        int64
	}{
		{
			name:        "paid delivery",
			orderAmount: 120000,
			want:        19900,
		},
		{
			name:        "free delivery",
			orderAmount: 200000,
			want:        0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := service.Calculate(tt.orderAmount)

			if got != tt.want {
				t.Fatalf("Calculate(%d) = %d, want %d", tt.orderAmount, got, tt.want)
			}
		})
	}
}
