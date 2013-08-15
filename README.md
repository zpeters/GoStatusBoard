GoStatusBoard
=============

A simple Go frontend to a json database. This is a simple "key store" that allows you to save statues of various services or objects.  Updates are automatically timestamped.  Contents of database and be dumped to a html file.  "Success" and "Fail" status are considered "special" status and will be highlighted with green or red in the html output

Usage
-----

    Usage of statusboard:
         [-d] Turn on debugging
         update OBJECT STATUS - Set objects status
         output - Dump current statuses
	 test COMMAND OBJECT SUCCESS_STATUS FAIL_STATUS - Run command COMMAND if exit code is 0 update OBJECT with SUCCESS_STATUS, otherwise OBJECT with FAIL_STATUS

Example
-------
    statusboard update httpd Running
    statusboard update smtpd Failed
    statusboard test SomethingThatFails Service Running Failed (updates database with "Service Failed" assuming it fails)
    statusboard test "ls -ahl" Listing Passed Failed (updates database with "Listing Passed" assuming it passed)

Crontab Examples (old version)
----------------
    00 06 * * * if someCommandThatReturnsAnErrorCode; then statusboard update "My Thing" "Success"; else statusboard update "My Thing" "Fail"; fi

Crontab Examples (new)
----------------
    00 06 * * * statusboard test someCommandThatReturnsAnErrorCode "My Thing" "Success" "Fail"

Script Example (a simple script that returns an error code)
--------------
```shell
#!/bin/sh
tar czf /tmp/mytar.tar.gz /myfiles
TAR_RESULT=$?
if [ ! $TAR_RESULT ]
then
	statusboard update "Tar Operation" "Fail"
else
	statusboard update "Tar Operation" "Success"
fi
```

HTML Output
-----------
    statusboard output > /var/www/html/status.html
