# downgrid

downgrid is an osu!direct alternative to download beatmaps on multiplayer lobbies. 

downgrid basically masks itself as a web browser and whenever you click to a link, it detects if the url is a beatmap url. If it is, then downgrid download the associated beatmap from chimu.moe. 

# How to install

- Download the latest executable from [releases page](https://github.com/YusufOzmen01/downgrid/releases)
- Put the executable somewhere you won't accidentally delete. (Make sure to put it to somewhere you have permission to such as your desktop.)
- **Run the executable as admin once.** This step is necessary to register downgrid as a web browser. After this, you don't need admin permissions unless you move the downgrid executable.
- Choose your web browser path when a file dialog opens. Whenever a non osu! beatmap url is detected, it will redirect that url to your actual web browser.
- After you choose your web browser's executable, head over to your system settings and change your default web browser as downgrid.
- You're done!

# How to compile
```sh
git clone https://github.com/YusufOzmen01/downgrid
cd downgrid
go mod download
go build main.go
```

# Notes
Because I'm not an experienced go programmer, any contribution or suggestion is welcome. Don't hesitate to help me :3

While downgrid might be parsing urls, downgrid does not send any information about you such as which url you're clicked to me or any server.