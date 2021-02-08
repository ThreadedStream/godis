
# What is godis?
  Godis is just another implementation of in-memory key-value store written in golang. It was conceieved with educational purposes in mind.
  Despite its obvious simplicity, it can be used for some production tasks, for example when you need some temporary storage for saving short living data.

  


## List of currently supported commands
  1) [SET](#SET), [HSET](#HSET), [MSET](#MSET) 
  2) [GET](#GET), [HGET](#HGET), [MGET](#MGET)
  3) [KEYS](#KEYS)
  4) [DEL](#DEL)
  5) [SAVE](#SAVE), [RESTORE](#RESTORE)
  6) [SIGNUP](#SIGNUP), [LOGIN](#LOGIN), [LOGOUT](#LOGOUT), [WHOAMI](#WHOAMI)


## SET
Associates an input key to some value. 
```bash
  SET key value
```

Or if you want to enable ttl for the given key

```bash
  SET key value 120 #Will be deleted after 2 minutes 
```

I should note that value acceptance range is limited to 3 types, namely strings, dictionaries, and lists.<br>

```bash
  SET key {key1:value1} #Associates given key with the provided dictionary value
```

```bash
  SET key [value1, value2,value3] #Similarly, associates given key with provided list of values
```

## HSET

HSET is similar to SET command -- distinction lies in the fact that now client has to provide field,
which enables access to the value. Apparently, field is stored in its SHA-256 hash representation.

```bash
  HSET key field value
```

## MSET
MSET allows a client to provide more key-value pairs, here is an example:

```bash
  MSET key1 value1 key2 value2 key3 value3
```

## GET
GET command simply retrieves value associated with the given key.

```bash
  SET key value
  # "OK"
  GET key
  # value
```

## HGET
HGET command provides client with a value only in case if proper key and value were passed:

```bash
  HSET key field value
  # "OK"
  HGET key field
  # value
```

## MGET
MGET returns values for multiple given keys

```bash
  MSET key1 value1 key2 value2 key3 value3
  # "OK"
  MGET key1, key2, key3
  # 1)value1
  # 2)value2
  # 3)value3
```

## KEYS
KEYS command returns all keys conforming to given pattern

```bash
  MSET key1 value1 key2 value2 key3 value3
  # "OK"
  KEYS key[1-9]*
  # 1)key1
  # 2)key2
  # 3)key3
```

## DEL
DEL command deletes entry containing a given key

```bash
  SET key value
  # "OK"
  KEYS key
  # "OK"
```

## SAVE
SAVE command allows a client to save current contents of key-value store on disk.
It accepts a single parameter being the path to the file. Command normally handles 
the case when file does not exist -- it will readily create one.

```bash
  SAVE store.txt
  # "OK"
```

## RESTORE
RESTORE command restores contents of previously saved store.

```bash
  RESTORE store.txt
  # "OK"
```

## SIGNUP
SIGNUP command allows you to save your credentials in database
```bash
  SIGNUP user password
  # "OK"
```
As a side note, all passwords are being hashed using sha-256 before being put into database.

## LOGIN
LOGIN command allows you to login into godis app using previously stored credentials
```bash
  SIGNUP user password
  # "OK"
  LOGIN user password
  # "OK"
```

## LOGOUT
LOGOUT command simply logs user out of godis
```bash
  LOGOUT 
  # "OK"
```

## WHOAMI
WHOAMI command spits out user's username. So putting it all together yields the following:
```bash
  
  SIGNUP user password
  # "OK"
  LOGIN user password
  # "OK"
  WHOAMI
  # user
  LOGOUT 
  # "OK"
  WHOAMI
  # "anonymous"
```


## How does it all work
Godis' architecture is very simple. The client and server are running on one machine, since 
they are a part of a single application (they are going to be decoupled soon). At first, application attempts to run a server, 
if it succeeds, then database initialization comes right after it.
Finally, once all needed components are up and running, app will trigger execution of a client pipe(fancy name for single component)
,which welcomes user with a message "You may start typing some commands". 

## Some notes for future self

The whole architecture design should be refactored. The first and foremost is to 
decouple client and server, in order for them to be independent components.

After accomplishing the latter, I should make an attempt to dockerize application, since
now it ends up being pretty wonky, by virtue of the fact that docker's passing a lot of EOFs when client starts, and 
the whole application stops reacting to user input. 
