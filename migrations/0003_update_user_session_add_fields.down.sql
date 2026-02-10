BEGIN;

ALTER TABLE user_sessions
	DROP COLUMN IF EXISTS operating_system;

ALTER TABLE user_sessions
	DROP COLUMN IF EXISTS client_device;

ALTER TABLE user_sessions
	DROP COLUMN IF EXISTS user_agent;

ALTER TABLE user_sessions
	ALTER COLUMN device_info DROP DEFAULT,
	ALTER COLUMN device_info SET DATA TYPE TEXT USING device_info::text;

COMMIT;
