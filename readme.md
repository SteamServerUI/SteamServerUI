# SteamServerUI - Your One-Stop Shop for Steam Server Shenanigans

![Go](https://img.shields.io/badge/Go-1.24.2-blue?logo=go&logoColor=white)
![Version](https://img.shields.io/badge/Version-v6%20Preview-orange?logo=github&logoColor=white)
![Issues](https://img.shields.io/github/issues/jacksonthemaster/StationeersServerUI?logo=github&logoColor=white)
![Stars](https://img.shields.io/github/stars/jacksonthemaster/StationeersServerUI?style=social&logo=github)
![Windows](https://img.shields.io/badge/Windows-supported-blue?logo=windows&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-supported-green?logo=linux&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-available-blue?logo=docker&logoColor=white)
![Last Commit](https://img.shields.io/github/last-commit/JacksonTheMaster/StationeersServerUI/v6-pre?logo=git&logoColor=white)

> **Note**: v6 is a work in progress. It currently IS able to run a Server successfully in a playable state, but the UI requires developer knowlege at some points to understand how to use it. Not recommended for production use, but technically feasible.


## 🚀 From Stationeers to Steam: The Great Servervolution

Once upon a time, I built **StationeersServerUI** (SSUI), a sleek, retro-themed UI to tame the wild beast that a Stationeers dedicated server is. It is glorious—automatic SteamCMD setups, one-click controls, Discord bots, and a backup system smarter than your average space engineer. But then, @mitoskalandiel dropped a galaxy-sized idea: *Why stop at Stationeers?* Why not make SSUI the ultimate overlord of *any* Steam server? And so, **SteamServerUI** was born, with @JacksonTheMaster and @mitoskalandiel leading the charge to generalize the chaos of server management.

**SteamServerUI (v6)** is the shiny, in-development evolution of SSUI, designed to run *any* Steam game server that can be wrangled with a `runfile`. Think Satisfactory, Project Zomboid, Stationeers, or even that obscure indie game you love (as long as you write a `runfile` for it). Meanwhile, Stationeers fans, fear not: **StationeersServerUI (v5)** remains a dedicated, maintained LTS version, chizzled in time as a rock-solid option for your spacefaring needs. v6 as **SteamServerUI** is a separate beast, and it won’t mess with v5’s vibe. It will probably be a while before v6 is released, and the details are still being worked out.
I am not sure if SteamServerUI will (should) move to a new repo later, but I will update this readme to reflect that once I know more about where I wanna go with this project and the Stationeers version.

> ⚠️ **Warning**: v6 is a *preview*. It’s like a prototype spaceship—cool, but expect a lot of loose bolts. This is meant for development, not production. Some Stationeers-specific features (like BackupManager and Discord) are currently in a state of… let’s call it “creative flux.” Non-breaking, but they complain with some noise.

## ✨ What’s the Big Deal? ✨

SteamServerUI takes the pain out of running Steam game servers by wrapping everything in a retro-styled web UI that’s equal parts nostalgic and powerful. The secret sauce? The **runfile**, a JSON config that tells SteamServerUI how to launch and manage your game server. Define the game, its executables, and its command-line args, and SteamServerUI handles the rest—downloading from SteamCMD, parsing args, and serving it all up in a UI that currently screams “I grew up playing DOS games.” Sorry not sorry. Looking for UI devs, btw.

### Games We’ve Tested (So Far)
- 🏭 **Satisfactory**: Build factories, crash servers, repeat.
- 🧟 **Project Zomboid**: Survive zombies, not server crashes.
- 🚀 **Stationeers**: Space is hard, server management isn’t.

### Games You *Could* Run (With a Runfile)

- 🔫 **Counter-Strike 2**: The king of FPS. Frag, flash, and manage servers with retro flair.
- ⚒️ **Rust**: Build bases, break trust, and run servers without breaking a sweat.
- 🧙 **ARK: Survival Evolved**: Dino-taming chaos. Config via `.ini` files for max roar.
- 🎖️ **Arma 3**: Military sims with less “operation: fix the server.” Needs `.cfg` tweaks.
- 💣 **Team Fortress 2**: Hats, rockets, and Source engine goodness, UI-managed.
- ⚔️ **Valheim**: Viking epicness. Skål to simple server setups!
- 🧟 **DayZ**: Survive zombies and server woes. Config files for extra grit.
- 🪓 **Garry’s Mod**: Sandbox insanity, from TTT to DarkRP, all UI-controlled.
- 🏹 **7 Days to Die**: Zombie hordes meet easy server launches. Tweak `.xml` for details.
- 🏰 **Conan Exiles**: Barbarian servers with more loincloths, less manual config.
- 🚀 **Space Engineers**: Build starships, not server scripts. Minimal args, in-game setup.
- 🪐 **No Man’s Sky**: Infinite worlds, one server. In-game config for galaxy hopping.
- 🧑‍🚀 **Astroneer**: Planet exploration without config file black holes.
- 🛡️ **Mount & Blade II: Bannerlord**: Lead armies, not error logs. `.cfg` for fine-tuning.
- 🦁 **Eco**: Save the planet, one server at a time. In-game or `.json` config.
- 🏎️ **Assetto Corsa**: Race servers that outpace your lap-time crashes.
- 🪖 **Squad**: Tactical FPS with less tactical server fiddling. `.cfg` tweaks needed.
- 🧝 **V Rising**: Vampire servers that don’t suck (except blood). In-game setup.
- 🏜️ **Hurtworld**: Outback survival, minus outback setup woes.
- 🏝️ **The Forest**: Cannibals are scary; server setup isn’t. Minimal args, in-game config.
- 🛠️ **Factorio**: Automate factories, not server maintenance. `.json` for extras.
- 🦕 **The Isle**: Dino servers where the only thing extinct is manual setup.
- 🏴‍☠️ **Blackwake**: Pirate battles, no need to walk the config plank.
- 🔪 **Dead by Daylight**: Scream at killers, not errors. In-game server settings.
- 🚗 **Wreckfest**: Smash cars, not keyboards over configs.
- 🏍️ **MX Bikes**: Dirt bike servers that don’t leave you in the mud.
- 🗡️ **Mordhau**: Medieval mayhem with modern server ease. `.ini` for details.
- 🦑 **Depth**: Sharks vs. divers, with UI-managed servers.
- 🏠 **Unturned**: Blocky survival with unblocky server management.
- 🛸 **Empyrion - Galactic Survival**: Conquer galaxies, not command lines.
- 🌌 **Stellaris**: Rule the stars, let SteamServerUI rule the server. In-game config.
- 🏞️ **Outpost Zero**: Sci-fi survival without sci-fi setup pain.
- 🦁 **Planet Zoo**: Manage animals, not server zoos. In-game or `.cfg` tweaks.
- 🏹 **Pavlov VR**: VR shootouts with non-virtual server simplicity.
- ⚽ **Rocket League**: Car soccer servers, UI-managed for epic goals.
- 🪖 **Hell Let Loose**: WWII battles with less server-side warfare. `.cfg` needed.
- 🧟 **Left 4 Dead 2**: Zombie co-op with Source engine server ease.
- 🏰 **Chivalry 2**: Medieval brawls, UI-managed for knightly glory.
- 🛠️ **Satisfactory**: Factory-building chaos, already tested and UI-approved!
- 🧟 **Project Zomboid**: Zombie survival with UI-managed servers, tested and true.
- 🚀 **Stationeers**: Space survival, our OG love. LTS version for diehards.
- 🪖 **Insurgency: Sandstorm**: Tactical firefights with less server friction.
- 🏜️ **Miscreated**: Post-apocalyptic survival, UI-managed for less doom.
- 🧙 **Barotrauma**: Submarine horrors with surface-level server ease.
- 🏹 **Blade & Sorcery**: VR swordfights with non-VR server simplicity.
- 🦁 **Zoo Tycoon**: Animal parks without server parkour. In-game config.
- 🏎️ **iRacing**: Sim racing servers that don’t spin out on setup.
- 🪓 **SCUM**: Island survival with less server struggle. `.ini` for details.
- 🧙 **Avorion**: Space sandbox servers, UI-managed for galactic fun.
- 🏰 **Stronghold: Crusader**: Castle sieges with modern server ease.
- 🛠️ **Stormworks: Build and Rescue**: Save lives, not servers. In-game config.
- 🦑 **Natural Selection 2**: Alien vs. marine servers, UI-managed for bites.
- 🏹 **Kādomon: Hyper Auto Battlers**: Auto-battler servers with auto-easy setup.
- …or *any* SteamCMD-compatible game, if you’re brave enough to write the runfile!


## 🛠️ The Runfile: Heart of SteamServerUI

The `runfile` is a JSON file that defines how to run a game server. It specifies the game’s Steam App ID, executables, and command-line arguments, which SteamServerUI uses to launch and manage the server. Here’s a sneak peek at a Satisfactory runfile:

```json
{
  "meta": {
    "name": "Satisfactory",
    "version": "1.1"
  },
  "architecture": "linux",
  "steam_app_id": "1690800",
  "windows_executable": "./FactoryServer.exe",
  "linux_executable": "./FactoryServer.sh",
  "args": {
    "basic": [
      {
        "flag": "-log",
        "default": "",
        "required": false,
        "requires_value": false,
        "description": "Forces the server to display logs in a window or terminal",
        "type": "bool",
        "ui_label": "Display Logs",
        "ui_group": "Basic",
        "weight": 10
      },
      ...
    ],
    "network": [...],
    "advanced": [...]
  }
}
```

SteamServerUI reads this, builds the command-line args, and serves them up in the UI for easy tweaking. Want to run Valheim instead? Select the `runValheim.ssui` file from the UI's runfile libary (not implemented), and you’re off to the races.SSUI handles the rest, no PhD in command-line wizardry required, writing a bash scrupt to run the game server will be a thing of the past.

### But what about ThisUnknownGame.exe?

The "plan" is to have a community-driven runfile library, where users can submit generalized, basic working versions of runfiles for games they love. This will 99% be a github branch or repo where the repo admins merge community runfiles into the main branch when tested and ready for broader use.

## 🌟 Features Currently Implemented 🌟

| 🚀 Zero Config | 🔄 Auto Updates | 🎮 One-Click Control | 🔒 Secure by Default | 🛠️ Mod Support |
|:-------------:|:---------------:|:-------------------:|:-------------------:|:--------------------:|
| Drop and run | SteamCMD updates | Start/stop with ease | JWT auth, TLS | BepInEx integration |

> **Note**: v6 is a work in progress. It currently IS able to run a Server successfully in a playable state, but the UI requires developer knowlege at some points to understand how to use it. Not recommended for production use, but technically feasible. 

## TL;DR - Get Started (If You Dare)

1. 📦 Grab the v6 branch with a git clone
2. 📁 build the project with `go build build/build.go` and copy the executable of your system to the root of the project (chmod +x) and execute it.
3. 🌐 Access the UI at `https://<<server-ip>>:8443`.

## Why You’ll (Eventually) Love It

- **Generalized Power**: One UI to rule *all* Steam servers (with the right runfile).
- **Community-Driven**: Built by @JacksonTheMaster, inspired by @mitoskalandiel, and all open to your ideas. 

## 🗺️ Documentation (v6)

There is currently NO documentation for version 6.
Earlier versions are documented in the [GitHub Wiki](https://github.com/JacksonTheMaster/StationeersServerUI/wiki).

## 🙌 Contributing

SteamServerUI is a community effort, and we’d love your input (but no pressure). Got a bug? [Open an issue](https://github.com/JacksonTheMaster/StationeersServerUI/issues). Got a runfile for your favorite game? Share it! See the [Contributing Guidelines](https://github.com/JacksonTheMaster/StationeersServerUI/wiki/Contributing) for details.
The License is there to protect this project, not to scare you away. It's Read-Source anyway!

Special thanks to:
- **@mitoskalandiel**: For the galaxy-brain idea to go beyond Stationeers, and providing me with a powerful Linux Server to test on!

## 📜 License

This project is licensed under the STATIONEERS SERVER UI LICENSE AGREEMENT (could update the name eventually) - see the [LICENSE](LICENSE) file for details.
