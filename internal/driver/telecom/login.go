package telecom

import (
	"encoding/hex"
	"errors"
	"net/url"
	"regexp"

	"github.com/bookstairs/bookhunter/internal/crypto"
	"github.com/bookstairs/bookhunter/internal/log"
)

func (t *Telecom) login(username, password string) error {
	// Query the rsa key.
	params, err := t.loginParams()
	if err != nil {
		return err
	}

	// Perform the login by the app API.
	app, err := t.appLogin(params, username, password)
	if err != nil {
		return err
	}

	// Access the login session.
	session, err := t.createSession(app.ToURL)
	if err != nil {
		return err
	}

	// Acquire the Ssk token.
	token, err := t.createToken(session.SessionKey)
	if err != nil {
		return err
	}

	// Refresh the cookies.
	err = t.refreshCookies(session.SessionKey)
	if err != nil {
		return err
	}

	// Save token into the driver client.
	appToken := &AppLoginToken{
		RsaPublicKey:            params.jRsaKey,
		SessionKey:              session.SessionKey,
		SessionSecret:           session.SessionSecret,
		FamilySessionKey:        session.FamilySessionKey,
		FamilySessionSecret:     session.FamilySessionSecret,
		AccessToken:             session.AccessToken,
		RefreshToken:            session.RefreshToken,
		SskAccessTokenExpiresIn: token.ExpiresIn,
		SskAccessToken:          token.AccessToken,
	}
	t.appToken = appToken

	return nil
}

func (t *Telecom) loginParams() (*AppLoginParams, error) {
	resp, err := t.client.R().
		SetQueryParam("appId", "8025431004").
		SetQueryParam("clientType", "10020").
		SetQueryParam("returnURL", "https://m.cloud.189.cn/zhuanti/2020/loginErrorPc/index.html").
		SetQueryParam("timeStamp", timeStamp()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Get(webPrefix + "/unifyLoginForPC.action")

	if err != nil {
		log.Debugf("login redirectURL occurs error: %s", err)
		return nil, err
	}
	content := resp.String()

	re := regexp.MustCompile("captchaToken' value='(.+?)'")
	captchaToken := re.FindStringSubmatch(content)[1]

	re = regexp.MustCompile("lt = \"(.+?)\"")
	lt := re.FindStringSubmatch(content)[1]

	re = regexp.MustCompile("returnUrl = '(.+?)'")
	returnURL := re.FindStringSubmatch(content)[1]

	re = regexp.MustCompile("paramId = \"(.+?)\"")
	paramID := re.FindStringSubmatch(content)[1]

	re = regexp.MustCompile("reqId = \"(.+?)\"")
	reqID := re.FindStringSubmatch(content)[1]

	// RSA key should be wrapped with the comments.
	re = regexp.MustCompile("j_rsaKey\" value=\"(.+?)\"")
	jRsaKey := "-----BEGIN PUBLIC KEY-----\n" + re.FindStringSubmatch(content)[1] + "\n-----END PUBLIC KEY-----"

	return &AppLoginParams{
		CaptchaToken: captchaToken,
		Lt:           lt,
		ReturnURL:    returnURL,
		ParamID:      paramID,
		ReqID:        reqID,
		jRsaKey:      jRsaKey,
	}, nil
}

func (t *Telecom) appLogin(params *AppLoginParams, username, password string) (*AppLoginResult, error) {
	rsaKey := params.jRsaKey
	rsaUsername, err := crypto.RsaEncrypt([]byte(rsaKey), []byte(username))
	if err != nil {
		return nil, err
	}
	rsaPassword, err := crypto.RsaEncrypt([]byte(rsaKey), []byte(password))
	if err != nil {
		return nil, err
	}

	// Start to perform login.
	resp, err := t.client.R().
		SetHeaders(map[string]string{
			"Content-Type":     "application/x-www-form-urlencoded",
			"Referer":          authPrefix + "/unifyAccountLogin.do",
			"Cookie":           "LT=" + params.Lt,
			"X-Requested-With": "XMLHttpRequest",
			"REQID":            params.ReqID,
			"lt":               params.Lt,
		}).
		SetFormData(map[string]string{
			"appKey":       "8025431004",
			"accountType":  "02",
			"userName":     "{RSA}" + hex.EncodeToString(crypto.Base64Encode(rsaUsername)),
			"password":     "{RSA}" + hex.EncodeToString(crypto.Base64Encode(rsaPassword)),
			"validateCode": "",
			"captchaToken": params.CaptchaToken,
			"returnUrl":    params.ReturnURL,
			"mailSuffix":   "@189.cn",
			"dynamicCheck": "FALSE",
			"clientType":   "10020",
			"cb_SaveName":  "1",
			"isOauth2":     "false",
			"state":        "",
			"paramId":      params.ParamID,
		}).
		SetResult(&AppLoginResult{}).
		Post(authPrefix + "/loginSubmit.do")
	if err != nil {
		return nil, err
	}

	// Check login result.
	res := resp.Result().(*AppLoginResult)
	if res.Result != 0 || res.ToURL == "" {
		return nil, errors.New("login failed in telecom disk")
	}

	return res, nil
}

func (t *Telecom) createSession(jumpURL string) (*AppSessionResp, error) {
	resp, err := t.client.R().
		SetHeader("Accept", "application/json;charset=UTF-8").
		SetQueryParams(map[string]string{
			"clientType":  "TELEMAC",
			"version":     "1.0.0",
			"channelId":   "web_cloud.189.cn",
			"redirectURL": url.QueryEscape(jumpURL),
		}).
		SetResult(&AppSessionResp{}).
		Get(apiPrefix + "/getSessionForPC.action")
	if err != nil {
		return nil, err
	}

	// Check the session result.
	res := resp.Result().(*AppSessionResp)
	if res.ResCode != 0 {
		return nil, errors.New("failed to acquire session")
	}

	return res, nil
}

func (t *Telecom) createToken(sessionKey string) (*AccessTokenResp, error) {
	timestamp := timeStamp()
	signParams := map[string]string{
		"Timestamp":  timestamp,
		"sessionKey": sessionKey,
		"AppKey":     "601102120",
	}
	resp, err := t.client.R().
		SetQueryParam("sessionKey", sessionKey).
		SetHeaders(map[string]string{
			"AppKey":    "601102120",
			"Signature": crypto.SignatureOfMd5(signParams),
			"Sign-Type": "1",
			"Accept":    "application/json",
			"Timestamp": timestamp,
		}).
		SetResult(&AccessTokenResp{}).
		Get(apiPrefix + "/open/oauth2/getAccessTokenBySsKey.action")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*AccessTokenResp), nil
}

func (t *Telecom) refreshCookies(sessionKey string) error {
	_, err := t.client.R().
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7").
		SetQueryParams(map[string]string{
			"sessionKey":  sessionKey,
			"redirectUrl": "main.action%%23recycle",
		}).
		Get(webPrefix + "/ssoLogin.action")

	return err
}
