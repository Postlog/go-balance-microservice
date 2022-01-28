create table if not exists balance
(
    user_id uuid           not null
        constraint balances_pk
            primary key,
    value   numeric(32, 2) not null
        constraint balances_balance_check
            check ((value)::double precision >= (0)::double precision)
);

create table if not exists transaction
(
    sender_id   uuid
        constraint transactions_balances_user_id_fk
            references balance
            on update cascade on delete cascade,
    receiver_id uuid
        constraint transactions_balances_user_id_fk_2
            references balance
            on update cascade on delete cascade,
    amount      numeric(32, 2) not null,
    description text           not null,
    id          serial
        constraint transactions_pk
            primary key,
    date        timestamp      not null,
    constraint both_not_null
        check ((sender_id IS NOT NULL) OR (receiver_id IS NOT NULL))
);


