Hi Silas,

Hope everything is well. I have done the tic-tac-toe server (https://github.com/minozihao/tic-tac-toe-server) and sent you a collaborator invite in Github. I've spent a few hours trying to understand the requirements and design the project. Not sure if I understand everything correctly, but here are my thoughts and assumptions.

I've decided to store sessions in a sync.Map (in-memory) initialized with the server init process (routing registration). A session has an id for auth and an Open/ActiveGame initially being nil. All finished games are stored in the in-memory cache with 1 minute default expiration for game state display. The host flow is as follows, the host creates a new session which returns a session id. The host is able to create a game in the session and he gets back a gameid and player id. The player id is used to validate players for play moves. If the host creates several games repeatedly in a session, the activeGame of the session is overwritten, previously created games in the session will get 'garbage-collected'. A game has a 3*3 gameboard represented by a 2-dim array. When another user joins the game, the user will receive a playerid. Both players play a move with sessionId, gameId, playerId, and specifying the row and column.

For simplicity and time concern, I am using uuid instead of generating tokens, and didn't implement any logging. Hope that is ok. Have Tried to get more tests coverage, still missing a lot, I focused on testing the core play move though.

I've attached the postman collection for your convenience.

you can run it locally or with docker by exposing and mapping port 8080
e.g. docker run -e 8080 -p 8080:8080 --name tictac image_id
