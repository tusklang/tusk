package types

type OmmThread struct {
  Channel   chan Returner
  Alive     bool
  returned  Returner
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
    return ot.returned//return none
  }

  getter := <- ot.Channel
  ot.returned = getter
  ot.Alive = false
  return getter
}

func (ot OmmThread) Format() string {

  if ot.Alive {
    return "{Alive Thread}"
  } else {
    return "{Dead Thread}"
  }

}

func (ot OmmThread) Type() string {
  return "thread"
}
