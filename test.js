import amqp from 'k6/x/amqp';
var CONFIG = require('./config.js');

export default function () {
  const connString = CONFIG.connString;
  const topic = CONFIG.topic;

  let sender = new amqp.Sender(connString,topic);
  sender.Connect();
  let receivers = [];
  for (let i = 0; i < 100; i++) {
    let receiver = new amqp.Receiver(connString, topic);
    receiver.Connect();
    receivers.push(receiver);
  }

  sender.Send("extra-super brilliant message2");
  receivers.forEach(receiver => {
    console.log(receiver.Receive());
    receiver.Disconnect();
  });

  sender.Disconnect();
}
