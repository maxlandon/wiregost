package transport

// Route - Should be probably be deleted, because we might not need it
// Route is somewhat different from a Chain object: the Route object does
// NOT keep state of any transport mechanism, which means that the Route is also the object
// that is stored in the database, displayed by users. They can also modify them, clone them,
// and use them as models. Then, the Route is parsed into a Chain object, which will actually
// implement the routed/proxied/forwarded communication.
type Route struct {
}
