# CBak (Configuration Backup)

Its a tool I made to take backup of my configuration files all over in my linux system, and restore them in my other computers.

## Built in flags:

    -C : takes a json file as input, which describes which paths that needs to be taken backup. (simple, tags are present for taking backups)
    -E : to extarct the backups
    -R : it takes the backed up file as input and restores the the files in that backup
    -o : tell your own output file name
    -p : custom path to take backup
    -v : version
    -h : help

It has its own compression methods built into it, and a own method to open the backup and view it. And it can't be opened using other tools such as ark.