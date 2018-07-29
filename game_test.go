package cardsagainstdiscord

import (
	"testing"
)

func TestNextCardCzar(t *testing.T) {
	players := []*Player{
		{ID: 1, Playing: true},
		{ID: 5, Playing: true},
		{ID: 2, Playing: true},
	}

	current := NextCardCzar(players, 0)
	if current != 1 {
		t.Error("Got ", current, " exected 1")
	}

	current = NextCardCzar(players, current)
	if current != 2 {
		t.Error("Got ", current, " exected 2")
	}

	current = NextCardCzar(players, current)
	if current != 5 {
		t.Error("Got ", current, " exected 5")
	}

	current = NextCardCzar(players, current)
	if current != 1 {
		t.Error("Got ", current, " exected 1")
	}
}

func TestNextCardCzar2(t *testing.T) {
	players := []*Player{
		{ID: 5, Playing: true},
		{ID: 1, Playing: true},
		{ID: 2, Playing: true},
	}

	current := NextCardCzar(players, 0)
	if current != 1 {
		t.Error("Got ", current, " exected 1")
	}

	current = NextCardCzar(players, current)
	if current != 2 {
		t.Error("Got ", current, " exected 2")
	}

	current = NextCardCzar(players, current)
	if current != 5 {
		t.Error("Got ", current, " exected 5")
	}

	current = NextCardCzar(players, current)
	if current != 1 {
		t.Error("Got ", current, " exected 1")
	}
}

func TestNextCardCzar3(t *testing.T) {
	players := []*Player{
		{ID: 5, Playing: true},
		{ID: 1, Playing: true},
		{ID: 2, Playing: true},
		{ID: 3, Playing: true},
	}

	current := NextCardCzar(players, 0)
	if current != 1 {
		t.Error("Got ", current, " exected 1")
	}

	current = NextCardCzar(players, current)
	if current != 2 {
		t.Error("Got ", current, " exected 2")
	}

	current = NextCardCzar(players, current)
	if current != 3 {
		t.Error("Got ", current, " exected 3")
	}

	current = NextCardCzar(players, current)
	if current != 5 {
		t.Error("Got ", current, " exected 5")
	}

	current = NextCardCzar(players, current)
	if current != 1 {
		t.Error("Got ", current, " exected 1")
	}
}
