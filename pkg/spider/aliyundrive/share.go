package aliyundrive

import (
	"github.com/bibliolater/bookhunter/pkg/log"
)

func (ali AliYunDrive) GetAnonymousShare(shareId string) (*GetShareInfoResponse, error) {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetBody(GetShareInfoRequest{ShareId: shareId}).
		SetResult(GetShareInfoResponse{}).
		SetError(ErrorResponse{}).
		Post(V2ShareLinkGetShareByAnonymous)
	if err != nil {
		return nil, err
	}
	response := downloadResp.Result().(*GetShareInfoResponse)
	return response, nil
}

func (ali AliYunDrive) GetShare(shareId string, shareToken string) (data chan *BaseShareFile, err error) {
	result := make(chan *BaseShareFile, 100)

	go func() {
		err = ali.fileList(shareToken, shareId, result)
		if err != nil {
			log.Fatal(err)
		}
		close(result)
	}()
	return result, nil
}

func (ali AliYunDrive) GetShredToken(shareId string, sharePwd string) (*GetShareTokenResponse, error) {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetBody(GetShareTokenRequest{ShareId: shareId, SharePwd: sharePwd}).
		SetResult(GetShareTokenResponse{}).
		SetError(ErrorResponse{}).
		Post(V2ShareLinkGetShareToken)
	if err != nil {
		return nil, err
	}
	response := downloadResp.Result().(*GetShareTokenResponse)
	return response, nil
}

func (ali AliYunDrive) fileList(shareToken string, shareId string, result chan *BaseShareFile) error {
	return ali.fileListByMarker(shareToken, shareId, "root", "", result)
}

func (ali AliYunDrive) fileListByMarker(shareToken string, shareId string, parentFileId string,
	marker string, result chan *BaseShareFile) error {

	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetHeader(xShareToken, shareToken).
		SetBody(GetShareFileListRequest{
			ShareId:        shareId,
			ParentFileId:   parentFileId,
			UrlExpireSec:   14400,
			OrderBy:        "name",
			OrderDirection: "DESC",
			Limit:          20,
			Marker:         marker,
		}).
		SetResult(GetShareFileListResponse{}).
		SetError(ErrorResponse{}).
		Post(V3FileList)
	if err != nil {
		return err
	}
	data := downloadResp.Result().(*GetShareFileListResponse)
	for _, item := range data.Items {
		if item.FileType == "folder" {
			err := ali.fileListByMarker(shareToken, shareId, item.FileId, "", result)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		result <- item
	}
	if data.NextMarker != "" {
		err := ali.fileListByMarker(shareToken, shareId, parentFileId, data.NextMarker, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ali AliYunDrive) GetFileDownloadUrl(shareToken string, shareId string, fileId string) (string, error) {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetHeader(xShareToken, shareToken).
		SetBody(GetShareLinkDownloadUrlRequest{
			ShareId:   shareId,
			FileId:    fileId,
			ExpireSec: 600,
		}).
		SetResult(GetShareLinkDownloadUrlResponse{}).
		SetError(ErrorResponse{}).
		Post(V2FileGetShareLinkDownloadUrl)
	if err != nil {
		return "", err
	}
	i := downloadResp.Result().(*GetShareLinkDownloadUrlResponse)
	return i.DownloadUrl, nil
}
