DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
    PERFORM cron.unschedule('purge-old-urls');
  END IF;
END $$;

ALTER TABLE short_urls DROP COLUMN IF EXISTS last_invokation;

DROP SEQUENCE IF EXISTS url_counter_seq;
