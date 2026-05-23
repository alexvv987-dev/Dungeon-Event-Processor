package config

import "testing"

func TestValidate(t *testing.T) {
	valid := Config{Floors: 2, Monsters: 2, OpenAt: "14:05:00", Duration: 2}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid config: %v", err)
	}

	cases := []Config{
		{Floors: 1, Monsters: 2, OpenAt: "14:00:00", Duration: 2},
		{Floors: 2, Monsters: 0, OpenAt: "14:00:00", Duration: 2},
		{Floors: 2, Monsters: 2, OpenAt: "14:00:00", Duration: 0},
		{Floors: 2, Monsters: 2, OpenAt: "bad", Duration: 2},
	}
	for _, c := range cases {
		if err := c.Validate(); err == nil {
			t.Errorf("expected error for config %+v", c)
		}
	}
}

func TestRegularFloors(t *testing.T) {
	c := Config{Floors: 4}
	if got := c.RegularFloors(); got != 3 {
		t.Errorf("RegularFloors() = %d, want 3", got)
	}
}
