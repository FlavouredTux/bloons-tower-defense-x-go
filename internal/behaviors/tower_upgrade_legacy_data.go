package behaviors

// generated upgrade chain data. do not edit by hand.

type legacyUpgradeEdge struct {
	NextObject string
	Cost       float64
	TowerCode  float64
}

var legacyUpgradeGraph = map[string][]legacyUpgradeEdge{
	"AHanger_0X": {
		{NextObject: "AHanger_1X", Cost: 900, TowerCode: 22.00},
	},
	"AHanger_1X": {
		{NextObject: "AHanger_2X", Cost: 1150, TowerCode: 22.10},
	},
	"AHanger_2X": {
		{NextObject: "AHanger_3L", Cost: 2500, TowerCode: 22.221},
		{NextObject: "AHanger_3L_Plus", Cost: 3600, TowerCode: 22.221},
		{NextObject: "AHanger_3M", Cost: 2400, TowerCode: 22.221},
		{NextObject: "AHanger_3R", Cost: 1900, TowerCode: 22.221},
	},
	"AHanger_3L": {
		{NextObject: "AHanger_4L", Cost: 4500, TowerCode: 22.31},
	},
	"AHanger_3L_Plus": {
		{NextObject: "AHanger_4L_Plus", Cost: 16300, TowerCode: 22.34},
	},
	"AHanger_3M": {
		{NextObject: "AHanger_4M", Cost: 17900, TowerCode: 22.32},
	},
	"AHanger_3R": {
		{NextObject: "AHanger_4R", Cost: 6000, TowerCode: 22.33},
	},
	"AHanger_4L": {
		{NextObject: "AHanger_5L", Cost: 9000, TowerCode: 22.41},
	},
	"AHanger_4L_Plus": {
		{NextObject: "AHanger_5L_Plus", Cost: 53300, TowerCode: 22.44},
	},
	"AHanger_4M": {
		{NextObject: "AHanger_5M", Cost: 53500, TowerCode: 22.42},
	},
	"AHanger_4R": {
		{NextObject: "AHanger_5R", Cost: 13000, TowerCode: 22.43},
	},
	"Acid_Shooter": {
		{NextObject: "Venomous_Monkey", Cost: 9000, TowerCode: 9.35},
	},
	"Adept_Wizard": {
		{NextObject: "Arcane_Wizard", Cost: 8500, TowerCode: 16.35},
	},
	"Advanced_Suction": {
		{NextObject: "Bloon_Juicer", Cost: 1600, TowerCode: 14.222},
		{NextObject: "Dense_Chipper", Cost: 4500, TowerCode: 14.222},
		{NextObject: "Razor_Shreder", Cost: 1200, TowerCode: 14.222},
		{NextObject: "Triple_Nozzle_Chipper", Cost: 2100, TowerCode: 14.222},
	},
	"Airburst_Sub": {
		{NextObject: "Assault_Wave_Sub", Cost: 2400, TowerCode: 7.31},
	},
	"Alchemical_Mastery": {
		{NextObject: "Diamond_Alchemist", Cost: 4400, TowerCode: 15.41},
	},
	"Arcane_Wizard": {
		{NextObject: "Pop_Elemental", Cost: 33000, TowerCode: 16.45},
	},
	"Arctic_Wind": {
		{NextObject: "Polar_Winds", Cost: 5400, TowerCode: 10.32},
	},
	"Army_of_Darkness": {
		{NextObject: "Dark_God", Cost: 13300, TowerCode: 16.42},
	},
	"Assault_Wave_Sub": {
		{NextObject: "Blockade_Sub", Cost: 6000, TowerCode: 7.41},
	},
	"Ballistic_Missile_Sub": {
		{NextObject: "First_Strike_Sub", Cost: 16000, TowerCode: 7.43},
	},
	"Banana_Factory": {
		{NextObject: "Banana_Replicator", Cost: 27000, TowerCode: 17.42},
	},
	"Banana_Farm": {
		{NextObject: "Banana_Plantation", Cost: 1200, TowerCode: 17.10},
	},
	"Banana_Plantation": {
		{NextObject: "Banana_Republic", Cost: 2800, TowerCode: 17.222},
		{NextObject: "Healthy_Bananas", Cost: 1400, TowerCode: 17.222},
		{NextObject: "Passive_Income", Cost: 1700, TowerCode: 17.222},
		{NextObject: "Rubberlust_Farm", Cost: 3000, TowerCode: 17.222},
	},
	"Banana_Republic": {
		{NextObject: "Banana_Factory", Cost: 9000, TowerCode: 17.32},
	},
	"Banana_Tree": {
		{NextObject: "Banana_Farm", Cost: 600, TowerCode: 17.00},
	},
	"Barbed_Darts_Sub": {
		{NextObject: "Twin_Guns", Cost: 530, TowerCode: 7.10},
	},
	"Barrel_Spin": {
		{NextObject: "Hydra_Rockets", Cost: 6600, TowerCode: 20.221},
		{NextObject: "Laser_Cannon", Cost: 5900, TowerCode: 20.221},
		{NextObject: "Unloader", Cost: 4300, TowerCode: 20.221},
		{NextObject: "Unloader_Plus", Cost: 5000, TowerCode: 20.221},
	},
	"Battery_Pack": {
		{NextObject: "Staggering_Explosions", Cost: 3000, TowerCode: 19.33},
	},
	"Battery_Plus": {
		{NextObject: "Carpet_Bombing", Cost: 7700, TowerCode: 19.36},
	},
	"Big_Bombs": {
		{NextObject: "Missile_Launcher", Cost: 450, TowerCode: 6.10},
	},
	"Big_One": {
		{NextObject: "Ion_Cannon", Cost: 25000, TowerCode: 19.42},
	},
	"Bigger_Piles": {
		{NextObject: "Faster_Piling", Cost: 1250, TowerCode: 21.10},
	},
	"Bionic_Boomer": {
		{NextObject: "Doublerangs", Cost: 2300, TowerCode: 3.33},
	},
	"Blacksmith_Village": {
		{NextObject: "Bloonbane_Village", Cost: 14000, TowerCode: 18.36},
	},
	"Blade_Shooter": {
		{NextObject: "Torque_Blades", Cost: 1330, TowerCode: 2.31},
	},
	"Blank_Tower_0M": {
		{NextObject: "Blank_Tower_1M", Cost: 100, TowerCode: 100001.00},
	},
	"Blank_Tower_1M": {
		{NextObject: "Blank_Tower_2M", Cost: 100, TowerCode: 100001.10},
	},
	"Blank_Tower_2M": {
		{NextObject: "Blank_Tower_3L", Cost: 100, TowerCode: 100000.222},
		{NextObject: "Blank_Tower_3L_Alt", Cost: 100, TowerCode: 100000.222},
		{NextObject: "Blank_Tower_3M", Cost: 100, TowerCode: 100000.222},
		{NextObject: "Blank_Tower_3M_Alt", Cost: 100, TowerCode: 100000.222},
		{NextObject: "Blank_Tower_3R", Cost: 100, TowerCode: 100000.222},
		{NextObject: "Blank_Tower_3R_Alt", Cost: 100, TowerCode: 100000.222},
	},
	"Blank_Tower_3L": {
		{NextObject: "Blank_Tower_4L", Cost: 100, TowerCode: 100001.31},
	},
	"Blank_Tower_3L_Alt": {
		{NextObject: "Blank_Tower_4L_Alt", Cost: 100, TowerCode: 100001.34},
	},
	"Blank_Tower_3M": {
		{NextObject: "Blank_Tower_4M", Cost: 100, TowerCode: 100001.32},
	},
	"Blank_Tower_3M_Alt": {
		{NextObject: "Blank_Tower_4M_Alt", Cost: 100, TowerCode: 100001.35},
	},
	"Blank_Tower_3R": {
		{NextObject: "Blank_Tower_4R", Cost: 100, TowerCode: 100001.33},
	},
	"Blank_Tower_3R_Alt": {
		{NextObject: "Blank_Tower_4R_Alt", Cost: 100, TowerCode: 100001.36},
	},
	"Blizzard_Wizard": {
		{NextObject: "Polar_Vortex_Mage", Cost: 41000, TowerCode: 10.46},
	},
	"Bloody_Sai_Ninja": {
		{NextObject: "Cursed_Katana_Ninja", Cost: 2900, TowerCode: 5.36},
	},
	"BloonX": {
		{NextObject: "Red_Shift", Cost: 16000, TowerCode: 23.43},
	},
	"Bloon_Buster_Cannon": {
		{NextObject: "Moab_Mauler", Cost: 1600, TowerCode: 6.31},
	},
	"Bloon_Dissolver": {
		{NextObject: "Bloon_Liquefier", Cost: 7200, TowerCode: 9.32},
	},
	"Bloon_Drain": {
		{NextObject: "Octo_Plasma", Cost: 3500, TowerCode: 23.31},
	},
	"Bloon_Impactor": {
		{NextObject: "Explosion_King", Cost: 7200, TowerCode: 6.42},
	},
	"Bloon_Juicer": {
		{NextObject: "Regurgitator", Cost: 3600, TowerCode: 14.31},
	},
	"Bloon_Liquefier": {
		{NextObject: "Moab_Poison", Cost: 16000, TowerCode: 9.42},
	},
	"Bloonbait_Farm": {
		{NextObject: "Bananabeam_Farm", Cost: 15000, TowerCode: 17.45},
	},
	"Bloonbane_Village": {
		{NextObject: "Fission_Village", Cost: 35000, TowerCode: 18.46},
	},
	"Bloonchipper": {
		{NextObject: "Succier_Chipper", Cost: 400, TowerCode: 14.00},
	},
	"Bloonjitzu": {
		{NextObject: "Ninja_God", Cost: 9500, TowerCode: 5.42},
	},
	"Bloontonium_Darts": {
		{NextObject: "Dart_Monkey_Gunner", Cost: 750, TowerCode: 1.31},
	},
	"Bloontonium_Reactor": {
		{NextObject: "Anti_Matter_Reactor", Cost: 9500, TowerCode: 7.42},
	},
	"Bloonzooka": {
		{NextObject: "RPG_Strike", Cost: 5500, TowerCode: 4.41},
	},
	"Bloonzooka_Plus": {
		{NextObject: "Railgun_Tank", Cost: 27000, TowerCode: 4.44},
	},
	"Bomb_Cannon": {
		{NextObject: "Big_Bombs", Cost: 500, TowerCode: 6.00},
	},
	"Bomb_Sprayer": {
		{NextObject: "Fire_Crackers", Cost: 5900, TowerCode: 2.35},
	},
	"Boomerang_Thrower": {
		{NextObject: "Multi_Pop_Thrower", Cost: 230, TowerCode: 3.00},
	},
	"Boost_Potions": {
		{NextObject: "Amplifier", Cost: 9000, TowerCode: 15.43},
	},
	"Brick_Layer": {
		{NextObject: "Moab_Crippler", Cost: 12500, TowerCode: 4.42},
	},
	"Camo_Radar": {
		{NextObject: "Monkey_Intelligence_Bureau", Cost: 4500, TowerCode: 18.31},
	},
	"Cannon_Ship": {
		{NextObject: "Harpoon_Ship", Cost: 2700, TowerCode: 11.33},
	},
	"Carpet_Bombing": {
		{NextObject: "Cross_Hyper_Beam", Cost: 37000, TowerCode: 19.46},
	},
	"Charge_Battery": {
		{NextObject: "Charge_Burst", Cost: 1500, TowerCode: 8.31},
	},
	"Charge_Burst": {
		{NextObject: "Charge_Overload", Cost: 3000, TowerCode: 8.41},
	},
	"Charge_Storage": {
		{NextObject: "Powerful_Charges", Cost: 330, TowerCode: 8.10},
	},
	"Charge_Tower": {
		{NextObject: "Charge_Storage", Cost: 220, TowerCode: 8.00},
	},
	"Cleansing_Foamer": {
		{NextObject: "Purification_Gun", Cost: 3500, TowerCode: 12.33},
	},
	"Cluster_Bombs": {
		{NextObject: "Bloon_Impactor", Cost: 3600, TowerCode: 6.32},
	},
	"Control_Sub": {
		{NextObject: "Mastermind_Sub", Cost: 23000, TowerCode: 7.44},
	},
	"Corrosive_Glue": {
		{NextObject: "Acid_Shooter", Cost: 2600, TowerCode: 9.222},
		{NextObject: "Bloon_Dissolver", Cost: 1800, TowerCode: 9.222},
		{NextObject: "Glue_Factory", Cost: 2500, TowerCode: 9.222},
		{NextObject: "Glue_Hose", Cost: 2000, TowerCode: 9.222},
	},
	"Crazy_Glue_Factory": {
		{NextObject: "Propeller_Glue", Cost: 7500, TowerCode: 9.43},
	},
	"Crows_Nest": {
		{NextObject: "Cannon_Ship", Cost: 1250, TowerCode: 11.221},
		{NextObject: "Destroyer", Cost: 2200, TowerCode: 11.221},
		{NextObject: "Dreadnaut_Ship", Cost: 2400, TowerCode: 11.221},
		{NextObject: "Swashbucklers", Cost: 900, TowerCode: 11.221},
	},
	"Cursed_Katana_Ninja": {
		{NextObject: "Cursed_Blade_Ninja", Cost: 19600, TowerCode: 5.46},
	},
	"Cursed_Pirate_Ship": {
		{NextObject: "Ghost_Ship", Cost: 18000, TowerCode: 11.44},
	},
	"Dark_Monkey": {
		{NextObject: "Nightmare", Cost: 44000, TowerCode: 24.41},
	},
	"Dart_Forest_Ranger": {
		{NextObject: "SMFC_Aficionado", Cost: 4500, TowerCode: 1.43},
	},
	"Dart_Monkey": {
		{NextObject: "Dart_Monkey_2", Cost: 160, TowerCode: 1.00},
	},
	"Dart_Monkey_2": {
		{NextObject: "Dart_Monkey_3", Cost: 180, TowerCode: 1.10},
	},
	"Dart_Monkey_3": {
		{NextObject: "Bloontonium_Darts", Cost: 210, TowerCode: 1.222},
		{NextObject: "Spike_o_Pult", Cost: 710, TowerCode: 1.222},
		{NextObject: "Spike_o_Pult_Plus", Cost: 950, TowerCode: 1.222},
		{NextObject: "Triple_Dart_Monkey", Cost: 440, TowerCode: 1.222},
	},
	"Dart_Monkey_Gunner": {
		{NextObject: "Dart_Tank", Cost: 1800, TowerCode: 1.41},
	},
	"Darter_Monkey": {
		{NextObject: "Darterer_Monkey", Cost: 100, TowerCode: 100001.00},
	},
	"Dartling_Gunner": {
		{NextObject: "Directed_Darts", Cost: 750, TowerCode: 20.00},
	},
	"Deadly_Precision": {
		{NextObject: "Brick_Layer", Cost: 6350, TowerCode: 4.32},
	},
	"Deep_Freeze": {
		{NextObject: "Ice_Shards", Cost: 2750, TowerCode: 10.33},
	},
	"Deep_Impact": {
		{NextObject: "Big_One", Cost: 9000, TowerCode: 19.32},
	},
	"Dense_Chipper": {
		{NextObject: "Gravity_Well", Cost: 12000, TowerCode: 14.35},
	},
	"Destroyer": {
		{NextObject: "Supreme_Battleship", Cost: 4700, TowerCode: 11.32},
	},
	"Directed_Darts": {
		{NextObject: "Barrel_Spin", Cost: 1200, TowerCode: 20.10},
	},
	"Distraction": {
		{NextObject: "Flash_Bombs", Cost: 2100, TowerCode: 5.31},
	},
	"Double_Shot": {
		{NextObject: "Bloonjitzu", Cost: 3000, TowerCode: 5.32},
	},
	"Doublerangs": {
		{NextObject: "Turbo_Charge", Cost: 3500, TowerCode: 3.43},
	},
	"Dragons_Breath": {
		{NextObject: "Pheonix_Flames", Cost: 6400, TowerCode: 16.31},
	},
	"Dreadnaut_Ship": {
		{NextObject: "Cursed_Pirate_Ship", Cost: 6000, TowerCode: 11.34},
	},
	"Drone_Engineer": {
		{NextObject: "Electroneer", Cost: 13500, TowerCode: 12.35},
	},
	"Droopy_Potions": {
		{NextObject: "Chemical_Engineer", Cost: 8000, TowerCode: 15.42},
	},
	"EMPs": {
		{NextObject: "Grid_Lock", Cost: 3900, TowerCode: 12.41},
	},
	"Electroneer": {
		{NextObject: "The_Machine", Cost: 27000, TowerCode: 12.45},
	},
	"Energy_Blast_Monkey": {
		{NextObject: "Super_Energy_Monkey", Cost: 55000, TowerCode: 24.36},
	},
	"Enhanced_Freeze": {
		{NextObject: "Perma_Frost", Cost: 390, TowerCode: 10.10},
	},
	"Even_More_Tacks": {
		{NextObject: "Blade_Shooter", Cost: 740, TowerCode: 2.222},
		{NextObject: "Bomb_Sprayer", Cost: 2200, TowerCode: 2.222},
		{NextObject: "Red_Hot_Tacks", Cost: 340, TowerCode: 2.222},
		{NextObject: "Tack_Sprayer", Cost: 630, TowerCode: 2.222},
	},
	"Explosion_Machine": {
		{NextObject: "Big_Bang_Machine", Cost: 31000, TowerCode: 6.46},
	},
	"Faster_Brewing": {
		{NextObject: "Healing_Potions", Cost: 1000, TowerCode: 15.223},
		{NextObject: "Paper_Potion_Monkey", Cost: 2200, TowerCode: 15.223},
		{NextObject: "Potent_Potions", Cost: 1900, TowerCode: 15.223},
		{NextObject: "Space_Potions", Cost: 3000, TowerCode: 15.223},
	},
	"Faster_Engineering": {
		{NextObject: "Cleansing_Foamer", Cost: 1100, TowerCode: 12.222},
		{NextObject: "Drone_Engineer", Cost: 1800, TowerCode: 12.222},
		{NextObject: "Shield_Buster", Cost: 750, TowerCode: 12.222},
		{NextObject: "Super_Nail_Gun", Cost: 550, TowerCode: 12.222},
	},
	"Faster_Launching": {
		{NextObject: "Ordinance", Cost: 500, TowerCode: 19.10},
	},
	"Faster_Piling": {
		{NextObject: "Moab_Shredr", Cost: 3400, TowerCode: 21.222},
		{NextObject: "Shield_Generator", Cost: 3500, TowerCode: 21.222},
		{NextObject: "Spikeball_Factory", Cost: 2400, TowerCode: 21.222},
		{NextObject: "Titanium_Spikes", Cost: 3000, TowerCode: 21.222},
	},
	"Faster_Shooting": {
		{NextObject: "Even_More_Tacks", Cost: 230, TowerCode: 2.10},
	},
	"Fire_Balls": {
		{NextObject: "Lightning_Rings", Cost: 1050, TowerCode: 16.10},
	},
	"Fire_Crackers": {
		{NextObject: "Fireworks_Shooter", Cost: 14400, TowerCode: 2.45},
	},
	"Fire_Strike": {
		{NextObject: "Nuke_Strike", Cost: 8500, TowerCode: 19.41},
	},
	"Flame_Jets": {
		{NextObject: "Ring_of_Fire", Cost: 3900, TowerCode: 2.42},
	},
	"Flash_Bombs": {
		{NextObject: "Mass_Distraction", Cost: 3800, TowerCode: 5.41},
	},
	"Flechette_Darts": {
		{NextObject: "Golden_Barrage", Cost: 17000, TowerCode: 20.41},
	},
	"Full_Metal_Jacket": {
		{NextObject: "Heat_Sniper", Cost: 1600, TowerCode: 4.10},
	},
	"Gamma_Rays": {
		{NextObject: "Omega_Rays", Cost: 33000, TowerCode: 23.42},
	},
	"Giga_Pops": {
		{NextObject: "Lightning_Bomb", Cost: 6300, TowerCode: 8.43},
	},
	"Glaive_King": {
		{NextObject: "Glaive_Lord", Cost: 7500, TowerCode: 3.42},
	},
	"Glaive_Ricochet": {
		{NextObject: "Glaive_King", Cost: 2700, TowerCode: 3.32},
	},
	"Glaive_Thrower": {
		{NextObject: "Bionic_Boomer", Cost: 1450, TowerCode: 3.221},
		{NextObject: "Glaive_Ricochet", Cost: 1200, TowerCode: 3.221},
		{NextObject: "Plasmarangs", Cost: 890, TowerCode: 3.221},
		{NextObject: "Plasmasaber_Thrower", Cost: 1900, TowerCode: 3.221},
	},
	"Glue_Factory": {
		{NextObject: "Crazy_Glue_Factory", Cost: 5000, TowerCode: 9.33},
	},
	"Glue_Gunner_L1": {
		{NextObject: "Piercing_Glue", Cost: 300, TowerCode: 9.00},
	},
	"Glue_Hose": {
		{NextObject: "Thick_Glue_Splatter", Cost: 4500, TowerCode: 9.31},
	},
	"Golden_Bolts": {
		{NextObject: "Golden_Shower_Shooter", Cost: 63000, TowerCode: 20.44},
	},
	"Grape_Shot": {
		{NextObject: "Crows_Nest", Cost: 350, TowerCode: 11.10},
	},
	"Gravity_Well": {
		{NextObject: "Singularity_Engine", Cost: 40000, TowerCode: 14.45},
	},
	"Hanger_0X": {
		{NextObject: "Hanger_1X", Cost: 440, TowerCode: 13.00},
	},
	"Hanger_1X": {
		{NextObject: "Hanger_2X", Cost: 550, TowerCode: 13.10},
	},
	"Hanger_2X": {
		{NextObject: "Hanger_3L", Cost: 900, TowerCode: 13.223},
		{NextObject: "Hanger_3M", Cost: 2800, TowerCode: 13.223},
		{NextObject: "Hanger_3R", Cost: 2700, TowerCode: 13.223},
		{NextObject: "Hanger_3R_plus", Cost: 3600, TowerCode: 13.223},
	},
	"Hanger_3L": {
		{NextObject: "Hanger_4L", Cost: 2700, TowerCode: 13.31},
	},
	"Hanger_3M": {
		{NextObject: "Hanger_4M", Cost: 6500, TowerCode: 13.32},
	},
	"Hanger_3R": {
		{NextObject: "Hanger_4R", Cost: 5000, TowerCode: 13.33},
	},
	"Hanger_3R_plus": {
		{NextObject: "Hanger_4R_plus", Cost: 7500, TowerCode: 13.36},
	},
	"Hanger_4L": {
		{NextObject: "Hanger_5L", Cost: 8100, TowerCode: 13.41},
	},
	"Hanger_4M": {
		{NextObject: "Hanger_5M", Cost: 17500, TowerCode: 13.42},
	},
	"Hanger_4R": {
		{NextObject: "Hanger_5R", Cost: 15500, TowerCode: 13.43},
	},
	"Hanger_4R_plus": {
		{NextObject: "Hanger_5R_plus", Cost: 21000, TowerCode: 13.46},
	},
	"Harpoon_Ship": {
		{NextObject: "MOAB_Takedown", Cost: 8100, TowerCode: 11.43},
	},
	"Healing_Potions": {
		{NextObject: "Boost_Potions", Cost: 2000, TowerCode: 15.33},
	},
	"Healthier_Bananas": {
		{NextObject: "Golden_Fruit", Cost: 6000, TowerCode: 17.41},
	},
	"Healthy_Bananas": {
		{NextObject: "Healthier_Bananas", Cost: 2600, TowerCode: 17.31},
	},
	"Heat_Sniper": {
		{NextObject: "Deadly_Precision", Cost: 3000, TowerCode: 4.221},
		{NextObject: "Semi_Automatic_Rifle", Cost: 3300, TowerCode: 4.221},
		{NextObject: "Shotgun_Plus", Cost: 2100, TowerCode: 4.221},
		{NextObject: "Tactical_Shotgun", Cost: 1700, TowerCode: 4.221},
	},
	"High_Energy_Beacon": {
		{NextObject: "Monkey_Energizer", Cost: 20000, TowerCode: 18.43},
	},
	"Hydra_Rockets": {
		{NextObject: "Rocket_Storm", Cost: 13900, TowerCode: 20.33},
	},
	"Hyper_Ultra_Beam": {
		{NextObject: "Hyper_Ultra_Beam", Cost: 75000, TowerCode: 20.52},
	},
	"Ice_Monkey": {
		{NextObject: "Enhanced_Freeze", Cost: 350, TowerCode: 10.00},
	},
	"Ice_Shards": {
		{NextObject: "Ice_Storm", Cost: 7700, TowerCode: 10.43},
	},
	"Ice_Wizard": {
		{NextObject: "Blizzard_Wizard", Cost: 10100, TowerCode: 10.36},
	},
	"Jungle_Drums": {
		{NextObject: "Blacksmith_Village", Cost: 3200, TowerCode: 18.223},
		{NextObject: "Camo_Radar", Cost: 2700, TowerCode: 18.223},
		{NextObject: "Monkey_Fort_Village", Cost: 3000, TowerCode: 18.223},
		{NextObject: "Village_Jams", Cost: 3000, TowerCode: 18.223},
	},
	"Katana_Ninja": {
		{NextObject: "Hidden_Monkey", Cost: 5600, TowerCode: 5.43},
	},
	"Laser_Cannon": {
		{NextObject: "Ray_of_Doom", Cost: 39000, TowerCode: 20.32},
	},
	"Laser_Vision": {
		{NextObject: "Plasma_Vision", Cost: 4600, TowerCode: 24.10},
	},
	"Life_Insurance": {
		{NextObject: "Insurance_Fraud", Cost: 7800, TowerCode: 17.43},
	},
	"Lightning_Rings": {
		{NextObject: "Adept_Wizard", Cost: 3500, TowerCode: 16.222},
		{NextObject: "Dragons_Breath", Cost: 3300, TowerCode: 16.222},
		{NextObject: "Necromancer", Cost: 1500, TowerCode: 16.222},
		{NextObject: "Whirlwind", Cost: 1600, TowerCode: 16.222},
	},
	"Loopy_Potions": {
		{NextObject: "Helium_Haze", Cost: 26000, TowerCode: 15.46},
	},
	"Machine_Gun": {
		{NextObject: "Supply_Drones", Cost: 13250, TowerCode: 4.43},
	},
	"Magnetic_Charge_Tower": {
		{NextObject: "Gravity_Bomb_Charger", Cost: 10500, TowerCode: 8.42},
	},
	"Magnetic_Field": {
		{NextObject: "Gravity_Bomb_Charger", Cost: 10500, TowerCode: 8.42},
	},
	"Masterangs": {
		{NextObject: "Megarang_Toss", Cost: 3000, TowerCode: 3.41},
	},
	"Mega_Charger": {
		{NextObject: "Mega_Mega_Charger", Cost: 21500, TowerCode: 8.44},
	},
	"Mega_Fruit_Cannon": {
		{NextObject: "MOAP_Cannon", Cost: 9000, TowerCode: 6.43},
	},
	"Missile_Launcher": {
		{NextObject: "Bloon_Buster_Cannon", Cost: 1200, TowerCode: 6.223},
		{NextObject: "Cluster_Bombs", Cost: 900, TowerCode: 6.223},
		{NextObject: "Pineapple_Launcher", Cost: 700, TowerCode: 6.223},
		{NextObject: "Pop_Cannon", Cost: 1500, TowerCode: 6.223},
	},
	"Moab_Mauler": {
		{NextObject: "Moab_Assassin_Cannon", Cost: 4200, TowerCode: 6.41},
	},
	"Moab_Shredr": {
		{NextObject: "Trvlr_Spikes", Cost: 7900, TowerCode: 21.33},
	},
	"Monkey_Alchemist": {
		{NextObject: "Poison_Alchemist", Cost: 600, TowerCode: 15.00},
	},
	"Monkey_Apprentice": {
		{NextObject: "Fire_Balls", Cost: 360, TowerCode: 16.00},
	},
	"Monkey_Buccaneer": {
		{NextObject: "Grape_Shot", Cost: 450, TowerCode: 11.00},
	},
	"Monkey_Engineer": {
		{NextObject: "Sentry_Turrets", Cost: 530, TowerCode: 12.00},
	},
	"Monkey_Fort_Village": {
		{NextObject: "Monkey_Town", Cost: 9000, TowerCode: 18.32},
	},
	"Monkey_Intelligence_Bureau": {
		{NextObject: "Call_to_Arms", Cost: 15000, TowerCode: 18.41},
	},
	"Monkey_Pirates": {
		{NextObject: "Pirate_Captain_Ship", Cost: 3500, TowerCode: 11.41},
	},
	"Monkey_Sub": {
		{NextObject: "Barbed_Darts_Sub", Cost: 330, TowerCode: 7.00},
	},
	"Monkey_Town": {
		{NextObject: "Monkey_Metropolis", Cost: 18000, TowerCode: 18.42},
	},
	"Monkey_Village": {
		{NextObject: "Tool_Sharpener", Cost: 800, TowerCode: 18.00},
	},
	"Mortar_Launcher": {
		{NextObject: "Faster_Launching", Cost: 360, TowerCode: 19.00},
	},
	"Multi_Pop_Thrower": {
		{NextObject: "Glaive_Thrower", Cost: 330, TowerCode: 3.10},
	},
	"Napalm_Launcher": {
		{NextObject: "Fire_Strike", Cost: 4400, TowerCode: 19.31},
	},
	"Necromancer": {
		{NextObject: "Army_of_Darkness", Cost: 3500, TowerCode: 16.32},
	},
	"Nega_Monkey": {
		{NextObject: "Dark_Monkey", Cost: 19900, TowerCode: 24.31},
	},
	"Ninja_Monkey": {
		{NextObject: "Sharp_Shurikens", Cost: 370, TowerCode: 5.00},
	},
	"Ninja_Training": {
		{NextObject: "Bloody_Sai_Ninja", Cost: 1600, TowerCode: 5.223},
		{NextObject: "Distraction", Cost: 500, TowerCode: 5.223},
		{NextObject: "Double_Shot", Cost: 1050, TowerCode: 5.223},
		{NextObject: "Sai_Ninja", Cost: 1050, TowerCode: 5.223},
	},
	"Octo_Plasma": {
		{NextObject: "Moab_Drain", Cost: 5500, TowerCode: 23.41},
	},
	"Orbital_Discharge": {
		{NextObject: "Magnetic_Charge_Tower", Cost: 3400, TowerCode: 8.32},
	},
	"Ordinance": {
		{NextObject: "Battery_Pack", Cost: 2250, TowerCode: 19.223},
		{NextObject: "Battery_Plus", Cost: 2550, TowerCode: 19.223},
		{NextObject: "Deep_Impact", Cost: 2600, TowerCode: 19.223},
		{NextObject: "Napalm_Launcher", Cost: 1300, TowerCode: 19.223},
	},
	"Paper_Potion_Monkey": {
		{NextObject: "Alchemical_Mastery", Cost: 2500, TowerCode: 15.31},
	},
	"Passive_Income": {
		{NextObject: "Life_Insurance", Cost: 2500, TowerCode: 17.33},
	},
	"Perma_Frost": {
		{NextObject: "Arctic_Wind", Cost: 2700, TowerCode: 10.223},
		{NextObject: "Deep_Freeze", Cost: 900, TowerCode: 10.223},
		{NextObject: "Ice_Wizard", Cost: 3100, TowerCode: 10.223},
		{NextObject: "Snowball_Thrower", Cost: 750, TowerCode: 10.223},
	},
	"Pheonix_Flames": {
		{NextObject: "Fire_God", Cost: 9300, TowerCode: 16.41},
	},
	"Piercing_Glue": {
		{NextObject: "Corrosive_Glue", Cost: 360, TowerCode: 9.10},
	},
	"Pineapple_Launcher": {
		{NextObject: "Mega_Fruit_Cannon", Cost: 2700, TowerCode: 6.33},
	},
	"Planetery_Spikes": {
		{NextObject: "Nebula", Cost: 25000, TowerCode: 21.42},
	},
	"Plasma_Monkey_": {
		{NextObject: "Quad_Plasma", Cost: 950, TowerCode: 23.00},
	},
	"Plasma_Vision": {
		{NextObject: "Energy_Blast_Monkey", Cost: 21000, TowerCode: 24.223},
		{NextObject: "Nega_Monkey", Cost: 9000, TowerCode: 24.223},
		{NextObject: "Robo_Monkey", Cost: 10000, TowerCode: 24.223},
		{NextObject: "Sun_Worshipper", Cost: 16500, TowerCode: 24.223},
	},
	"Plasmapunch_Monkey": {
		{NextObject: "Doomfist_Monkey", Cost: 60000, TowerCode: 23.44},
	},
	"Plasmarangs": {
		{NextObject: "Masterangs", Cost: 2000, TowerCode: 3.31},
	},
	"Plasmasaber_Knight": {
		{NextObject: "Projectile_Master", Cost: 17000, TowerCode: 3.44},
	},
	"Plasmasaber_Thrower": {
		{NextObject: "Plasmasaber_Knight", Cost: 5500, TowerCode: 3.34},
	},
	"Plasmawhip_Monkey": {
		{NextObject: "Plasmapunch_Monkey", Cost: 16000, TowerCode: 23.34},
	},
	"Poison_Alchemist": {
		{NextObject: "Faster_Brewing", Cost: 430, TowerCode: 15.10},
	},
	"Polar_Winds": {
		{NextObject: "Absolute_Zero", Cost: 21000, TowerCode: 10.42},
	},
	"Pop_Cannon": {
		{NextObject: "Explosion_Machine", Cost: 6300, TowerCode: 6.36},
	},
	"Potent_Plasma": {
		{NextObject: "Bloon_Drain", Cost: 1800, TowerCode: 23.221},
		{NextObject: "Plasmawhip_Monkey", Cost: 4500, TowerCode: 23.221},
		{NextObject: "Radiation", Cost: 3000, TowerCode: 23.221},
		{NextObject: "Solar_Flare", Cost: 3350, TowerCode: 23.221},
	},
	"Potent_Potions": {
		{NextObject: "Droopy_Potions", Cost: 4100, TowerCode: 15.32},
	},
	"Powerful_Charges": {
		{NextObject: "Charge_Battery", Cost: 750, TowerCode: 8.221},
		{NextObject: "Orbital_Discharge", Cost: 1100, TowerCode: 8.221},
		{NextObject: "Super_Charge_Tower", Cost: 2000, TowerCode: 8.221},
		{NextObject: "Tesla_Coil", Cost: 1700, TowerCode: 8.221},
	},
	"Purification_Gun": {
		{NextObject: "Bloon_Containment_Unit", Cost: 7500, TowerCode: 12.43},
	},
	"Quad_Core_Rotors": {
		{NextObject: "Super_Wide_Funnel", Cost: 12000, TowerCode: 14.42},
	},
	"Quad_Plasma": {
		{NextObject: "Potent_Plasma", Cost: 1150, TowerCode: 23.10},
	},
	"Radiation": {
		{NextObject: "BloonX", Cost: 6000, TowerCode: 23.33},
	},
	"Ray_of_Doom": {
		{NextObject: "Hyper_Ultra_Beam", Cost: 85000, TowerCode: 20.42},
	},
	"Razor_Shreder": {
		{NextObject: "Quad_Core_Rotors", Cost: 3300, TowerCode: 14.32},
	},
	"Red_Hot_Tacks": {
		{NextObject: "Flame_Jets", Cost: 1550, TowerCode: 2.32},
	},
	"Regurgitator": {
		{NextObject: "Vampire_Blender", Cost: 9000, TowerCode: 14.41},
	},
	"Robo_Monkey": {
		{NextObject: "Technological_Terror", Cost: 27000, TowerCode: 24.33},
	},
	"Rocket_Storm": {
		{NextObject: "Bloon_Area_Denial_System", Cost: 21000, TowerCode: 20.43},
	},
	"Rubberlust_Farm": {
		{NextObject: "Bloonbait_Farm", Cost: 3000, TowerCode: 17.35},
	},
	"Sai_Ninja": {
		{NextObject: "Katana_Ninja", Cost: 1900, TowerCode: 5.33},
	},
	"Semi_Automatic_Rifle": {
		{NextObject: "Machine_Gun", Cost: 5700, TowerCode: 4.33},
	},
	"Sentry_Turrets": {
		{NextObject: "Faster_Engineering", Cost: 350, TowerCode: 12.10},
	},
	"Sharp_Shurikens": {
		{NextObject: "Ninja_Training", Cost: 480, TowerCode: 5.10},
	},
	"Shield_Buster": {
		{NextObject: "EMPs", Cost: 2250, TowerCode: 12.31},
	},
	"Shield_Generator": {
		{NextObject: "Ultra_Forcefield_Generator", Cost: 11500, TowerCode: 21.35},
	},
	"Shotgun_Plus": {
		{NextObject: "Bloonzooka_Plus", Cost: 5900, TowerCode: 4.34},
	},
	"Singularity_Engine": {
		{NextObject: "Singularity_Engine", Cost: 40000, TowerCode: 14.55},
	},
	"Smart_Sub": {
		{NextObject: "Control_Sub", Cost: 6500, TowerCode: 7.34},
	},
	"Sniper_Monkey": {
		{NextObject: "Full_Metal_Jacket", Cost: 500, TowerCode: 4.00},
	},
	"Snowball_Thrower": {
		{NextObject: "Snowmound_Cannon", Cost: 2000, TowerCode: 10.31},
	},
	"Snowmound_Cannon": {
		{NextObject: "Freezerburn_Cannon_XL", Cost: 4800, TowerCode: 10.41},
	},
	"Solar_Flare": {
		{NextObject: "Gamma_Rays", Cost: 11900, TowerCode: 23.32},
	},
	"Space_Potions": {
		{NextObject: "Loopy_Potions", Cost: 8300, TowerCode: 15.36},
	},
	"Spike_Assault_Rifle": {
		{NextObject: "Spike_Mini_Gun", Cost: 9500, TowerCode: 1.45},
	},
	"Spike_Factory": {
		{NextObject: "Bigger_Piles", Cost: 700, TowerCode: 21.00},
	},
	"Spike_o_Pult": {
		{NextObject: "Triple_Pult", Cost: 1600, TowerCode: 1.32},
	},
	"Spike_o_Pult_Plus": {
		{NextObject: "Spike_Assault_Rifle", Cost: 3300, TowerCode: 1.35},
	},
	"Spikeball_Factory": {
		{NextObject: "Spiked_Mines", Cost: 8500, TowerCode: 21.31},
	},
	"Spiked_Mines": {
		{NextObject: "Spike_Wall", Cost: 10000, TowerCode: 21.41},
	},
	"Staggering_Explosions": {
		{NextObject: "Pop_and_Awe", Cost: 10500, TowerCode: 19.43},
	},
	"Succier_Chipper": {
		{NextObject: "Advanced_Suction", Cost: 500, TowerCode: 14.10},
	},
	"Sun_Temple": {
		{NextObject: "Sun_God", Cost: 600000, TowerCode: 24.42},
	},
	"Sun_Worshipper": {
		{NextObject: "Sun_Temple", Cost: 90000, TowerCode: 24.32},
	},
	"Super_Charge_Tower": {
		{NextObject: "Mega_Charger", Cost: 7000, TowerCode: 8.34},
	},
	"Super_Energy_Monkey": {
		{NextObject: "Arcane_Guardian_Monkey", Cost: 155000, TowerCode: 24.46},
	},
	"Super_Monkey": {
		{NextObject: "Laser_Vision", Cost: 2400, TowerCode: 24.00},
	},
	"Super_Nail_Gun": {
		{NextObject: "Super_Sentries", Cost: 1200, TowerCode: 12.32},
	},
	"Super_Sentries": {
		{NextObject: "Omega_Tech", Cost: 6350, TowerCode: 12.42},
	},
	"Support_Sub": {
		{NextObject: "Bloontonium_Reactor", Cost: 3000, TowerCode: 7.32},
	},
	"Supreme_Battleship": {
		{NextObject: "Aircraft_Carrier", Cost: 14000, TowerCode: 11.42},
	},
	"Swashbucklers": {
		{NextObject: "Monkey_Pirates", Cost: 1800, TowerCode: 11.31},
	},
	"Tack_Shooter": {
		{NextObject: "Faster_Shooting", Cost: 170, TowerCode: 2.00},
	},
	"Tack_Sprayer": {
		{NextObject: "Tack_Storm", Cost: 2300, TowerCode: 2.33},
	},
	"Tack_Storm": {
		{NextObject: "Tack_Typhoon", Cost: 5600, TowerCode: 2.43},
	},
	"Tactical_Shotgun": {
		{NextObject: "Bloonzooka", Cost: 3300, TowerCode: 4.31},
	},
	"Technological_Terror": {
		{NextObject: "Annihilator", Cost: 54000, TowerCode: 24.43},
	},
	"Tempest_Tornado": {
		{NextObject: "Wind_God", Cost: 11300, TowerCode: 16.43},
	},
	"Tesla_Coil": {
		{NextObject: "Giga_Pops", Cost: 3500, TowerCode: 8.33},
	},
	"Thick_Glue_Splatter": {
		{NextObject: "Glue_Striker", Cost: 6000, TowerCode: 9.41},
	},
	"Titanium_Spikes": {
		{NextObject: "Planetery_Spikes", Cost: 7800, TowerCode: 21.32},
	},
	"Tool_Sharpener": {
		{NextObject: "Jungle_Drums", Cost: 1300, TowerCode: 18.10},
	},
	"Torpedo_Sub": {
		{NextObject: "Ballistic_Missile_Sub", Cost: 2550, TowerCode: 7.33},
	},
	"Torque_Blades": {
		{NextObject: "Blade_Maelstrom", Cost: 2800, TowerCode: 2.41},
	},
	"Triple_Dart_Monkey": {
		{NextObject: "Dart_Forest_Ranger", Cost: 850, TowerCode: 1.33},
	},
	"Triple_Nozzle_Chipper": {
		{NextObject: "Turbo_Sucker", Cost: 4200, TowerCode: 14.33},
	},
	"Triple_Pult": {
		{NextObject: "Juggernaut", Cost: 3200, TowerCode: 1.42},
	},
	"Trvlr_Spikes": {
		{NextObject: "Spike_Storm", Cost: 13000, TowerCode: 21.43},
	},
	"Turbo_Sucker": {
		{NextObject: "Supa_Vac", Cost: 10500, TowerCode: 14.43},
	},
	"Twin_Guns": {
		{NextObject: "Airburst_Sub", Cost: 950, TowerCode: 7.221},
		{NextObject: "Smart_Sub", Cost: 2300, TowerCode: 7.221},
		{NextObject: "Support_Sub", Cost: 1000, TowerCode: 7.221},
		{NextObject: "Torpedo_Sub", Cost: 1050, TowerCode: 7.221},
	},
	"Ultra_Forcefield_Generator": {
		{NextObject: "Lockdown_Factory", Cost: 42500, TowerCode: 21.45},
	},
	"Unloader": {
		{NextObject: "Flechette_Darts", Cost: 8500, TowerCode: 20.31},
	},
	"Unloader_Plus": {
		{NextObject: "Golden_Bolts", Cost: 15000, TowerCode: 20.34},
	},
	"Venomous_Monkey": {
		{NextObject: "King_Cobra", Cost: 28000, TowerCode: 9.45},
	},
	"Village_Jams": {
		{NextObject: "High_Energy_Beacon", Cost: 6000, TowerCode: 18.33},
	},
	"Whirlwind": {
		{NextObject: "Tempest_Tornado", Cost: 5500, TowerCode: 16.33},
	},
}

func legacyLinearUpgradeRuleFor(objectName string) (linearUpgradeRule, bool) {
	edges, ok := legacyUpgradeGraph[objectName]
	if !ok || len(edges) == 0 {
		return linearUpgradeRule{}, false
	}

	// linear stages are represented by tower code fraction .00 and .10.
	for _, e := range edges {
		frac1000 := int((e.TowerCode-float64(int(e.TowerCode)))*1000.0 + 0.5)
		if frac1000 == 0 || frac1000 == 100 {
			// TowerVal for the new instance is the code stored on the next object's own edges.
			nextVal, _ := legacyTowerCodeForObject(e.NextObject)
			return linearUpgradeRule{Cost: e.Cost, NextObject: e.NextObject, TowerVal: nextVal}, true
		}
	}
	return linearUpgradeRule{}, false
}

func legacyTowerCodeForObject(objectName string) (float64, bool) {
	edges, ok := legacyUpgradeGraph[objectName]
	if !ok || len(edges) == 0 {
		return 0, false
	}
	return edges[0].TowerCode, true
}

func legacyUpgradeEdgesFor(objectName string) []legacyUpgradeEdge {
	edges := legacyUpgradeGraph[objectName]
	if len(edges) == 0 {
		return nil
	}
	out := make([]legacyUpgradeEdge, len(edges))
	copy(out, edges)
	return out
}

func legacyLinearCostsFor(objectName string, maxSteps int) []float64 {
	if maxSteps <= 0 {
		return nil
	}

	costs := make([]float64, 0, maxSteps)
	cur := objectName
	for i := 0; i < maxSteps; i++ {
		rule, ok := legacyLinearUpgradeRuleFor(cur)
		if !ok {
			break
		}
		costs = append(costs, rule.Cost)
		cur = rule.NextObject
	}
	return costs
}
