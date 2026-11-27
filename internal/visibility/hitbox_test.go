package visibility

import (
	"strings"
	"testing"
)

func TestStanceOf(t *testing.T) {
	tests := []struct {
		eyesHeight float64
		isAirborne bool
		want       Stance
	}{
		{64, false, StanceStanding},
		{46, false, StanceCrouching},
		{64, true, StanceAirborne},
		{46, true, StanceAirborne},
	}
	for _, test := range tests {
		if got := StanceOf(test.eyesHeight, test.isAirborne); got != test.want {
			t.Errorf("StanceOf(%v, %v) = %v, want %v", test.eyesHeight, test.isAirborne, got, test.want)
		}
	}
}

func TestHitboxTemplates(t *testing.T) {
	seen := make(map[int]bool)
	for stance := Stance(0); stance < stanceCount; stance++ {
		for _, isMoving := range []bool{false, true} {
			for weaponKind := WeaponKind(0); weaponKind < weaponKindCount; weaponKind++ {
				index := HitboxTemplateIndex(stance, isMoving, weaponKind)
				if index < 0 || index >= HitboxTemplateCount || seen[index] {
					t.Fatalf("invalid or duplicated index %d for %v %v %v", index, stance, isMoving, weaponKind)
				}
				seen[index] = true

				gotStance, gotIsMoving, gotWeaponKind := SplitHitboxTemplateIndex(index)
				if gotStance != stance || gotIsMoving != isMoving || gotWeaponKind != weaponKind {
					t.Errorf("SplitHitboxTemplateIndex(%d) = %v %v %v, want %v %v %v", index, gotStance, gotIsMoving,
						gotWeaponKind, stance, isMoving, weaponKind)
				}
				name := HitboxTemplateName(index)
				if strings.HasPrefix(name, ",") || strings.HasSuffix(name, " ") {
					t.Errorf("template %d has no name for its stance or weapon kind: %q", index, name)
				}

				template := GetHitboxTemplate(stance, isMoving, weaponKind)
				if template.EyesHeight < 40 || template.EyesHeight > 70 || len(template.Points) < 10 {
					t.Errorf("template %s looks wrong: eyes height %v, %d points", HitboxTemplateName(index),
						template.EyesHeight, len(template.Points))
				}
			}
		}
	}
}
