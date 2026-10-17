package api

import (
	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

// noProjectileID is used when a shot or a grenade event can't be linked to a projectile.
// The parser assigns unique IDs from 1, and it must not be empty because PostgreSQL's COPY reads an empty CSV field as NULL.
const noProjectileID = "0"

// findGrenadeProjectileID returns the unique ID of the live projectile thrown with the given grenade.
// A player may have thrown several projectiles with the same grenade (weapon instances are re-used
// across throws), in that case the closest one to the event position is the right one.
func findGrenadeProjectileID(analyzer *Analyzer, grenadeUniqueID common.UniqueID, eventPosition r3.Vector) string {
	projectileID := noProjectileID
	bestDistance := -1.0
	bestEntityID := -1
	for _, projectile := range analyzer.parser.GameState().GrenadeProjectiles() {
		if projectile.WeaponInstance.UniqueID() != grenadeUniqueID {
			continue
		}

		distance := projectile.Position().Sub(eventPosition).Norm2()
		entityID := projectile.Entity.ID()
		isCloser := bestDistance == -1 || distance < bestDistance
		// On an exact distance tie either projectile is an equally valid match, so prefer the lower entity ID
		// to keep the result deterministic regardless of iteration order.
		isSameDistanceWithLowerEntityID := distance == bestDistance && entityID < bestEntityID
		if isCloser || isSameDistanceWithLowerEntityID {
			bestDistance = distance
			bestEntityID = entityID
			projectileID = projectile.UniqueID().String()
		}
	}

	return projectileID
}
