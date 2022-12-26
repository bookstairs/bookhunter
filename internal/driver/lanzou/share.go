package lanzou

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/log"
)

func (l *Lanzou) ResolveShareURL(shareURL, pwd string) ([]ResponseData, error) {
	shareURL = strings.TrimSpace(shareURL)
	// 移除url前部的主机
	rawURL, _ := url.Parse(shareURL)
	parsedURI := rawURL.RequestURI()

	if l.IsFileURL(shareURL) {
		fileShareURL, err := l.resolveFileShareURL(parsedURI, pwd)
		if err != nil {
			return nil, err
		}
		return []ResponseData{*fileShareURL}, err
	} else if l.IsDirURL(shareURL) {
		return l.resolveFileItemShareURL(parsedURI, pwd)
	} else {
		log.Warnf("Unexpected share url, try to download by using directory share API. %s", shareURL)
		return l.resolveFileItemShareURL(parsedURI, pwd)
	}
}

func (l *Lanzou) IsDirURL(shareURL string) bool {
	return dirURLRe.MatchString(shareURL)
}

func (l *Lanzou) IsFileURL(shareURL string) bool {
	return fileURLRe.MatchString(shareURL)
}

func (l *Lanzou) removeNotes(html string) string {
	html = htmlNoteRe.ReplaceAllString(html, "")
	html = jsNoteRe.ReplaceAllString(html, "$1")
	return html
}

func (l *Lanzou) resolveFileShareURL(parsedURI, pwd string) (*ResponseData, error) {
	resp, err := l.R().Get(parsedURI)
	if err != nil {
		return nil, err
	}
	firstPage := resp.String()
	firstPage = l.removeNotes(firstPage)

	// 参考https://github.com/zaxtyson/LanZouCloud-API 中对acwScV2的处理
	if strings.Contains(firstPage, "acw_sc__v2") {
		// 在页面被过多访问或其他情况下，有时候会先返回一个加密的页面，其执行计算出一个acw_sc__v2后放入页面后再重新访问页面才能获得正常页面
		// 若该页面进行了js加密，则进行解密，计算acw_sc__v2，并加入cookie
		acwScV2 := l.calcAcwScV2(firstPage)
		l.SetCookie(&http.Cookie{
			Name:  "acw_sc__v2",
			Value: acwScV2,
		})
		log.Infof("Set Cookie: acw_sc__v2=%v", acwScV2)
		get, _ := l.R().Get(parsedURI)
		firstPage = get.String()
	}

	if strings.Contains(firstPage, "文件取消") || strings.Contains(firstPage, "文件不存在") {
		return nil, fmt.Errorf("文件不存在 %v", parsedURI)
	}

	// Share with password
	if strings.Contains(firstPage, "id=\"pwdload\"") ||
		strings.Contains(firstPage, "id=\"passwddiv\"") {
		if pwd == "" {
			return nil, fmt.Errorf("缺少密码 %v", parsedURI)
		}
		return l.ParsePasswordShare(parsedURI, pwd, firstPage)
	} else if find2Re.MatchString(firstPage) {
		lanzouDom, err := l.ParseAnonymousShare(parsedURI, firstPage)
		return lanzouDom, err
	}
	return nil, fmt.Errorf("解析页面失败")
}

func (l *Lanzou) ParsePasswordShare(parsedURI string, pwd string, firstPage string) (*ResponseData, error) {
	allString := find1Re.FindStringSubmatch(firstPage)
	urlpath := allString[1]
	params := allString[2] + pwd

	result := &Dom{}
	query, _ := url.ParseQuery(params)
	_, err := l.R().
		SetHeader("referer", l.BaseURL+parsedURI).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetResult(result).
		SetFormDataFromValues(query).
		Post(urlpath)
	if err != nil {
		return nil, err
	}
	return l.parseDom(result)
}

func (l *Lanzou) ParseAnonymousShare(parsedURI string, firstPage string) (*ResponseData, error) {
	allString := find2Re.FindStringSubmatch(firstPage)

	dom, err := l.R().Get(allString[1])
	if err != nil {
		return nil, err
	}
	data := make(map[string]string)
	fnDom := l.removeNotes(dom.String())
	var re = regexp.MustCompile(`(?m)var\s+(\w+)\s+=\s+'(.*)';`)
	for _, match := range re.FindAllStringSubmatch(fnDom, -1) {
		data[match[1]] = match[2]
	}
	title := l.extractRegex(find2TitleRe, firstPage)

	var formRe = regexp.MustCompile(`(?m)('(\w+)':([\w']+),?)`)

	fromData := make(map[string]string)
	for _, match := range formRe.FindAllStringSubmatch(fnDom, -1) {
		k := match[2]
		v := match[3]

		if v == "1" {
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			v = strings.TrimLeft(strings.TrimRight(v, "'"), "'")
		} else {
			v = data[v]
		}
		fromData[k] = v
	}

	result := &Dom{}
	_, err = l.R().
		SetHeader("origin", l.BaseURL).
		SetHeader("referer", l.BaseURL+parsedURI).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetResult(result).
		SetFormData(fromData).
		Post("/ajaxm.php")
	if err != nil {
		return nil, err
	}
	lanzouDom, err := l.parseDom(result)
	if lanzouDom != nil {
		lanzouDom.Name = title
	}
	return lanzouDom, err
}

func (l *Lanzou) parseDom(result *Dom) (*ResponseData, error) {
	if result.Zt != 1 {
		return nil, fmt.Errorf("解析直链失败")
	}

	var header = map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Referer":         "https://" + l.Host,
	}

	request := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).
		R()
	rr, err := request.SetHeaders(header).
		Get(result.Dom + "/file/" + result.URL.(string))
	if rr.StatusCode() != 302 && err != nil {
		log.Fatalf("解析链接失败 %v", err)
	}

	if strings.Contains(rr.String(), "网络异常") {
		log.Fatalf("访问过多，被限制，解限功能待实现 %v", err)
	}

	location := rr.Header().Get("location")

	title, _ := result.Inf.(string)
	return &ResponseData{
		Name: title,
		URL:  location,
	}, nil
}

func (l *Lanzou) calcAcwScV2(htmlText string) string {
	arg1Re := regexp.MustCompile(`arg1='([0-9A-Z]+)'`)
	arg1 := l.extractRegex(arg1Re, htmlText)
	acwScV2 := l.hexXor(l.unbox(arg1), "3000176000856006061501533003690027800375")
	return acwScV2
}

func (l *Lanzou) unbox(arg string) string {
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

func (l *Lanzou) hexXor(arg, args string) string {
	a := min(len(arg), len(args))
	res := ""
	for idx := 0; idx < a; idx += 2 {
		v1, _ := strconv.ParseInt(arg[idx:idx+2], 16, 32)
		v2, _ := strconv.ParseInt(args[idx:idx+2], 16, 32)
		//		v to lower case hex
		v3 := fmt.Sprintf("%02x", v1^v2)
		res += v3
	}
	return res
}

func (l *Lanzou) extractRegex(reg *regexp.Regexp, str string) string {
	matches := reg.FindStringSubmatch(str)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func (l *Lanzou) resolveFileItemShareURL(parsedURI, pwd string) ([]ResponseData, error) {
	resp, _ := l.R().Get(parsedURI)
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
	_, err := l.R().SetFormData(formData).SetResult(result).Post("/filemoreajax.php")
	if err != nil {
		return nil, err
	}

	if result.Zt != 1 {
		log.Warnf("lanzou 文件列表解析失败  %v  %v %v", parsedURI, pwd, result.Info)
		return []ResponseData{}, nil
	}

	item, ok := result.Text.([]interface{})

	if !ok {
		log.Warnf("lanzou 文件列表解析失败  %v  %v %v", parsedURI, pwd, result)
		return nil, errors.New("lanzou 文件列表解析失败")
	}

	data := make([]ResponseData, len(item))
	for i, d := range item {
		file := d.(map[string]interface{})

		respData, err := l.resolveFileShareURL("/"+file["id"].(string), pwd)
		if err != nil {
			return nil, err
		}
		respData.Name = file["name_all"].(string)
		data[i] = *respData
	}
	return data, nil
}
