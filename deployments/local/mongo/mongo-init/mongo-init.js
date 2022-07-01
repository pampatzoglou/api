conn = new Mongo();
db = conn.getDB(process.env.MONGODB_DATABASE);
db.auth(process.env.MONGODB_USERNAME, process.env.MONGODB_PASSWORD)
db.log.insertOne({"message": "Database created."});

db.createCollection('shop')
db.shop.createIndex({ "address.zip": 1 }, { unique: false });
//id: ID!
///name: String!
//banner:
//  base64: String
//  url: String
//homepage: String
//location:
//  street: String!
//  number: String!
//  town: String!
//  coordinates:
///    latitude: String
//    longitude: String
//telephone: [Int]

db.shop.createIndex({ "address.zip": 1 }, { unique: false });
db.shop.insertOne({ "address": { "city": "Paris", "zip": "123" }, "name": "Mike", "phone": "1234" });
db.shop.insertOne({ "address": { "city": "Marsel", "zip": "321" }, "name": "Helga", "phone": "4321" });


db.createCollection('owner')
db.createCollection('item')
db.createCollection('client')



