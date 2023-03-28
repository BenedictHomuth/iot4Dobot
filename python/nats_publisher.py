import asyncio
import json
import time
from nats.aio.client import Client as NATS

async def publish_message(msg):

    event = {}
    event['data'] = "Python: " + str(msg)
    event_json = json.dumps(event).encode("utf-8")

    # Establish a connection to the NATS server
    nc = NATS()
    await nc.connect("nats://localhost:4222")

    # Publish a message on the "my-topic" subject
    await nc.publish("roboPos", event_json)

    # Wait for the message to be delivered
    # await asyncio.sleep(1)

    # Close the connection to the NATS server
    await nc.close()

counter = 0
while True:
    asyncio.run(publish_message(counter))
    counter += 1
    time.sleep(1)