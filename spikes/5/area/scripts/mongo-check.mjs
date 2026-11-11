import { MongoClient } from 'mongodb';

const uri = process.env.MONGODB_URI;
if (!uri) {
  throw new Error('MONGODB_URI environment variable is required');
}

const client = new MongoClient(uri);

try {
  await client.connect();
  const db = client.db('area');
  const probes = db.collection('probes');
  const testedAt = new Date().toISOString();

  await probes.insertOne({ testedAt });
  const found = await probes.findOne({ testedAt });

  if (!found) {
    throw new Error('inserted document was not read back');
  }

  console.log(JSON.stringify({ status: 'ok', testedAt: found.testedAt }));
} finally {
  await client.close();
}
