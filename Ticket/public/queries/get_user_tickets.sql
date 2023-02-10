select * from
(
	SELECT filtered_purchase.*, flight.flight_id
	from
	(select * from purchase where purchase.corresponding_user_id = $1) as filtered_purchase
	INNER JOIN flight
	ON flight.flight_serial = filtered_purchase.flight_serial
) as filtered_purchase
INNER JOIN available_offers
ON available_offers.flight_id = filtered_purchase.flight_id
