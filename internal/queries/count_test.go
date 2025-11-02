package queries

import "testing"

func TestGetStatesCount(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test 1",
			args: args{
				id: 1,
			},
			want:    7,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStatesCount(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatesCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStatesCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
