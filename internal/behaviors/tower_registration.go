package behaviors

import "btdx/internal/engine"

// RegisterExtraTowerBehaviors registers all additional towers and their projectiles.
func RegisterExtraTowerBehaviors(im *engine.InstanceManager) {
	// towers
	im.RegisterBehavior("Tack_Shooter", func() engine.InstanceBehavior { return &TackShooterBehavior{} })
	im.RegisterBehavior("Boomerang_Thrower", func() engine.InstanceBehavior { return &BoomerangThrowerBehavior{} })
	im.RegisterBehavior("Sniper_Monkey", func() engine.InstanceBehavior { return &SniperMonkeyBehavior{} })
	im.RegisterBehavior("Ninja_Monkey", func() engine.InstanceBehavior { return &NinjaMonkeyBehavior{} })
	im.RegisterBehavior("Bomb_Cannon", func() engine.InstanceBehavior { return &BombCannonBehavior{} })
	// banana Farm and all upgrade objects share BananaFarmBehavior
	for _, bfName := range []string{
		"Banana_Farm",
		"Banana_Plantation",
		"Banana_Republic",
		"Healthy_Bananas",
		"Passive_Income",
		"Rubberlust_Farm",
		"Banana_Factory",
		"Banana_Replicator",
	} {
		n := bfName
		im.RegisterBehavior(n, func() engine.InstanceBehavior { return &BananaFarmBehavior{} })
	}
	im.RegisterBehavior("Ice_Monkey", func() engine.InstanceBehavior { return &IceMonkeyBehavior{} })
	im.RegisterBehavior("Glue_Gunner_L1", func() engine.InstanceBehavior { return &GlueGunnerBehavior{} })
	im.RegisterBehavior("Charge_Tower", func() engine.InstanceBehavior { return &ChargeTowerBehavior{} })
	// All charge tower upgrade names share the same behavior
	for _, name := range []string{
		"Charge_Storage", "Powerful_Charges",
		"Charge_Battery", "Charge_Burst", "Charge_Overload",
		"Orbital_Discharge", "Magnetic_Charge_Tower",
		"Tesla_Coil", "Giga_Pops", "Lightning_Bomb",
		"Super_Charge_Tower", "Mega_Charger", "Mega_Mega_Charger",
	} {
		n := name
		im.RegisterBehavior(n, func() engine.InstanceBehavior { return &ChargeTowerBehavior{} })
	}

	// projectiles
	im.RegisterBehavior("Tack", func() engine.InstanceBehavior { return &TackBehavior{} })
	im.RegisterBehavior("Blade", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 18} })
	im.RegisterBehavior("Red_Hot_Tack", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 16} })
	im.RegisterBehavior("Water_Tack", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 20} })
	im.RegisterBehavior("Storm_Tack", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 22} })
	im.RegisterBehavior("Firecracker", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 22} })
	im.RegisterBehavior("Flame_Jet", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 26} })
	im.RegisterBehavior("Firework_I", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Firework_II", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Firework_III", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Firework_IV", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Torque_Blade", func() engine.InstanceBehavior { return &SpinningBladeBehavior{} })
	im.RegisterBehavior("RoF", func() engine.InstanceBehavior { return &RingOfFireBehavior{} })
	im.RegisterBehavior("Boomerang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Plasmarang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Masterang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Ricochet_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("King_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Lord_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Turbo_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Megarang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Extra_pop", func() engine.InstanceBehavior { return &ExtraPopBehavior{} })
	im.RegisterBehavior("Super_Glaive_Proj", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Turbo_Glaive_Proj", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Glaive_Lord_Proj", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Sniper_Dart", func() engine.InstanceBehavior { return &SniperDartBehavior{} })
	im.RegisterBehavior("Shotgun_Slug", func() engine.InstanceBehavior { return &ShotgunSlugBehavior{} })
	im.RegisterBehavior("Bloonzooka_Shot", func() engine.InstanceBehavior { return &BloonzookaShotBehavior{} })
	im.RegisterBehavior("Frag", func() engine.InstanceBehavior { return &FragBehavior{} })
	im.RegisterBehavior("RPG_Projectile", func() engine.InstanceBehavior { return &RPGProjectileBehavior{} })
	im.RegisterBehavior("Shuriken", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Distraction_Shot", func() engine.InstanceBehavior { return &DistractionShotBehavior{} })
	im.RegisterBehavior("Flash_Bomb_Proj", func() engine.InstanceBehavior { return &FlashBombProjBehavior{} })
	im.RegisterBehavior("Flash", func() engine.InstanceBehavior { return &FlashBehavior{} })
	im.RegisterBehavior("Sai", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Alt_Sai", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Katana", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Cursed_Katana", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Cursed_Blade", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Crouching_Blade", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Golden_Ninja_Star", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Bomb", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("Cluster_Bomb", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("Impact_Bomb", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("King_Bomb", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("Cluster", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("Impactor", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("Small_Explosion", func() engine.InstanceBehavior { return &SmallExplosionBehavior{} })
	im.RegisterBehavior("Ice_Aura", func() engine.InstanceBehavior { return &IceAuraBehavior{} })
	im.RegisterBehavior("Glue_Glob", func() engine.InstanceBehavior { return &GlueGlobBehavior{} })
	// charge Tower projectiles
	im.RegisterBehavior("Charge_Proj", func() engine.InstanceBehavior { return &ChargeProjBehavior{hitRadius: 10} })
	im.RegisterBehavior("Powerful_Charge", func() engine.InstanceBehavior { return &ChargeProjBehavior{hitRadius: 10} })
	im.RegisterBehavior("Burst_Charge", func() engine.InstanceBehavior { return &BurstChargeBehavior{} })
	im.RegisterBehavior("Orbital_Charge", func() engine.InstanceBehavior { return &OrbitalChargeBehavior{} })
	im.RegisterBehavior("Magnetic_Charge", func() engine.InstanceBehavior { return &OrbitalChargeBehavior{} })
	im.RegisterBehavior("Small_Energy", func() engine.InstanceBehavior { return &EnergyProjBehavior{hitRadius: 22} })
	im.RegisterBehavior("Big_Energy", func() engine.InstanceBehavior { return &EnergyProjBehavior{hitRadius: 32} })
	im.RegisterBehavior("Super_Charge", func() engine.InstanceBehavior { return &ChargeProjBehavior{hitRadius: 14} })
	im.RegisterBehavior("Mega_Charge", func() engine.InstanceBehavior { return &MegaChargeBehavior{} })
	im.RegisterBehavior("Mega_Mega_Charge", func() engine.InstanceBehavior { return &MegaMegaChargeBehavior{} })
	im.RegisterBehavior("Mega_Proj", func() engine.InstanceBehavior { return &ChargeProjBehavior{hitRadius: 10} })
	im.RegisterBehavior("Mega_Mega_Proj", func() engine.InstanceBehavior { return &ChargeProjBehavior{hitRadius: 10} })
	im.RegisterBehavior("Energy_Bomb", func() engine.InstanceBehavior { return &EnergyBombBehavior{} })
	// monkey Sub — tower forms
	im.RegisterBehavior("Monkey_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Barbed_Darts_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Twin_Guns", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Torpedo_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Ballistic_Missile_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("First_Strike_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Airburst_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Assault_Wave_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Blockade_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Support_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Bloontonium_Reactor", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Anti_Matter_Reactor", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	im.RegisterBehavior("Smart_Sub", func() engine.InstanceBehavior { return &MonkeySubBehavior{} })
	// monkey Sub — projectiles
	im.RegisterBehavior("Monkey_Sub_Dart", func() engine.InstanceBehavior { return &SubHomingDartBehavior{hitRadius: 12} })
	im.RegisterBehavior("Barbed_Dart", func() engine.InstanceBehavior { return &SubHomingDartBehavior{hitRadius: 12} })
	im.RegisterBehavior("Airburst_Dart", func() engine.InstanceBehavior {
		return &AirburstDartBehavior{SubHomingDartBehavior: SubHomingDartBehavior{hitRadius: 12}}
	})
	im.RegisterBehavior("Airwave_Dart", func() engine.InstanceBehavior {
		return &AirwaveDartBehavior{SubHomingDartBehavior: SubHomingDartBehavior{hitRadius: 12}}
	})
	im.RegisterBehavior("Torpedo", func() engine.InstanceBehavior { return &TorpedoProjectileBehavior{} })
	im.RegisterBehavior("Ballistic_Missile", func() engine.InstanceBehavior { return &BallisticMissileBehavior{} })
	im.RegisterBehavior("First_Strike_Missile", func() engine.InstanceBehavior { return &FirstStrikeMissileBehavior{} })
	im.RegisterBehavior("Moab_Explosion", func() engine.InstanceBehavior { return &MoabExplosionBehavior{} })
	// monkey Sub — energy burst objects
	im.RegisterBehavior("Sub_Energy", func() engine.InstanceBehavior { return &SubEnergyBehavior{} })
	im.RegisterBehavior("Bloontonium_Energy", func() engine.InstanceBehavior { return &SubEnergyBehavior{} })
}
