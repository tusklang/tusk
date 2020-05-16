#ifndef READDIR_HPP_
#define READDIR_HPP_

#include <vector>
#include <algorithm>
#include <deque>
#include <dirent.h>
using namespace std;

vector<string> read_dir(string dir_str) {

  DIR *dir;
  struct dirent *en;
  dir = opendir(dir_str.c_str());

  deque<string> dirs;

  if (dir) {

     while ((en = readdir(dir)) != NULL) dirs.push_back(en->d_name);

     closedir(dir);
  }

  //remove the first two dirs
  //the first two dirs are always . and ..
  dirs.pop_front();
  dirs.pop_front();
  ////////////////////////////////////////

  vector<string> dirv;

  copy(dirs.begin(), dirs.end(), back_inserter(dirv));

  vector<string> copiedDirV;

  for (string it : dirv) copiedDirV.push_back(dir_str + it);

  return copiedDirV;
}

#endif
