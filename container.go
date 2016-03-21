// The Beanstalk Container makes it easier to deal with an item
// received from a tube. It packages the connection and the item
// identifier along with the body of the item. Then the user can
// decide what to do with the item w/o having to manually track
// which connection originated the item, which is necessary to
// properly disposition it.
package beanstalk

import (
	"time"
)

// Any object that ends up in a container should implement this to
// get out of it. Alternatively, see Container.Body().
type Item interface {
   FromByteArray( []byte ) error
}

// The container of some Item. The user should call either
// Delete() or Release() once done with the Item. If neither,
// then the behavior is governed by the tube and/or beanstalk.
type Container interface {
   // Returns the Beanstalk ID of the item.
   ID() uint64

   // Returns the body of the item received.
   Body() []byte

   // Performs the conversion necessary to obtain the item via
   // the FromByteArray() method.
   Access( i Item ) error

   // Tells beanstalk that the Item was consumed.
   Delete() error

   // Tells beanstalk that the Item is being abandoned.
   Release( priority uint32, delay time.Duration ) error

   // Tells beanstalk that the Item is still being worked.
   Touch() error
}
