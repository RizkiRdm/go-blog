CREATE TABLE IF NOT EXISTS `blog_categories` (
    `id_blog` INT NULL,
    `id_category` INT NULL,
    PRIMARY KEY (`id_blog`)
);
ALTER TABLE `blog_categories`
ADD FOREIGN KEY (`id`) REFERENCES `blogs`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `blog_categories`
ADD FOREIGN KEY (`id_category`) REFERENCES `categories`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;