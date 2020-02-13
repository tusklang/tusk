const http = require('http')
, express = require('express')
, url = require('url')
, app = express();

const server = http.createServer(app);

app.use(express.static(__dirname + '/'));

app.get('/', (req, res) => {
  res.sendFile('/index.html');
});

app.get('/docs', (req, res) => {
  var page = url.parse(req.url).query.split('=')[1];

  if (!page) return res.redirect('/');

  res.sendFile(__dirname + '/docs/' + page);
});

server.listen(process.env.PORT || 80);
