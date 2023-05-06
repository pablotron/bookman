SELECT id,
       name,
       author,
       ts_rank_cd(ts_vec, websearch_to_tsquery('english', @q)) AS rank
  FROM bookman.books
 WHERE websearch_to_tsquery('english', @q) @@ ts_vec
 ORDER BY rank DESC;
