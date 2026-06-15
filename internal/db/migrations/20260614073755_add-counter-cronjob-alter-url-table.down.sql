ALTER TABLE short_urls DROP COLUMN IF EXISTS last_invokation;

DROP SEQUENCE IF EXISTS url_counter_seq;
