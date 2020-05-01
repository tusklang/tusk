const fs = require('fs');

module.exports = dir => {

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

      return JSON.stringify(files);
    } catch (e) {

      return 'Error: cannot import directory ' + dir.slice(0, -3);
    }
  } else if (dir.endsWith('*')) {

    try {
      var files = fs.readdirSync(dir.slice(0, -1))

      files = files.filter(f => f.endsWith('.omm') && fs.lstatSync(dir.slice(0, -1) + f).isFile());

      return JSON.stringify(files.map(f => fs.readFileSync(dir.slice(0, -1) + f, 'utf8')));
    } catch {

      return 'Error: cannot import directory ' + dir.slice(0, -1);
    }
  } else {

    try {
      return JSON.stringify(
        [fs.readFileSync(dir, 'utf8')]
      );
    } catch (e) {

      return 'Error: cannot import ' + dir;
    }
  }
}
