CREATE TABLE category (
    id serial,
    primary key(id),
    name text,
    red smallint,
    green smallint,
    blue smallint,
    create_date timestamp,
    last_update timestamp,
    active boolean,
    inner_transaction boolean
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
    transaction_type text
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