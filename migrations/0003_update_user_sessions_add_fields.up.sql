BEGIN;

ALTER TABLE user_sessions
	ADD COLUMN client_device VARCHAR(100) DEFAULT 'Unknown';

ALTER TABLE user_sessions
	ADD COLUMN operating_system VARCHAR(100) DEFAULT 'Windows 10';

ALTER TABLE user_sessions
	ADD COLUMN user_agent VARCHAR(100) DEFAULT 'Unknown';

ALTER TABLE user_sessions
	ALTER COLUMN device_info SET DATA TYPE JSONB USING device_info::JSONB,
	ALTER COLUMN device_info SET DEFAULT '{}';

COMMIT;
