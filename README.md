Simple Go applications created as a way to both get more familiar with Go and learn how to make apps that interact with AWS tools, like DynamoDB.  Most of the database interaction stuff was written using github copilot, to get past the "where do I start" and tons of googling that I'd need to learn this.


# Current state

Multiple applications:
- client-game - this is the main "game" that a player would play.  Simple controls, move a character around a grid
- server-game - this is the server-side version that a player wouldn't launch themselves.  Controls NPCs and randomly moves them around the grid
- entity-creator - helper app to create NPCs and Players without needing to manually go into the dynamodb dashboard


# How to
Note: this really only runs properly from my own laptop, due to how I set up the AWS stuff.  Requires logging into AWS (via "aws login") to communicate with the database.

## build
./buildall.sh

## run
Run these in two separate terminals
./bin/client-game player-1
./bin/server-game 

### controls (client-game)
Use WASD to move. The player's current position is saved to DynamoDB.
Arrow keys (and other keys like CMD) appear to cause the app to crash.  Not sure why, but ignoring for now, because this is just a learning exercise and I care more about the AWS/DB stuff.

Other players (and eventually NPCs) are loaded at startup and refreshed from the DB every 1 second while the game is running.

### controls (server-game)
No controls right now, basically just a view of the game.  NPCs move every 2s, using a random direction, and are currently allowed to occupy spaces that other players are in.

### entity creator
Uses command-line args to set up an object to place in the DB
(Note: right now just creates NPCs, need additional flag for players)
./bin/entity-creator --id npc-1 --health 100 --x 0 --y 0

# Notes
This is in a super rough state and unlikely I'll make it actually nice and well-designed, since this is really just a learning experiment.  Some stuff that would be nice to do, if I get the time:
- rename some of the go files so they're specific to the module they're in
- rename some modules and paths so it makes a bit more sense (i.e. common should be awstestcommon or something)
- split up how db is accessed between player and npc so they can actually have fields with different names, or rename it to generic "entity" with a type flag
- update how data is actually pulled from the db.  Need a better way than "every 1s" but also don't want to pull everything on every Update call, which is overkill
