package spider

import "testing"

func TestGenerateUrl(t *testing.T) {
	type args struct {
		base  string
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "trim \"/\" suffix",
			args: args{
				base:  "http://www.baidu.com/",
				paths: []string{"/path1", "/path2"},
			},
			want: "http://www.baidu.com/path1/path2",
		},
		{
			name: "add schema protocol",
			args: args{
				base:  "www.baidu.com",
				paths: []string{"/path1", "/path2"},
			},
			want: "http://www.baidu.com/path1/path2",
		},
		{
			name: "trim request path suffix",
			args: args{
				base:  "http://www.baidu.com/",
				paths: []string{"/path1/", "/path2/"},
			},
			want: "http://www.baidu.com/path1/path2",
		},
		{
			name: "add prefix to request path",
			args: args{
				base:  "http://www.baidu.com/",
				paths: []string{"path1", "path2"},
			},
			want: "http://www.baidu.com/path1/path2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateUrl(tt.args.base, tt.args.paths...); got != tt.want {
				t.Errorf("GenerateUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
