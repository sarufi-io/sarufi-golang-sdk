package sarufi

import "testing"

func TestValidationError_Error(t *testing.T) {
	t.Parallel()
	type fields struct {
		Detail []ValidationErrorDetail
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal validation error",
			fields: fields{
				Detail: []ValidationErrorDetail{
					{
						Loc:  []string{"username", "body"},
						Msg:  "field required",
						Type: "value_error.missing",
					},
				},
			},
			want: "validation error: location: [username,body], message: field required, error type: value_error.missing",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &ValidationError{
				Detail: tt.fields.Detail,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
