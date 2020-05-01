const fs = require('fs')
, read = require('./read');

var stdinBuffer = fs.readFileSync(0)
, d = stdinBuffer.toString().trim();

console.log(read(d));
