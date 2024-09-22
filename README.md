# Tic-Tac-Toe Server
Build a tic-tac-toe game server with a REST/JSON API in Go.

Requirements:

* players should be able to create a new session (no authentication needed)
* all other endpoints should be authenticated by a session ID (`Authorization: Bearer <session-id>` or whatever you prefer)
* players should be able to start a new game
* players should be able to list open games (those with only a single player in them)
* players should be able to join an open game
* players should only be allowed to have a single open/active game per session
* players should only be allowed to play legal moves
* players should be able to see the win/draw status of the game for 1 minute after it finishes
* the server should be able to handle 100's of open games and 1,000's of active games

API (use JSON for both input and output):

```
POST   /session                create a new session (should return a token/ID which can be used for authentication)
GET    /session                get current session ID and active/open game ID
DELETE /session                end session
POST   /games                  create a new game (sets the current game ID for authenticated session)
GET    /games                  list open games
GET    /games/<game-id>        get the game state
POST   /games/<game-id>/join   join an open game (sets the current game ID for the authenticated session)
POST   /games/<game-id>/play   play a legal move
DELETE /games/<game-id>        end game
```

The game state should include:
* the board with the positions marked with X's and O's
* the game players with their session ID and if they are X or O
* the current players turn
* when the game finishes, who won or if it was a draw

Store sessions and games in memory and apply appropriate locking to ensure changes to the data are safe (don't worry about optimizing locking or anything like that).

How to prioritize work:
* completeness and correctness over everything else
* attempt to make the API user experience nice (return useful and actionable error messages when possible)
* write clean and testable code
* don't prioritize tests and project scaffolding (logging, configuration, etc.) unless they help you more efficiently complete the project