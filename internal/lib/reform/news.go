package reform

//reform:news
type ReformNews struct {
	Id      int    `reform:"news_id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}
