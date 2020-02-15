-- # of Transacions in the last x days
-- first day = $1 (ex: 2020-01-01)
-- last day = $1 (ex: 2020-02-16)
SELECT COUNT(*) FROM transactions WHERE transaction_date >= '2020-01-01'
AND transaction_date <= '2020-02-16' AND active='t';

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