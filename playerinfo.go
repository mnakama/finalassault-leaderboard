package main

type playerInfo struct {
	Twitch  string
	Youtube string
	Steam   string
	VRML    string
}

var playerLookup = map[uint]playerInfo{
	6474: { // flameboy
		Steam: "profiles/76561198939943059",
	},
	6475: { // Shiftguns
		Steam: "id/Shiftguns",
	},
	6484: { // SadlyItsBradley
		Twitch:  "sadlyitsbradley",
		Youtube: "user/BronyBrad",
		Steam:   "id/Bradllez",
		VRML:    "ZWRjMm94djc4bmM90",
	},
	6487: { // Heathen
		Twitch:  "heathenist",
		Youtube: "channel/UCBKvirlCz8mEBerLEbkrWSA",
		Steam:   "id/KDLGates",
	},
	6500: { // Milo.[HUN]
		Youtube: "channel/UCOofPnTfop6nRRc_CW9-xfQ",
		Steam:   "profiles/76561197970773588",
	},
	6607: { // Bastarducci
		Steam: "profiles/76561198030045413",
	},
	6803: { // Splaticus
		Youtube: "channel/UCnPGJInTNIj1Qq4qCMdHqHg",
	},
	7010: { // photogineer
		Twitch:  "photogineer",
		Steam:   "profiles/76561198168719260",
		Youtube: "channel/UCIC4xwPkibJXB-ExvTNju1g",
	},
	7227: { // NotSporks
		Twitch:  "maybesporks",
		Youtube: "channel/UCOrsHqpF_umj3vyr9J05mMg",
	},
	7741: { // PhaserLock Alex
		Twitch: "dav3schneider",
	},
	7888: { // Beastrick
		Youtube: "channel/UCzmnvJ9oWukHV7vYvYcRrZA",
		Steam:   "profiles/76561197981208885",
		VRML:    "MTYrMVFpUmRlOVE90",
	},
	7944: { // MasterShadow
		Youtube: "channel/UCF9RAKFGIBThbFhYKFDOiWA",
	},
	8074: { // Gnuman
		Twitch:  "gnum4n",
		Youtube: "channel/UC6y_Bdmk_Vj26Zhj2xunt1g",
		Steam:   "id/gnum4n",
		VRML:    "UjNPUkNLWkNpcnM90",
	},
	9063: { // Naoko
		Twitch: "naokomoon",
		VRML:   "NllHdlZtcWx1b1k90",
	},
	9588: { // Manello
		VRML: "T1NmbjA4akoyNEE90",
	},
	10247: { // Diesel
		VRML: "Ump6MFdFTHpCTzQ90",
	},
	14237: { // Excel
		Twitch: "hps_excel",
	},
	28352: { // CrazierRex
		Steam:   "id/CrazierRex",
		Youtube: "user/CrazierRex1/videos",
	},
}
