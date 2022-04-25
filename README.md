# Go Orderbook

A fast (4m TPS) orderbook in Go.

### Features

- Market orders
- Limit orders (GTC)
- Add, Cancel, Amend (quantity)
- Best price and depth of each side
- Last filled price

### Todo

- Persistence / restore from snapshot
- Add more execution styles (FOK, IOC etc)
- Test and finish inbound messages
- Implement outbound messages

I'm not actively working on this but plan to pick it up again some time.

### I saw a Linked List, ew

Linked lists are misunderstood data structures that are powerful if combined with hashmaps.

Add an order:
Adding to a linked list is an optimal operation that doesn't require any re-allocations (could possibly be further optimized using a mempool).

Cancel order:
This is the part that people misunderstand about linked lists; we don't iterate to search. We keep a reference to the list's node by order ID in a hash map and so we can look up a node in O(1). After that a cancel is simply repointing the adjacent nodes. Compare this to a VecDeque (ring buffer) for instance which would require a binary search to the order and then a memshift of all items before or after.

Matching:
Matching happens <10% of the time and an order is typically matched in a small amount of iterations, thus the loss in performance compared to iterating contiguous memory is negligable compared to the gains in add and cancel.

Combined with a balanced tree of levels we are able to have a performant algorithm that is simple and memory efficient. All other "hyper" performant algorithms I could think of came with notable sacrafices (complexity, memory, or attack vectors).
