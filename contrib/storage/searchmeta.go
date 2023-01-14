package storage

import "github.com/zander-84/seagull/contract"

type searchMeta struct {
	cursor       bool
	doPagination bool
	page         int
	pageSize     int

	maxPage     int
	maxPageSize int

	doCount bool
}

func NewSearchMete() contract.SearchMeta {
	out := new(searchMeta)
	out.maxPage = 10000
	out.maxPageSize = 100
	out.cursor = false
	out.doPagination = true
	return out
}

func (s *searchMeta) UseCursor(cursor bool) contract.SearchMeta {
	s.cursor = cursor
	return s
}

func (s *searchMeta) IsCursor() bool {
	return s.cursor
}

func (s *searchMeta) UsePage(page bool) contract.SearchMeta {
	s.doPagination = page
	return s
}

func (s *searchMeta) IsPage() bool {
	return s.doPagination
}

func (s *searchMeta) SetPage(page int) contract.SearchMeta {
	s.page = page
	return s
}

func (s *searchMeta) SetPageSize(pageSize int) contract.SearchMeta {
	s.pageSize = pageSize
	return s
}
func (s *searchMeta) SetMaxPage(maxPage int) contract.SearchMeta {
	s.maxPage = maxPage
	return s
}

func (s *searchMeta) SetMaxPageSize(maxPageSize int) contract.SearchMeta {
	s.maxPageSize = maxPageSize
	return s
}

func (s *searchMeta) Page() int {
	if s.page < 1 {
		return 1
	} else if s.page > s.maxPage {
		return s.maxPage
	} else {
		return s.page
	}
}

func (s *searchMeta) PageSize() int {
	if s.pageSize < 1 {
		return 10
	} else if s.pageSize > s.maxPageSize {
		return s.maxPageSize
	} else {
		return s.pageSize
	}
}

func (s *searchMeta) Offset() int {
	return (s.Page() - 1) * s.PageSize()
}

func (s *searchMeta) IsCount() bool {
	return s.doCount
}

func (s *searchMeta) UseCount(cnt bool) contract.SearchMeta {
	s.doCount = cnt
	return s
}
