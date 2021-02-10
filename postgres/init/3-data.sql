BEGIN;

-- Settings
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

--
-- Statistics
--

-- Numbers
INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    'Total Balance',
    'SELECT SUM(a.balance) FROM (SELECT balance FROM accounts WHERE active=True) AS a;',
    NOW(),
    NOW(),
    NOW(),
    'number',
    'total_balance',
    't',
    '',
    '',
    '',
    ''
);

INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    '# of Transactions last 30 days',
    'SELECT COUNT(*) FROM transactions WHERE transaction_date >= NOW() - interval ''30'' day AND transaction_date <= NOW() + interval ''1'' day AND active=True;',
    NOW(),
    NOW(),
    NOW(),
    'number',
    'transaction_count',
    'f',
    '',
    '',
    '',
    ''
);

INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    'Total income last 30 days',
    'SELECT SUM(amount) FROM transactions WHERE transaction_date >= NOW() - interval ''30'' day
    AND transaction_date <= NOW() + interval ''1'' day AND active=''t'' AND account_id IS NULL;',
    NOW(),
    NOW(),
    NOW(),
    'number',
    'total_income',
    't',
    '',
    '',
    '',
    ''
);

INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    'Total expenses last 30 days',
    'SELECT SUM(amount) FROM transactions WHERE transaction_date >= NOW() - interval ''30'' day
    AND transaction_date <= NOW() + interval ''1'' day AND active=''t'' AND to_account IS NULL;',
    NOW(),
    NOW(),
    NOW(),
    'number',
    'total_expenses',
    't',
    '',
    '',
    '',
    ''
);

-- Graphs
INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    'Balance per day',
    'SELECT json_object_agg(a.name, a.money_per_day) FROM (
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
    ) AS a;',
    NOW(),
    NOW(),
    NOW(),
    'bar',
    'balance_per_day',
    't',
    '',
    '',
    '',
    ''
);

INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    'Total amount per category',
    'SELECT json_object_agg(b.name, b.obj) FROM (
        SELECT a.name,json_build_object(''hex'', a.hex, ''value'', a.sum) obj FROM (
            SELECT c.name,SUM(t.amount),c.hex FROM categories AS c
            JOIN transactions AS t ON c.id=t.category_id
            WHERE t.active=true AND c.active=true
            GROUP BY c.name,c.hex
        ) AS a
    ) AS b;',
    NOW(),
    NOW(),
    NOW(),
    'pie',
    'total_category_amount',
    't',
    '',
    '',
    '',
    ''
);

INSERT INTO statistics (active, name, compute_query, create_date, last_update, execution_date, visualisation, external_id, monetary, description, keys, value, suffix) VALUES (
    't',
    'Amount per Category, last 30 days',
    'SELECT json_object_agg(b.name, b.obj) FROM (
        SELECT a.name,json_build_object(''hex'', a.hex, ''value'', a.sum) obj FROM (
            SELECT c.name,SUM(t.amount),c.hex FROM categories AS c
            JOIN transactions AS t ON c.id=t.category_id
            WHERE t.active=true AND c.active=true AND t.transaction_date >= NOW() - interval ''30 days''
            GROUP BY c.name,c.hex
        ) AS a
    ) AS b;',
    NOW(),
    NOW(),
    NOW(),
    'pie',
    'past_category_amount',
    't',
    '',
    '',
    '',
    ''
);

COMMIT;