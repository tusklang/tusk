package types

type OmmThread struct {
  Channel    chan Returner
  WasJoined  bool
  Returned   Returner
}

func (ot *OmmThread) FromGoType(val chan Returner) {
  ot.Channel = val
  ot.WasJoined = true
}

func (ot OmmThread) ToGoType() chan Returner {
  return ot.Channel
}

func (ot *OmmThread) WaitFor() Returner {

  if ot.WasJoined { //if it is not alive
    return ot.Returned //return the already given value
  }

  getter := <- ot.Channel
  ot.Returned = getter
  ot.WasJoined = true
  return getter
}

func (ot OmmThread) Format() string {

  if ot.WasJoined {
    return "{Joined Thread}"
  } else {
    return "{Detached Thread}"
  }

}

func (ot OmmThread) Type() string {
  return "thread"
}

func (ot OmmThread) TypeOf() string {
  return ot.Type()
}
