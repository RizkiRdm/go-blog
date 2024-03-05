-- Relasi many to many antara Blog dan Tag
CREATE TABLE blog_tag (
  id_blog INT FOREIGN KEY REFERENCES blog(id_blog),
  id_tag INT FOREIGN KEY REFERENCES tag(id_tag),
  PRIMARY KEY (id_blog, id_tag)
);
