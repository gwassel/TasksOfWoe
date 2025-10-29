CREATE TYPE external_source AS ENUM ('tg');

CREATE TABLE user (
    id BIGSERIAL PRIMARY KEY,
    external_user_id BIGINT NOT NULL,
    external_source external_source NOT NULL
); 

CREATE UNIQUE INDEX user__external_user_id_external_source__idx ON user(external_user_id, external_source);