package aliyundrive

import (
	"encoding/json"
	"fmt"

	"github.com/bibliolater/bookhunter/pkg/log"
)

func (ali AliYunDrive) GetAuthorizationToken() string {
	token, err := ali.getToken()
	if err != nil {
		log.Fatal(err)
	}
	return token.AccessToken
}

func (ali AliYunDrive) getToken() (*TokenResponse, error) {
	at, exist := ali.Cache[AccessTokenPrefix+ali.RefreshToken]
	if exist {
		resp := &TokenResponse{}
		if err := json.Unmarshal([]byte(at), &resp); err != nil {
			return nil, err
		}
		return resp, nil
	}

	rt, exist := ali.Cache[RefreshTokenPrefix+ali.RefreshToken]
	if exist {
		ali.RefreshToken = rt
	}
	resp, err := ali.Client.R().
		SetHeader(ContentType, ContentTypeJSON).
		SetBody(TokenRequest{GrantType: "refresh_token", RefreshToken: ali.RefreshToken}).
		SetResult(TokenResponse{}).
		SetError(ErrorResponse{}).
		Post(V2AccountToken)
	if err != nil {
		return nil, err
	}
	if e, ok := resp.Error().(*ErrorResponse); ok {
		return nil, fmt.Errorf("token: %s, %s", e.Code, e.Message)
	}
	res := resp.Result().(*TokenResponse)
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	ali.Cache[AccessTokenPrefix+ali.RefreshToken] = string(marshal)
	ali.Cache[RefreshTokenPrefix+ali.RefreshToken] = res.RefreshToken
	return res, nil
}
