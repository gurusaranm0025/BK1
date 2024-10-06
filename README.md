# CBak (Configuration Backup) v0.6.8-alpha

Its a tool I made to take backup of my configuration files all over in my linux system, and restore them in my other computers.

## Built in flags:

    -C : takes a json file as input, which describes which paths that needs to be taken backup. 
    (simple, tags are present for taking backups)
    -E : to extarct the backups
    -R : it takes the backed up file as input and restores the the files in that backup
    -o : tell your own output file name
    -p : custom paths to take backup, for providing more than one path then separate each one with a ','. 
    -v : version
    -h : help

It has its own compression methods built into it, and a own method to open the backup and view it. And it can't be opened using other tools such as ark.

## Installation

### Step 1:
To install the dependencies **make and go**.

#### For Arch Linux:

    sudo pacman -S make go

#### For Ubuntu/Debian:

    sudo apt update
    sudo apt install build-essential golang

#### For CentOS/RHEL 7:

    sudo yum install make golang

#### For CentOS/RHEL 8 and Fedora:

    sudo dnf install make golang

### For OpenSUSE:

    sudo zypper install make go

### For Alpine linux:

    sudo apk add make go

### Step 2:
Clone this repository, where you feel safe,

    git clone https://github.com/gurusaranm0025/CBak.git
    cd CBak

then run ,

    sudo make

for installing the package.

## Backup.conf

Its a json file, from which you can tell which files and folders u want to take backup.

Below is a simple example,

        {
            "BackupName": "testing",
            "BackupPaths": [
                "~/Downloads/ZenBrowser"
            ],
            "Tags": ["hypr","rofi","dunst"]
        }

Here, 

    BackupName  - defines the backup file name.

    BackupPaths - it takes various paths of the files and folders which you want to take backup.

    Tags        - tags are built-in flags which already points towards certain location which may contain config files. There are five available tags right now.

The available tags are,

    hypr    - points to ~/.config/hypr
    rofi    - points to ~/.config/rofi
    waybar  - points to ~/.config/waybar
    wlogout - points to ~/.config/wlogout
    dunst   - points to ~/.config/dunst

