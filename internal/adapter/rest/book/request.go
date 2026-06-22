package book

type CreateBookReq struct {
	Title  string `json:"title" binding:"required" example:"The Great Gatsby"`
	Author string `json:"author" binding:"required" example:"John Doe"`
}

type GetBooksReq struct {
	Page    int `form:"page,default=1" binding:"min=1" example:"1"`
	PerPage int `form:"per_page,default=10" binding:"min=1,max=100" example:"10"`
}

type UpdateBookReq struct {
	Title  string `json:"title" binding:"required" example:"The Great Gatsby"`
	Author string `json:"author" binding:"required" example:"John Doe"`
}
