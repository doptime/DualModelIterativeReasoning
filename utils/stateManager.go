package utils

import "strings"

type StateManager struct{}

func NewStateManager() *StateManager {
	return &StateManager{}
}

// TransformState performs transformations like adding or removing tokens.
func (sm *StateManager) TransformState(state string, operation string, token string) string {
	tokens := strings.Fields(state)
	switch operation {
	case "add":
		tokens = append(tokens, token)
	case "remove":
		for i, t := range tokens {
			if t == token {
				tokens = append(tokens[:i], tokens[i+1:]...)
				break
			}
		}
	}
	return strings.Join(tokens, " ")
}

// CompareStates checks if two states are equivalent.
func (sm *StateManager) CompareStates(stateA, stateB string) bool {
	return stateA == stateB
}
