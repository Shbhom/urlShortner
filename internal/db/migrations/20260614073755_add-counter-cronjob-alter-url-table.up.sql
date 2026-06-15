-- counter for generating unique short codes
CREATE SEQUENCE IF NOT EXISTS url_counter_seq START 1;

ALTER TABLE short_urls
ADD COLUMN IF NOT EXISTS last_invokation TIMESTAMP WITH TIME ZONE DEFAULT NOW();
