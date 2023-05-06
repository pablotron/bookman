INSERT INTO bookman.books(name, author, body) VALUES (
  @name,
  'Unknown Author',
  @body
) RETURNING id;
