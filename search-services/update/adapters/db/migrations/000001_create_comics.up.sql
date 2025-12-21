CREATE TABLE IF NOT EXISTS comics (
    comics_id integer PRIMARY KEY,
    image_url text NOT NULL,
    words_total integer NOT NULL CHECK (words_total > 0)

);

CREATE TABLE IF NOT EXISTS key_words (
    key_word text,
    comics_id integer REFERENCES comics (comics_id),
    PRIMARY KEY (key_word, comics_id)
);

CREATE OR REPLACE VIEW db_stats AS
SELECT words_total."total", words_unique."unique", comics_fetched."fetched"
FROM (SELECT COUNT(*) "total"
    FROM key_words) words_total
CROSS JOIN (SELECT COUNT(DISTINCT key_word) "unique"
    FROM key_words) words_unique
CROSS JOIN (SELECT COUNT(*) "fetched"
    FROM comics) comics_fetched;