#include <iostream>
#include <vector>
#include <map>
#include "structs.h"
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

json hashIndex(json val, json indexes);
json arrayIndex(json val, json indexes);
json expressionIndex(json val, json indexes);
json math(json exp, const json calc_params, json vars, const string dir, int line);
Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line);
json indexesCalc(json val, json indexes, json calc_params, int line);

json hashIndex(json val, json indexes, json calc_params, int line) {

  json inner = val[0];

  inner.erase(inner.begin());
  inner.erase(inner.end());

  int bCnt = 0
  , glCnt = 0
  , cur = false;
  vector<vector<string>> hash {
    {"", ""}
  };

  for (int i = 0; i < inner.size(); i++) {
    if (inner[i] == "[:") glCnt++;
    if (inner[i] == ":]") glCnt--;
    if (inner[i] == "[") bCnt++;
    if (inner[i] == "]") bCnt--;

    if (inner[i] == "newlineN") continue;

    if (glCnt == 0 && bCnt == 0 && inner[i] == ",") {

      vector<string> pusher {"", ""};

      hash.push_back(pusher);
      continue;
    }

    if (glCnt == 0 && bCnt == 0 && inner[i] == ":") {
      cur = !cur;
      continue;
    }

    hash[hash.size() - 1][cur]+=inner[i];
  }

  map<string, string> hashMap;

  for (int i = 0; i < hash.size(); i++)
    hashMap.insert(pair<string, string>(hash[i][0].substr(0, hash[i][0].length() - 1).substr(1), hash[i][1]));

  json index = indexes[0];

  //removing [ and ]
  index.erase(index.begin());
  index.erase(index.end());

  string indexD = index[0].dump();

  //erase the quotes
  indexD.erase(0, 3);
  indexD.erase(indexD.length() - 3, 3);

  string value = hashMap[indexD];

  cout << indexes.dump() << endl;

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return json::parse("[\"" + value + "\"]");

  char* datatype = GetType(&value[0]);

  string valueN;

  if (strcmp(datatype, "hash") == 0) {

    string valueCopy = value.substr(2, value.length() - 4);

    vector<string> hashCopied {""};

    int bCnt = 0
    , glCnt = 0;

    for (int i = 0; i < valueCopy.length(); i++) {
      if (valueCopy[i] == '{') glCnt++;
      if (valueCopy[i] == '}') glCnt--;
      if (valueCopy[i] == '[') bCnt++;
      if (valueCopy[i] == ']') bCnt--;

      string curC(1, valueCopy[i]);

      if (bCnt == 0 && glCnt == 0 && valueCopy[i] == ':') {
        hashCopied.push_back(":");
        hashCopied.push_back("");

        continue;
      }

      if (bCnt == 0 && glCnt == 0 && valueCopy[i] == ',') {
        hashCopied.push_back(",");
        hashCopied.push_back("");

        continue;
      }

      hashCopied[hashCopied.size() - 1]+=valueCopy[i];
    }

    hashCopied.insert(hashCopied.begin(), "[:");
    hashCopied.push_back(":]");

    string puts = "";

    for (int i = 0; i < hashCopied.size(); i++) puts+=("\"" + hashCopied[i] + "\",");

    valueN = puts;

    cout << valueN << endl;
  } else if (strcmp(datatype, "array") == 0) {

    string valCopy = value.substr(1, value.length() - 2);
    vector<string> valVect {""};

    int bCnt = 0
    , glCnt = 0;

    for (int i = 0; i < valCopy.length(); i++) {
      if (valCopy[i] == '{') glCnt++;
      if (valCopy[i] == '}') glCnt--;
      if (valCopy[i] == '[') bCnt++;
      if (valCopy[i] == ']') bCnt--;

      string cur(1, valCopy[i]);

      if (bCnt == 0 && glCnt == 0 && valCopy[i] == ',') {
        valVect.push_back(",");
        valVect.push_back("");
        continue;
      } else valVect[valVect.size() - 1]+=cur;
    }

    string valNPut = "";

    for (int i = 0; i < valVect.size(); i++) valNPut+=("\"" + valVect[i] + "\",");

    valueN = "\"[\"," + valNPut + "\"]\"";
  } else if (strcmp(datatype, "string") == 0) {
    valueN = "\"\\\'" + value.substr(1, value.length() - 2) + "\\\'\"";
  }

  json valueR = json::parse("[[" + valueN + "]]");

  return indexesCalc(valueR, indexes, calc_params, line);
}

json arrayIndex(json val, json indexes, json calc_params, int line) {

  json inner = val[0];

  inner.erase(inner.begin());
  inner.erase(inner.end());

  int bCnt = 0
  , glCnt = 0;
  vector<string> arr {""};

  for (int i = 0; i < inner.size(); i++) {
    if (inner[i] == "[:") glCnt++;
    if (inner[i] == ":]:") glCnt--;
    if (inner[i] == "[") bCnt++;
    if (inner[i] == "]") bCnt--;

    if (inner[i] == "newlineN") continue;

    if (glCnt == 0 && bCnt == 0 && inner[i] == ",") {
      arr.push_back("");
      continue;
    }

    arr[arr.size() - 1] = arr[arr.size() - 1] + inner[i].dump();
  }

  json arrJ = arr;

  json index = indexes[0];

  //erase the quotes
  index.erase(index.begin());
  index.erase(index.end());

  string indexD = index[0].dump();

  //erase the quotes
  indexD.erase(0, 1);
  indexD.erase(indexD.length() - 1, 1);

  char* indexC = &(indexD)[0];
  char* cp = &(calc_params.dump())[0];
  char* zero = "0";

  //getting the position of the index
  while (IsLessC(zero, indexC)) {

    if (arrJ.size() == 0) return json::parse("[\'undefined\']");

    arrJ.erase(arrJ.begin());

    indexC = &(Subtract(indexC, "1", cp, line))[0];
  }

  string first = arrJ[0].dump();
  first.erase(0, 3);
  first.erase(first.length() - 3, 3);
  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return json::parse("[\'" + first + "\']");

  json firstR = json::parse("[[" + first + "]]");

  return indexesCalc(firstR, indexes, calc_params, line);
}

json expressionIndex(json valJ, json indexes, json calc_params, int line) {

  string val = valJ[0][0].dump()
  , expVal = "";

  for (int i = 3; i < val.length() - 3; i++) expVal+=val[i];

  json index = indexes[0];

  //erase the quotes
  index.erase(index.begin());
  index.erase(index.end());

  string indexD = index[0].dump();

  //erase the quotes
  indexD.erase(0, 1);
  indexD.erase(indexD.length() - 1, 1);

  char* indexC = &(indexD)[0];
  char* cp = &(calc_params.dump())[0];
  char* zero = "0";

  //erase the quotes
  val.erase(0, 3);
  val.erase(val.length() - 3, 3);

  //getting the position of the index
  while (IsLessC(zero, indexC)) {

    if (val.length() == 0) return json::parse("[\'undefined\']");

    val = val.substr(1);

    indexC = &(Subtract(indexC, "1", cp, line))[0];
  }

  string first(1, val[0]);

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return json::parse("[\'\\\'" + first + "\\\'\']");

  json nIndex = indexes
  , nCP = calc_params;
  int nL = line;

  json firstR = json::parse("[[\'\\\'" + first + "\\\'\']]");

  return indexesCalc(firstR, nIndex, nCP, nL);
}


json indexesCalc(json val, json indexes, json calc_params, int line) {

  string getTypeVal = val[0][0];

  char* datatype = GetType(&getTypeVal[0]);

  if (strcmp(datatype, "hash") == 0) return hashIndex(val, indexes, calc_params, line);
  else if (strcmp(datatype, "array") == 0) return arrayIndex(val, indexes, calc_params, line);
  else if (strcmp(datatype, "string") == 0) return expressionIndex(val, indexes, calc_params, line);

  return json::parse("[\'undefined\']");
}
