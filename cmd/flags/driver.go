package flags

var (
	Driver = "telecom"

	// Aliyun Drive

	RefreshToken = ""

	// Telecom Cloud

	TelecomUsername = ""
	TelecomPassword = ""
)

// NewDriverProperties which would be used in driver.New
func NewDriverProperties() map[string]string {
	return map[string]string{
		"driver":          Driver,
		"refreshToken":    RefreshToken,
		"telecomUsername": Username,
		"telecomPassword": Password,
	}
}
