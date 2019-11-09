# finalassault-leaderboard
A better leaderboard for Final Assault

It's faster, javascript-free and framework-free, _and_ it has more features!

## Features
- Link to your name in the leaderboard to show off your rank (Example: <http://gnuman.games/finalassault/leaderboard#Gnuman>)
- Colorful Icons!
- Links to steam, twitch, youtube to see profiles vids. Check out top players' content to improve your game
- Fixes some unicode player names that are corrupted in the source data
(PhaserLock Interactive is using the wrong encoding somewhere)
- Super fast and lightweight

## Why?
The official leaderboard at https://phasermm.com runs on AngularJS and has a lot of extra javascript code. It used
to take about 10+ seconds to render the page. Their site is a bit better now, though.

If you like to compete, you'll probably want to refresh the page often and see results immediately. That's what this
project is for =)

### Potential future improvements
- Allow players to upload their own links (need a way to verify they are who they say they are. Does steam have OAuth?)
- Link to Oculus/Vive profiles somehow? If they even have a profile system?

### Stuff I wish PhaserLock would add
- Steam player ID, so I could automatically link to steam profiles
- Fix the text encoding, so I don't have to correct unicode player names
- More APIs, like active PvP battles, recent battle outcomes, players in PvP queue, etc.
