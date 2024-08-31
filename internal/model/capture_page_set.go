package model

type PaginatedCapturePageSet struct {
	CapturePages []CapturePageSet `json:"capture_page_set"`
	Pagination   *Pagination      `json:"pagination"`
}

type CapturePageSetFilters struct {
	IdsIn                  []int `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}
