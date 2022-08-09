export TOKEN=`cat .token`
export GUILD=`cat .guild`
go run cmd/runner.go -token=$TOKEN -guild=$GUILD