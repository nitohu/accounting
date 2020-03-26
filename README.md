# Accounting

This is an accounting application which is developed in go.

You can create accounts, transactions, categorise them and have statistics which whom you can analyze your finances.
The server (master) of this application is written in Golang.

## Requirements

### Installations

* Golang [https://golang.org/dl/](Download)
* Postgresql [https://www.postgresql.org/download/](Download)

Recommendation:

* If you aren't familiar with Postgres or SQL in general, it might be helpful to install pgAdmin:

[https://www.pgadmin.org/download/](Download)

### Set up your database

**!!! If you are manually inserting records into the database, please always make sure none of the fields in the record are NULL !!!**

Create the database inside postgres:

```sql
> CREATE DATABASE accounting;
-- Connect to the database
> \c accounting
```

Open the `db.sql` file and execute the SQL commands for creating the tables.

You can now insert the queries for the statistics in the statistics table. Please make sure none of the fields are NULL, use an empty string or a placeholder instead.

Please use the following external IDs and types of visualisation for the statistics to display them properly on the dashboard:

Name | External Identifier | Visualisation
--- | --- | ---
Total Number of Transactions last 30 days | transaction_count     | number
 Total Balance                             | total_balance         | number
 Balance per day                           | balance_per_day       | bar
 Total Amount per Category                 | total_category_amount | pie
 Amount per Category, last 30 days         | past_category_amount  | pie
 Total expenses last 30 days               | total_expenses        | number
 Total income last 30 days                 | total_income          | number
 Average balance per day, total            | total_avg_balance     | number

### Set a master password

You need to set a master password for the application, so that you can login. You need to insert the record manually into the statistics table.

Use the following SQL command for inserting the record into the settings table:

```sql
INSERT INTO settings (name, password, email, last_update, calc_interval, calc_uom, currency, session_key, salary_date) VALUES (
    'your name',
    '8C6976E5B5410415BDE908BD4DEE15DFB167A9C873FC4BB8A81F6F2AB448A918',
    'your email',
    NOW(),
    30,
    'minutes',
    '€',
    '',
    NOW()
);
```
The password in this query is `admin`.

If you want to set your own password, please encrypt it with Sha256.
