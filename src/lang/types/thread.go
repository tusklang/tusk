package types

type OmmThread struct {
  Channel chan Returner
  Alive   bool
}

func (ot *OmmThread) FromGoType(val chan Returner) {
  ot.Channel = val
  ot.Alive = true
}

func (ot OmmThread) ToGoType() chan Returner {
  return ot.Channel
}

func (ot OmmThread) IsAlive() bool {
  return ot.Alive
}

func (ot OmmThread) WaitFor() Returner {

  if !ot.Alive { //if it is not alive
    return Returner{} //return none
  }

  defer func() {
    ot.Alive = false
  }() //set the thread to killed once the function finishes

  getter := <- ot.Channel
  return getter
}

func (_ OmmThread) ValueFunc() {}
