import asyncio
from nats.aio.client import Client as NATS

async def run(loop):
    # Verbindung zu NATS herstellen
    nc = NATS()
    await nc.connect(servers=["nats://demo.nats.io:4222"], loop=loop)

    # Nachricht an ein Topic senden
    topic = "my_topic"
    message = "Hello, NATS!"
    await nc.publish(topic, message.encode())

    # Verbindung schließen
    await nc.close()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))