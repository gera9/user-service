package utils

import (
	"reflect"
	"testing"
)

func TestUpdateQueryBuilder(t *testing.T) {
	type args struct {
		table string
		o     any
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
		want2 []any
	}{
		{
			name: "Test 1",
			args: args{
				table: "users",
				o: struct {
					Username string
					Email    string
					Password string
				}{
					Username: "test",
					Email:    "test",
					Password: "test",
				},
			},
			want:  "UPDATE users SET Username = $1, Email = $2, Password = $3",
			want1: 4,
			want2: []any{"test", "test", "test"},
		},
		{
			name: "Test 2",
			args: args{
				table: "users",
				o: struct {
					Username string
					Email    string
					Password string
				}{
					Username: "test",
					Email:    "",
					Password: "test",
				},
			},
			want:  "UPDATE users SET Username = $1, Password = $2",
			want1: 3,
			want2: []any{"test", "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := UpdateQueryBuilder(tt.args.table, tt.args.o)
			if got != tt.want {
				t.Errorf("UpdateQueryBuilder() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UpdateQueryBuilder() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("UpdateQueryBuilder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
