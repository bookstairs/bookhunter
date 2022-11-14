package aliyun

import "fmt"

// AuthToken will return the token by the given refreshToken.
func (ali *Aliyun) AuthToken() (string, error) {
	// Get the token from the cache.
	token := ali.cachedToken()
	if token != nil {
		return token.AccessToken, nil
	}

	// Refresh the access token by the refresh token.
	resp, err := ali.client.R().
		SetBody(&TokenReq{GrantType: "refresh_token", RefreshToken: ali.refreshToken}).
		SetResult(&TokenResp{}).
		SetResult(&ErrorResp{}).
		Post("https://auth.aliyundrive.com/v2/account/token")
	if err != nil {
		return "", err
	}

	if e, ok := resp.Error().(*ErrorResp); ok {
		return "", fmt.Errorf("token: %s, %s", e.Code, e.Message)
	}
	token = resp.Result().(*TokenResp)
	ali.cacheToken(token)

	return token.AccessToken, nil
}
