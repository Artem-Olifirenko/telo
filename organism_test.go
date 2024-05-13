package organism

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrganism(t *testing.T) {
	o := NewOrganism()
	assert.True(t, o.IsAlive())
	assert.False(t, o.IsReady())

	o.Ready()

	assert.True(t, o.IsReady())

	o.Die()

	assert.False(t, o.IsAlive())
	assert.True(t, o.IsReady())
}

func TestOrganism_GrowLimb_One(t *testing.T) {
	o := NewOrganism()
	limb := o.GrowLimb("limb")

	assert.True(t, limb.IsAlive())
	assert.False(t, limb.IsReady())
	assert.Equal(t, "limb", limb.Name())

	limb.Ready()

	assert.True(t, limb.IsAlive())
	assert.True(t, limb.IsReady())
	assert.True(t, o.IsAlive())
	assert.False(t, o.IsReady())

	o.Ready()

	assert.True(t, limb.IsAlive())
	assert.True(t, limb.IsReady())
	assert.True(t, o.IsAlive())
	assert.True(t, o.IsReady())

	limb.Die()

	assert.False(t, limb.IsAlive())
	assert.True(t, limb.IsReady())
	assert.False(t, o.IsAlive())
	assert.True(t, o.IsReady())

	o.Die()

	assert.False(t, limb.IsAlive())
	assert.True(t, limb.IsReady())
	assert.False(t, o.IsAlive())
	assert.True(t, o.IsReady())
}

func TestOrganism_DeadLimbs(t *testing.T) {
	tests := []struct {
		name      string
		organism  *Organism
		deadLimbs []*Limb
		isAlive   bool
	}{
		{
			name: "organism have one dead limb",
			organism: func() *Organism {
				o := NewOrganism()
				o.GrowLimb("alive_limb")
				deadLimb := o.GrowLimb("custom_limb")
				deadLimb.Die()

				return o
			}(),
			deadLimbs: []*Limb{
				{name: "custom_limb"},
			},
		},
		{
			name: "organism died",
			organism: func() *Organism {
				o := NewOrganism()
				o.GrowLimb("alive_limb1")
				o.GrowLimb("alive_limb2")
				o.Die()

				return o
			}(),
			deadLimbs: []*Limb{
				{name: "core"},
			},
		},
		{
			name: "organism is alive",
			organism: func() *Organism {
				o := NewOrganism()
				o.GrowLimb("alive_limb1")
				o.GrowLimb("alive_limb2")

				return o
			}(),
			isAlive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isAlive, tt.organism.IsAlive())
			assert.Equal(t, tt.deadLimbs, tt.organism.DeadLimbs())
		})
	}
}

func TestOrganism_NotReadyLimbs(t *testing.T) {
	tests := []struct {
		name          string
		organism      *Organism
		notReadyLimbs []*Limb
		isReady       bool
	}{
		{
			name: "not ready: empty organism",
			organism: func() *Organism {
				o := NewOrganism()

				return o
			}(),
			notReadyLimbs: []*Limb{
				{name: "core", isAlive: true},
			},
		},
		{
			name: "not ready: organism with one limb",
			organism: func() *Organism {
				o := NewOrganism()
				o.GrowLimb("first_limb")

				return o
			}(),
			notReadyLimbs: []*Limb{
				{name: "core", isAlive: true},
				{name: "first_limb", isAlive: true},
			},
		},
		{
			name: "not ready: organism with one ready limb",
			organism: func() *Organism {
				o := NewOrganism()
				firstLimb := o.GrowLimb("first_limb")
				firstLimb.Ready()

				return o
			}(),
			notReadyLimbs: []*Limb{
				{name: "core", isAlive: true},
			},
		},
		{
			name: "not ready: ready organism with one not ready limb",
			organism: func() *Organism {
				o := NewOrganism()
				o.GrowLimb("first_limb")
				o.Ready()

				return o
			}(),
			notReadyLimbs: []*Limb{
				{name: "first_limb", isAlive: true},
			},
		},
		{
			name: "ready: ready organism with ready limbs",
			organism: func() *Organism {
				o := NewOrganism()
				o.GrowLimb("first_limb").Ready()
				o.GrowLimb("second_limbs").Ready()
				o.GrowLimb("third_limb").Ready()
				o.Ready()

				return o
			}(),
			isReady: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isReady, tt.organism.IsReady())
			assert.Equal(t, tt.notReadyLimbs, tt.organism.NotReadyLimbs())
		})
	}
}
