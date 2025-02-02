package p2p

// All ETH68 message codes.
// We will be only supporting the newest ETH protocol.

const ETH = 68

const (
	// Base protocol msg codes: 0x00...0x10 (0-16)
	// Only 4 are used for now.
	HANDSHAKE = 0x00
	DISC      = 0x01
	PING      = 0x02
	PONG      = 0x03

	// Extended protocol msg codes: 0x10...0x1C (16-28)
	STATUS                        = 0x10
	NEW_BLOCK_HASHES              = 0x11
	TRANSACTIONS                  = 0x12
	GET_BLOCK_HEADERS             = 0x13
	BLOCK_HEADERS                 = 0x14
	GET_BLOCK_BODIES              = 0x15
	BLOCK_BODIES                  = 0x16
	NEW_BLOCK                     = 0x17
	NEW_POOLED_TRANSACTION_HASHES = 0x18
	GET_POOLED_TRANSACTIONS       = 0x19
	POOLED_TRANSACTIONS           = 0x1A
	GET_RECEIPTS                  = 0x1B
	RECEIPTS                      = 0x1C
)
