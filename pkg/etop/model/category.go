package model

type GetEtopCategoryQuery struct {
	CategoryID int64
	Status     *Status3

	Result *EtopCategory
}

type GetEtopCategoriesQuery struct {
	Status *Status3

	Result struct {
		Categories []*EtopCategory
	}
}

type CreateEtopCategoryCommand struct {
	Category *EtopCategory

	Result *EtopCategory
}
