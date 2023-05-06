SELECT id,
       name,
       author,
       body
  FROM bookman.books
 WHERE id = @id;
