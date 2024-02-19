CREATE TABLE IF NOT EXISTS News_Categories (
 	news_id INT NOT NULL,
	categories_id INT NOT NULL,
	PRIMARY KEY(news_id, categories_id),
	FOREIGN KEY(news_id) REFERENCES News (news_id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(categories_id) REFERENCES Categories(categories_id) ON DELETE CASCADE ON UPDATE CASCADE
)