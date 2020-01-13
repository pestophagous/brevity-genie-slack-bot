package chat

type chatSortByLastActivity []ReadonlyChat

func (cs chatSortByLastActivity) Len() int { return len(cs) }

func (cs chatSortByLastActivity) Less(lhs, rhs int) bool {
	return cs[lhs].EndTime().Before(cs[rhs].EndTime())
}

func (cs chatSortByLastActivity) Swap(lhs, rhs int) { cs[lhs], cs[rhs] = cs[rhs], cs[lhs] }
