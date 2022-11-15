package aliyun

import (
	"io"

	"github.com/bookstairs/bookhunter/internal/log"
)

// AnonymousShare will try to access the share without the user information.
func (ali *Aliyun) AnonymousShare(shareID string) (*ShareInfoResp, error) {
	token, err := ali.AuthToken()
	if err != nil {
		return nil, err
	}

	resp, err := ali.client.R().
		SetAuthToken(token).
		SetBody(&ShareInfoReq{ShareID: shareID}).
		SetResult(&ShareInfoResp{}).
		SetError(&ErrorResp{}).
		Post("https://api.aliyundrive.com/adrive/v2/share_link/get_share_by_anonymous")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*ShareInfoResp), nil
}

func (ali *Aliyun) Share(shareID string, shareToken string) ([]ShareFile, error) {
	return ali.listShareFiles(&listShareFilesParam{
		shareToken:   shareToken,
		shareID:      shareID,
		parentFileID: "root",
		marker:       "",
	})
}

func (ali *Aliyun) listShareFiles(param *listShareFilesParam) ([]ShareFile, error) {
	token, err := ali.AuthToken()
	if err != nil {
		return nil, err
	}

	resp, err := ali.client.R().
		SetAuthToken(token).
		SetHeader("x-share-token", param.shareToken).
		SetBody(&ShareFileListReq{
			ShareID:        param.shareID,
			ParentFileID:   param.parentFileID,
			URLExpireSec:   14400,
			OrderBy:        "name",
			OrderDirection: "DESC",
			Limit:          20,
			Marker:         param.marker,
		}).
		SetResult(&ShareFileListResp{}).
		SetError(&ErrorResp{}).
		Post("https://api.aliyundrive.com/adrive/v3/file/list")
	if err != nil {
		return nil, err
	}

	res := resp.Result().(*ShareFileListResp)
	var files []ShareFile

	for _, item := range res.Items {
		if item.FileType == "folder" {
			list, err := ali.listShareFiles(&listShareFilesParam{
				shareToken:   param.shareToken,
				shareID:      param.shareID,
				parentFileID: item.FileID,
				marker:       "",
			})
			if err != nil {
				return nil, err
			}

			files = append(files, list...)
		} else {
			files = append(files, *item)
		}
	}

	if res.NextMarker != "" {
		list, err := ali.listShareFiles(&listShareFilesParam{
			shareToken:   param.shareToken,
			shareID:      param.shareID,
			parentFileID: param.parentFileID,
			marker:       res.NextMarker,
		})
		if err != nil {
			return nil, err
		}

		files = append(files, list...)
	}

	return files, nil
}

func (ali *Aliyun) ShareToken(shareID string, sharePwd string) (*ShareTokenResp, error) {
	token, err := ali.AuthToken()
	if err != nil {
		return nil, err
	}

	resp, err := ali.client.R().
		SetAuthToken(token).
		SetBody(&ShareTokenReq{ShareID: shareID, SharePwd: sharePwd}).
		SetResult(&ShareTokenResp{}).
		SetError(&ErrorResp{}).
		Post("https://auth.aliyundrive.com/v2/share_link/get_share_token")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*ShareTokenResp), nil
}

func (ali *Aliyun) DownloadURL(shareToken string, shareID string, fileID string) (string, error) {
	token, err := ali.AuthToken()
	if err != nil {
		return "", err
	}

	resp, err := ali.client.R().
		SetAuthToken(token).
		SetHeader("x-share-token", shareToken).
		SetBody(&ShareLinkDownloadURLReq{
			ShareID: shareID,
			FileID:  fileID,
			// Only ten minutes valid
			ExpireSec: 600,
		}).
		SetResult(&ShareLinkDownloadURLResp{}).
		SetError(&ErrorResp{}).
		Post("https://api.aliyundrive.com/v2/file/get_share_link_download_url")
	if err != nil {
		return "", err
	}

	res := resp.Result().(*ShareLinkDownloadURLResp)
	return res.DownloadURL, nil
}

func (ali *Aliyun) DownloadFile(downloadURL string) (io.ReadCloser, int64, error) {
	log.Debugf("Start to download file from aliyun drive: %s", downloadURL)

	resp, err := ali.client.R().
		SetDoNotParseResponse(true).
		Get(downloadURL)
	response := resp.RawResponse

	return response.Body, response.ContentLength, err
}
