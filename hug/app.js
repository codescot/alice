exports.lambdaHandler = async (event, context) => {
    context.callbackWaitsForEmptyEventLoop = false;
    const body = JSON.parse(event.body);
    
    const from = body.from;
    const to = body.to;

    return `${from} gives ${to} an awkwardly long and warm hug. <3`;
};
