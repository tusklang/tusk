//file to remove whitespace from a string (not within a string)

module.exports = str => {

  //a new string to return
  var nStr = '';

  //quote type
  var qType = null
  //if it is an escaped char
  , escaped = false;

  for (let i of str) {

    //if it is an escaped character
    if (!escaped && i == '\\') {
      nStr+='\\';
      escaped = true;
      continue;
    }

    //if the current character is not escaped
    //if the current character is a quote then set qType to the current quote type if qType is null
    //otherwise is qType is not null, set qType to the current quote type
    if (!escaped && /\'|\"|\`/.test(i)) qType = qType ? null : i;

    //if it not in a quote
    if (!qType) {

      //if the current char is whitespace then continue
      if (/\s/.test(i)) continue;

      //otherwise add the current char to nStr
      nStr+=i;

    } else nStr+=i; //if it is a string, then add the current char no matter what
  }

  //just using .trim() for extra precautions
  return nStr.trim();
}
