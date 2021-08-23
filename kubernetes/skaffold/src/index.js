const express = require('express');

const app = express();

app.get('/', (request, response) => {
  response.status(200).send({ message: 'Hello Skaffold!' });
});

app.listen(4000, () => {
  console.log('Server is running!');
});
