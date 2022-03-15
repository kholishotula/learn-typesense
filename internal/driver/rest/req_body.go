package rest

type createReqBody struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Year    int    `json:"year"`
	Summary string `json:"summary"`
}

func (rb createReqBody) Validate() error {
	if len(rb.Title) == 0 {
		return NewErrBadRequest("missing `title`")
	}
	if len(rb.Author) == 0 {
		return NewErrBadRequest("missing `author`")
	}
	if rb.Year == 0 {
		return NewErrBadRequest("missing `year`")
	}
	return nil
}
