BEGIN;

INSERT INTO settings (name, password, email, last_update, calc_interval, calc_uom, currency, session_key, salary_date) VALUES (
    'your name',
    '8C6976E5B5410415BDE908BD4DEE15DFB167A9C873FC4BB8A81F6F2AB448A918',
    'admin',
    NOW(),
    30,
    'minutes',
    '€',
    '',
    NOW()
);

COMMIT;