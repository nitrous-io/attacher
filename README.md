Attacher
========

Provides unix socket access to lxc-attach command

## Testing

`sudo ./attacher -s /tmp/attacher.sock`
`sudo netcat -U /tmp/attacher.sock`

## Protocol

```
u s
<-- {"connected": true}
--> {"attach": "597713fb3a55156a3e32e16be90559b779deb2cc520bf35079f406e6e63bfcb9", "user": "action"}
<-- {"attached": true}
<-- ["d", "action@foobar4:/$ "]
--> ["d", "echo hello world\n"]
<-- ["d", "echo hello world\r\nhello world\r\naction@foobar4:~$ "]
--> ["disconnect", ""]
```

```
<-- {"connected", true}
--> {"some_unknown_command": "foobar"}
<-- {"error": "wtf"}
```

- - -
Copyright (c) 2014 Irrational Industries, Inc. d.b.a. Nitrous.IO. All Rights Reserved.
