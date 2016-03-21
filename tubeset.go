package beanstalk

import (
	"time"
)

// TubeSet represents a set of tubes on the server connected to by Conn.
// Name names the tubes represented.
type TubeSet struct {
	Conn *Conn
	Name map[string]bool
}

// NewTubeSet returns a new TubeSet representing the given names.
func NewTubeSet(c *Conn, name ...string) *TubeSet {
	ts := &TubeSet{c, make(map[string]bool)}
	for _, s := range name {
		ts.Name[s] = true
	}
	return ts
}

// Reserve reserves and returns a job from one of the tubes in t. If no
// job is available before time timeout has passed, Reserve returns a
// ConnError recording ErrTimeout.
//
// Typically, a client will reserve a job, perform some work, then delete
// the job with Conn.Delete.
func (t *TubeSet) Reserve(timeout time.Duration) (id uint64, body []byte, err error) {
	r, err := t.Conn.cmd(nil, t, nil, "reserve-with-timeout", dur(timeout))
	if err != nil {
		return 0, nil, err
	}
	body, err = t.Conn.readResp(r, true, "RESERVED %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}


// Reserve reserves and returns a job from one of the tubes in t.
//
// Typically, a client will reserve a job, perform some work, then delete
// the job with Conn.Delete.
func (t *TubeSet) ReserveNoTimeout() (id uint64, body []byte, err error) {
	r, err := t.Conn.cmd(nil, t, nil, "reserve")
	if err != nil {
		return 0, nil, err
	}
	body, err = t.Conn.readResp(r, true, "RESERVED %d", &id)
	if err != nil {
		return 0, nil, err
	}
	return id, body, nil
}

// The beanstalk.Container implementation.
type box struct {
   conn *Conn
   id uint64
   body []byte
}

func (b *box) ID() uint64 { return b.id }

func (b *box) Body() []byte { return b.body }

func (b *box) Access( i Item ) error { return i.FromByteArray( b.body ) }

func (b *box) Delete() error { return b.conn.Delete( b.id ) }

func (b *box) Release( priority uint32, delay time.Duration ) error {
   return b.conn.Release( b.id, priority, delay )
}

func (b *box) Touch() error { return b.conn.Touch( b.id ) }

func (t *TubeSet) ReserveItem( timeout time.Duration ) (container Container, err error) {
   an_id,the_body,err := t.Reserve( timeout )
   if err != nil {
      return nil,err
   }
   b := &box{
      conn : t.Conn,
      id : an_id,
      body : the_body,
   }
   return b,nil
}

func (t *TubeSet) ReserveItemNoTimeout() (container Container, err error) {
   an_id,the_body,err := t.ReserveNoTimeout()
   if err != nil {
      return nil,err
   }
   b := &box{
      conn : t.Conn,
      id : an_id,
      body : the_body,
   }
   return b,nil
}

