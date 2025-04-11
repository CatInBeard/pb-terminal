# pb-terminal
A terminal emulator for PocketBook e-readers running on Linux.
## Description
A terminal emulator for PocketBook e-readers, written in Go. It uses the libinkview library implementation by dennwc: https://github.com/dennwc/inkview. The project is built using the Docker image from Skeeve: https://github.com/Skeeve/SDK_6.3.0.
## Installation
To install the application, connect your e-reader to your computer via a cable and copy the release file to the `applications` folder. On Linux, you can build the application from source using the command `make build` and automatically install it on the connected e-reader using the command `make instal`
## License
The project is licensed under the MIT license.
## Warning
Using a terminal emulator on an e-reader can pose risks when used. The terminal provides access to system files and commands that can potentially harm your e-reader or disrupt its operation. If you are not familiar with the command line and do not know what you are doing, you may accidentally:
* Delete important system files or settings
* Change system configuration, leading to unstable operation or crashes
* Install unverified or malicious software
* Access sensitive information or settings
Please use the terminal emulator with caution and only if you are confident in your actions. If you are unsure, it is better not to use the terminal or seek help from an experienced user.
## Disclaimer
The developer of the terminal emulator disclaims any warranty, express or implied, including, but not limited to, the warranties of merchantability, fitness for a particular purpose, and non-infringement. You use the terminal emulator at your own risk. No liability is assumed for any damages or losses arising from the use of the terminal emulator.