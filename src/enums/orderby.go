package enums

// Параметр отвечает за фильтрацию по полям
// swagger:enum OrderBy
type OrderBy string

const (
	ASCENDING  OrderBy = "asc"
	DESCENDING OrderBy = "desc"
)
