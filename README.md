![](./img/banner.jpg)

# ArozOS Launcher

The ArozOS default launcher for over the air (OTA) updates



The ArozOS Launcher is the default launcher that ship with the ArozOS Raspberry Pi image since v1.119. The launcher is a very basic program that start arozos and update the arozos if an update folder exists. 

## Supported Platforms

The launcher only support a limited number of platforms including

- Windows (amd64)
- Linux (amd64, armv6/7, arm64)
- MacOS (amd64)



## Build

Require Go 1.17 or above

```
git clone https://github.com/aroz-online/launcher
cd ./launcher
go mod tidy
go build
```

For building for cross platforms, see ```build.bat```

## Usage

### Automatic Net Install (for fresh installation)

(This feature is only added since Launcher v1.3 (for ArozOS v2.011))

If you have a clean Linux installation and want to install ArozOS from Github Release, you can use the following command to download and start the launcher. The launcher will download the release of ArozOS from Github and start it for you.

Here is an example for using the Launcher net install function

```bash
wget https://github.com/aroz-online/launcher/releases/download/v2.010/launcher_linux_arm64
mv launcher_linux_arm64 launcher
sudo chmod 775 ./launcher
./launcher
```

*Depending on your network speed, the web.tar.gz download might take some time, sometime up to 8+ minutes on low end SBCs.*

### Folder Structure (for adding OTA Update to existing ArozOS systems)

Put the launcher binary inside the arozos source root (the src folder if you are cloning from Github repo) with a executable binary next to it in the same directory, and start the launcher with the start command.

Here is an example folder structure that will launch correctly

- web/
- system/
- arozos
- launcher



### Start Command

All parameter passed to the launcher will be passed to the arozos binary which the launcher will pick from wildcard arozos_* or ./arozos (or arozos.exe if you are running it on windows).  Example:

```
./launcher -port 80 -hostname "MyServer"
```

If you are already having a startup script that calls to arozos binary, you can simply switch out the ```./arozos``` to ```./launcher```. Here is an example of such switch

```
# Original
./arozos -port 80 -tls=true -tls_port 433 -hostname "MyServer"

# Using launcher
./launcher -port 80 -tls=true -tls_port 433 -hostname "MyServer"
```

### Port Occupy

This launcher will listen to ```127.0.0.1:25576``` for arozos launcher check. If you have another application that is using this port, you might not be able to use arozos with this launcher.

### Auto Configuration Backup

This launcher will backup the following files and migrate them to the newer version of arozos. 

```
system/bridge.json
system/dev.uuid
system/cron.json
system/storage.json
web/SystemAO/vendor/*
```

**By design the database files will not be overwritten by updates. However, we still do recommend users to backup their important data and database before applying an OTA update from their vendor.**

### Changing Update Download Source

**The update download was not done by the launcher.** The launcher only launch arozos, backup old arozos folder structure and apply updates (overwrite the current arozos with contents inside the ```updates/``` folder) if required folder structure exists. See ```arozos/src/system/update.json``` for more details.

## License

MIT
