package types

const (
	DefaultPage     int = 1
	DefaultPageSize int = 25
)

type PageContent interface{}

type PageDto struct {
	Page    int           `json:"page"`
	Pages   int           `json:"pages"`
	Size    int           `json:"size"`
	Records int           `json:"records"`
	Content []PageContent `json:"content"`
}

func NewPageDto(pagination Pagination, content []PageContent) PageDto {
	return PageDto{
		Page:    pagination.Page,
		Pages:   pagination.Pages,
		Size:    pagination.Size,
		Records: pagination.Records,
		Content: content,
	}
}

func EmptyPageDto() PageDto {
	return PageDto{}
}

type Pagination struct {
	Page    int
	Pages   int
	Size    int
	Records int
	Limit   int
	Offset  int
}

func NewPagination(page, size, all int) Pagination {
	pages := 0

	if size > 0 && all > 0 {
		pages = all / size

		if all%size != 0 {
			pages = pages + 1
		}
	}

	return Pagination{
		Page:    page,
		Pages:   pages,
		Size:    size,
		Records: all,
		Limit:   size,
		Offset:  (page - 1) * size,
	}
}

func (p Pagination) IsValid() bool {
	if p.Pages == 0 || p.Page <= 0 || p.Size <= 0 || p.Page > p.Pages {
		return false
	}

	return true
}
