CREATE TABLE pengguna (
  id_pengguna INT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  biodata TEXT,
  foto_profil VARCHAR(255),
  peran ENUM('admin', 'penulis', 'pembaca') DEFAULT 'pembaca'
);