package service

import (
	"errors"
	engine "github.com/zhangmingkai4315/go-microservice-project/engine"
	"strings"
)

type inMemoryRepository struct {
	matches []engine.Match
}

func NewInMemoryRepository() *inMemoryRepository {
	repo := &inMemoryRepository{}
	repo.matches = []engine.Match{}
	return repo
}

func (repo *inMemoryRepository) addMatch(match engine.Match) (err error) {
	repo.matches = append(repo.matches, match)
	return err
}

func (repo *inMemoryRepository) getMatches() (matches []engine.Match, err error) {
	matches = repo.matches
	return
}

func (repo *inMemoryRepository) getMatch(id string) (match engine.Match, err error) {
	found := false
	for _, target := range repo.matches {
		if strings.Compare(target.ID, id) == 0 {
			match = target
			found = true
		}
	}
	if !found {
		err = errors.New("Could not found match in repos")
	}
	return match, err
}
