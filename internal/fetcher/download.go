package fetcher

// Service defines the link ISP.
type Service string

const (
	Aliyun  Service = "aliyun"  // aliyundrive.com
	Lanzou  Service = "lanzou"  // lanzou.com
	Telecom Service = "telecom" // cloud.189.cn
	Direct  Service = "direct"  // Directly download from the link.
)

// Download will access the download link and return the io response.
func Download(service Service, url string, c *Config) error {
	// TODO Implement this method now.
	return nil
}
