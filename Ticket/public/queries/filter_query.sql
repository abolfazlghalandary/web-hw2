select *
from available_offers
where
$0 > $1 and
origin = $2 and
destination = $3 and
date_trunc('day', departure_local_time) = to_date($4,'YYYY-MM-DD') and
date_trunc('day', arrival_local_time) = to_date($5,'YYYY-MM-DD')
