using System;
using System.Collections.Generic;
using System.Threading.Tasks;

using Amazon.Lambda.Core;
using Amazon.Lambda.APIGatewayEvents;
using MongoDB.Driver;
using System.Linq;
using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using Newtonsoft.Json;

[assembly: LambdaSerializer(typeof(Amazon.Lambda.Serialization.Json.JsonSerializer))]

namespace Quotes
{
    public class Quote
    {
        [BsonId]
        public ObjectId Id { get; set; }

        [BsonElement("number")]
        public int Number { get; set; }

        [BsonElement("text")]
        public string Text { get; set; }

        [BsonElement("game")]
        public string Game { get; set; }

        [BsonElement("date")]
        public DateTime Date { get; set; }

        public override string ToString()
        {
            return $"Quote #{Number}: {Text} [{Game}] [{Date:dd-MM-yyyy}]";
        }
    }

    public class Function
    {
        private static string MongoDbConnection = Environment.GetEnvironmentVariable("MONGODB_CONNECTION");

        private static async Task<string> GetRandomQuote()
        {
            var client = new MongoClient(MongoDbConnection);
            var database = client.GetDatabase("quotes");
            var collection = database.GetCollection<Quote>("twitchuserid");

            var emptyFilter = Builders<Quote>.Filter.Empty;
            long count = await collection.CountDocumentsAsync(emptyFilter);

            var random = new Random();
            var nextInt = random.Next(Convert.ToInt32(count));
            var randomQuoteIdFilter = Builders<Quote>.Filter.Eq(nameof(Quote.Number).ToLower(), nextInt);

            var quotes = await collection.FindAsync<Quote>(randomQuoteIdFilter);
            var quote = quotes
                .ToList()
                .Single();

            return quote.ToString();
        }

        private static async Task<string> GetQuoteByNumber(int number)
        {
            var client = new MongoClient(MongoDbConnection);
            var database = client.GetDatabase("quotes");
            var collection = database.GetCollection<Quote>("twitchuserid");

            var randomQuoteIdFilter = Builders<Quote>.Filter.Eq(nameof(Quote.Number).ToLower(), number);

            var quotes = await collection.FindAsync<Quote>(randomQuoteIdFilter);
            var quote = quotes
                .ToList()
                .Single();

            return quote.ToString();
        }

        public async Task<APIGatewayProxyResponse> FunctionHandler(APIGatewayProxyRequest apigProxyEvent, ILambdaContext context)
        {
            string quote;

            if (apigProxyEvent.QueryStringParameters != null &&
                apigProxyEvent.QueryStringParameters.ContainsKey("number") &&
                int.TryParse(apigProxyEvent.QueryStringParameters["number"], out var number))
            {
                quote = await GetQuoteByNumber(number);
            }
            else
            {
                quote = await GetRandomQuote();
            }

            return new APIGatewayProxyResponse
            {
                Body = quote ?? "quote not found",
                StatusCode = 200,
                Headers = new Dictionary<string, string> { { "Content-Type", "text/plain" } }
            };
        }
    }
}
