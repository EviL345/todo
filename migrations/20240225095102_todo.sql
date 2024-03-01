-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    login varchar(255) NOT NULL,
    password varchar(255) NOT NULL
    );

CREATE TABLE IF NOT EXISTS tasks (
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL,
    title varchar(255) NOT NULL,
    description varchar(255),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
