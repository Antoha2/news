package reform

//reform:people
type reformNews struct {
	Id      int    `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}
