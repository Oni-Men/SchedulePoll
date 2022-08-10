package poll

type PollManager struct {
	polls map[string]*Poll
}

var manager *PollManager = CreatePollManager()

func CreatePollManager() *PollManager {
	pm := PollManager{
		polls: make(map[string]*Poll),
	}

	return &pm
}

func AddPoll(poll *Poll) {
	manager.polls[poll.ID] = poll
}

func GetPoll(id string) *Poll {
	return manager.polls[id]
}

func RemovePoll(id string) {
	delete(manager.polls, id)
}
