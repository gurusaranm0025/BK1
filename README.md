# CBak (Configuration Backup) v0.6.5-alpha

Its a tool I made to take backup of my configuration files all over in my linux system, and restore them in my other computers.

## Built in flags:

    -C : takes a json file as input, which describes which paths that needs to be taken backup. 
    (simple, tags are present for taking backups)
    -E : to extarct the backups
    -R : it takes the backed up file as input and restores the the files in that backup
    -o : tell your own output file name
    -p : custom path to take backup
    -v : version
    -h : help

It has its own compression methods built into it, and a own method to open the backup and view it. And it can't be opened using other tools such as ark.

## Installation

### Step 1:
To install the only dependency **make**.

#### For Arch Linux:

    sudo pacman -S make

#### For Ubuntu/Debian:

    sudo apt update
    sudo apt install build-essential

#### For CentOS/RHEL 7:

    sudo yum install make

#### For CentOS/RHEL 8 and Fedora:

    sudo dnf install make

### Step 2:
Clone this repository, where you feel safe,

    git clone https://github.com/gurusaranm0025/CBak.git
    cd CBak

then run ,

    sudo make

for installing the package.

