package lanzou

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/log"
)

type Drive struct {
	Client  *resty.Client
	BaseURL string
}

func NewDefaultDrive(client *resty.Client) *Drive {
	baseURL := "https://lanzoux.com"
	return NewDrive(client, baseURL)
}

func NewDrive(client *resty.Client, baseURL string) *Drive {
	client.
		SetBaseURL(baseURL).
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Referer", baseURL)
	return &Drive{
		Client:  client,
		BaseURL: baseURL,
	}
}

type Response struct {
	Code int64        `json:"code"`
	Data ResponseData `json:"data"`
	Msg  string       `json:"msg"`
}

type ResponseData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (l Drive) ResolveShareURL(shareURL string, pwd string) (*Response, error) {
	// 移除url前部的主机
	rawURL, _ := url.Parse(shareURL)
	parsedURI := rawURL.RequestURI()
	if l.IsFileURL(shareURL) {
		return l.resolveFileShareURL(parsedURI, pwd)
	} else if l.IsDirURL(shareURL) {
		return l.resolveFileItemShareURL(parsedURI, pwd)
	}
	return nil, fmt.Errorf("unsupport shareURL %v", shareURL)
}

func (l Drive) IsDirURL(shareURL string) bool {
	return dirURLRe.MatchString(shareURL)
}
func (l Drive) IsFileURL(shareURL string) bool {
	return fileURLRe.MatchString(shareURL)
}

func (l Drive) removeNotes(html string) string {
	html = htmlNoteRe.ReplaceAllString(html, "")
	html = jsNoteRe.ReplaceAllString(html, "$1")
	return html
}

func (l Drive) resolveFileShareURL(parsedURI string, pwd string) (*Response, error) {
	get, _ := l.Client.R().Get(parsedURI)

	firstPage := get.String()
	firstPage = l.removeNotes(firstPage)

	// 参考https://github.com/zaxtyson/LanZouCloud-API 中对acwScV2的处理
	if strings.Contains(firstPage, "acw_sc__v2") {
		// 在页面被过多访问或其他情况下，有时候会先返回一个加密的页面，其执行计算出一个acw_sc__v2后放入页面后再重新访问页面才能获得正常页面
		//	# 若该页面进行了js加密，则进行解密，计算acw_sc__v2，并加入cookie
		acwScV2 := l.calcAcwScV2(firstPage)
		l.Client.SetCookie(&http.Cookie{
			Name:  "acw_sc__v2",
			Value: acwScV2,
		})
		log.Infof("Set Cookie: acw_sc__v2=%v", acwScV2)
		get, _ := l.Client.R().Get(parsedURI)
		firstPage = get.String()
	}

	if strings.Contains(firstPage, "文件取消") || strings.Contains(firstPage, "文件不存在") {
		return nil, fmt.Errorf("文件不存在 %v", parsedURI)
	}

	// share with password
	allString := find1Re.FindStringSubmatch(firstPage)
	if len(allString) == 3 {
		urlpath := allString[1]
		params := allString[2] + pwd

		result := &LanzouyDom{}

		query, _ := url.ParseQuery(params)

		_, _ = l.Client.R().
			SetHeader("referer", l.BaseURL+parsedURI).
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetResult(result).
			SetFormDataFromValues(query).
			Post(urlpath)
		return l.parseLanzouDom(result)
	}

	// Share without password
	allString = find2Re.FindStringSubmatch(firstPage)
	if len(allString) == 2 {
		dom, _ := l.Client.R().Get(allString[1])

		data := make(map[string]string)

		var re = regexp.MustCompile(`(?m)var\s+(\w+)\s+=\s+'(.*)';`)
		for _, match := range re.FindAllStringSubmatch(dom.String(), -1) {
			data[match[1]] = match[2]
		}
		title := l.extractRegex(find2TitleRe, firstPage)

		result := &LanzouyDom{}
		_, _ = l.Client.R().
			SetHeader("origin", l.BaseURL).
			SetHeader("referer", l.BaseURL+parsedURI).
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetResult(result).
			SetFormData(map[string]string{
				"action":     "downprocess",
				"signs":      data["ajaxdata"],
				"sign":       data["msigns"],
				"websign":    data["wsigns"],
				"websignkey": data["cwebsignkeyc"],
				"ves":        "1",
			}).
			Post("/ajaxm.php")

		lanzouDom, err := l.parseLanzouDom(result)
		if lanzouDom != nil {
			lanzouDom.Data.Name = title
		}
		return lanzouDom, err
	}

	return nil, fmt.Errorf("解析页面失败")
}

func (l Drive) parseLanzouDom(result *LanzouyDom) (*Response, error) {
	var header = map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Referer":         "https://lanzous.com",
	}
	if result.Zt == 1 {
		request := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
		rr, err := request.SetHeaders(header).
			Get(result.Dom + "/file/" + result.URL)
		if rr.StatusCode() != 302 && err != nil {
			log.Fatal("解析链接失败", err)
		}

		if strings.Contains(rr.String(), "网络异常") {
			log.Fatal("访问过多，被限制，解限功能待实现", err)
		}

		location := rr.Header().Get("location")

		title, _ := result.Inf.(string)
		return &Response{
			Code: 200,
			Data: ResponseData{
				Name: title,
				URL:  location,
			},
		}, nil
	} else {
		return nil, fmt.Errorf("解析直链失败")
	}
}

func (l Drive) calcAcwScV2(htmlText string) string {
	arg1Re := regexp.MustCompile(`arg1='([0-9A-Z]+)'`)
	arg1 := l.extractRegex(arg1Re, htmlText)
	acwScV2 := l.hexXor(l.unsbox(arg1), "3000176000856006061501533003690027800375")
	return acwScV2
}

func (l Drive) unsbox(arg string) string {
	v1 := []int{15, 35, 29, 24, 33, 16, 1, 38, 10, 9, 19, 31, 40, 27, 22, 23, 25, 13, 6, 11,
		39, 18, 20, 8, 14, 21, 32, 26, 2, 30, 7, 4, 17, 5, 3, 28, 34, 37, 12, 36}
	v2 := make([]string, len(v1))
	for idx, v3 := range arg {
		for idx2, in := range v1 {
			if in == (idx + 1) {
				v2[idx2] = string(v3)
			}
		}
	}
	return strings.Join(v2, "")
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func (l Drive) hexXor(arg string, args string) string {
	a := min(len(arg), len(args))
	res := ""
	for idx := 0; idx < a; idx += 2 {
		v1, _ := strconv.ParseInt(arg[idx:idx+2], 16, 32)
		v2, _ := strconv.ParseInt(args[idx:idx+2], 16, 32)
		//		v to lowcase hex
		v3 := fmt.Sprintf("%02x", v1^v2)
		res += v3
	}
	return res
}

type LanzouyDom struct {
	Zt  int         `json:"zt"`
	Dom string      `json:"dom"`
	URL string      `json:"url"`
	Inf interface{} `json:"inf"`
}

type FileList struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
	Text []struct {
		Icon    string `json:"icon"`
		T       int    `json:"t"`
		ID      string `json:"id"`
		NameAll string `json:"name_all"`
		Size    string `json:"size"`
		Time    string `json:"time"`
		Duan    string `json:"duan"`
		PIco    int    `json:"p_ico"`
	} `json:"text"`
}

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

func (l Drive) extractRegex(reg *regexp.Regexp, str string) string {
	matches := reg.FindStringSubmatch(str)

	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func (l Drive) resolveFileItemShareURL(parsedURI string, pwd string) (*Response, error) {
	resp, _ := l.Client.R().Get(parsedURI)
	str := resp.String()
	formData := map[string]string{
		"lx":  l.extractRegex(lxReg, str),
		"fid": l.extractRegex(fidReg, str),
		"uid": l.extractRegex(uidReg, str),
		"pg":  "1",
		"rep": l.extractRegex(repReg, str),
		"t":   l.extractRegex(regexp.MustCompile("var "+l.extractRegex(tVar, str)+" = '(\\d+)';"), str),
		"k":   l.extractRegex(regexp.MustCompile("var "+l.extractRegex(kVar, str)+" = '(\\S+)';"), str),
		"up":  l.extractRegex(upReg, str),
		"ls":  l.extractRegex(lsReg, str),
		"pwd": pwd,
	}

	result := &FileList{}
	_, _ = l.Client.R().SetFormData(formData).SetResult(result).Post("/filemoreajax.php")

	if len(result.Text) > 0 {
		u := ""

		for _, file := range result.Text {
			if strings.Contains(file.NameAll, "epub") {
				u = "/" + file.ID
			}
		}

		if u == "" {
			u = "/" + result.Text[0].ID
		}

		return l.resolveFileShareURL(u, pwd)
	}

	return nil, fmt.Errorf("获取连接失败")
}
