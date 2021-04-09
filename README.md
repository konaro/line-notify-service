# Line Notify Service

This is an API service for line-notify (https://notify-bot.line.me/zh_TW)

## Using Technique

- Go:
  - Use Go routine to push notifications synchronous.
  - Use JWT for admin authentication.
- Sqlite
  - Use file base DB to store client tokens.



**As a client (message receiver), the only step is using line OAuth to registry access token.**

**The Path is  {host}/oauth/line**



## APIs

There are multiple APIs for administrators such as login, update password, list/delete tokens, push notify.