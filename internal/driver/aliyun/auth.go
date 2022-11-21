package aliyun

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/skip2/go-qrcode"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
)

var (
	// Token will be refreshed before this time.
	acceleratedExpirationDuration = 10 * time.Minute
	headerAuthorization           = http.CanonicalHeaderKey("Authorization")
)

// RefreshToken can only be called after the New method.
func (ali *Aliyun) RefreshToken() string {
	return ali.authentication.refreshToken
}

func newAuthentication(c *client.Config, refreshToken string) (*authentication, error) {
	// Create session file.
	path, err := c.ConfigPath()
	if err != nil {
		return nil, err
	}
	sessionFile := filepath.Join(path, "session.db")
	open, err := os.OpenFile(sessionFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Exit(err)
	}
	_ = open.Close()

	cl, err := client.New(c)
	if err != nil {
		return nil, err
	}

	return &authentication{
		Client:       cl,
		sessionFile:  sessionFile,
		refreshToken: refreshToken,
	}, nil
}

type authentication struct {
	*client.Client
	sessionFile  string
	refreshToken string
	tokenCache   *TokenResp
}

func (auth *authentication) authenticationHook() resty.PreRequestHook {
	return func(_ *resty.Client, req *http.Request) error {
		if req.Header.Get("x-empty-content-type") != "" {
			req.Header.Del("x-empty-content-type")
			req.Header.Set("content-type", "")
		}

		// Refresh the token and set it on the header.
		token, err := auth.accessToken()
		if err != nil {
			return err
		}
		req.Header.Set(headerAuthorization, "Bearer "+token)

		return nil
	}
}

func (auth *authentication) Auth() error {
	// Check the given token if it's valid.
	if auth.validateRefreshToken() {
		return nil
	}

	// The given token isn't valid. Try to load the token from file.
	b, err := os.ReadFile(auth.sessionFile)
	if err != nil {
		return err
	}
	auth.refreshToken = string(b)
	if auth.validateRefreshToken() {
		return nil
	}

	// If the token doesn't exist. We will get the refreshToken by QR code.
	if err := auth.login(); err != nil {
		return err
	}
	_ = auth.validateRefreshToken()

	return nil
}

func (auth *authentication) login() error {
	// Set the OAuth2 request into cookies.
	_, err := auth.R().
		SetQueryParams(map[string]string{
			"client_id":     "25dzX3vbYqktVxyX",
			"redirect_uri":  "https://www.aliyundrive.com/sign/callback",
			"response_type": "code",
			"login_type":    "custom",
			"state":         `{"origin":"https://www.aliyundrive.com"}`,
		}).
		Get("https://auth.aliyundrive.com/v2/oauth/authorize")
	if err != nil {
		return err
	}

	// Get the QR code.
	resp, err := auth.R().
		SetQueryParams(map[string]string{
			"appName":     "aliyun_drive",
			"fromSite":    "52",
			"appEntrance": "web",
			"isMobile":    "false",
			"lang":        "zh_CN",
			"_bx-v":       "2.0.31",
			"returnUrl":   "",
			"bizParams":   "",
		}).
		SetResult(&QRCodeResp{}).
		Get("https://passport.aliyundrive.com/newlogin/qrcode/generate.do")
	if err != nil {
		return err
	}
	qr := resp.Result().(*QRCodeResp).Content.Data
	if msg := qr.TitleMsg; msg != "" {
		return errors.New(msg)
	}

	// Print the QR code into console.
	code, _ := qrcode.New(qr.CodeContent, qrcode.Low)
	fmt.Println()
	fmt.Println(code.ToSmallString(false))
	log.Info("Use Aliyun Drive App to scan this QR code.")

	// Wait for the scan result.
	scanned := false
	for {
		// NEW / SCANED / EXPIRED / CANCELED / CONFIRMED
		resp, err := auth.queryQRCode(qr.T, qr.Ck)
		if err != nil {
			return err
		}
		res := resp.Content.Data

		switch res.QrCodeStatus {
		case "NEW":
		case "SCANED":
			if !scanned {
				log.Info("Scan success. Confirm login on your mobile.")
			}
			scanned = true
		case "EXPIRED":
			return fmt.Errorf("the QR code expired")
		case "CANCELED":
			return fmt.Errorf("user canceled login")
		case "CONFIRMED":
			biz := res.BizAction.PdsLoginResult
			if err := auth.confirmLogin(biz.AccessToken); err != nil {
				return err
			}
			auth.refreshToken = biz.RefreshToken
			return nil
		default:
			return fmt.Errorf("%v", res)
		}

		// Sleep one second.
		time.Sleep(time.Second)
	}
}

func (auth *authentication) queryQRCode(t int64, ck string) (*QueryQRCodeResp, error) {
	resp, err := auth.R().
		SetQueryParams(map[string]string{
			"appName":  "aliyun_drive",
			"fromSite": "52",
			"_bx-v":    "2.0.31",
		}).
		SetFormData(map[string]string{
			"t":           strconv.FormatInt(t, 10),
			"ck":          ck,
			"appName":     "aliyun_drive",
			"appEntrance": "web",
			"isMobile":    "false",
			"lang":        "zh_CN",
			"returnUrl":   "",
			"fromSite":    "52",
			"bizParams":   "",
			"navlanguage": "zh-CN",
			"navPlatform": "MacIntel",
		}).
		SetResult(&QueryQRCodeResp{}).
		Post("https://passport.aliyundrive.com/newlogin/qrcode/query.do")
	if err != nil {
		return nil, err
	}

	res := resp.Result().(*QueryQRCodeResp)
	if res.Content.Data.BizExt != "" {
		bs, _ := base64.StdEncoding.DecodeString(res.Content.Data.BizExt)
		_ = json.Unmarshal(bs, &res.Content.Data.BizAction)
	}

	return res, nil
}

func (auth *authentication) confirmLogin(accessToken string) error {
	resp, err := auth.R().
		SetBody(map[string]string{
			"token": accessToken,
		}).
		SetResult(&ConfirmLoginResp{}).
		Post("https://auth.aliyundrive.com/v2/oauth/token_login")
	if err != nil {
		return err
	}
	jump := resp.Result().(*ConfirmLoginResp).Goto

	_, err = auth.R().Get(jump)
	if err != nil {
		return err
	}

	gotoURL, _ := url.Parse(jump)
	if gotoURL != nil && gotoURL.Query().Has("code") {
		code := gotoURL.Query().Get("code")
		_, err := auth.getToken(code)
		if err != nil {
			return err
		} else {
			log.Info("Successfully login.")
		}
	}

	return nil
}

func (auth *authentication) getToken(code string) (*TokenResp, error) {
	resp, err := auth.R().
		SetBody(map[string]string{
			"code":      code,
			"loginType": "normal",
			"deviceId":  "aliyundrive",
		}).
		SetResult(&TokenResp{}).
		Get("https://api.aliyundrive.com/token/get")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*TokenResp), nil
}

func (auth *authentication) validateRefreshToken() bool {
	if auth.refreshToken != "" {
		_, err := auth.accessToken()
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}
	return false
}

// accessToken will return the token by the given refreshToken.
// You can call this method for automatically refreshing the access token.
func (auth *authentication) accessToken() (string, error) {
	// Get the token from the cache or refreshToken.
	if auth.tokenCache == nil || time.Now().Add(acceleratedExpirationDuration).After(auth.tokenCache.ExpireTime) {
		// Refresh the access token by the refresh token.
		resp, err := auth.R().
			SetBody(&TokenReq{GrantType: "refresh_token", RefreshToken: auth.refreshToken}).
			SetResult(&TokenResp{}).
			SetError(&ErrorResp{}).
			Post("https://auth.aliyundrive.com/v2/account/token")
		if err != nil {
			return "", err
		}
		if e, ok := resp.Error().(*ErrorResp); ok {
			return "", fmt.Errorf("token: %s, %s", e.Code, e.Message)
		}

		// Cache the token.
		auth.tokenCache = resp.Result().(*TokenResp)
		auth.refreshToken = auth.tokenCache.RefreshToken

		// Persist the refresh token if it's valid.
		if err := os.WriteFile(auth.sessionFile, []byte(auth.refreshToken), 0o644); err != nil {
			return "", err
		}
	}

	return auth.tokenCache.AccessToken, nil
}
