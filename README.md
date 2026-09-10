Simple Go application(s) created as a way to both get more familiar with Go and learn how to make apps that interact with AWS tools, like DynamoDB.  Used a lot of AI to create this, to get past the "where do I start?" phase.


Current state: simple "game" made with ebiten that takes player id as a command-line arg (pretend it's like a login), loads all players from DynamoDB at startup, and displays them on the map. The active player is green; other players are amber. Movement updates the active player both in game and in the DB..

Example run:
go run . player-1

Use WASD to move. The player's current position is saved to DynamoDB.
Arrow keys (and other keys like CMD) appear to cause the app to crash.  Not sure why, but ignoring for now, because this is just a learning exercise and I care more about the AWS/DB stuff.

Other players are loaded at startup and refreshed from the DB every 1 second while the game is running.
