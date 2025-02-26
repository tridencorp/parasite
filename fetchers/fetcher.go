package fetchers

// Fetchers are responsible for fetching and validating data from the network.
//
// Since the blockchain is a decentralized and public network, we cannot be 100%
// sure that the data we receive is valid. We must cross-check it with other peers,
// and this will be the fetcher's job.
type Fetcher struct {
	PeerCount  uint8
	MatchCount uint8
}
