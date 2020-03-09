-- # of Transacions in the last x days
-- first day = $1 (ex: 2020-01-01)
-- last day = $1 (ex: 2020-02-16)
SELECT COUNT(*) FROM transactions WHERE transaction_date >= NOW() - interval '30' day
AND transaction_date <= NOW() + interval '1' day AND active='t';

-- Total balance
SELECT SUM(a.balance) FROM (
    SELECT balance FROM accounts WHERE active=True
) AS a;

-- Get the number of days as a decimal since the start date defined from settings
SELECT 
    CASE 
        WHEN EXTRACT(epoch FROM AGE(start_date, NOW()))/86400 >= 0
        THEN EXTRACT(epoch FROM AGE(start_date, NOW()))/86400
        ELSE EXTRACT(epoch FROM AGE(start_date, NOW()))/86400 * -1
    END delta_start_date
FROM settings LIMIT 1;

-- Average balance per day, per account
SELECT
    acc.id,
    acc.balance / b.delta_start_date as money_per_day,
    acc.name,
    acc.balance
FROM accounts AS acc
JOIN (
    SELECT 
        CASE 
            WHEN EXTRACT(epoch FROM AGE(start_date, NOW()))/86400 >= 0
            THEN EXTRACT(epoch FROM AGE(start_date, NOW()))/86400
            ELSE EXTRACT(epoch FROM AGE(start_date, NOW()))/86400 * -1
        END delta_start_date
    FROM settings LIMIT 1
) AS b ON 1=1
WHERE acc.active=True;

-- Total expenses last 30 days 
SELECT SUM(amount) FROM transactions WHERE transaction_date >= NOW() - interval '30' day
AND transaction_date <= NOW() + interval '1' day AND active='t'
AND to_account IS NULL;

-- Total income last 30 days
SELECT SUM(amount) FROM transactions WHERE transaction_date >= NOW() - interval '30' day
AND transaction_date <= NOW() + interval '1' day AND active='t' AND from_account IS NULL;

-- Average value moved per account
SELECT 
    SUM(amount) / COUNT(account_id)
FROM transactions WHERE transaction_date >= '2020-01-01'
AND transaction_date <= '2020-02-27' AND active='t';
