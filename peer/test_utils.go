package peer

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"io"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/ltcsuite/lnd/chainntnfs"
	"github.com/ltcsuite/lnd/channeldb"
	"github.com/ltcsuite/lnd/channelnotifier"
	"github.com/ltcsuite/lnd/htlcswitch"
	"github.com/ltcsuite/lnd/input"
	"github.com/ltcsuite/lnd/keychain"
	"github.com/ltcsuite/lnd/lntest/channels"
	"github.com/ltcsuite/lnd/lntest/mock"
	"github.com/ltcsuite/lnd/lnwallet"
	"github.com/ltcsuite/lnd/lnwallet/chainfee"
	"github.com/ltcsuite/lnd/lnwire"
	"github.com/ltcsuite/lnd/netann"
	"github.com/ltcsuite/lnd/queue"
	"github.com/ltcsuite/lnd/shachain"
	"github.com/ltcsuite/ltcd/btcec/v2"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/stretchr/testify/require"
)

const (
	broadcastHeight = 100

	// timeout is a timeout value to use for tests which need to wait for
	// a return value on a channel.
	timeout = time.Second * 5

	// testCltvRejectDelta is the minimum delta between expiry and current
	// height below which htlcs are rejected.
	testCltvRejectDelta = 13
)

var (
	testKeyLoc = keychain.KeyLocator{Family: keychain.KeyFamilyNodeKey}
)

// noUpdate is a function which can be used as a parameter in createTestPeer to
// call the setup code with no custom values on the channels set up.
var noUpdate = func(a, b *channeldb.OpenChannel) {}

// createTestPeer creates a channel between two nodes, and returns a peer for
// one of the nodes, together with the channel seen from both nodes. It takes
// an updateChan function which can be used to modify the default values on
// the channel states for each peer.
func createTestPeer(t *testing.T, notifier chainntnfs.ChainNotifier,
	publTx chan *wire.MsgTx, updateChan func(a, b *channeldb.OpenChannel),
	mockSwitch *mockMessageSwitch) (
	*Brontide, *lnwallet.LightningChannel, error) {

	nodeKeyLocator := keychain.KeyLocator{
		Family: keychain.KeyFamilyNodeKey,
	}
	aliceKeyPriv, aliceKeyPub := btcec.PrivKeyFromBytes(
		channels.AlicesPrivKey,
	)
	aliceKeySigner := keychain.NewPrivKeyMessageSigner(
		aliceKeyPriv, nodeKeyLocator,
	)
	bobKeyPriv, bobKeyPub := btcec.PrivKeyFromBytes(
		channels.BobsPrivKey,
	)

	channelCapacity := ltcutil.Amount(10 * 1e8)
	channelBal := channelCapacity / 2
	aliceDustLimit := ltcutil.Amount(200)
	bobDustLimit := ltcutil.Amount(1300)
	csvTimeoutAlice := uint32(5)
	csvTimeoutBob := uint32(4)
	isAliceInitiator := true

	prevOut := &wire.OutPoint{
		Hash:  channels.TestHdSeed,
		Index: 0,
	}
	fundingTxIn := wire.NewTxIn(prevOut, nil, nil)

	aliceCfg := channeldb.ChannelConfig{
		ChannelConstraints: channeldb.ChannelConstraints{
			DustLimit:        aliceDustLimit,
			MaxPendingAmount: lnwire.MilliSatoshi(rand.Int63()),
			ChanReserve:      ltcutil.Amount(rand.Int63()),
			MinHTLC:          lnwire.MilliSatoshi(rand.Int63()),
			MaxAcceptedHtlcs: uint16(rand.Int31()),
			CsvDelay:         uint16(csvTimeoutAlice),
		},
		MultiSigKey: keychain.KeyDescriptor{
			PubKey: aliceKeyPub,
		},
		RevocationBasePoint: keychain.KeyDescriptor{
			PubKey: aliceKeyPub,
		},
		PaymentBasePoint: keychain.KeyDescriptor{
			PubKey: aliceKeyPub,
		},
		DelayBasePoint: keychain.KeyDescriptor{
			PubKey: aliceKeyPub,
		},
		HtlcBasePoint: keychain.KeyDescriptor{
			PubKey: aliceKeyPub,
		},
	}
	bobCfg := channeldb.ChannelConfig{
		ChannelConstraints: channeldb.ChannelConstraints{
			DustLimit:        bobDustLimit,
			MaxPendingAmount: lnwire.MilliSatoshi(rand.Int63()),
			ChanReserve:      ltcutil.Amount(rand.Int63()),
			MinHTLC:          lnwire.MilliSatoshi(rand.Int63()),
			MaxAcceptedHtlcs: uint16(rand.Int31()),
			CsvDelay:         uint16(csvTimeoutBob),
		},
		MultiSigKey: keychain.KeyDescriptor{
			PubKey: bobKeyPub,
		},
		RevocationBasePoint: keychain.KeyDescriptor{
			PubKey: bobKeyPub,
		},
		PaymentBasePoint: keychain.KeyDescriptor{
			PubKey: bobKeyPub,
		},
		DelayBasePoint: keychain.KeyDescriptor{
			PubKey: bobKeyPub,
		},
		HtlcBasePoint: keychain.KeyDescriptor{
			PubKey: bobKeyPub,
		},
	}

	bobRoot, err := chainhash.NewHash(bobKeyPriv.Serialize())
	if err != nil {
		return nil, nil, err
	}
	bobPreimageProducer := shachain.NewRevocationProducer(*bobRoot)
	bobFirstRevoke, err := bobPreimageProducer.AtIndex(0)
	if err != nil {
		return nil, nil, err
	}
	bobCommitPoint := input.ComputeCommitmentPoint(bobFirstRevoke[:])

	aliceRoot, err := chainhash.NewHash(aliceKeyPriv.Serialize())
	if err != nil {
		return nil, nil, err
	}
	alicePreimageProducer := shachain.NewRevocationProducer(*aliceRoot)
	aliceFirstRevoke, err := alicePreimageProducer.AtIndex(0)
	if err != nil {
		return nil, nil, err
	}
	aliceCommitPoint := input.ComputeCommitmentPoint(aliceFirstRevoke[:])

	aliceCommitTx, bobCommitTx, err := lnwallet.CreateCommitmentTxns(
		channelBal, channelBal, &aliceCfg, &bobCfg, aliceCommitPoint,
		bobCommitPoint, *fundingTxIn, channeldb.SingleFunderTweaklessBit,
		isAliceInitiator, 0,
	)
	if err != nil {
		return nil, nil, err
	}

	dbAlice, err := channeldb.Open(t.TempDir())
	if err != nil {
		return nil, nil, err
	}
	t.Cleanup(func() {
		require.NoError(t, dbAlice.Close())
	})

	dbBob, err := channeldb.Open(t.TempDir())
	if err != nil {
		return nil, nil, err
	}
	t.Cleanup(func() {
		require.NoError(t, dbBob.Close())
	})

	estimator := chainfee.NewStaticEstimator(12500, 0)
	feePerKw, err := estimator.EstimateFeePerKW(1)
	if err != nil {
		return nil, nil, err
	}

	// TODO(roasbeef): need to factor in commit fee?
	aliceCommit := channeldb.ChannelCommitment{
		CommitHeight:  0,
		LocalBalance:  lnwire.NewMSatFromSatoshis(channelBal),
		RemoteBalance: lnwire.NewMSatFromSatoshis(channelBal),
		FeePerKw:      ltcutil.Amount(feePerKw),
		CommitFee:     feePerKw.FeeForWeight(input.CommitWeight),
		CommitTx:      aliceCommitTx,
		CommitSig:     bytes.Repeat([]byte{1}, 71),
	}
	bobCommit := channeldb.ChannelCommitment{
		CommitHeight:  0,
		LocalBalance:  lnwire.NewMSatFromSatoshis(channelBal),
		RemoteBalance: lnwire.NewMSatFromSatoshis(channelBal),
		FeePerKw:      ltcutil.Amount(feePerKw),
		CommitFee:     feePerKw.FeeForWeight(input.CommitWeight),
		CommitTx:      bobCommitTx,
		CommitSig:     bytes.Repeat([]byte{1}, 71),
	}

	var chanIDBytes [8]byte
	if _, err := io.ReadFull(crand.Reader, chanIDBytes[:]); err != nil {
		return nil, nil, err
	}

	shortChanID := lnwire.NewShortChanIDFromInt(
		binary.BigEndian.Uint64(chanIDBytes[:]),
	)

	aliceChannelState := &channeldb.OpenChannel{
		LocalChanCfg:            aliceCfg,
		RemoteChanCfg:           bobCfg,
		IdentityPub:             aliceKeyPub,
		FundingOutpoint:         *prevOut,
		ShortChannelID:          shortChanID,
		ChanType:                channeldb.SingleFunderTweaklessBit,
		IsInitiator:             isAliceInitiator,
		Capacity:                channelCapacity,
		RemoteCurrentRevocation: bobCommitPoint,
		RevocationProducer:      alicePreimageProducer,
		RevocationStore:         shachain.NewRevocationStore(),
		LocalCommitment:         aliceCommit,
		RemoteCommitment:        aliceCommit,
		Db:                      dbAlice.ChannelStateDB(),
		Packager:                channeldb.NewChannelPackager(shortChanID),
		FundingTxn:              channels.TestFundingTx,
	}
	bobChannelState := &channeldb.OpenChannel{
		LocalChanCfg:            bobCfg,
		RemoteChanCfg:           aliceCfg,
		IdentityPub:             bobKeyPub,
		FundingOutpoint:         *prevOut,
		ChanType:                channeldb.SingleFunderTweaklessBit,
		IsInitiator:             !isAliceInitiator,
		Capacity:                channelCapacity,
		RemoteCurrentRevocation: aliceCommitPoint,
		RevocationProducer:      bobPreimageProducer,
		RevocationStore:         shachain.NewRevocationStore(),
		LocalCommitment:         bobCommit,
		RemoteCommitment:        bobCommit,
		Db:                      dbBob.ChannelStateDB(),
		Packager:                channeldb.NewChannelPackager(shortChanID),
	}

	// Set custom values on the channel states.
	updateChan(aliceChannelState, bobChannelState)

	aliceAddr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 18555,
	}

	if err := aliceChannelState.SyncPending(aliceAddr, 0); err != nil {
		return nil, nil, err
	}

	bobAddr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 18556,
	}

	if err := bobChannelState.SyncPending(bobAddr, 0); err != nil {
		return nil, nil, err
	}

	aliceSigner := input.NewMockSigner(
		[]*btcec.PrivateKey{aliceKeyPriv}, nil,
	)
	bobSigner := input.NewMockSigner(
		[]*btcec.PrivateKey{bobKeyPriv}, nil,
	)

	alicePool := lnwallet.NewSigPool(1, aliceSigner)
	channelAlice, err := lnwallet.NewLightningChannel(
		aliceSigner, aliceChannelState, alicePool,
	)
	if err != nil {
		return nil, nil, err
	}
	_ = alicePool.Start()
	t.Cleanup(func() {
		require.NoError(t, alicePool.Stop())
	})

	bobPool := lnwallet.NewSigPool(1, bobSigner)
	channelBob, err := lnwallet.NewLightningChannel(
		bobSigner, bobChannelState, bobPool,
	)
	if err != nil {
		return nil, nil, err
	}
	_ = bobPool.Start()
	t.Cleanup(func() {
		require.NoError(t, bobPool.Stop())
	})

	chainIO := &mock.ChainIO{
		BestHeight: broadcastHeight,
	}
	wallet := &lnwallet.LightningWallet{
		WalletController: &mock.WalletController{
			RootKey:               aliceKeyPriv,
			PublishedTransactions: publTx,
		},
	}

	// If mockSwitch is not set by the caller, set it to the default as the
	// caller does not need to control it.
	if mockSwitch == nil {
		mockSwitch = &mockMessageSwitch{}
	}

	nodeSignerAlice := netann.NewNodeSigner(aliceKeySigner)

	const chanActiveTimeout = time.Minute

	chanStatusMgr, err := netann.NewChanStatusManager(&netann.ChanStatusConfig{
		ChanStatusSampleInterval: 30 * time.Second,
		ChanEnableTimeout:        chanActiveTimeout,
		ChanDisableTimeout:       2 * time.Minute,
		DB:                       dbAlice.ChannelStateDB(),
		Graph:                    dbAlice.ChannelGraph(),
		MessageSigner:            nodeSignerAlice,
		OurPubKey:                aliceKeyPub,
		OurKeyLoc:                testKeyLoc,
		IsChannelActive:          func(lnwire.ChannelID) bool { return true },
		ApplyChannelUpdate: func(*lnwire.ChannelUpdate,
			*wire.OutPoint, bool) error {

			return nil
		},
	})
	if err != nil {
		return nil, nil, err
	}
	if err = chanStatusMgr.Start(); err != nil {
		return nil, nil, err
	}

	errBuffer, err := queue.NewCircularBuffer(ErrorBufferSize)
	if err != nil {
		return nil, nil, err
	}

	var pubKey [33]byte
	copy(pubKey[:], aliceKeyPub.SerializeCompressed())

	cfgAddr := &lnwire.NetAddress{
		IdentityKey: aliceKeyPub,
		Address:     aliceAddr,
		ChainNet:    wire.SimNet,
	}

	interceptableSwitchNotifier := &mock.ChainNotifier{
		EpochChan: make(chan *chainntnfs.BlockEpoch, 1),
	}
	interceptableSwitchNotifier.EpochChan <- &chainntnfs.BlockEpoch{
		Height: 1,
	}

	interceptableSwitch, err := htlcswitch.NewInterceptableSwitch(
		&htlcswitch.InterceptableSwitchConfig{
			CltvRejectDelta:    testCltvRejectDelta,
			CltvInterceptDelta: testCltvRejectDelta + 3,
			Notifier:           interceptableSwitchNotifier,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	// TODO(yy): change ChannelNotifier to be an interface.
	channelNotifier := channelnotifier.New(dbAlice.ChannelStateDB())
	require.NoError(t, channelNotifier.Start())
	t.Cleanup(func() {
		require.NoError(t, channelNotifier.Stop(),
			"stop channel notifier failed")
	})

	cfg := &Config{
		Addr:              cfgAddr,
		PubKeyBytes:       pubKey,
		ErrorBuffer:       errBuffer,
		ChainIO:           chainIO,
		Switch:            mockSwitch,
		ChanActiveTimeout: chanActiveTimeout,
		InterceptSwitch:   interceptableSwitch,
		ChannelDB:         dbAlice.ChannelStateDB(),
		FeeEstimator:      estimator,
		Wallet:            wallet,
		ChainNotifier:     notifier,
		ChanStatusMgr:     chanStatusMgr,
		Features:          lnwire.NewFeatureVector(nil, lnwire.Features),
		DisconnectPeer:    func(b *btcec.PublicKey) error { return nil },
		ChannelNotifier:   channelNotifier,
	}

	alicePeer := NewBrontide(*cfg)
	alicePeer.remoteFeatures = lnwire.NewFeatureVector(nil, lnwire.Features)

	chanID := lnwire.NewChanIDFromOutPoint(channelAlice.ChannelPoint())
	alicePeer.activeChannels.Store(chanID, channelAlice)

	alicePeer.wg.Add(1)
	go alicePeer.channelManager()

	return alicePeer, channelBob, nil
}

// mockMessageSwitch is a mock implementation of the messageSwitch interface
// used for testing without relying on a *htlcswitch.Switch in unit tests.
type mockMessageSwitch struct {
	links []htlcswitch.ChannelUpdateHandler
}

// BestHeight currently returns a dummy value.
func (m *mockMessageSwitch) BestHeight() uint32 {
	return 0
}

// CircuitModifier currently returns a dummy value.
func (m *mockMessageSwitch) CircuitModifier() htlcswitch.CircuitModifier {
	return nil
}

// RemoveLink currently does nothing.
func (m *mockMessageSwitch) RemoveLink(cid lnwire.ChannelID) {}

// CreateAndAddLink currently returns a dummy value.
func (m *mockMessageSwitch) CreateAndAddLink(cfg htlcswitch.ChannelLinkConfig,
	lnChan *lnwallet.LightningChannel) error {

	return nil
}

// GetLinksByInterface returns the active links.
func (m *mockMessageSwitch) GetLinksByInterface(pub [33]byte) (
	[]htlcswitch.ChannelUpdateHandler, error) {

	return m.links, nil
}

// mockUpdateHandler is a mock implementation of the ChannelUpdateHandler
// interface. It is used in mockMessageSwitch's GetLinksByInterface method.
type mockUpdateHandler struct {
	cid lnwire.ChannelID
}

// newMockUpdateHandler creates a new mockUpdateHandler.
func newMockUpdateHandler(cid lnwire.ChannelID) *mockUpdateHandler {
	return &mockUpdateHandler{
		cid: cid,
	}
}

// HandleChannelUpdate currently does nothing.
func (m *mockUpdateHandler) HandleChannelUpdate(msg lnwire.Message) {}

// ChanID returns the mockUpdateHandler's cid.
func (m *mockUpdateHandler) ChanID() lnwire.ChannelID { return m.cid }

// Bandwidth currently returns a dummy value.
func (m *mockUpdateHandler) Bandwidth() lnwire.MilliSatoshi { return 0 }

// EligibleToForward currently returns a dummy value.
func (m *mockUpdateHandler) EligibleToForward() bool { return false }

// MayAddOutgoingHtlc currently returns nil.
func (m *mockUpdateHandler) MayAddOutgoingHtlc(lnwire.MilliSatoshi) error { return nil }

// ShutdownIfChannelClean currently returns nil.
func (m *mockUpdateHandler) ShutdownIfChannelClean() error { return nil }

type mockMessageConn struct {
	t *testing.T

	// MessageConn embeds our interface so that the mock does not need to
	// implement every function. The mock will panic if an unspecified function
	// is called.
	MessageConn

	// writtenMessages is a channel that our mock pushes written messages into.
	writtenMessages chan []byte

	readMessages   chan []byte
	curReadMessage []byte

	// writeRaceDetectingCounter is incremented on any function call
	// associated with writing to the connection. The race detector will
	// trigger on this counter if a data race exists.
	writeRaceDetectingCounter int

	// readRaceDetectingCounter is incremented on any function call
	// associated with reading from the connection. The race detector will
	// trigger on this counter if a data race exists.
	readRaceDetectingCounter int
}

func newMockConn(t *testing.T, expectedMessages int) *mockMessageConn {
	return &mockMessageConn{
		t:               t,
		writtenMessages: make(chan []byte, expectedMessages),
		readMessages:    make(chan []byte, 1),
	}
}

// SetWriteDeadline mocks setting write deadline for our conn.
func (m *mockMessageConn) SetWriteDeadline(time.Time) error {
	m.writeRaceDetectingCounter++
	return nil
}

// Flush mocks a message conn flush.
func (m *mockMessageConn) Flush() (int, error) {
	m.writeRaceDetectingCounter++
	return 0, nil
}

// WriteMessage mocks sending of a message on our connection. It will push
// the bytes sent into the mock's writtenMessages channel.
func (m *mockMessageConn) WriteMessage(msg []byte) error {
	m.writeRaceDetectingCounter++
	select {
	case m.writtenMessages <- msg:
	case <-time.After(timeout):
		m.t.Fatalf("timeout sending message: %v", msg)
	}

	return nil
}

// assertWrite asserts that our mock as had WriteMessage called with the byte
// slice we expect.
func (m *mockMessageConn) assertWrite(expected []byte) {
	select {
	case actual := <-m.writtenMessages:
		require.Equal(m.t, expected, actual)

	case <-time.After(timeout):
		m.t.Fatalf("timeout waiting for write: %v", expected)
	}
}

func (m *mockMessageConn) SetReadDeadline(t time.Time) error {
	m.readRaceDetectingCounter++
	return nil
}

func (m *mockMessageConn) ReadNextHeader() (uint32, error) {
	m.readRaceDetectingCounter++
	m.curReadMessage = <-m.readMessages
	return uint32(len(m.curReadMessage)), nil
}

func (m *mockMessageConn) ReadNextBody(buf []byte) ([]byte, error) {
	m.readRaceDetectingCounter++
	return m.curReadMessage, nil
}

func (m *mockMessageConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockMessageConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockMessageConn) Close() error {
	return nil
}
