package telecom

type (
	AppLoginParams struct {
		CaptchaToken string
		Lt           string
		ReturnURL    string
		ParamID      string
		ReqID        string
		jRsaKey      string
	}

	AppLoginResult struct {
		Result int    `json:"result"`
		Msg    string `json:"msg"`
		ToURL  string `json:"toUrl"`
	}

	AppSessionResp struct {
		ResCode             int    `json:"res_code"`
		ResMessage          string `json:"res_message"`
		AccessToken         string `json:"accessToken"`
		FamilySessionKey    string `json:"familySessionKey"`
		FamilySessionSecret string `json:"familySessionSecret"`
		GetFileDiffSpan     int    `json:"getFileDiffSpan"`
		GetUserInfoSpan     int    `json:"getUserInfoSpan"`
		IsSaveName          string `json:"isSaveName"`
		KeepAlive           int    `json:"keepAlive"`
		LoginName           string `json:"loginName"`
		RefreshToken        string `json:"refreshToken"`
		SessionKey          string `json:"sessionKey"`
		SessionSecret       string `json:"sessionSecret"`
	}

	AccessTokenResp struct {
		// The expiry time for token, default 30 days.
		ExpiresIn   int64  `json:"expiresIn"`
		AccessToken string `json:"accessToken"`
	}

	AppLoginToken struct {
		SessionKey              string `json:"sessionKey"`
		SessionSecret           string `json:"sessionSecret"`
		FamilySessionKey        string `json:"familySessionKey"`
		FamilySessionSecret     string `json:"familySessionSecret"`
		AccessToken             string `json:"accessToken"`
		RefreshToken            string `json:"refreshToken"`
		SskAccessToken          string `json:"sskAccessToken"`
		SskAccessTokenExpiresIn int64  `json:"sskAccessTokenExpiresIn"`
		RsaPublicKey            string `json:"rsaPublicKey"`
	}
)
