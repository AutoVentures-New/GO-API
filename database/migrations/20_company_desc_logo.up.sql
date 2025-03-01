ALTER TABLE companies ADD COLUMN description TEXT NULL AFTER updated_at;

ALTER TABLE companies ADD COLUMN logo VARCHAR(255) NULL AFTER status;
