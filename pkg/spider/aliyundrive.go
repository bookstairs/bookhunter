package spider

import (
	"strings"

	"github.com/bibliolater/bookhunter/pkg/spider/aliyundrive"
)

type AliYunConfig struct {
	RefreshToken string
}

var AliyunConfig = &AliYunConfig{}

var drive *aliyundrive.AliYunDrive

// ResolveAliYunDrive reclusive translate the telecom link to a direct download link.
func ResolveAliYunDrive(client *Client, url, passcode string, formats ...string) ([]string, error) {
	return resolveAliYunDrive(client, url, passcode, formats)
}

func resolveAliYunDrive(client *Client, shareUrl string, sharePwd string, formats []string) ([]string, error) {
	shareId := strings.TrimPrefix(shareUrl, "https://www.aliyundrive.com/s/")
	sharePwd = strings.TrimSpace(sharePwd)

	if drive == nil {
		drive = NewAliYunDrive(client, AliyunConfig)
	}
	return resolveShare(drive, shareId, sharePwd, formats)
}

func resolveShare(drive *aliyundrive.AliYunDrive, shareId string, sharePwd string, formats []string) ([]string, error) {
	token, err := drive.GetShredToken(shareId, sharePwd)
	if err != nil {
		return nil, err
	}
	shareFiles, err := drive.GetShare(shareId, token.ShareToken)
	if err != nil {
		return nil, err
	}
	var links []string
	for item := range shareFiles {
		for _, format := range formats {
			if strings.EqualFold(item.FileExtension, format) {
				url, err := drive.GetFileDownloadUrl(token.ShareToken, shareId, item.FileId)
				if err != nil {
					return nil, err
				}
				links = append(links, url)
			}
		}
	}
	return links, nil
}

func NewAliYunDrive(c *Client, aliConfig *AliYunConfig) *aliyundrive.AliYunDrive {
	client := c.client
	client.SetPreRequestHook(aliyundrive.HcHook)
	return &aliyundrive.AliYunDrive{
		Client:       c.client,
		RefreshToken: aliConfig.RefreshToken,
		Cache:        map[string]string{},
	}
}
