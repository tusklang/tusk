const fs = require('fs');

var stdinBuffer = fs.readFileSync(0)
, dir = stdinBuffer.toString().trim();

if (dir.endsWith("*")) {

  try {
    fs.readdir(dir.slice(0, -1), (err, files) => {
      files = files.filter(f => f.endsWith('.omm'));

      console.log(JSON.stringify(files.map(f => fs.readFileSync(dir.slice(0, -1) + f, 'utf8'))));
    });
  } catch {
    console.log('Error: cannot import directory ' + dir.slice(0, -1));
  }
} else {

  try {
    console.log(JSON.stringify(
      [fs.readFileSync(dir, 'utf8')]
    ));
  } catch {
    console.log('Error: cannot import ' + dir);
  }
}
