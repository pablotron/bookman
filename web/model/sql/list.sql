SELECT id,
       name,
       author,
       0.0 AS rank

  FROM bookman.books

 ORDER BY LOWER(name)
