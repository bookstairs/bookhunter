package telecom

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bookstairs/bookhunter/internal/log"
)

// ShareFiles will resolve the telecom-shared link.
func (t *Telecom) ShareFiles(accessURL, accessCode string) (*ShareInfo, []ShareFile, error) {
	log.Debugf("Download telecom file from %s -- %s", accessURL, accessCode)

	// Get share info.
	info, err := t.shareInfo(accessURL)
	if err != nil {
		return nil, nil, err
	}

	// Reclusive get the shared files.
	var files []ShareFile
	if info.IsFolder {
		files, err = t.listShareFolders(accessCode, info.FileID, info.FileID, info.ShareID, info.ShareMode)
	} else {
		files, err = t.listShareFiles(accessCode, info.FileID, info.ShareID, info.ShareMode)
	}
	if err != nil {
		return nil, nil, err
	}

	return info, files, nil
}

func (t *Telecom) DownloadURL(shareCode, shareID, fileID string) (string, error) {
	resp, err := t.client.R().
		SetHeaders(map[string]string{
			"accept":  "application/json;charset=UTF-8",
			"origin":  "https://cloud.189.cn",
			"Referer": "https://cloud.189.cn/web/share?code=" + shareCode,
		}).
		SetQueryParams(map[string]string{
			"fileId":  fileID,
			"shareId": shareID,
			"dt":      "1",
		}).
		SetResult(&ShareLink{}).
		Get(webPrefix + "/api/open/file/getFileDownloadUrl.action")
	if err != nil {
		return "", err
	}
	link := resp.Result().(*ShareLink)

	return link.FileDownloadURL, nil
}

func (t *Telecom) DownloadFile(url string) (io.ReadCloser, error) {
	resp, err := t.client.R().
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.RawBody(), nil
}

// ShareCode extract the share code.
func (t *Telecom) ShareCode(accessURL string) (string, error) {
	if idx := strings.LastIndex(accessURL, "/"); idx > 0 {
		rs := []rune(accessURL)
		return string(rs[idx+1:]), nil
	} else {
		return "", fmt.Errorf("invalid share link, couldn't find share code: %s", accessURL)
	}
}

func (t *Telecom) shareInfo(accessURL string) (*ShareInfo, error) {
	shareCode, err := t.ShareCode(accessURL)
	if err != nil {
		return nil, err
	}
	resp, err := t.client.R().
		SetHeaders(map[string]string{
			"accept":  "application/json;charset=UTF-8",
			"origin":  "https://cloud.189.cn",
			"Referer": "https://cloud.189.cn/web/share?code=" + shareCode,
		}).
		SetQueryParam("shareCode", shareCode).
		SetResult(&ShareInfo{}).
		Get(webPrefix + "/api/open/share/getShareInfoByCode.action")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*ShareInfo), nil
}

func (t *Telecom) listShareFiles(code, fileID string, shareID int64, mode int) ([]ShareFile, error) {
	resp, err := t.client.R().
		SetQueryParams(map[string]string{
			"fileId":     fileID,
			"shareId":    strconv.FormatInt(shareID, 10),
			"shareMode":  strconv.Itoa(mode),
			"accessCode": code,
			"isFolder":   "false",
			"iconOption": "5",
			"pageNum":    "1",
			"pageSize":   "10",
		}).
		SetResult(&ShareFiles{}).
		Get(webPrefix + "/api/open/share/listShareDir.action")
	if err != nil {
		return nil, err
	}
	res := resp.Result().(*ShareFiles)

	return res.FileListAO.FileList, nil
}

func (t *Telecom) listShareFolders(code, fileID, shareDirFileID string, shareID int64, mode int) ([]ShareFile, error) {
	resp, err := t.client.R().
		SetQueryParams(map[string]string{
			"fileId":         fileID,
			"shareDirFileId": shareDirFileID,
			"shareId":        strconv.FormatInt(shareID, 10),
			"shareMode":      strconv.Itoa(mode),
			"accessCode":     code,
			"isFolder":       "true",
			"orderBy":        "lastOpTime",
			"descending":     "true",
			"iconOption":     "5",
			"pageNum":        "1",
			"pageSize":       "60",
		}).
		SetResult(&ShareFiles{}).
		Get(webPrefix + "/api/open/share/listShareDir.action")
	if err != nil {
		return nil, err
	}

	res := resp.Result().(*ShareFiles).FileListAO
	files := res.FileList

	for _, folder := range res.FolderList {
		id := strconv.FormatInt(folder.ID, 10)
		children, err := t.listShareFolders(code, id, shareDirFileID, shareID, mode)
		if err != nil {
			return nil, err
		}

		files = append(files, children...)
	}

	return files, nil
}
