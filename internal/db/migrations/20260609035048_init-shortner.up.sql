CREATE TABLE short_urls(
    shortnedKey varchar(16) PRIMARY KEY,
    url TEXT NOT NULL,
    Created_At TIMESTAMPTZ DEFAULT NOW()
)
