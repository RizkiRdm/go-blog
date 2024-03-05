CREATE TABLE blog (
  id_blog INT PRIMARY KEY AUTO_INCREMENT,
  thumbnail VARCHAR(100),
  judul VARCHAR(100) NOT NULL,
  body TEXT NOT NULL,
  tanggal_posting DATETIME NOT NULL,
  id_pengguna INT FOREIGN KEY REFERENCES pengguna(id_pengguna)
);