package telecom

type (
	// All the request and response structs for the telecom login.

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

type (
	// All the request and response structs for the telecom share link query.

	ShareInfo struct {
		ResCode        int    `json:"res_code"`
		ResMessage     string `json:"res_message"`
		AccessCode     string `json:"accessCode"`
		ExpireTime     int    `json:"expireTime"`
		ExpireType     int    `json:"expireType"`
		FileID         string `json:"fileId"`
		FileName       string `json:"fileName"`
		FileSize       int    `json:"fileSize"`
		IsFolder       bool   `json:"isFolder"`
		NeedAccessCode int    `json:"needAccessCode"`
		ShareDate      int64  `json:"shareDate"`
		ShareID        int64  `json:"shareId"`
		ShareMode      int    `json:"shareMode"`
		ShareType      int    `json:"shareType"`
	}

	ShareFile struct {
		CreateDate string `json:"createDate"`
		FileCata   int    `json:"fileCata"`
		ID         int64  `json:"id"`
		LastOpTime string `json:"lastOpTime"`
		Md5        string `json:"md5"`
		MediaType  int    `json:"mediaType"`
		Name       string `json:"name"`
		Rev        string `json:"rev"`
		Size       int64  `json:"size"`
		StarLabel  int    `json:"starLabel"`
	}

	ShareFolder struct {
		CreateDate   string `json:"createDate"`
		FileCata     int    `json:"fileCata"`
		FileListSize int    `json:"fileListSize"`
		ID           int64  `json:"id"`
		LastOpTime   string `json:"lastOpTime"`
		Name         string `json:"name"`
		ParentID     int64  `json:"parentId"`
		Rev          string `json:"rev"`
		StarLabel    int    `json:"starLabel"`
	}

	ShareFiles struct {
		ResCode    int    `json:"res_code"`
		ResMessage string `json:"res_message"`
		ExpireTime int    `json:"expireTime"`
		ExpireType int    `json:"expireType"`
		FileListAO struct {
			Count        int           `json:"count"`
			FileList     []ShareFile   `json:"fileList"`
			FileListSize int64         `json:"fileListSize"`
			FolderList   []ShareFolder `json:"folderList"`
		} `json:"fileListAO"`
		LastRev int64 `json:"lastRev"`
	}

	ShareLink struct {
		ResCode         int    `json:"res_code"`
		ResMessage      string `json:"res_message"`
		FileDownloadURL string `json:"fileDownloadUrl"`
	}
)
