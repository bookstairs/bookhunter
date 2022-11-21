package aliyun

import "time"

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type QRCodeResp struct {
	Content struct {
		Data struct {
			TitleMsg    string `json:"title_msg"`
			T           int64  `json:"t"`
			CodeContent string `json:"codeContent"`
			Ck          string `json:"ck"`
			ResultCode  int    `json:"resultCode"`
		} `json:"data"`
	} `json:"content"`
}

type QueryQRCodeResp struct {
	Content struct {
		Data struct {
			QrCodeStatus         string `json:"qrCodeStatus"`
			ResultCode           int    `json:"resultCode"`
			LoginResult          string `json:"loginResult"`
			LoginSucResultAction string `json:"loginSucResultAction"`
			BizAction            struct {
				PdsLoginResult struct {
					Role           string        `json:"role"`
					IsFirstLogin   bool          `json:"isFirstLogin"`
					NeedLink       bool          `json:"needLink"`
					LoginType      string        `json:"loginType"`
					NickName       string        `json:"nickName"`
					NeedRpVerify   bool          `json:"needRpVerify"`
					Avatar         string        `json:"avatar"`
					AccessToken    string        `json:"accessToken"`
					UserName       string        `json:"userName"`
					UserID         string        `json:"userId"`
					DefaultDriveID string        `json:"defaultDriveId"`
					ExistLink      []interface{} `json:"existLink"`
					ExpiresIn      int           `json:"expiresIn"`
					ExpireTime     time.Time     `json:"expireTime"`
					RequestID      string        `json:"requestId"`
					DataPinSetup   bool          `json:"dataPinSetup"`
					State          string        `json:"state"`
					TokenType      string        `json:"tokenType"`
					DataPinSaved   bool          `json:"dataPinSaved"`
					RefreshToken   string        `json:"refreshToken"`
					Status         string        `json:"status"`
				} `json:"pds_login_result"`
			} `json:"-"`
			St          string `json:"st"`
			LoginType   string `json:"loginType"`
			BizExt      string `json:"bizExt"`
			LoginScene  string `json:"loginScene"`
			AppEntrance string `json:"appEntrance"`
			Smartlock   bool   `json:"smartlock"`
		} `json:"data"`
	} `json:"content"`
}

type ConfirmLoginResp struct {
	Goto string `json:"goto"`
}

type TokenReq struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResp struct {
	DefaultSboxDriveID string    `json:"default_sbox_drive_id"`
	Role               string    `json:"role"`
	DeviceID           string    `json:"device_id"`
	UserName           string    `json:"user_name"`
	NeedLink           bool      `json:"need_link"`
	ExpireTime         time.Time `json:"expire_time"`
	PinSetup           bool      `json:"pin_setup"`
	NeedRpVerify       bool      `json:"need_rp_verify"`
	Avatar             string    `json:"avatar"`
	TokenType          string    `json:"token_type"`
	AccessToken        string    `json:"access_token"`
	DefaultDriveID     string    `json:"default_drive_id"`
	DomainID           string    `json:"domain_id"`
	RefreshToken       string    `json:"refresh_token"`
	IsFirstLogin       bool      `json:"is_first_login"`
	UserID             string    `json:"user_id"`
	NickName           string    `json:"nick_name"`
	State              string    `json:"state"`
	ExpiresIn          int       `json:"expires_in"`
	Status             string    `json:"status"`
}

type ShareInfoReq struct {
	ShareID string `json:"share_id"`
}

type ShareInfoResp struct {
	Avatar             string          `json:"avatar"`
	CreatorID          string          `json:"creator_id"`
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
	FileID        string `json:"file_id"`
	Thumbnail     string `json:"thumbnail"`
	FileType      string `json:"type"`
}

type ShareLinkDownloadURLReq struct {
	ShareID   string `json:"share_id"`
	FileID    string `json:"file_id"`
	ExpireSec int    `json:"expire_sec"`
}

type ShareLinkDownloadURLResp struct {
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
}

type ShareTokenReq struct {
	ShareID  string `json:"share_id"`
	SharePwd string `json:"share_pwd"`
}

type ShareTokenResp struct {
	ShareToken string `json:"share_token"`
	ExpireTime string `json:"expire_time"`
	ExpiresIn  int    `json:"expires_in"`
}

type ShareFileListReq struct {
	ShareID               string `json:"share_id"`
	Starred               bool   `json:"starred"`
	All                   bool   `json:"all"`
	Category              string `json:"category"`
	Fields                string `json:"fields"`
	ImageThumbnailProcess string `json:"image_thumbnail_process"`
	Limit                 int    `json:"limit"`
	Marker                string `json:"marker"`
	OrderBy               string `json:"order_by"`
	OrderDirection        string `json:"order_direction"`
	ParentFileID          string `json:"parent_file_id"`
	Status                string `json:"status"`
	FileType              string `json:"type"`
	URLExpireSec          int    `json:"url_expire_sec"`
	VideoThumbnailProcess string `json:"video_thumbnail_process"`
}

type ShareFileListResp struct {
	Items      []*ShareFile `json:"items"`
	NextMarker string       `json:"next_marker"`
}

type ShareFile struct {
	ShareID       string   `json:"share_id"`
	Name          string   `json:"name"`
	Size          int      `json:"size"`
	Creator       string   `json:"creator"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	DownloadURL   int      `json:"download_url"`
	URL           int      `json:"url"`
	FileExtension string   `json:"file_extension"`
	FileID        string   `json:"file_id"`
	Thumbnail     string   `json:"thumbnail"`
	ParentFileID  string   `json:"parent_file_id"`
	FileType      string   `json:"type"`
	UpdatedAt     string   `json:"updated_at"`
	CreatedAt     string   `json:"created_at"`
	Selected      string   `json:"selected"`
	MimeExtension string   `json:"mime_extension"`
	MimeType      string   `json:"mime_type"`
	PunishFlag    int      `json:"punish_flag"`
	ActionList    []string `json:"action_list"`
	DriveID       string   `json:"drive_id"`
	DomainID      string   `json:"domain_id"`
	RevisionID    string   `json:"revision_id"`
}

// listShareFilesParam is used in file list query context.
type listShareFilesParam struct {
	shareToken   string
	shareID      string
	parentFileID string
	marker       string
}
