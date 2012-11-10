GoStatusBoard
=============

A simple Go frontend to a json database. This is a simple "key store" that allows you to save statues of various services or objects.  Updates are automatically timestamped.  Contents of database and be dumped to a html file.  "Success" and "Fail" status are considered "special" status and will be highlighted with green or red in the html output

Usage
-----

    Usage of statusboard:
         [-d] Turn on debugging
         update OBJECT STATUS - Set objects status
         output - Dump current statuses

Example
-------
    statusboard update httpd Running
    statusboard update smtpd Failed

Crontab Examples
----------------
    00 06 * * * if someCommandThatReturnsAnErrorCode; then statusboard update "My Thing" "Success"; else statusboard update "My Thing" "Fail"; fi

Script Example
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
