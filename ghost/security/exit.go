package security

// Exit - The exit function to call when exiting the implant process. It takes care of all
// cleanup and security checks necessary for an as-secure-as-possible exit.
// As well, the behavior of this function can vary depending on ghost implant state, or on
// details decided by Wiregost users, or even be used to automate the process of persistence.
// This function can be called from anywhere in a ghost implant code.
// It also must not be misused with the RPC Exit command (for users): here, this function does not
// check permissions and other details.
func Exit() {

}
