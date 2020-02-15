CREATE TABLE settings (
    name text,
    password text,
    email text,
    last_update timestamp,
    start_date timestamp,
    calc_interval int,
    calc_uom text,
    currency text,
    session_key text,
);

CREATE TABLE accounts (
    id serial,
    primary key(id),
    name text,
    active boolean,
    balance float,
    balance_forecast float,
    iban text,
    account_holder text,
    bank_code text,
    account_nr text,
    bank_name text,
    bank_type text,
    create_date timestamp,
    last_update timestamp
);

CREATE TABLE categories (
    id serial,
    primary key(id),
    name text,
    create_date timestamp,
    last_update timestamp,
    hex text
);

CREATE TABLE statistics (
    id serial,
    primary key(id),
    active boolean,
    name text,
    compute_query text,
    create_date timestamp,
    last_update timestamp,
    description text,
    chart_data text,
    keys text,
    value float,
    execution_date timestamp
);

CREATE TABLE transactions (
    id serial,
    primary key(id),
    name text,
    active boolean,
    transaction_date TIMESTAMP,
    last_update TIMESTAMP,
    create_date TIMESTAMP,
    amount float,
    account_id int references accounts(id),
    to_account int references accounts(id),
    transaction_type text,
    dest_booked boolean,
    origin_booked boolean,
    description text,
    category_id references categories(id)
);
