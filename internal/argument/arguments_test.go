package argument

import "testing"

func TestHideSensitive(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Content which length is [0, 6)", args: args{content: "12345"}, want: "1***5"},
		{name: "Content which length is [6, 9)", args: args{content: "12345678"}, want: "12****78"},
		{name: "Content which length is [9, +∞)", args: args{content: "1234567890"}, want: "123****890"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HideSensitive(tt.args.content); got != tt.want {
				t.Errorf("HideSensitive() = %v, want %v", got, tt.want)
			}
		})
	}
}
