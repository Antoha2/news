package reform

//reform:news_categories
type NewsCategories struct {
	NewsId       int `reform:"news_id,pk"`
	CategoriesId int `reform:"categories_id"`
}
