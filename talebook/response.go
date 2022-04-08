package talebook

const (
	SuccessStatus      = "ok"
	BookNotFoundStatus = "not_found"
)

// CommonResponse is the base response for all the requests.
type CommonResponse struct {
	Err string `json:"err"`
}

// LoginResponse is used for login action.
type LoginResponse struct {
	CommonResponse
	Msg string `json:"msg"`
}

// BookListResponse is used to return recent books.
type BookListResponse struct {
	CommonResponse
	Msg   string `json:"msg"`
	Title string `json:"title"`
	Total int64  `json:"total"`
	Books []struct {
		ID int64 `json:"id"`
	} `json:"books"`
}

// BookResponse stands for default book information
type BookResponse struct {
	CommonResponse
	Msg          string `json:"msg"`
	KindleSender string `json:"kindle_sender"`
	Book         struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Files []struct {
			Format string `json:"format"`
			Size   int64  `json:"size"`
			Href   string `json:"href"`
		} `json:"files"`
	} `json:"book"`
}
