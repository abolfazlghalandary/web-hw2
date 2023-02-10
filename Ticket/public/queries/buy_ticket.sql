INSERT INTO purchase (
    corresponding_user_id,
    title,
    first_name,
    last_name,
    flight_serial,
    offer_price,
    offer_class,
    transaction_id,
    transaction_result
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
