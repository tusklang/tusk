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

    if v.Type == "function" {

      var fn = v.Value.(OmmFunc)
      fn.Body = insert_garbage_collectors(fn.Body)
      actions[k].Value = fn

      continue
    }
    if v.Type == "proto" {

      prototype := v.Value.(OmmProto)

      for i := range prototype.Static {

        var val = *prototype.Static[i]

        if val.Type() == "function" {
          var fn = val.(OmmFunc)
          fn.Body = insert_garbage_collectors(fn.Body)
          *prototype.Static[i] = fn
        }

      }

      for i := range prototype.Instance {

        var val = *prototype.Instance[i]

        if val.Type() == "function" {
          var fn = val.(OmmFunc)
          fn.Body = insert_garbage_collectors(fn.Body)
          *prototype.Static[i] = fn
        }

      }

      actions[k].Value = prototype

      continue
    }

    //perform insert_garbage_collectors on all of the sub actions
    actions[k].ExpAct = insert_garbage_collectors(v.ExpAct)
    actions[k].First = insert_garbage_collectors(v.First)
    actions[k].Second = insert_garbage_collectors(v.Second)

    //also do it for the (runtime) arrays and hashes
    for i := range v.Array {
      v.Array[i] = insert_garbage_collectors(v.Array[i])
    }
    for i := range v.Hash {
      v.Hash[i] = insert_garbage_collectors(v.Hash[i])
    }
    ////////////////////////////////////////////////

    /////////////////////////////////////////////////////////////

  }

  return inserted
}
