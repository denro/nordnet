/*
	The nordnet package is a set of packages for communicating with the NEXT API.

	The library contains the following packages:

	The api package provides a wrapper to the REST-API.

	The feed package is an implementation for subscribing to the real-time events.

	The util package contans all models used by the packages as well as a function for generating credentials.

*/

package nordnet

import (
	_ "github.com/denro/nordnet/api"
	_ "github.com/denro/nordnet/feed"
	_ "github.com/denro/nordnet/util"
)
