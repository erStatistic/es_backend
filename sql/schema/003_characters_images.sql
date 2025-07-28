-- +goose Up
CREATE TABLE characters_images (
    character_id INT NOT NULL,
    image_id INT NOT NULL,
    PRIMARY KEY (character_id, image_id),
    FOREIGN KEY (character_id) REFERENCES characters (id),
    FOREIGN KEY (image_id) REFERENCES images (id)
);

-- +goose Down
DROP TABLE characters_images;
