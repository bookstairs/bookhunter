package aliyundrive

import (
	"github.com/bookstairs/bookhunter/pkg/log"
)

func (ali AliYunDrive) GetAnonymousShare(shareID string) (*GetShareInfoResponse, error) {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetBody(GetShareInfoRequest{ShareID: shareID}).
		SetResult(GetShareInfoResponse{}).
		SetError(ErrorResponse{}).
		Post(V2ShareLinkGetShareByAnonymous)
	if err != nil {
		return nil, err
	}
	response := downloadResp.Result().(*GetShareInfoResponse)
	return response, nil
}

func (ali AliYunDrive) GetShare(shareID string, shareToken string) (data chan *BaseShareFile, err error) {
	result := make(chan *BaseShareFile, 100)

	go func() {
		err = ali.fileList(shareToken, shareID, result)
		if err != nil {
			log.Fatal(err)
		}
		close(result)
	}()
	return result, nil
}

func (ali AliYunDrive) GetShredToken(shareID string, sharePwd string) (*GetShareTokenResponse, error) {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetBody(GetShareTokenRequest{ShareID: shareID, SharePwd: sharePwd}).
		SetResult(GetShareTokenResponse{}).
		SetError(ErrorResponse{}).
		Post(V2ShareLinkGetShareToken)
	if err != nil {
		return nil, err
	}
	response := downloadResp.Result().(*GetShareTokenResponse)
	return response, nil
}

func (ali AliYunDrive) fileList(shareToken string, shareID string, result chan *BaseShareFile) error {
	return ali.fileListByMarker(FileListParam{
		shareToken:   shareToken,
		shareID:      shareID,
		parentFileID: "root",
		marker:       "",
	}, result)
}

func (ali AliYunDrive) fileListByMarker(param FileListParam, result chan *BaseShareFile) error {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetHeader(xShareToken, param.shareToken).
		SetBody(GetShareFileListRequest{
			ShareID:        param.shareID,
			ParentFileID:   param.parentFileID,
			URLExpireSec:   14400,
			OrderBy:        "name",
			OrderDirection: "DESC",
			Limit:          20,
			Marker:         param.marker,
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
			err := ali.fileListByMarker(FileListParam{
				shareToken:   param.shareToken,
				shareID:      param.shareID,
				parentFileID: item.FileID,
				marker:       "",
			}, result)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		result <- item
	}
	if data.NextMarker != "" {
		err := ali.fileListByMarker(FileListParam{
			shareToken:   param.shareToken,
			shareID:      param.shareID,
			parentFileID: param.parentFileID,
			marker:       data.NextMarker,
		}, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ali AliYunDrive) GetFileDownloadURL(shareToken string, shareID string, fileID string) (string, error) {
	downloadResp, err := ali.Client.R().
		SetAuthToken(ali.GetAuthorizationToken()).
		SetHeader(xShareToken, shareToken).
		SetBody(GetShareLinkDownloadURLRequest{
			ShareID:   shareID,
			FileID:    fileID,
			ExpireSec: 600,
		}).
		SetResult(GetShareLinkDownloadURLResponse{}).
		SetError(ErrorResponse{}).
		Post(V2FileGetShareLinkDownloadURL)
	if err != nil {
		return "", err
	}
	i := downloadResp.Result().(*GetShareLinkDownloadURLResponse)
	return i.DownloadURL, nil
}
