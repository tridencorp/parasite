```
@@@@@@@    @@@@@@   @@@@@@@    @@@@@@    @@@@@@   @@@  @@@@@@@  @@@@@@@@
@@@@@@@@  @@@@@@@@  @@@@@@@@  @@@@@@@@  @@@@@@@   @@@  @@@@@@@  @@@@@@@@
@@!  @@@  @@!  @@@  @@!  @@@  @@!  @@@  !@@       @@!    @@!    @@!
!@!  @!@  !@!  @!@  !@!  @!@  !@!  @!@  !@!       !@!    !@!    !@!
@!@@!@!   @!@!@!@!  @!@!!@!   @!@!@!@!  !!@@!!    !!@    @!!    @!!!:!
!!@!!!    !!!@!!!!  !!@!@!    !!!@!!!!   !!@!!!   !!!    !!!    !!!!!:
!!:       !!:  !!!  !!: :!!   !!:  !!!       !:!  !!:    !!:    !!:
:!:       :!:  !:!  :!:  !:!  :!:  !:!      !:!   :!:    :!:    :!:
 ::       ::   :::  ::   :::  ::   :::  :::: ::    ::     ::     :: ::::
 :         :   : :   :   : :   :   : :  :: : :    :       :     : :: ::

```

# PARASITE

Small P2P/RPC server for ETH and BTC.

# GOAL
The goal is to have something small, easy to understand, with fast TCP api and much more simpler database design than Geth - without stupid Trie (don't confused it with Merkle Tree which is actually neat !).

# INSTALL
```go
go mod init parasite
go mod tidy

go build ./cmd/parasite
```

# RUN
```go
./parasite
```
# CURRENT STATE
It’s under heavy development. I’m switching between many projects for Triden right now, so progress will come — baby steps.