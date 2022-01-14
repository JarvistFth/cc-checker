This is a "key_value_store" chaincode.
The chain-code support four different function in order to achieve a voting session. The function inputs are as follows: 

- put(k,v) : it store a couple key 'k' and value 'v'
- get(k):    get a value associated to a key 'k'
- getAll():  retrieve all the keys stored
- delete(k): delete a couple by referencing the key 'k'