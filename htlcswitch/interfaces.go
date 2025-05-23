package htlcswitch

import (
	"github.com/ltcsuite/lnd/channeldb"
	"github.com/ltcsuite/lnd/channeldb/models"
	"github.com/ltcsuite/lnd/invoices"
	"github.com/ltcsuite/lnd/lnpeer"
	"github.com/ltcsuite/lnd/lntypes"
	"github.com/ltcsuite/lnd/lnwallet/chainfee"
	"github.com/ltcsuite/lnd/lnwire"
	"github.com/ltcsuite/lnd/record"
	"github.com/ltcsuite/ltcd/wire"
)

// InvoiceDatabase is an interface which represents the persistent subsystem
// which may search, lookup and settle invoices.
type InvoiceDatabase interface {
	// LookupInvoice attempts to look up an invoice according to its 32
	// byte payment hash.
	LookupInvoice(lntypes.Hash) (invoices.Invoice, error)

	// NotifyExitHopHtlc attempts to mark an invoice as settled. If the
	// invoice is a debug invoice, then this method is a noop as debug
	// invoices are never fully settled. The return value describes how the
	// htlc should be resolved. If the htlc cannot be resolved immediately,
	// the resolution is sent on the passed in hodlChan later. The eob
	// field passes the entire onion hop payload into the invoice registry
	// for decoding purposes.
	NotifyExitHopHtlc(payHash lntypes.Hash, paidAmount lnwire.MilliSatoshi,
		expiry uint32, currentHeight int32,
		circuitKey models.CircuitKey, hodlChan chan<- interface{},
		payload invoices.Payload) (invoices.HtlcResolution, error)

	// CancelInvoice attempts to cancel the invoice corresponding to the
	// passed payment hash.
	CancelInvoice(payHash lntypes.Hash) error

	// SettleHodlInvoice settles a hold invoice.
	SettleHodlInvoice(preimage lntypes.Preimage) error

	// HodlUnsubscribeAll unsubscribes from all htlc resolutions.
	HodlUnsubscribeAll(subscriber chan<- interface{})
}

// packetHandler is an interface used exclusively by the Switch to handle
// htlcPacket and pass them to the link implementation.
type packetHandler interface {
	// handleSwitchPacket handles the switch packets. These packets might
	// be forwarded to us from another channel link in case the htlc
	// update came from another peer or if the update was created by user
	// initially.
	//
	// NOTE: This function should block as little as possible.
	handleSwitchPacket(*htlcPacket) error
}

// dustHandler is an interface used exclusively by the Switch to evaluate
// whether a link has too much dust exposure.
type dustHandler interface {
	// getDustSum returns the dust sum on either the local or remote
	// commitment.
	getDustSum(remote bool) lnwire.MilliSatoshi

	// getFeeRate returns the current channel feerate.
	getFeeRate() chainfee.SatPerKWeight

	// getDustClosure returns a closure that can evaluate whether a passed
	// HTLC is dust.
	getDustClosure() dustClosure
}

// scidAliasHandler is an interface that the ChannelLink implements so it can
// properly handle option_scid_alias channels.
type scidAliasHandler interface {
	// attachFailAliasUpdate allows the link to properly fail incoming
	// HTLCs on option_scid_alias channels.
	attachFailAliasUpdate(failClosure func(
		sid lnwire.ShortChannelID,
		incoming bool) *lnwire.ChannelUpdate)

	// getAliases fetches the link's underlying aliases. This is used by
	// the Switch to determine whether to forward an HTLC and where to
	// forward an HTLC.
	getAliases() []lnwire.ShortChannelID

	// isZeroConf returns whether or not the underlying channel is a
	// zero-conf channel.
	isZeroConf() bool

	// negotiatedAliasFeature returns whether the option-scid-alias feature
	// bit was negotiated.
	negotiatedAliasFeature() bool

	// confirmedScid returns the confirmed SCID for a zero-conf channel.
	confirmedScid() lnwire.ShortChannelID

	// zeroConfConfirmed returns whether or not the zero-conf channel has
	// confirmed.
	zeroConfConfirmed() bool
}

// ChannelUpdateHandler is an interface that provides methods that allow
// sending lnwire.Message to the underlying link as well as querying state.
type ChannelUpdateHandler interface {
	// HandleChannelUpdate handles the htlc requests as settle/add/fail
	// which sent to us from remote peer we have a channel with.
	//
	// NOTE: This function MUST be non-blocking (or block as little as
	// possible).
	HandleChannelUpdate(lnwire.Message)

	// ChanID returns the channel ID for the channel link. The channel ID
	// is a more compact representation of a channel's full outpoint.
	ChanID() lnwire.ChannelID

	// Bandwidth returns the amount of milli-satoshis which current link
	// might pass through channel link. The value returned from this method
	// represents the up to date available flow through the channel. This
	// takes into account any forwarded but un-cleared HTLC's, and any
	// HTLC's which have been set to the over flow queue.
	Bandwidth() lnwire.MilliSatoshi

	// EligibleToForward returns a bool indicating if the channel is able
	// to actively accept requests to forward HTLC's. A channel may be
	// active, but not able to forward HTLC's if it hasn't yet finalized
	// the pre-channel operation protocol with the remote peer. The switch
	// will use this function in forwarding decisions accordingly.
	EligibleToForward() bool

	// MayAddOutgoingHtlc returns an error if we may not add an outgoing
	// htlc to the channel, taking the amount of the htlc to add as a
	// parameter.
	MayAddOutgoingHtlc(lnwire.MilliSatoshi) error

	// ShutdownIfChannelClean shuts the link down if the channel state is
	// clean. This can be used with dynamic commitment negotiation or coop
	// close negotiation which require a clean channel state.
	ShutdownIfChannelClean() error
}

// ChannelLink is an interface which represents the subsystem for managing the
// incoming htlc requests, applying the changes to the channel, and also
// propagating/forwarding it to htlc switch.
//
//	abstraction level
//	     ^
//	     |
//	     | - - - - - - - - - - - - Lightning - - - - - - - - - - - - -
//	     |
//	     | (Switch)		     (Switch)		       (Switch)
//	     |  Alice <-- channel link --> Bob <-- channel link --> Carol
//	     |
//	     | - - - - - - - - - - - - - TCP - - - - - - - - - - - - - - -
//	     |
//	     |  (Peer) 		     (Peer)	                (Peer)
//	     |  Alice <----- tcp conn --> Bob <---- tcp conn -----> Carol
//	     |
type ChannelLink interface {
	// TODO(roasbeef): modify interface to embed mail boxes?

	// Embed the packetHandler interface.
	packetHandler

	// Embed the ChannelUpdateHandler interface.
	ChannelUpdateHandler

	// Embed the dustHandler interface.
	dustHandler

	// Embed the scidAliasHandler interface.
	scidAliasHandler

	// IsUnadvertised returns true if the underlying channel is
	// unadvertised.
	IsUnadvertised() bool

	// ChannelPoint returns the channel outpoint for the channel link.
	ChannelPoint() *wire.OutPoint

	// ShortChanID returns the short channel ID for the channel link. The
	// short channel ID encodes the exact location in the main chain that
	// the original funding output can be found.
	ShortChanID() lnwire.ShortChannelID

	// UpdateShortChanID updates the short channel ID for a link. This may
	// be required in the event that a link is created before the short
	// chan ID for it is known, or a re-org occurs, and the funding
	// transaction changes location within the chain.
	UpdateShortChanID() (lnwire.ShortChannelID, error)

	// UpdateForwardingPolicy updates the forwarding policy for the target
	// ChannelLink. Once updated, the link will use the new forwarding
	// policy to govern if it an incoming HTLC should be forwarded or not.
	UpdateForwardingPolicy(models.ForwardingPolicy)

	// CheckHtlcForward should return a nil error if the passed HTLC details
	// satisfy the current forwarding policy fo the target link. Otherwise,
	// a LinkError with a valid protocol failure message should be returned
	// in order to signal to the source of the HTLC, the policy consistency
	// issue.
	CheckHtlcForward(payHash [32]byte, incomingAmt lnwire.MilliSatoshi,
		amtToForward lnwire.MilliSatoshi,
		incomingTimeout, outgoingTimeout uint32,
		heightNow uint32, scid lnwire.ShortChannelID) *LinkError

	// CheckHtlcTransit should return a nil error if the passed HTLC details
	// satisfy the current channel policy.  Otherwise, a LinkError with a
	// valid protocol failure message should be returned in order to signal
	// the violation. This call is intended to be used for locally initiated
	// payments for which there is no corresponding incoming htlc.
	CheckHtlcTransit(payHash [32]byte, amt lnwire.MilliSatoshi,
		timeout uint32, heightNow uint32) *LinkError

	// Stats return the statistics of channel link. Number of updates,
	// total sent/received milli-satoshis.
	Stats() (uint64, lnwire.MilliSatoshi, lnwire.MilliSatoshi)

	// Peer returns the representation of remote peer with which we have
	// the channel link opened.
	Peer() lnpeer.Peer

	// AttachMailBox delivers an active MailBox to the link. The MailBox may
	// have buffered messages.
	AttachMailBox(MailBox)

	// Start/Stop are used to initiate the start/stop of the channel link
	// functioning.
	Start() error
	Stop()
}

// ForwardingLog is an interface that represents a time series database which
// keep track of all successfully completed payment circuits. Every few
// seconds, the switch will collate and flush out all the successful payment
// circuits during the last interval.
type ForwardingLog interface {
	// AddForwardingEvents is a method that should write out the set of
	// forwarding events in a batch to persistent storage. Outside
	// sub-systems can then query the contents of the log for analysis,
	// visualizations, etc.
	AddForwardingEvents([]channeldb.ForwardingEvent) error
}

// TowerClient is the primary interface used by the daemon to backup pre-signed
// justice transactions to watchtowers.
type TowerClient interface {
	// RegisterChannel persistently initializes any channel-dependent
	// parameters within the client. This should be called during link
	// startup to ensure that the client is able to support the link during
	// operation.
	RegisterChannel(lnwire.ChannelID) error

	// BackupState initiates a request to back up a particular revoked
	// state. If the method returns nil, the backup is guaranteed to be
	// successful unless the justice transaction would create dust outputs
	// when trying to abide by the negotiated policy.
	BackupState(chanID *lnwire.ChannelID, stateNum uint64) error
}

// InterceptableHtlcForwarder is the interface to set the interceptor
// implementation that intercepts htlc forwards.
type InterceptableHtlcForwarder interface {
	// SetInterceptor sets a ForwardInterceptor.
	SetInterceptor(interceptor ForwardInterceptor)

	// Resolve resolves an intercepted packet.
	Resolve(res *FwdResolution) error
}

// ForwardInterceptor is a function that is invoked from the switch for every
// incoming htlc that is intended to be forwarded. It is passed with the
// InterceptedForward that contains the information about the packet and a way
// to resolve it manually later in case it is held.
// The return value indicates if this handler will take control of this forward
// and resolve it later or let the switch execute its default behavior.
type ForwardInterceptor func(InterceptedPacket) error

// InterceptedPacket contains the relevant information for the interceptor about
// an htlc.
type InterceptedPacket struct {
	// IncomingCircuit contains the incoming channel and htlc id of the
	// packet.
	IncomingCircuit models.CircuitKey

	// OutgoingChanID is the destination channel for this packet.
	OutgoingChanID lnwire.ShortChannelID

	// Hash is the payment hash of the htlc.
	Hash lntypes.Hash

	// OutgoingExpiry is the absolute block height at which the outgoing
	// htlc expires.
	OutgoingExpiry uint32

	// OutgoingAmount is the amount to forward.
	OutgoingAmount lnwire.MilliSatoshi

	// IncomingExpiry is the absolute block height at which the incoming
	// htlc expires.
	IncomingExpiry uint32

	// IncomingAmount is the amount of the accepted htlc.
	IncomingAmount lnwire.MilliSatoshi

	// CustomRecords are user-defined records in the custom type range that
	// were included in the payload.
	CustomRecords record.CustomSet

	// OnionBlob is the onion packet for the next hop
	OnionBlob [lnwire.OnionPacketSize]byte

	// AutoFailHeight is the block height at which this intercept will be
	// failed back automatically.
	AutoFailHeight int32
}

// InterceptedForward is passed to the ForwardInterceptor for every forwarded
// htlc. It contains all the information about the packet which accordingly
// the interceptor decides if to hold or not.
// In addition this interface allows a later resolution by calling either
// Resume, Settle or Fail.
type InterceptedForward interface {
	// Packet returns the intercepted packet.
	Packet() InterceptedPacket

	// Resume notifies the intention to resume an existing hold forward. This
	// basically means the caller wants to resume with the default behavior for
	// this htlc which usually means forward it.
	Resume() error

	// Settle notifies the intention to settle an existing hold
	// forward with a given preimage.
	Settle(lntypes.Preimage) error

	// Fail notifies the intention to fail an existing hold forward with an
	// encrypted failure reason.
	Fail(reason []byte) error

	// FailWithCode notifies the intention to fail an existing hold forward
	// with the specified failure code.
	FailWithCode(code lnwire.FailCode) error
}

// htlcNotifier is an interface which represents the input side of the
// HtlcNotifier which htlc events are piped through. This interface is intended
// to allow for mocking of the htlcNotifier in tests, so is unexported because
// it is not needed outside of the htlcSwitch package.
type htlcNotifier interface {
	// NotifyForwardingEvent notifies the HtlcNotifier than a htlc has been
	// forwarded.
	NotifyForwardingEvent(key HtlcKey, info HtlcInfo,
		eventType HtlcEventType)

	// NotifyIncomingLinkFailEvent notifies that a htlc has failed on our
	// incoming link. It takes an isReceive bool to differentiate between
	// our node's receives and forwards.
	NotifyLinkFailEvent(key HtlcKey, info HtlcInfo,
		eventType HtlcEventType, linkErr *LinkError, incoming bool)

	// NotifyForwardingFailEvent notifies the HtlcNotifier that a htlc we
	// forwarded has failed down the line.
	NotifyForwardingFailEvent(key HtlcKey, eventType HtlcEventType)

	// NotifySettleEvent notifies the HtlcNotifier that a htlc that we
	// committed to as part of a forward or a receive to our node has been
	// settled.
	NotifySettleEvent(key HtlcKey, preimage lntypes.Preimage,
		eventType HtlcEventType)

	// NotifyFinalHtlcEvent notifies the HtlcNotifier that the final outcome
	// for an htlc has been determined.
	NotifyFinalHtlcEvent(key models.CircuitKey,
		info channeldb.FinalHtlcInfo)
}
