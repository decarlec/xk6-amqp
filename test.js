import amqp from 'k6/x/amqp';
import { randomBytes } from 'k6/crypto';

export default function () {
  const connString = __ENV.AMQP_CONN_STRING;
  const topic = __ENV.AMQP_TOPIC;

  let sender = new amqp.Sender(connString,topic);
  sender.Connect();
  let receivers = [];
  for (let i = 0; i < 20; i++) {
    let receiver = new amqp.Receiver(connString, topic);
    receiver.Connect();
    receivers.push(receiver);
  }

  sender.Send(new Uint32Array(randomBytes(1024)));
  receivers.forEach(receiver => {
    receiver.Receive();
  });

  receivers.forEach(receiver => {
    receiver.Disconnect();
  });

  sender.Disconnect();
}
