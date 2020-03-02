#include <iostream>
#include <vector>
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

json hashIndex(json val, json indexes);
json arrayIndex(json val, json indexes);
json expressionIndex(json val, json indexes);
json indexesCalc(json val, json indexes, json calc_params, int line);
json math(json exp, const json calc_params, json vars, const string dir, int line);

json hashIndex(json val, json indexes, json calc_params, int line) {
  return "[]"_json;
}

json arrayIndex(json val, json indexes, json calc_params, int line) {

  json inner = val[0];

  inner.erase(inner.begin());
  inner.erase(inner.end());

  int bCnt = 0
  , cbCnt = 0;
  vector<string> arr {""};

  for (int i = 0; i < inner.size(); i++) {
    if (inner[i] == "{") cbCnt++;
    if (inner[i] == "}") cbCnt--;
    if (inner[i] == "[") bCnt++;
    if (inner[i] == "]") bCnt--;

    if (cbCnt == 0 && bCnt == 0 && inner[i] == ",") {
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

    if (arrJ.size() == 0) return json::parse("[\"undefined\"]");

    arrJ.erase(arrJ.begin());

    indexC = &(Subtract(indexC, "1", cp, line))[0];
  }

  string first = arrJ[0].dump();
  first.erase(0, 3);
  first.erase(first.length() - 3, 3);
  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return json::parse("[\"" + first + "\"]");

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

    if (val.length() == 0) return json::parse("[\"undefined\"]");

    val = val.substr(1);

    indexC = &(Subtract(indexC, "1", cp, line))[0];
  }

  string first(1, val[0]);

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return json::parse("[\"\\\"" + first + "\\\"\"]");

  json nIndex = indexes
  , nCP = calc_params;
  int nL = line;

  json firstR = json::parse("[[\"\\\"" + first + "\\\"\"]]");

  return indexesCalc(firstR, nIndex, nCP, nL);
}


json indexesCalc(json val, json indexes, json calc_params, int line) {

  string valDump = val.dump();

  char* datatype = GetType(&valDump[0]);

  if (strcmp(datatype, "hash") == 0) return hashIndex(val, indexes, calc_params, line);
  else if (strcmp(datatype, "array") == 0) return arrayIndex(val, indexes, calc_params, line);
  else if (strcmp(datatype, "string") == 0) return expressionIndex(val, indexes, calc_params, line);

  return json::parse("[\"undefined\"]");
}
