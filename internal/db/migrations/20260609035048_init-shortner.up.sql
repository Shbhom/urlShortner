CREATE TABLE shornted_url(
    shortnedKey varchar(16) PRIMARY KEY,
    url TEXT NOT NULL,
    Created_At TIMESTAMPTZ DEFAULT NOW()
)
