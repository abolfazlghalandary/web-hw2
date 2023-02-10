const { Pool, Client } = require('pg')
const fs = require('fs');

const clientData = {
    user: 'postgres',
    host: 'localhost',
    database: 'dev',
    password: 'dev',
    port: 5432,
}

async function filter_flights(params) {

    let {origin, destination, arrivalDate, departureDate, number, cls} = params

    const client = new Client(clientData)
    await client.connect()
    let rawQuery = fs.readFileSync('public/queries/filter_query.sql', 'utf8');

    let queryInputs = [number, origin, destination, departureDate, arrivalDate]
    rawQuery = rawQuery.replace('$0', cls + '_class_free_capacity')
    let result = await client.query(rawQuery, queryInputs);
    await client.end();

    return result.rows;
}

async function validate_buy(params){
    let {
        number,
        flight_id,
        offer_class
    } = params

    const client = new Client(clientData)
    await client.connect()

    let rawQuery0 = `select flight_serial from flight where flight_id = '${flight_id}'`
    let result0 = await client.query(rawQuery0);
    if(!result0.rows.length){
        return 'invalid flight id'
    }
    let flight_serial = result0.rows[0]['flight_serial']

    let rawQuery1 = fs.readFileSync('public/queries/get_flight_capacity.sql', 'utf8');
    rawQuery1 = rawQuery1.replace("$0", offer_class  + '_class_free_capacity')
    rawQuery1 = rawQuery1.replace("$1", flight_id)

    let rawQuery2 = fs.readFileSync('public/queries/get_bought_tickets_count.sql', 'utf8');


    let capacity_result = await client.query(rawQuery1)
    let bought_number_result = await client.query(rawQuery2, [flight_serial, offer_class])
    await client.end()

    let capacity = capacity_result.rows[0][offer_class + "_class_free_capacity"]
    let bought_number =  bought_number_result.rows[0]['count'];

    return capacity - bought_number >= number ? null : 'there is no enough capacity';
}

async function buy_ticket(params,userId) {
    let {
        title,
        first_name,
        last_name,
        flight_id,
        offer_price,
        offer_class,
        transaction_id,
        transaction_result
    } = params

    const client = new Client(clientData)

    await client.connect()

    let rawQuery1 = `select flight_serial from flight where flight_id = '${flight_id}'`
    let result1 = await client.query(rawQuery1);
    if(!result1.rows.length){
        return 'invalid flight id'
    }
    let flight_serial = result1.rows[0]['flight_serial']


    let queryInputs = [
        userId,
        title,
        first_name,
        last_name,
        flight_serial,
        offer_price,
        offer_class,
        transaction_id,
        transaction_result
    ]
    let rawQuery2 = fs.readFileSync('public/queries/buy_ticket.sql', 'utf8');
    await client.query(rawQuery2, queryInputs);

    await client.end();
    return ''
}


async function get_user_tickets(user_id){
    
    const client = new Client(clientData)
    await client.connect()
    let rawQuery = fs.readFileSync('public/queries/get_user_tickets.sql', 'utf8');

    let result = await client.query(rawQuery, [user_id]);
    await client.end();

    return result.rows;
}

module.exports = { filter_flights, buy_ticket, validate_buy, get_user_tickets };

