CREATE TABLE space (
    id BIGSERIAL PRIMARY KEY,
    external_space_id BIGINT NOT NULL,
    external_source external_source NOT NULL
); 

CREATE UNIQUE INDEX user__external_space_id_external_source__idx ON space (external_space_id, external_source);