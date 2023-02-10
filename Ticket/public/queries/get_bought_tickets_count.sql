select
    count(*)
from
    purchase
where
    flight_serial = $1 and offer_class = $2 and transaction_result = 1
