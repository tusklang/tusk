package compiler

import . "lang/types"

/*
  Omm's garbage collector works like so:

    var main: fn() {
      var a: 3
      log a
    }

  Now we have to remove the variable `a`
  Omm does this by deferring it (in go) to delete it at the end of the group

    var main: fn() {
      var a: 3
      defer delete a ;defer is a go keyword to call a function after a `return`
      log a
      return a
    }
*/

func insert_garbage_collectors(actions []Action) []Action {

  var inserted []Action

  for k, v := range actions {

    inserted = append(inserted, v)

    if v.Type == "var" || v.Type == "declare" {
      inserted = append(inserted, Action{
        Type: "delete",
        Name: v.Name,
      })
    }

    //perform insert_garbage_collectors on all of the sub actions
    actions[k].ExpAct = insert_garbage_collectors(v.ExpAct)
    actions[k].First = insert_garbage_collectors(v.First)
    actions[k].Second = insert_garbage_collectors(v.Second)
    /////////////////////////////////////////////////////////////

  }

  return inserted
}
