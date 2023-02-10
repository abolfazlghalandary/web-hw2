let express = require('express');
var request = require('request');

const {filter_flights, buy_ticket, validate_buy, get_user_tickets} = require("../public/javascripts/logic");
let router = express.Router();

router.post('/filter_flights', async function (req, res, next) {
  res.json({
    flights: await filter_flights(req.body)
  });
});

router.put('/buy_ticket', async function (req, res, next) {
  const token = req.header("Authorization");
  let userId;
  request.get(
    {
      url:"http://localhost:8080/auth/user_info",
      headers : {"Authorization":token}},
      async function(error, response,body){
      userId = JSON.parse(body).user.ID;
      res.status(await buy_ticket(req.body) ? 300 : 200).send()
  })
});

router.get('/validate_buy', async function (req, res, next) {
  let message = await validate_buy(req.body)
  res.status(message ? 400 : 200).send(message);
});

router.get('/get_user_tickets', async function (req, res, next) {
  res.json({
    tickets : await get_user_tickets(req.query)
  })
});

module.exports = router;
