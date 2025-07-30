# BrightSign Player SCP

This is a command line tool written in go.  It uses the Diagnostic Web Server (DWS) APIs (described in bs-api-docs-20250614).  It is similar to the "scp" command in linux, but it will be called "bscp" instead.  The host portion of the designation can be an IP address or a hostname.  Similar to scp, it uses the format of host:full-path and writes the file to full path.  It will use a similar CLI as scp but it uses the file upload API to copy the file to the player.  It will never try to use SSH keys though, instead it will always prompt for a password.  That password will be the DWS password.

The tool will then make the API call to upload the file.  It will always copy to the SD card (e.g. to /storage/sd) and then will make an API call to list the files and ensure that the file is there.  If yes, it returns success otherwise it prints an error and errors out.

Create the file system layout according to go best practices.  Also build unit tests for all functions.
