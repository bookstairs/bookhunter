package lanzou

import "regexp"

var (
	lxReg  = regexp.MustCompile(`'lx':(\d+),`)
	fidReg = regexp.MustCompile(`'fid':(\d+),`)
	uidReg = regexp.MustCompile(`'uid':'(\d+)',`)
	repReg = regexp.MustCompile(`'rep':'(\d+)',`)
	upReg  = regexp.MustCompile(`'up':(\d+),`)
	lsReg  = regexp.MustCompile(`'ls':(\d+),`)
	tVar   = regexp.MustCompile(`'t':(\S+),`)
	kVar   = regexp.MustCompile(`'k':(\S+),`)

	dirURLRe  = regexp.MustCompile(`(?m)https?://[a-zA-Z0-9-]*?\.?lanzou[a-z]\.com/(/s/)?b[a-zA-Z0-9]{7,}/?`)
	fileURLRe = regexp.MustCompile(`(?m)https?://[a-zA-Z0-9-]*?\.?lanzou[a-z]\.com/(/s/)?i[a-zA-Z0-9]{7,}/?`)

	find1Re      = regexp.MustCompile(`(?m)url\s:\s+'(.*?)',\n\t+data\s:\s+'(.*?)'\+pwd,`)
	find2Re      = regexp.MustCompile(`(?m)<iframe.*?src="(/fn\?\w{10,})"\s.*>`)
	find2TitleRe = regexp.MustCompile(`(?m)<title>(.*?)\s-\s蓝奏云</title>`)

	htmlNoteRe = regexp.MustCompile(`(?m)<!--.+?-->|\s+//\s*.+`)
	jsNoteRe   = regexp.MustCompile(`(?m)(.+?[,;])\s*//.+`)
)

type (
	ResponseData struct {
		Name string `json:"name"`
		Size string `json:"size"`
		URL  string `json:"url"`
	}

	Dom struct {
		Zt  int    `json:"zt"`
		Dom string `json:"dom"`
		// URL 可能为string或int
		URL interface{} `json:"url"`
		// Inf 可能为string或int
		Inf interface{} `json:"inf"`
	}

	FileList struct {
		Zt   int    `json:"zt"`
		Info string `json:"info"`
		// Text 可能为int或数组
		Text interface{} `json:"text"`
	}

	FileItem []struct {
		Icon    string `json:"icon"`
		T       int    `json:"t"`
		ID      string `json:"id"`
		NameAll string `json:"name_all"`
		Size    string `json:"size"`
		Time    string `json:"time"`
		Duan    string `json:"duan"`
		PIco    int    `json:"p_ico"`
	}
)
