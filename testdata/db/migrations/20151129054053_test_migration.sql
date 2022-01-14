-- migrate:up
select 1;
select * from posts;
insert into users (id, name) values (1, 'alic' ||
                                        'e;');
-- some comment with ;

select 1;

CREATE OR REPLACE FUNCTION public.edit_data_subscription_purshases (
)
    RETURNS void AS
$body$
DECLARE
    user_ids BIGINT[];
    max_expiration_at TIMESTAMPTZ;
    cur_user_id BIGINT;
    count_expiration_at BIGINT;
    count_expiration_at_is_null BIGINT;
    sp_id BIGINT;
    sp_ids BIGINT[];

BEGIN
    -- Checking and initializing column expiration_at

    user_ids := ARRAY(SELECT DISTINCT user_id::BIGINT
                      FROM subscription.subscription_purchases ORDER BY user_id);

    FOREACH cur_user_id IN ARRAY user_ids
        LOOP
            count_expiration_at :=
                    (SELECT count(*) FROM subscription.subscription_purchases
                     WHERE user_id = cur_user_id AND expiration_at IS NOT NULL)::BIGINT;

            count_expiration_at_is_null :=
                    (SELECT count(*) FROM subscription.subscription_purchases
                     WHERE user_id = cur_user_id AND expiration_at IS NULL)::BIGINT;

            IF count_expiration_at <> 0 AND count_expiration_at_is_null <> 0
            THEN
                max_expiration_at :=
                        (SELECT max(expiration_at) FROM subscription.subscription_purchases
                         WHERE user_id = cur_user_id AND expiration_at IS NOT NULL)::TIMESTAMPTZ;

                sp_ids :=
                        ARRAY(SELECT id FROM subscription.subscription_purchases
                              WHERE user_id = cur_user_id AND expiration_at IS NULL ORDER BY id);

                FOREACH sp_id IN ARRAY sp_ids
                    LOOP
                        max_expiration_at := max_expiration_at + interval '1 month';
                        UPDATE subscription.subscription_purchases
                        SET expiration_at = max_expiration_at WHERE id = sp_id;
                    END LOOP;
            END IF;
        END LOOP;

    -- Checking and initializing column start_at

    UPDATE subscription.subscription_purchases
    SET start_at = (expiration_at::DATE - interval '1 month' + interval '1 day')::TIMESTAMPTZ
    WHERE start_at IS NULL;

END;
$body$
    LANGUAGE 'plpgsql'
    VOLATILE
    CALLED ON NULL INPUT
    SECURITY INVOKER
    PARALLEL UNSAFE
    COST 100;

create table users (
                       id integer,
                       name varchar(255)
);

select 1;

-- migrate:down
drop table users;
