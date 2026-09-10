Simple Go application(s) created as a way to both get more familiar with Go and learn how to make apps that interact with AWS tools, like DynamoDB.  Used a lot of AI to create this, to get past the "where do I start?" phase.


Current state: simple application that can do some basic operations on Player objects.  All requests are done as command-line args (for now).  Next step: add ebiten engine and make it do these requests actively when running.

Example commands for now:
go run . create-player --id player-1 --health 100 --x 0 --y 0
go run . move-player --id player-1 --x 12 --y 8
go run . damage-player --id player-1 --amount 25
go run . get-player --id player-1
