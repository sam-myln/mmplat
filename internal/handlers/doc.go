/*
Package handlers consist of Auth, Handler parts.
Handler responsible for routing, binding m-wares, etc.
Notably, w/ Auth, process consists of following steps:
1. show login page, receive POST
2. set `userAuth` session
3. redirect to totp page
4. receive POST w/ totp
5. set `userVerif` redirect
6. finally, set finalToken
*/
package handlers
