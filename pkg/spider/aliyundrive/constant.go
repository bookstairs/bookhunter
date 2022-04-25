package aliyundrive

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

const (
	Authorization = "authorization"
	ContentType   = "content-type"
	UserAgent     = "User-Agent"

	ContentTypeJSON = "application/json"

	Bearer = "bearer"

	BaseURL = "https://auth.aliyundrive.com"
	ApiURL  = "https://api.aliyundrive.com"

	AccessTokenPrefix  = "at:"
	RefreshTokenPrefix = "rt:"

	xEmptyContentType = "x-empty-content-type"
	xShareToken       = "x-share-token"

	V2AccountToken                 = BaseURL + "/v2/account/token"
	V2ShareLinkGetShareToken       = BaseURL + "/v2/share_link/get_share_token"
	V2FileGetShareLinkDownloadUrl  = ApiURL + "/v2/file/get_share_link_download_url"
	V3FileList                     = ApiURL + "/adrive/v3/file/list"
	V2ShareLinkGetShareByAnonymous = BaseURL + "/adrive/v2/share_link/get_share_by_anonymous"
)

type AliYunDrive struct {
	Client       *resty.Client
	RefreshToken string
	Cache        map[string]string
}

func HcHook(_ *resty.Client, req *http.Request) error {
	if req.Header.Get(xEmptyContentType) != "" {
		req.Header.Del(xEmptyContentType)
		req.Header.Set(ContentType, "")
	}
	return nil
}
