# AltWebServer
Inspired by [AltStore](https://github.com/altstoreio/AltStore) and [AltServer-Linux](https://github.com/NyaMisty/AltServer-Linux), 
AltWebServer is a web server helps people side load IPAs on their iDevices.

# Requrements
* [libimobiledevice](https://github.com/libimobiledevice/libimobiledevice)
* [idevice_installer](https://github.com/libimobiledevice/ideviceinstaller)
* [usbmuxd](https://github.com/libimobiledevice/usbmuxd)
* [netmuxd](https://github.com/jkcoxson/netmuxd)
* [AltServer-Linux](https://github.com/NyaMisty/AltServer-Linux)

# Usage 
1. install libimobiledevice, idevice_installer, usbmuxd and netmuxd on your computer
2. pair your iDevice with usbmuxd, then stop usbmuxd
3. run netmuxd
4. download alterser-linux executable on your computer and make it executable
5. mv config.json.example config.json
6. config your Apple ID, password, altserver_path and anisett_url(optional) in your config.json. **Note: Your Apple ID, and it's password will be saved as plaintext on your computer, make sure your computer is not compromised or your Apple ID might get STOLEN**
7. run ./altwebserver
