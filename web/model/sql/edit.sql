UPDATE bookman.books
   SET name = @name,
       author = @author
 WHERE id = @id;
