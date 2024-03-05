-- tabel relasi many to many
CREATE TABLE blog_kategori (
  id_blog INT FOREIGN KEY REFERENCES blog(id_blog),
  id_kategori INT FOREIGN KEY REFERENCES kategori(id_kategori),
  PRIMARY KEY (id_blog, id_kategori)
);