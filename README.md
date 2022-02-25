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
