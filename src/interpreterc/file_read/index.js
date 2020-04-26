const fs = require('fs')
, path = require('path');

var stdinBuffer = fs.readFileSync(0)
, dir = stdinBuffer.toString().trim();

//nested directory get
if (dir.endsWith('*->')) {

  try {

    var fileGet = (dir) => fs.readdirSync(dir).filter(f =>  f.endsWith('.omm')).map(f => dir + (dir.endsWith('/') ? '' : '/') + f).filter(f => fs.lstatSync(f).isFile());
    var dirGet = (dir) => fs.readdirSync(dir).map(d => dir + (dir.endsWith('/') ? '' : '/') + d).filter(d => fs.lstatSync(d).isDirectory()).map(getAll);

    function getAll(dir) {
      var dirs = dirGet(dir)
      , files = fileGet(dir);

      return [...dirs, ...files];
    }

    var files = getAll(dir.slice(0, -3)).flat(Infinity).map(f => fs.readFileSync(f, 'utf8'));

    console.log(JSON.stringify(files));
  } catch (e) {
    console.log(e)
    console.log('Error: cannot import directory ' + dir.slice(0, -3));
    process.exit(1);
  }
} else if (dir.endsWith('*')) {

  try {
    fs.readdir(dir.slice(0, -1), (err, files) => {

      try {
        if (err) {
          console.log('Error: cannot import directory ' + dir.slice(0, -1));
          process.exit(1);
        }

        files = files.filter(f => f.endsWith('.omm') && fs.lstatSync(dir.slice(0, -1) + f).isFile());

        console.log(JSON.stringify(files.map(f => fs.readFileSync(dir.slice(0, -1) + f, 'utf8'))));
      } catch {
        console.log('Error: cannot import directory ' + dir.slice(0, -1));
        process.exit(1);
      }
    });
  } catch {
    console.log('Error: cannot import directory ' + dir.slice(0, -1));
    process.exit(1);
  }
} else {

  try {
    console.log(JSON.stringify(
      [fs.readFileSync(dir, 'utf8')]
    ));
  } catch {
    console.log('Error: cannot import ' + dir);
    process.exit(1);
  }
}
