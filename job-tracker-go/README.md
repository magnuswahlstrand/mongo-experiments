

## Todo

* [ ] Deduplicate jobs
* [*] List jobs
* [ ] Run tests on push
* [ ] Progress state machine
* [ ] Locking mechanism
  * https://github.com/qiniu/qmgo/pull/95/files
```json
{
  "query": {
    "user_id": "Obj()",
    "state": "CART"
  },
  "update": {
    "$set": {
      "state": "PRE-AUTHORIZE"
    }
  },
  "new": true
}

```
* [ ] Unlock stuck jobs
* [ ] Build CLI
* [x] Use test-containers
