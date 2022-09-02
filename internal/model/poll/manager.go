package poll

type PollManager struct {
	polls map[string]*Poll
}

func CreatePollManager() *PollManager {
	pm := PollManager{
		polls: make(map[string]*Poll),
	}

	return &pm
}

func (pm *PollManager) AddPoll(poll *Poll) {
	pm.polls[poll.ID] = poll
}

func (pm *PollManager) GetPoll(id string) *Poll {
	return pm.polls[id]
}

func (pm *PollManager) RemovePoll(id string) {
	delete(pm.polls, id)
}
