package aliyundrive

import (
	"time"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	DefaultSboxDriveId string    `json:"default_sbox_drive_id"`
	Role               string    `json:"role"`
	DeviceId           string    `json:"device_id"`
	UserName           string    `json:"user_name"`
	NeedLink           bool      `json:"need_link"`
	ExpireTime         time.Time `json:"expire_time"`
	PinSetup           bool      `json:"pin_setup"`
	NeedRpVerify       bool      `json:"need_rp_verify"`
	Avatar             string    `json:"avatar"`
	TokenType          string    `json:"token_type"`
	AccessToken        string    `json:"access_token"`
	DefaultDriveId     string    `json:"default_drive_id"`
	DomainId           string    `json:"domain_id"`
	RefreshToken       string    `json:"refresh_token"`
	IsFirstLogin       bool      `json:"is_first_login"`
	UserId             string    `json:"user_id"`
	NickName           string    `json:"nick_name"`
	State              string    `json:"state"`
	ExpiresIn          int       `json:"expires_in"`
	Status             string    `json:"status"`
}

type GetShareInfoRequest struct {
	ShareId string `json:"share_id"`
}

type GetShareInfoResponse struct {
	Avatar             string          `json:"avatar"`
	CreatorId          string          `json:"creator_id"`
	CreatorName        string          `json:"creator_name"`
	CreatorPhone       string          `json:"creator_phone"`
	Expiration         string          `json:"expiration"`
	UpdatedAt          string          `json:"updated_at"`
	ShareName          string          `json:"share_name"`
	FileCount          int             `json:"file_count"`
	FileInfos          []ShareItemInfo `json:"file_infos"`
	Vip                string          `json:"vip"`
	DisplayName        string          `json:"display_name"`
	IsFollowingCreator bool            `json:"is_following_creator"`
}

type ShareItemInfo struct {
	Category      string `json:"category"`
	FileExtension string `json:"file_extension"`
	FileId        string `json:"file_id"`
	Thumbnail     string `json:"thumbnail"`
	FileType      string `json:"type"`
}

type GetShareLinkDownloadUrlRequest struct {
	ShareId   string `json:"share_id"`
	FileId    string `json:"file_id"`
	ExpireSec int    `json:"expire_sec"`
}

type GetShareLinkDownloadUrlResponse struct {
	DownloadUrl string `json:"download_url"`
	Url         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
}

type GetShareTokenRequest struct {
	ShareId  string `json:"share_id"`
	SharePwd string `json:"share_pwd"`
}

type GetShareTokenResponse struct {
	ShareToken string `json:"share_token"`
	ExpireTime string `json:"expire_time"`
	ExpiresIn  int    `json:"expires_in"`
}

type GetShareFileListRequest struct {
	ShareId               string `json:"share_id"`
	Starred               bool   `json:"starred"`
	All                   bool   `json:"all"`
	Category              string `json:"category"`
	Fields                string `json:"fields"`
	ImageThumbnailProcess string `json:"image_thumbnail_process"`
	Limit                 int    `json:"limit"`
	Marker                string `json:"marker"`
	OrderBy               string `json:"order_by"`
	OrderDirection        string `json:"order_direction"`
	ParentFileId          string `json:"parent_file_id"`
	Status                string `json:"status"`
	FileType              string `json:"type"`
	UrlExpireSec          int    `json:"url_expire_sec"`
	VideoThumbnailProcess string `json:"video_thumbnail_process"`
}

type GetShareFileListResponse struct {
	Items      []*BaseShareFile `json:"items"`
	NextMarker string           `json:"next_marker"`
}

type BaseShareFile struct {
	ShareId       string   `json:"share_id"`
	Name          string   `json:"name"`
	Size          int      `json:"size"`
	Creator       string   `json:"creator"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	DownloadUrl   int      `json:"download_url"`
	Url           int      `json:"url"`
	FileExtension string   `json:"file_extension"`
	FileId        string   `json:"file_id"`
	Thumbnail     string   `json:"thumbnail"`
	ParentFileId  string   `json:"parent_file_id"`
	FileType      string   `json:"type"`
	UpdatedAt     string   `json:"updated_at"`
	CreatedAt     string   `json:"created_at"`
	Selected      string   `json:"selected"`
	MimeExtension string   `json:"mime_extension"`
	MimeType      string   `json:"mime_type"`
	PunishFlag    int      `json:"punish_flag"`
	ActionList    []string `json:"action_list"`
	DriveId       string   `json:"drive_id"`
	DomainId      string   `json:"domain_id"`
	RevisionId    string   `json:"revision_id"`
}

type FileListParam struct {
	shareToken   string
	shareId      string
	parentFileId string
	marker       string
}
