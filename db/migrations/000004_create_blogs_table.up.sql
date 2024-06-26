CREATE TABLE IF NOT EXISTS `blogs` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    title VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    thumbnail VARCHAR(100),
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
ALTER TABLE `blogs`
ADD FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE
SET NULL ON UPDATE CASCADE;