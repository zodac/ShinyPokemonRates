CREATE TABLE IF NOT EXISTS shiny_pokemon (
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    pokemon varchar(255) NOT NULL,
    id int,
    seen int,
    found int,

    PRIMARY KEY (timestamp, id)
);

CREATE INDEX index_timestamp ON shiny_pokemon(timestamp)