package fetcher

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
)

const (
	textBookMetadata = "https://s-file-2.ykt.cbern.com.cn/zxx/ndrs/resources/tch_material/version/data_version.json"
	downloadLinkTmpl = "https://r1-ndr.ykt.cbern.com.cn/edu_product/esp/assets_document/%s.pkg/pdf.pdf"
)

func newK12Service(config *Config) (service, error) {
	c, err := client.New(config.Config)
	if err != nil {
		return nil, err
	}

	return &k12Service{Client: c, config: config, metadata: map[int64]TextBook{}}, nil
}

type k12Service struct {
	*client.Client
	config   *Config
	metadata map[int64]TextBook
}

func (k *k12Service) size() (int64, error) {
	resp, err := k.R().
		SetResult(&MetadataSource{}).
		Get(textBookMetadata)
	if err != nil {
		return 0, err
	}

	res := resp.Result().(*MetadataSource)
	maxID := int64(0)

	for _, url := range strings.Split(res.URLs, ",") {
		books, err := k.downloadMetadata(url)
		if err != nil {
			return 0, err
		}

		for idx := range books {
			book := books[idx]
			id := book.CustomProperties.ResSort
			k.metadata[id] = book
			if id > maxID {
				maxID = id
			}
		}
	}

	return maxID, nil
}

func (k *k12Service) downloadMetadata(url string) ([]TextBook, error) {
	resp, err := k.R().SetResult([]TextBook{}).Get(url)
	if err != nil {
		return nil, err
	}

	return *resp.Result().(*[]TextBook), nil
}

func (k *k12Service) formats(id int64) (map[file.Format]driver.Share, error) {
	book, ok := k.metadata[id]
	if !ok {
		return map[file.Format]driver.Share{}, nil
	}

	tags := make(map[string]string, len(book.TagList))
	for _, tag := range book.TagList {
		tags[tag.TagID] = tag.TagName
	}

	subPath := ""
	for _, id := range strings.Split(book.TagPaths[0], "/") {
		if name := tags[id]; name != "" {
			if subPath == "" {
				subPath = name
			} else {
				subPath = filepath.Join(subPath, name)
			}
		}
	}

	return map[file.Format]driver.Share{
		file.PDF: {
			FileName:   book.Title + "." + book.CustomProperties.Format,
			SubPath:    subPath,
			Size:       book.CustomProperties.Size,
			URL:        book.ID,
			Properties: nil,
		},
	}, nil
}

func (k *k12Service) fetch(_ int64, _ file.Format, share driver.Share, writer file.Writer) error {
	resp, err := k.R().SetDoNotParseResponse(true).Get(fmt.Sprintf(downloadLinkTmpl, share.URL))
	if err != nil {
		return err
	}

	body := resp.RawBody()
	defer func() { _ = body.Close() }()

	// Save the download content info files.
	_, err = io.Copy(writer, body)
	return err
}

type MetadataSource struct {
	URLs string `json:"urls"`
}

type TextBook struct {
	ID               string `json:"id"`
	CustomProperties struct {
		AutoFillThumb bool `json:"auto_fill_thumb"`
		ExtProperties struct {
			CatalogType string `json:"catalog_type"`
			LibraryID   string `json:"library_id"`
			SubCatalog  string `json:"sub_catalog"`
		} `json:"ext_properties"`
		Format         string   `json:"format"`
		Height         string   `json:"height"`
		IsTop          int      `json:"is_top"`
		Providers      []string `json:"providers"`
		ResSort        int64    `json:"res_sort"`
		Resolution     string   `json:"resolution"`
		Size           int64    `json:"size"`
		SysTransStatus string   `json:"sys_trans_status"`
		Thumbnails     []string `json:"thumbnails"`
		Width          string   `json:"width"`
	} `json:"custom_properties"`
	ResourceTypeCode string `json:"resource_type_code"`
	Language         string `json:"language"`
	Provider         string `json:"provider"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
	TagList          []struct {
		TagID          string `json:"tag_id"`
		TagName        string `json:"tag_name"`
		TagDimensionID string `json:"tag_dimension_id"`
		OrderNum       int    `json:"order_num"`
	} `json:"tag_list"`
	Status               string   `json:"status"`
	CreateContainerID    string   `json:"create_container_id"`
	ResourceTypeCodeName string   `json:"resource_type_code_name"`
	OnlineTime           string   `json:"online_time"`
	TagPaths             []string `json:"tag_paths"`
	ContainerID          string   `json:"container_id"`
	TenantID             string   `json:"tenant_id"`
	Title                string   `json:"title"`
	Label                []string `json:"label"`
	Description          string   `json:"description"`
	ProviderList         []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"provider_list"`
}
