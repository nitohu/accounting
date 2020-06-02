--
-- Numbers
--

-- # of Transacions in the last x days
SELECT COUNT(*) FROM transactions WHERE transaction_date >= NOW() - interval '30 days'
AND transaction_date <= NOW() + interval '1 day' day AND active='t';

-- Total balance
SELECT SUM(a.balance) FROM (
    SELECT balance FROM accounts WHERE active=True
) AS a;

-- Total expenses last 30 days 
SELECT SUM(amount) FROM transactions WHERE transaction_date >= NOW() - interval '30' day
AND transaction_date <= NOW() + interval '1' day AND active='t'
AND to_account IS NULL;

-- Total income last 30 days
SELECT SUM(amount) FROM transactions WHERE transaction_date >= NOW() - interval '30' day
AND transaction_date <= NOW() + interval '1' day AND active='t' AND account_id IS NULL;

-- Average value moved per account
SELECT 
    SUM(amount) / COUNT(account_id)
FROM transactions WHERE transaction_date >= NOW() - interval '30' day
AND transaction_date <= NOW() + interval '1' day AND active='t';

--
-- Graphs
--

-- Average balance per day, per account
SELECT json_object_agg(a.name, a.money_per_day) FROM (
    SELECT
        acc.id,
        CASE
            WHEN b.delta_salary_date > 1
            THEN acc.balance / b.delta_salary_date 
            ELSE acc.balance * b.delta_salary_date
        END money_per_day,
        acc.name,
        acc.balance
    FROM accounts AS acc
    JOIN (
        SELECT 
            CASE 
                WHEN EXTRACT(epoch FROM AGE(salary_date, NOW()))/86400 >= 0
                THEN EXTRACT(epoch FROM AGE(salary_date, NOW()))/86400
                ELSE EXTRACT(epoch FROM AGE(salary_date, NOW()))/86400 * -1
            END delta_salary_date
        FROM settings LIMIT 1
    ) AS b ON 1=1
    WHERE acc.active=True AND acc.balance > 0
) AS a;

-- Money spent per category, total
SELECT json_object_agg(b.name, b.obj) FROM (
    SELECT a.name,json_build_object('hex', a.hex, 'value', a.sum) obj FROM (
        SELECT c.name,SUM(t.amount),c.hex FROM categories AS c
        JOIN transactions AS t ON c.id=t.category_id
        WHERE t.active=true AND c.active=true
        GROUP BY c.name,c.hex
    ) AS a
) AS b;

-- Money spent per category, last 30 days
SELECT json_object_agg(b.name, b.obj) FROM (
    SELECT a.name,json_build_object('hex', a.hex, 'value', a.sum) obj FROM (
        SELECT c.name,SUM(t.amount),c.hex FROM categories AS c
        JOIN transactions AS t ON c.id=t.category_id
        WHERE t.active=true AND c.active=true AND t.transaction_date >= NOW() - interval '30 days'
        GROUP BY c.name,c.hex
    ) AS a
) AS b;

-- Average Money spent per category, last 30 days
SELECT json_object_agg(b.name, b.amount) FROM (
    SELECT a.name,a.amount/30 amount FROM (
        SELECT c.name,SUM(t.amount) amount FROM categories AS c
        JOIN transactions AS t ON c.id=t.category_id
        WHERE t.active=true AND c.active=true AND t.transaction_date >= NOW() - interval '30 days'
        GROUP BY c.name
    ) as a
) as b;

-- Money spent per category, per month => Ext. ID: category_spent_monthly
select json_object_agg(
	CONCAT(a.month, '/', a.year),
	a.amount
) from (
	SELECT to_char(date_trunc('month', t.transaction_date), 'YYYY') AS year,
		   to_char(date_trunc('month', t.transaction_date), 'Mon') AS month,
		   sum(t.amount) AS amount
	FROM transactions AS t
	GROUP BY date_trunc('month', t.transaction_date)
) as a;
