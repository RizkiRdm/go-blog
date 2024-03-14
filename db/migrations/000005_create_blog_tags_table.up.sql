CREATE TABLE IF NOT EXISTS `blog_tags` (
    `id_tag` INT NULL,
    `id_blog` INT NULL,
    PRIMARY KEY (`id_tag`)
);
ALTER TABLE `blog_tags`
ADD FOREIGN KEY (`id_blog`) REFERENCES `blogs`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `blog_tags`
ADD FOREIGN KEY (`id_tag`) REFERENCES `tags`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;