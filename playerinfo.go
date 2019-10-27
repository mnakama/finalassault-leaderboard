package main

type playerInfo struct {
	Twitch  string
	Youtube string
	VRML    string
}

var playerLookup = map[uint]playerInfo{
	6487: { // Heathen
		Twitch:  "heathenist",
		Youtube: "channel/UCBKvirlCz8mEBerLEbkrWSA",
	},
	6484: { // SadlyItsBradley
		Twitch:  "sadlyitsbradley",
		Youtube: "user/BronyBrad",
		VRML:    "ZWRjMm94djc4bmM90",
	},
	6500: { // Milo.[HUN]
		Youtube: "channel/UCOofPnTfop6nRRc_CW9-xfQ",
	},
	6803: { // Splaticus
		Youtube: "channel/UCnPGJInTNIj1Qq4qCMdHqHg",
	},
	7010: { // photogineer
		Twitch: "photogineer",
	},
	7227: { // NotSporks
		Twitch: "maybesporks",
	},
	7741: { // PhaserLock Alex
		Twitch: "dav3schneider",
	},
	7888: { // Beastrick
		Youtube: "channel/UCzmnvJ9oWukHV7vYvYcRrZA",
		VRML:    "MTYrMVFpUmRlOVE90",
	},
	7944: { // MasterShadow
		Youtube: "channel/UCF9RAKFGIBThbFhYKFDOiWA",
	},
	8074: { // Gnuman
		Twitch:  "gnum4n",
		Youtube: "channel/UC6y_Bdmk_Vj26Zhj2xunt1g",
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
}
