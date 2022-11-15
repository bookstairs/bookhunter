package lanzou

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bookstairs/bookhunter/internal/client"
)

func TestParseLanzouUrl(t *testing.T) {
	type args struct {
		url string
		pwd string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test lanzoue.com",
			args: args{
				url: "https://tianlangbooks.lanzoue.com/iI0lc0fj3mpa",
				pwd: "tlsw",
			},
		},
		{
			name: "test lanzoui",
			args: args{
				url: "https://tianlangbooks.lanzoui.com/i9CTIws9s6d",
				pwd: "tlsw",
			},
		},
		{
			name: "test lanzouf.com",
			args: args{
				url: "https://tianlangbooks.lanzouf.com/ic7HY05ejl2h",
				pwd: "tlsw",
			},
		},
		{
			name: "test lanzoup.com",
			args: args{
				url: "https://tianlangbooks.lanzoup.com/i4q4Chcm2cf",
				pwd: "tlsw",
			},
		},
		{
			name: "test lanzouy.com",
			args: args{
				url: "https://fast8.lanzouy.com/ibZCg0b8tibi",
				pwd: "",
			},
		},
		{
			name: "test sobook lanzouy.com",
			args: args{
				url: "https://sobooks.lanzoum.com/b03phl3te",
				pwd: "htuj",
			},
		},
		{
			name: "test sobook lanzouy.com 2",
			args: args{
				url: "https://sobooks.lanzoum.com/ihOex0fiodri",
				pwd: "",
			},
		},
	}

	drive, _ := NewDrive(&client.Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := drive.ResolveShareURL(tt.args.url, tt.args.pwd)
			assert.NoError(t, err, "Failed to resolve link")
			assert.Equal(t, int64(200), response.Code, "Failed to resolve link: "+response.Msg)
			assert.NotEmpty(t, response.Data)
			assert.NotNil(t, response.Data)
			assert.NotNil(t, response.Data.URL)
			assert.NotNil(t, response.Data.Name)
		})
	}
}
