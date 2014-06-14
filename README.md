Attacher
========

Provides socket access to lxc-attach command

## Testing

`sudo ./attacher -s /tmp/attacher.sock`
`sudo netcat -U /tmp/attacher.sock`

## Protocol

```
--> {"attach": "597713fb3a55156a3e32e16be90559b779deb2cc520bf35079f406e6e63bfcb9", "user": "action"}\n
<-- {"attach": "success"}\n
<-- ["d", "action@foobar4:/$ "]\n
--> ["d", "echo hello world\n"]\n
<-- ["d", "echo hello world\r\nhello world\r\naction@foobar4:~$ "]
--> ["disconnect", ""]\n
```

```
--> {"some_unknown_command": "foobar"}
<-- {"error": "wtf"}
```

- - -
Copyright (c) 2014 Irrational Industries, Inc. d.b.a. Nitrous.IO. All Rights Reserved.
