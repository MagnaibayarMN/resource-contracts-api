package queries

import "testing"

func TestGetAnnotationsCount(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "GetAnnotationsCount",
			want:    162,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAnnotationsCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnnotationsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAnnotationsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
