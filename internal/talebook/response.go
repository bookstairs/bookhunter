package talebook

const (
	SuccessStatus      = "ok"
	BookNotFoundStatus = "not_found"
)

// CommonResp is the base response for all the requests.
type CommonResp struct {
	Err string `json:"err"`
}

// LoginResp is used for login action.
type LoginResp struct {
	CommonResp
	Msg string `json:"msg"`
}

// BookResp stands for default book information
type BookResp struct {
	CommonResp
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

// BooksResp is used to return recent books.
type BooksResp struct {
	CommonResp
	Msg   string `json:"msg"`
	Title string `json:"title"`
	Total int64  `json:"total"`
	Books []struct {
		ID int64 `json:"id"`
	} `json:"books"`
}
