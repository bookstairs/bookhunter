package telecom

import (
	"net/http"

	"github.com/tickstep/cloudpan189-api/cloudpan"

	"github.com/bookstairs/bookhunter/internal/log"
)

func (t *Telecom) login(username, password string) error {
	appToken, e := cloudpan.AppLogin(username, password)
	if e != nil {
		return e
	}

	// Refresh the cookies.
	if loginUser := cloudpan.RefreshCookieToken(appToken.SessionKey); loginUser != "" {
		t.client.SetCookie(&http.Cookie{
			Name:   "COOKIE_LOGIN_USER",
			Value:  loginUser,
			Domain: "cloud.189.cn",
			Path:   "/",
		})
	}

	log.Info("Login into telecom success.")

	// Save token into the driver client.
	t.appToken = &AppLoginToken{
		RsaPublicKey:            appToken.RsaPublicKey,
		SessionKey:              appToken.SessionKey,
		SessionSecret:           appToken.SessionSecret,
		FamilySessionKey:        appToken.FamilySessionKey,
		FamilySessionSecret:     appToken.FamilySessionSecret,
		AccessToken:             appToken.AccessToken,
		RefreshToken:            appToken.RefreshToken,
		SskAccessTokenExpiresIn: appToken.SskAccessTokenExpiresIn,
		SskAccessToken:          appToken.AccessToken,
	}

	return nil
}
