package assets

// This file contains a build-time server configuration (address, certificates, user, etc)
// When a console is compiled from the server, it may decide which values to inject

// HasBuiltinServer - If this is different from "", it means we use the values below
var HasBuiltinServer string

// ServerLHost - Host of server
var ServerLHost string

// ServerLPort - Port on which to contact server
var ServerLPort string

// ServerUser - Username
var ServerUser string

// ServerCACertificate - CA Certificate
var ServerCACertificate string

// ServerCertificate - CA Certificate
var ServerCertificate string

// ServerPrivateKey - Private key
var ServerPrivateKey string

// Token - A unique number for this client binary
var Token string
