#include <vector>
#include <map>
#include "structs.h"
#include "json.hpp"
using namespace std;
using json = nlohmann::json;

json hashIndex(json valJ, json indexes, json calc_params, int line);
json arrayIndex(json valJ, json indexes, json calc_params, int line);
json expressionIndex(json val, json indexes);
json math(json exp, const json calc_params, json vars, const string dir, int line);
Returner parser(const json actions, const json calc_params, json vars, const string dir, bool groupReturn, int line);
json indexesCalc(json val, json indexes, json calc_params, int line);

json hashIndex(json valJ, json indexes, json calc_params, int line) {

  valJ = valJ[0];

  valJ.erase(valJ.begin());
  valJ.erase(valJ.end());

  vector<json> vals { json::parse("[]") };

  int bCnt = 0
  , glCnt = 0;

  for (int i = 0; i < valJ.size(); i++) {
    if (valJ[i] == "[") bCnt++;
    if (valJ[i] == "]") bCnt--;
    if (valJ[i] == "[:") glCnt++;
    if (valJ[i] == ":]") glCnt--;

    if (valJ[i] == "newlineN") continue;

    if (valJ[i] == "," && bCnt == 0 && glCnt == 0) {
      vals.push_back(json::parse("[]"));
      continue;
    }

    vals[vals.size() - 1].push_back(valJ[i]);
  }

  map<string, json> valMap;

  for (int i = 0; i < vals.size(); i++) {

    int bCnt = 0
    , glCnt = 0;

    pair<string, json> curVal = pair<string, json>("", "[]"_json);
    int cur = 0;

    for (int o = 0; o < vals[i].size(); o++) {
      if (vals[i][o] == "[") bCnt++;
      if (vals[i][o] == "]") bCnt--;
      if (vals[i][o] == "[:") glCnt++;
      if (vals[i][o] == ":]") glCnt--;

      if (vals[i][o] == ":" && bCnt == 0 && glCnt == 0) {
        cur = 1;
        continue;
      }

      if (cur == 0) curVal.first+=vals[i][o];
      else curVal.second.push_back(vals[i][o]);
    }

    if (curVal.first.rfind("\'", 0) == 0 || curVal.first.rfind("\"", 0) == 0 || curVal.first.rfind("`", 0) == 0)
      curVal.first = curVal.first.substr(1, curVal.first.length() - 2);

    valMap.insert(curVal);
  }

  json _curIndex = indexes[0];

  _curIndex.erase(_curIndex.begin());
  _curIndex.erase(_curIndex.end());

  string curIndex = (string) _curIndex[0];

  if (curIndex.rfind("\'", 0) == 0 || curIndex.rfind("\"", 0) == 0 || curIndex.rfind("`", 0) == 0)
    curIndex = curIndex.substr(1, curIndex.length() - 2);

  json val = valMap[curIndex];

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return val;

  return indexesCalc(json::parse("[" + val.dump() + "]"), indexes, calc_params, line);
}

json arrayIndex(json valJ, json indexes, json calc_params, int line) {
  valJ = valJ[0];

  valJ.erase(valJ.begin());
  valJ.erase(valJ.end());

  vector<json> vals { json::parse("[]") };

  int bCnt = 0
  , glCnt = 0;

  for (int i = 0; i < valJ.size(); i++) {
    if (valJ[i] == "[") bCnt++;
    if (valJ[i] == "]") bCnt--;
    if (valJ[i] == "[:") glCnt++;
    if (valJ[i] == ":]") glCnt++;

    if (valJ[i] == "newlineN") continue;

    if (valJ[i] == "," && bCnt == 0 && glCnt == 0) {
      vals.push_back(json::parse("[]"));
      continue;
    }

    vals[vals.size() - 1].push_back(valJ[i]);
  }

  json _curIndex = indexes[0];

  _curIndex.erase(_curIndex.begin());
  _curIndex.erase(_curIndex.end());

  char* curIndex = &(_curIndex[0].dump())[0];
  char* zero = "0";

  while (IsLessC(zero, curIndex)) {
    vals.erase(vals.begin());

    string cp = calc_params.dump();

    curIndex = Subtract(curIndex, "1", &cp[0], line);
  }

  json first = vals[0];

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return first;

  return indexesCalc(json::parse("[" + first.dump() + "]"), indexes, calc_params, line);
}

json expressionIndex(json valJ, json indexes, json calc_params, int line) {

  cout << indexes << endl;

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

  val.erase(0, 2);
  val.erase(val.length() - 2, 2);

  //getting the position of the index
  while (IsLessC(zero, indexC)) {

    if (val.length() == 0) return json::parse("[\"undefined\"]");

    val = val.substr(1);

    indexC = &(Subtract(indexC, "1", cp, line))[0];
  }

  string first(1, val[0]);

  indexes.erase(indexes.begin());

  if (indexes.size() == 0) return json::parse("[\"\'" + first + "\'\"]");

  json nIndex = indexes
  , nCP = calc_params;
  int nL = line;

  json firstR = json::parse("[[\"\\\'" + first + "\\\'\"]]");

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
