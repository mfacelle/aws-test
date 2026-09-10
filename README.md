Simple Go application(s) created as a way to both get more familiar with Go and learn how to make apps that interact with AWS tools, like DynamoDB.  Used a lot of AI to create this, to get past the "where do I start?" phase.


Current state: simple "game" made with ebiten that takes player id as a command-line arg (pretend it's like a login), and then displays player location and allows for movement that updates dynamoDB.

Example run:
go run . player-1

Use WASD to move. The player's current position is saved to DynamoDB.
Arrow keys (and other keys like CMD) appear to cause the app to crash.  Not sure why, but ignoring for now, because this is just a learning exercise and I care more about the AWS/DB stuff.

Next step: load all players from DB and display them (and update in real time!?)
