package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDIDDocControllerChanged(t *testing.T) {
	setup := Setup()

	//Init did
	aliceKeys, aliceDid, _ := setup.InitDid("did:cheqd:test:alice")
	bobKeys, _, _ := setup.InitDid("did:cheqd:test:bob")

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.Controller = append(updatedDidDoc.Controller, "did:cheqd:test:bob")
	receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, ConcatKeys(aliceKeys, bobKeys))

	// check
	require.NotEqual(t, aliceDid.Controller, receivedDid.Controller)
	require.NotEqual(t, []string{"did:cheqd:test:alice", "did:cheqd:test:bob"}, receivedDid.Controller)
	require.Equal(t, []string{"did:cheqd:test:bob"}, receivedDid.Controller)
}

func TestDIDDocVerificationMethodChangedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, aliceDid, _ := setup.InitDid("did:cheqd:test:alice")
	bobKeys, _, _ := setup.InitDid("did:cheqd:test:bob")

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod[0].Type = "new"
	_, err := setup.SendUpdateDid(updatedDidDoc, bobKeys)

	// check
	require.Error(t, err)
	require.Equal(t, "signature did:cheqd:test:alice not found: invalid signature detected", err.Error())
}

func TestDIDDocVerificationMethodControllerChangedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, aliceDid, _ := setup.InitDid("did:cheqd:test:alice")
	bobKeys, _, _ := setup.InitDid("did:cheqd:test:bob")

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod[0].Controller = "did:cheqd:test:bob"
	_, err := setup.SendUpdateDid(updatedDidDoc, bobKeys)

	// check
	require.Error(t, err)
	require.Equal(t, "signature did:cheqd:test:alice not found: invalid signature detected", err.Error())
}

func TestDIDDocControllerChangedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, aliceDid, _ := setup.InitDid("did:cheqd:test:alice")
	bobKeys, _, _ := setup.InitDid("did:cheqd:test:bob")

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.Controller = append(updatedDidDoc.Controller, "did:cheqd:test:bob")
	_, err := setup.SendUpdateDid(updatedDidDoc, bobKeys)

	// check
	require.Error(t, err)
	require.Equal(t, "signature did:cheqd:test:alice not found: invalid signature detected", err.Error())
}

func TestDIDDocVerificationMethodDeletedWithoutOldSignature(t *testing.T) {
	setup := Setup()

	//Init did
	_, bodDidDoc, _ := setup.InitDid("did:cheqd:test:bob")

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:alice")

	aliceDid.VerificationMethod = append(aliceDid.VerificationMethod, &types.VerificationMethod{
		Id:                 "did:cheqd:test:alice#key-2",
		Controller:         "did:cheqd:test:bob",
		PublicKeyMultibase: bodDidDoc.VerificationMethod[0].PublicKeyMultibase,
	})

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:alice#key-1": privKey}
	_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
	_, err := setup.SendUpdateDid(updatedDidDoc, aliceKeys)

	// check
	require.Error(t, err)
	require.Equal(t, "signature did:cheqd:test:bob not found: invalid signature detected", err.Error())
}

func TestDIDDocVerificationMethodDeleted(t *testing.T) {
	setup := Setup()

	//Init did
	bobKeys, bodDidDoc, _ := setup.InitDid("did:cheqd:test:bob")

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	aliceDid := setup.CreateDid(pubKey, "did:cheqd:test:alice")

	aliceDid.VerificationMethod = append(aliceDid.VerificationMethod, &types.VerificationMethod{
		Id:                 "did:cheqd:test:alice#key-2",
		Controller:         "did:cheqd:test:bob",
		PublicKeyMultibase: bodDidDoc.VerificationMethod[0].PublicKeyMultibase,
	})

	aliceKeys := map[string]ed25519.PrivateKey{"did:cheqd:test:alice#key-1": privKey}
	_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

	updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
	updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
	receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, ConcatKeys(aliceKeys, bobKeys))

	// check
	require.NotEqual(t, len(aliceDid.VerificationMethod), len(receivedDid.VerificationMethod))
	require.True(t, reflect.DeepEqual(aliceDid.VerificationMethod[0], receivedDid.VerificationMethod[0]))
}

/*func TestHandler_UpdateDidInvalidSignature(t *testing.T) {
	setup := Setup()

	_, did, _ := setup.InitDid("did:cheqd:test:alice")

	// query Did
	receivedDid, _, _ := setup.Keeper.GetDid(&setup.Ctx, did.Id)

	//Init priv key
	newPubKey, newPrivKey, _ := ed25519.GenerateKey(rand.Reader)

	// add new Did
	didMsgUpdate := setup.UpdateDid(receivedDid, newPubKey)
	dataUpdate, _ := ptypes.NewAnyWithValue(didMsgUpdate)
	_, err := setup.Handler(setup.Ctx, setup.WrapRequest(newPrivKey, dataUpdate, make(map[string]string)))
	require.Error(t, err)
	require.Equal(t, "did:cheqd:test:alice: invalid signature detected", err.Error())
}*/
