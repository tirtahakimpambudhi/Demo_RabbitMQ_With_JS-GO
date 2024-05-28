import amqp from "amqplib";
export const subscribe = async () => {
    const connection = await amqp.connect(process.env.MESSAGE_BROKER_CLOUD);
    const channel = await connection.createChannel();

    process.on('SIGINT', async () => {
        console.log("Received SIGINT. Closing connection...");
        await channel.close();
        await connection.close();
        process.exit(0);
    });

    process.on('SIGTERM', async () => {
        console.log("Received SIGTERM. Closing connection...");
        await channel.close();
        await connection.close();
        process.exit(0);
    });

    await channel.consume(process.env.QUEUE_PROCESS_SERVICE, (message) => {
        if (message) {
            const data = message.content.toString();
            console.table([{ queue: process.env.QUEUE_PROCESS_SERVICE, message: JSON.parse(data) }]);
        }
    }, { noAck: true });
};
