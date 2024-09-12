# Quick Start

## Install Obsidian Plugin

### Installing from the Plugin Marketplace

> As the plugin is not yet available on the official community marketplace, this method is currently unavailable.

1. Open "Settings" and click on "Third-party Plugins" on the left side of the window.
2. If you are using plugins for the first time, you will be prompted to disable "Safe Mode."
3. Click "Disable Safe Mode" in the "Third-party Plugins" section.
4. After disabling Safe Mode, click the "Browse" button next to "Community Plugin Market."
5. In the "Browse" interface, search for `BetaX NAS Sync` to find the plugin, then click on the search result.
6. On the `BetaX NAS Sync` plugin details page, click the "Install" button. The installation is now complete.

### Manual Installation from Repository

1. Visit the [Releases](https://github.com/skye-z/ons/releases) page.
2. Download the latest version's `main.js` and `manifest.json`.
3. Open your `Obsidian` vault directory, append `/.obsidian` to the path, and press Enter.
4. Create a `plugins` directory inside `.obsidian`, then navigate into the `plugins` directory.
5. Create an `ons` directory within `plugins` and place the downloaded `main.js` and `manifest.json` files in it.
6. Restart `Obsidian` to see `BetaX NAS Sync` in the "Installed Plugins." The installation is now complete.

## Install NAS Service

### Installation via App Store

> Currently, there is no supported vendor; future plans include supporting Synology SPK installation.

### Installation via Docker Container

1. Pull the `skyezhang/ons-nas` image.
2. Map the container's `/app/vault` directory to your NAS's `Obsidian` vault directory with read/write permissions.
3. Open all ports on the UDP side and open port `9892` on the TCP side.
4. Once configured, start the container.

> Note: This service uses UDP hole punching technology, so all ports need to be open on the UDP side, and port `9892` is for the web interface.
