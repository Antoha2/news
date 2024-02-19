CREATE TABLE News (
  id int NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL
) 

-- --------------------------------------------------------

--
-- Структура таблицы `NewsCategories`
--

CREATE TABLE NewsCategories (
  NewsId BIGINT NOT NULL,
  CategoryId BIGINT NOT NULL
)

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `News`
--
ALTER TABLE News
  ADD PRIMARY KEY (Id);

--
-- Индексы таблицы `NewsCategories`
--
ALTER TABLE NewsCategories
  ADD PRIMARY KEY (NewsId,CategoryId);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `News`
--
ALTER TABLE News
  MODIFY Id BIGINT NOT NULL AUTO_INCREMENT;