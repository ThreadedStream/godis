
# What is godis?
  Godis is just another implementation of in-memory key-value store written in golang. It was conceieved with educational purposes in mind.
  Despite its obvious simplicity, it can be used for some production tasks, for example when you need some temporary storage for saving short living data.

  
## What i have achieved so far
  <ul> 
    <li>Basic functionality of server is implemented and tested</li>
    <li>Client part has been finished, although it needs some improvements</li>
  </ul>

## List of currently supported commands
  1) [SET](#SET), [HSET](#HSET), [MSET](#MSET) 
  2) [GET](#GET), [HGET](#HGET), [MGET](#MGET)
  3) [KEYS](#KEYS)
  4) [DEL](#DEL)


## SET
Associates an input key to some value. 
```bash
  SET key value
```

Or if you want to enable ttl for the given key

```bash
  SET key value 120 #Will be deleted after 2 minutes 
```

## HSET
