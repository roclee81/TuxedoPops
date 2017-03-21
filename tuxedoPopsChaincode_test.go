package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/protobuf/proto"
	txcache "github.com/skuchain/TuxedoPops/TXCache"
	"github.com/skuchain/TuxedoPops/TuxedoPopsTX"

	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcec"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Notes fromessage Testing popcode
// Public Key: 02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc
// Private Key: 94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20
// Hyperledger address hex 74ded2036e988fc56e3cff77a40c58239591e921
// Hyperledger address Base58: 8sDMfw2Ti7YumfTkbf7RHMgSSSxuAmMFd2GS9wnjkUoX

// Notes fromessage Testing popcode2
// Public Key: 02cb6d65b04c4b84502015f918fe549e95cad4f3b899359a170d4d7d438363c0ce
// Private Key: 60977f22a920c9aa18d58d12cb5e90594152d7aa724bcce21484dfd0f4490b58
// Hyperledger address hex 10734390011641497f489cb475743b8e50d429bb
// Hyperledger address Base58: EHxhLN3Ft4p9jPkR31MJMEMee9G

//Owner1 key
// Public Key: 0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7
// Private Key: 7142c92e6eba38de08980eeb55b8c98bb19f8d417795adb56b6c4d25da6b26c5

// Owner2 key
// Public Key: 02e138b25db2e74c54f8ca1a5cf79e2d1ed6af5bd1904646e7dc08b6d7b0d12bfd
// Private Key: b18b7d3082b3ff9438a7bf9f5f019f8a52fb64647ea879548b3ca7b551eefd65
func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "", args)
	if err != nil {
		fmt.Println("INIT", args, "failed", err)
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInvoke("1", "invoke", args)
	if err != nil {
		fmt.Println("invoke", args, "failed", err)
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes, err := stub.MockQuery("balance", []string{name})

	if err != nil {
		fmt.Println("Query for address (", name, ") failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query for address (", name, ") failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("Query value for address (", name, ") was not", value, "as expected instead", string(bytes))
		t.FailNow()
	}
}

func mint(t *testing.T, stub *shim.MockStub, counterSeed string) {
	createArgs := TuxedoPopsTX.CreateTX{}
	createArgs.Address = "74ded2036e988fc56e3cff77a40c58239591e921"
	createArgs.Amount = 10
	pubKeyBytes, err := hex.DecodeString("03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26")
	if err != nil {
		fmt.Println(err)
	}
	createArgs.CreatorPubKey = pubKeyBytes
	hexCreatorSig := generateCreateSig(counterSeed, 10, "Test Asset", "Test Data", "74ded2036e988fc56e3cff77a40c58239591e921", "7ff1ac3d9dfc56315ee610d0a15609d13c399cf9c92ba2e32e7b1d25ea5c9494")

	createArgs.CreatorSig, err = hex.DecodeString(hexCreatorSig)
	if err != nil {
		fmt.Println(err)
	}
	createArgs.Data = "Test Data"
	createArgs.Type = "Test Asset"
	createArgBytes, err := proto.Marshal(&createArgs)
	createArgBytesStr := hex.EncodeToString(createArgBytes)
	_, err = stub.MockInvoke("3", "create", []string{createArgBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

//altMint takes in a key struct which holds input of private key, public key, and address
func altMint(t *testing.T, stub *shim.MockStub, keys *keyInfo) {
	createArgs := TuxedoPopsTX.CreateTX{}
	createArgs.Address = keys.address
	createArgs.Amount = 10
	createArgs.Data = "Test Data"
	createArgs.Type = "Test Asset"

	pubKeyBytes, err := hex.DecodeString(keys.pubKeyStr)
	if err != nil {
		fmt.Println(err)
	}
	createArgs.CreatorPubKey = pubKeyBytes
	hexCreatorSig := generateCreateSig(keys.counter, 10, "Test Asset", "Test Data", keys.address, keys.privKeyStr)

	createArgs.CreatorSig, err = hex.DecodeString(hexCreatorSig)
	if err != nil {
		fmt.Println(err)
	}

	createArgBytes, err := proto.Marshal(&createArgs)
	createArgBytesStr := hex.EncodeToString(createArgBytes)
	_, err = stub.MockInvoke("3", "create", []string{createArgBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

//altMint takes in a key struct which holds input of private key, public key, and address
func altMint1(t *testing.T, stub *shim.MockStub, user *keyInfo, popcode *keyInfo, data string, createdType string, amount int) {
	createArgs := TuxedoPopsTX.CreateTX{}
	createArgs.Address = popcode.address
	createArgs.Amount = int32(amount)
	createArgs.Data = data
	createArgs.Type = createdType

	creatorPubKeyBytes, err := hex.DecodeString(user.pubKeyStr)
	if err != nil {
		fmt.Println(err)
	}
	createArgs.CreatorPubKey = creatorPubKeyBytes
	popcode.counter, err = getCounter(stub, popcode)
	if err != nil {
		t.Errorf("error getting counterseed in altMint: (%s)", err.Error())
		t.FailNow()
	}
	hexCreatorSig := generateCreateSig(popcode.counter, amount, createdType, data, popcode.address, user.privKeyStr)

	createArgs.CreatorSig, err = hex.DecodeString(hexCreatorSig)
	if err != nil {
		fmt.Println(err)
	}

	createArgBytes, err := proto.Marshal(&createArgs)
	createArgBytesStr := hex.EncodeToString(createArgBytes)
	_, err = stub.MockInvoke("3", "create", []string{createArgBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

func generateCombineSig(counter string, combine TuxedoPopsTX.Combine, amount int, data string, privateKeyStr string) string {
	privKeyByte, _ := hex.DecodeString(privateKeyStr)

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := counter
	message += ":" + combine.Recipe
	for _, source := range combine.GetSources() {
		message += ":" + strconv.FormatInt(int64(source.Idx()), 10)
		message += ":" + strconv.FormatInt(int64(source.Amount()), 10)
	}
	message += ":" + strconv.FormatInt(int64(amount), 10) + ":" + data

	fmt.Printf("\n\ncombine message: (%s)\n\n", message)
	messageBytes := sha256.Sum256([]byte(message))
	fmt.Println(message)

	sig, _ := privKey.Sign(messageBytes[:])

	return hex.EncodeToString(sig.Serialize())
}

func generateRecipeSig(recipeName string, createdType string,
	ingredients []*TuxedoPopsTX.Ingredient, privateKeyStr string) string {

	privKeyByte, _ := hex.DecodeString(privateKeyStr)

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := recipeName + ":" + createdType
	for _, ingredient := range ingredients {
		message += ":" + strconv.FormatInt(int64(ingredient.Numerator), 10) + ":" +
			strconv.FormatInt(int64(ingredient.Denominator), 10) + ":" + ingredient.Type
	}
	fmt.Printf("Signed Message: (%s)\n\n\n", message)
	messageBytes := sha256.Sum256([]byte(message))
	sig, _ := privKey.Sign(messageBytes[:])
	return hex.EncodeToString(sig.Serialize())
}

func registerRecipe(t *testing.T, stub *shim.MockStub) {
	recipeArgs := TuxedoPopsTX.Recipe{}
	recipeArgs.RecipeName = "test recipe"
	recipeArgs.CreatedType = "B"
	recipeArgs.CreatorPubKey, _ = hex.DecodeString("02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
	test := make([]*TuxedoPopsTX.Ingredient, 1)
	test[0] = new(TuxedoPopsTX.Ingredient)
	test[0].Denominator = 1
	test[0].Numerator = 1
	test[0].Type = "Test Asset"

	recipeArgs.Ingredients = test

	sigHex := generateRecipeSig(recipeArgs.RecipeName, recipeArgs.CreatedType,
		recipeArgs.Ingredients, "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")

	var err error
	recipeArgs.CreatorSig, err = hex.DecodeString(sigHex)
	if err != nil {
		fmt.Printf("error decoding creator signature in register recipe. ERR: (%v)", err.Error())
		t.FailNow()
	}
	recipeArgsBytes, err := proto.Marshal(&recipeArgs)
	if err != nil {
		fmt.Printf("error marshalling recipeArgs in registerRecipe. ERR: (%s)\n", err.Error())
		t.FailNow()
	}
	recipeArgsBytesStr := hex.EncodeToString(recipeArgsBytes)
	_, err = stub.MockInvoke("4", "recipe", []string{recipeArgsBytesStr})
	if err != nil {
		fmt.Println(err)
		t.Errorf("error invoking recipe: (%v)", err.Error())
	}
}

func generateCreateSig(CounterSeedStr string, amount int, assetType string, data string, addr string, privateKeyStr string) string {
	privKeyByte, _ := hex.DecodeString(privateKeyStr)

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := CounterSeedStr + ":" + addr + ":" + strconv.FormatInt(int64(amount), 10) + ":" + assetType + ":" + data
	fmt.Println("Signed Message")
	fmt.Println(message)
	messageBytes := sha256.Sum256([]byte(message))
	sig, _ := privKey.Sign(messageBytes[:])
	return hex.EncodeToString(sig.Serialize())
}

func possess(t *testing.T, stub *shim.MockStub, counterSeed string, idx int) {
	transferArgs := TuxedoPopsTX.TransferOwners{}
	transferArgs.Address = "74ded2036e988fc56e3cff77a40c58239591e921"
	transferArgs.Data = "Test possess"
	transferArgs.PopcodePubKey, _ = hex.DecodeString("02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
	ownerBytes, _ := hex.DecodeString("0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7")
	transferArgs.Owners = [][]byte{ownerBytes}
	transferArgs.Output = int32(idx)
	ownerHex := hex.EncodeToString(ownerBytes)
	hexPossessSig := generatePossessSig(counterSeed, idx, "Test possess", ownerHex, "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")
	var err error
	transferArgs.PopcodeSig, err = hex.DecodeString(hexPossessSig)
	transferArgsBytes, _ := proto.Marshal(&transferArgs)
	transferArgsBytesStr := hex.EncodeToString(transferArgsBytes)

	_, err = stub.MockInvoke("4", "transfer", []string{transferArgsBytesStr})
	if err != nil {
		fmt.Println(err)
	}
}

func altPossess(t *testing.T, stub *shim.MockStub, popcode *keyInfo,
	prevOwners []*keyInfo, newOwners []*keyInfo, idx int, data string) {

	transferArgs := TuxedoPopsTX.TransferOwners{}
	transferArgs.Address = popcode.address
	transferArgs.Data = data
	transferArgs.PopcodePubKey, _ = hex.DecodeString(popcode.pubKeyStr)
	ownerBytes, _ := hex.DecodeString(newOwners[0].pubKeyStr)
	transferArgs.Owners = [][]byte{ownerBytes}
	var err error
	popcode.counter, err = getCounter(stub, popcode)
	if err != nil {
		t.Errorf("error getting counter in altPossess: (%v)\n", err.Error())
		t.FailNow()
	}
	ownerHex := hex.EncodeToString(ownerBytes)

	fmt.Printf("\n\nlen prevOwners: (%d)\nprevOwners: (%v)\n", len(prevOwners), prevOwners)
	transferArgs.PrevOwnerSigs = make([][]byte, len(prevOwners))

	for i, owner := range prevOwners {
		if owner != nil {
			hexPrevOwnerSig := generatePossessSig(popcode.counter, idx, data, ownerHex, owner.privKeyStr)
			transferArgs.PrevOwnerSigs[i], err = hex.DecodeString(hexPrevOwnerSig)
			if err != nil {
				t.Errorf("error decoding hexPrevOwnerSig:\ni=(%d)\nowner = (%v)\nerr: (%v)\n", i, owner, err.Error())
				t.FailNow()
			}

		}
	}
	fmt.Printf("\n\nlen prevOwnerSigs: (%d)\nprevOwnerSigs: (%v)\n", len(transferArgs.PrevOwnerSigs), transferArgs.PrevOwnerSigs)

	transferArgs.Output = int32(idx)
	hexPossessSig := generatePossessSig(popcode.counter, idx, data, ownerHex, popcode.privKeyStr)
	transferArgs.PopcodeSig, err = hex.DecodeString(hexPossessSig)
	transferArgsBytes, _ := proto.Marshal(&transferArgs)
	transferArgsBytesStr := hex.EncodeToString(transferArgsBytes)

	_, err = stub.MockInvoke("4", "transfer", []string{transferArgsBytesStr})
	if err != nil {
		t.Errorf("POSSESS ERROR: (%v)", err.Error())
		t.FailNow()
	}
}

func generatePossessSig(CounterSeedStr string, outputIdx int, data string, newOwnersHex string, privateKeyStr string) string {
	privKeyByte, _ := hex.DecodeString(privateKeyStr)

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := CounterSeedStr + ":" + strconv.FormatInt(int64(outputIdx), 10) + ":" + data
	newOwnersTmp, err := hex.DecodeString(newOwnersHex)
	if err != nil {
		fmt.Println(err)
	}
	newOwners := [][]byte{newOwnersTmp}

	for _, newO := range newOwners {
		message += ":"
		message += hex.EncodeToString(newO)
	}
	// fmt.Printf("Signed message %s \n", message)
	mDigest := sha256.Sum256([]byte(message))
	sig, _ := privKey.Sign(mDigest[:])
	return hex.EncodeToString(sig.Serialize())
}

// func generatePossessSig(CounterSeedStr string, outputIdx int, data string, newOwners [][]byte, privateKeyStr string) []byte {
// 	privKeyByte, _ := hex.DecodeString(privateKeyStr)

// 	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

// 	message := CounterSeedStr + ":" + strconv.FormatInt(int64(outputIdx), 10) + ":" + data

// 	for _, newO := range newOwners {
// 		message += ":"
// 		message += hex.EncodeToString(newO)
// 	}
// 	// fmt.Printf("Signed message %s \n", message)
// 	mDigest := sha256.Sum256([]byte(message))
// 	sig, _ := privKey.Sign(mDigest[:])
// 	return sig.Serialize()
// }

func unitize(t *testing.T, stub *shim.MockStub, counterSeed string) {
	unitizeArgs := TuxedoPopsTX.Unitize{}
	unitizeArgs.Data = "Test Unitize"
	unitizeArgs.DestAddress = "10734390011641497f489cb475743b8e50d429bb"
	unitizeArgs.DestAmounts = []int32{10}
	unitizeArgs.SourceAddress = "74ded2036e988fc56e3cff77a40c58239591e921"
	unitizeArgs.SourceOutput = 0
	unitizeArgs.PopcodePubKey, _ = hex.DecodeString("02ca4a8c7dc5090f924cde2264af240d76f6d58a5d2d15c8c5f59d95c70bd9e4dc")
	ownerSig := generateUnitizeSig(counterSeed, unitizeArgs.DestAddress, 0, []int{10}, unitizeArgs.Data, "7142c92e6eba38de08980eeb55b8c98bb19f8d417795adb56b6c4d25da6b26c5")
	unitizeArgs.OwnerSigs = [][]byte{ownerSig}
	unitizeArgs.PopcodeSig = generateUnitizeSig(counterSeed, unitizeArgs.DestAddress, 0, []int{10}, unitizeArgs.Data, "94d7fe7308a452fdf019a0424d9c48ba9b66bdbca565c6fa3b1bf9c646ebac20")
	unitizeArgsBytes, _ := proto.Marshal(&unitizeArgs)
	unitizeArgsBytesStr := hex.EncodeToString(unitizeArgsBytes)

	_, err := stub.MockInvoke("4", "unitize", []string{unitizeArgsBytesStr})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func altUnitize(t *testing.T, stub *shim.MockStub, sourcePopcode *keyInfo, destPopcode *keyInfo, owners []*keyInfo, data string, amounts []int32, output int32) {
	unitizeArgs := TuxedoPopsTX.Unitize{}
	unitizeArgs.Data = data
	unitizeArgs.DestAddress = destPopcode.address
	unitizeArgs.DestAmounts = amounts

	unitizeArgs.SourceAddress = sourcePopcode.address
	unitizeArgs.SourceOutput = output
	unitizeArgs.PopcodePubKey, _ = hex.DecodeString(sourcePopcode.pubKeyStr)
	var err error
	sourcePopcode.counter, err = getCounter(stub, sourcePopcode)
	if err != nil {
		t.Errorf("Error getting counter: (%v)\n", err.Error())
		t.FailNow()
	}

	if len(owners) > 1 {
		t.Errorf("length of owners slice is larger than 1. Length: (%d)\nCurrently this test works with a maximum of one owner\n", len(owners))
		t.FailNow()
	}

	intAmounts := make([]int, len(amounts))
	for i, amount := range amounts {
		intAmounts[i] = int(amount)
	}
	ownerSig := generateUnitizeSig(sourcePopcode.counter, unitizeArgs.DestAddress, int(output), intAmounts, unitizeArgs.Data, owners[0].privKeyStr)
	unitizeArgs.OwnerSigs = [][]byte{ownerSig}
	unitizeArgs.PopcodeSig = generateUnitizeSig(sourcePopcode.counter, unitizeArgs.DestAddress, int(output), intAmounts, unitizeArgs.Data, sourcePopcode.privKeyStr)
	unitizeArgsBytes, _ := proto.Marshal(&unitizeArgs)
	unitizeArgsBytesStr := hex.EncodeToString(unitizeArgsBytes)

	_, err = stub.MockInvoke("4", "unitize", []string{unitizeArgsBytesStr})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func generateUnitizeSig(CounterSeedStr string, destAddr string, outputIdx int, amounts []int, data string, privateKeyStr string) []byte {
	privKeyByte, _ := hex.DecodeString(privateKeyStr)

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyByte)

	message := CounterSeedStr + ":" + destAddr + ":" + data + ":" + strconv.FormatInt(int64(outputIdx), 10)

	for _, amount := range amounts {
		message += ":" + strconv.FormatInt(int64(amount), 10)
	}
	fmt.Printf("\n\nunitize message: (%s)\n\n", message)

	mDigest := sha256.Sum256([]byte(message))
	sig, _ := privKey.Sign(mDigest[:])
	return sig.Serialize()
}

func checkCombine(t *testing.T, stub *shim.MockStub) {
	txCache := txcache.TXCache{}
	txCacheBytes, err := stub.GetState("TxCache")
	if err != nil {
		fmt.Println(err)
	}
	proto.Unmarshal(txCacheBytes, &txCache)

	//create a new set of keys
	keys := new(keyInfo)
	keys.privKeyStr, err = newPrivateKeyString()
	if err != nil {
		fmt.Printf("error generating private key: %v", err)
	}
	keys.pubKeyStr, err = newPubKeyString(keys.privKeyStr)
	if err != nil {
		fmt.Printf("error generating public key: %v", err)
	}
	keys.address = newAddress(keys.pubKeyStr)
	keys.counter, err = getCounter(stub, keys)
	if err != nil {
		t.Errorf("error retrieving counterseed: (%v)", err.Error())
	}

	//mint transaction with keys and counterseed
	altMint(t, stub, keys)

	keys.counter, err = getCounter(stub, keys)
	if err != nil {
		t.Errorf("error retrieving counterseed: (%v)", err.Error())
	}

	// registerRecipe(t, stub)

	//perform combination
	combineArgs := TuxedoPopsTX.Combine{}
	combineArgs.Address = keys.address
	//Sources
	combineArgs.Sources = make([]*TuxedoPopsTX.CombineSources, 1)
	combineArgs.Sources[0] = new(TuxedoPopsTX.CombineSources)
	combineArgs.Sources[0].SourceAmount = 10
	combineArgs.Sources[0].SourceOutput = 0

	combineArgs.Amount = 10
	combineArgs.Recipe = "test recipe"
	combineArgs.Data = "test data"

	creatorPrivKey, _ := newPrivateKeyString()
	creatorPubKeyStr, _ := newPubKeyString(creatorPrivKey)

	combineArgs.CreatorPubKey, _ = hex.DecodeString(creatorPubKeyStr)
	if err != nil {
		fmt.Printf("error generating private key: %v", err.Error())
	}

	combineArgs.CreatorSig, err = hex.DecodeString(generateCombineSig(keys.counter, combineArgs, int(combineArgs.Amount), combineArgs.Data, creatorPrivKey))
	if err != nil {
		fmt.Printf("Error decoding creator sig string in checkCombine. ERR: (%s)", err.Error())
		t.FailNow()
	}

	combineArgs.OwnerSigs = make([][]byte, 0)
	combineArgs.PopcodePubKey, _ = hex.DecodeString(keys.pubKeyStr)
	combineArgs.PopcodeSig, err = hex.DecodeString(generateCombineSig(keys.counter, combineArgs, int(combineArgs.Amount), combineArgs.Data, keys.privKeyStr))
	if err != nil {
		fmt.Printf("Error decoding creator sig string in checkCombine. ERR: (%s)", err.Error())
		t.FailNow()
	}

	combineArgsBytes, _ := proto.Marshal(&combineArgs)
	combineArgsBytesStr := hex.EncodeToString(combineArgsBytes)

	_, err = stub.MockInvoke("4", "combine", []string{combineArgsBytesStr})
	if err != nil {
		fmt.Printf("\nError invoking combine in checkCombine. ERR: (%s)", err.Error())
		t.FailNow()
	}
}

/*
	//To create new private and public keys
	privKeyString, err := newPrivateKeyString()
	if err != nil {
		fmt.Println(err)
	}
	pubKeyString, err := newPubKeyString(privKeyString)
	if err != nil {
		fmt.Println(err)
	}
*/
//generates and returns SHA256 private key string
func newPrivateKeyString() (string, error) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", fmt.Errorf("Error generating private key\n")
	}
	privKeyBytes := privKey.Serialize()
	privKeyString := hex.EncodeToString(privKeyBytes)
	return privKeyString, nil
}

//generates and returns SHA256 public key string fromessage private key string input
func newPubKeyString(privKeyString string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privKeyString)
	if err != nil {
		return "", fmt.Errorf("error decoding private key string (%s)", privKeyString)
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyBytes := pubKey.SerializeCompressed()
	pubkKeyString := hex.EncodeToString(pubKeyBytes)
	return pubkKeyString, nil
}

//generates and returns first forty characters of sha256 hash of public key string
func newAddress(pubKeyStr string) string {
	pubKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		fmt.Printf("error decoding pubkeystring (%s)", pubKeyStr)
	}
	hasher := sha256.New()
	hasher.Write(pubKeyBytes)
	hashedPubKeyBytes := []byte{}
	hashedPubKeyBytes = hasher.Sum(hashedPubKeyBytes)
	hashedPubKeyString := hex.EncodeToString(hashedPubKeyBytes[0:20])
	address := hashedPubKeyString
	return address
}

func generateKeys() (*keyInfo, error) {
	var err error
	keys := new(keyInfo)
	keys.privKeyStr, err = newPrivateKeyString()
	if err != nil {
		fmt.Printf("error generating private key: %v", err.Error())
		return nil, fmt.Errorf("error generating private key: %v", err.Error())
	}
	keys.pubKeyStr, err = newPubKeyString(keys.privKeyStr)
	if err != nil {
		fmt.Printf("error generating public key: %v", err.Error())
		return nil, fmt.Errorf("error generating public key: %v", err.Error())
	}
	keys.address = newAddress(keys.pubKeyStr)
	return keys, nil
}

func getCounter(stub *shim.MockStub, keys *keyInfo) (string, error) {
	//query balance to get counterSeed
	bytes, err := stub.MockQuery("balance", []string{keys.address})
	if err != nil {
		return "", fmt.Errorf("balance query failure on address: (%s)\n", keys.address)
	}
	balanceResult := make(map[string]string)
	json.Unmarshal(bytes, &balanceResult)
	return balanceResult["Counter"], nil
}

/*
	checkCounterSeedChange creates 150 popcodes and checks that the counterseed changes at the appropriate time.
*/
func checkCounterSeedChange(t *testing.T, stub *shim.MockStub) {
	originalCounterseed, err := stub.GetState("CounterSeed")
	if err != nil {
		t.Error("error retrieving counterseed through counterseed query")
	}
	txCache := txcache.TXCache{}
	txCacheBytes, err := stub.GetState("TxCache")
	if err != nil {
		fmt.Println(err)
	}
	proto.Unmarshal(txCacheBytes, &txCache)
	//create up to 150 popcodes
	for i := len(txCache.Cache); i < 150; i++ {
		//create a new set of keys
		keys, err := generateKeys()
		if err != nil {
			t.Errorf("error generating keys: (%v)\n", err.Error())
		}

		keys.counter, err = getCounter(stub, keys)
		if err != nil {
			t.Errorf("error retrieving counterseed: (%v)", err.Error())
			t.FailNow()
		}

		//mint transaction with keys and counterseed
		altMint(t, stub, keys)

		//check counterseed
		counterseed, err := stub.GetState("CounterSeed")
		if err != nil {
			t.Error("error retrieving counterseed through call to getState")
		}

		txCache := txcache.TXCache{}
		txCacheBytes, err := stub.GetState("TxCache")
		if err != nil {
			fmt.Println(err)
		}
		proto.Unmarshal(txCacheBytes, &txCache)

		fmt.Printf("\n\nCOUNTERSEEDSTRING: (%s)\ni: (%d)\nTXCACHELEN: (%d)\n\n\n", hex.EncodeToString(counterseed), i, len(txCache.Cache))

		//check for correct counterSeed value
		if (i < 101) && (hex.EncodeToString(counterseed) != hex.EncodeToString(originalCounterseed)) {
			t.Errorf("\nCounterseed got:\n(%s)\nwant:\n(%s)\n", hex.EncodeToString(counterseed), hex.EncodeToString(originalCounterseed))
			t.FailNow()
		}
		if expected := sha256.Sum256(originalCounterseed); i > 101 && (hex.EncodeToString(counterseed) != hex.EncodeToString(expected[:])) {
			t.Errorf("\nCounterseed got:\n(%s)\nwant:\n(%s)\n", hex.EncodeToString(counterseed), hex.EncodeToString(expected[:]))
			t.FailNow()
		}
	}
}

func generateUsers(stub *shim.MockStub) (*users, error) {
	users := new(users)
	var err error
	users.user1, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user1.counter, err = getCounter(stub, users.user1)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	users.user2, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user2.counter, err = getCounter(stub, users.user2)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	users.user3, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user3.counter, err = getCounter(stub, users.user3)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	users.user4, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user4.counter, err = getCounter(stub, users.user4)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	users.user5, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user5.counter, err = getCounter(stub, users.user5)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	users.user6, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user6.counter, err = getCounter(stub, users.user6)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	users.user7, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generateUsers: (%v)\n", err.Error())
	}
	users.user7.counter, err = getCounter(stub, users.user7)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	return users, nil
}

func generatePopcodes(stub *shim.MockStub) (*popcodes, error) {
	popcodes := new(popcodes)
	var err error
	popcodes.popcode1, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode1.counter, err = getCounter(stub, popcodes.popcode1)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	popcodes.popcode2, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode2.counter, err = getCounter(stub, popcodes.popcode2)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	popcodes.popcode3, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode3.counter, err = getCounter(stub, popcodes.popcode3)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	popcodes.popcode4, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode4.counter, err = getCounter(stub, popcodes.popcode4)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	popcodes.popcode5, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode5.counter, err = getCounter(stub, popcodes.popcode5)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	popcodes.popcode6, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode6.counter, err = getCounter(stub, popcodes.popcode6)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	popcodes.popcode7, err = generateKeys()
	if err != nil {
		return nil, fmt.Errorf("error generating keys in generatePopcodes: (%v)\n", err.Error())
	}
	popcodes.popcode7.counter, err = getCounter(stub, popcodes.popcode7)
	if err != nil {
		return nil, fmt.Errorf("error generating counterSeed in generateUsers: (%v)\n", err.Error())
	}

	return popcodes, nil
}

type keyInfo struct {
	privKeyStr string
	pubKeyStr  string
	address    string
	counter    string
}

type test struct {
	t        *testing.T
	stub     *shim.MockStub
	users    *users
	popcodes *popcodes
}

type users struct {
	user1 *keyInfo
	user2 *keyInfo
	user3 *keyInfo
	user4 *keyInfo
	user5 *keyInfo
	user6 *keyInfo
	user7 *keyInfo
}

type popcodes struct {
	popcode1 *keyInfo
	popcode2 *keyInfo
	popcode3 *keyInfo
	popcode4 *keyInfo
	popcode5 *keyInfo
	popcode6 *keyInfo
	popcode7 *keyInfo
}

type possessInfo struct {
	t          *testing.T
	stub       *shim.MockStub
	popcode    *keyInfo
	prevOwners []*keyInfo
	newOwners  []*keyInfo
	idx        int
	data       string
}

func TestPopcodeChaincode(t *testing.T) {
	bst := new(tuxedoPopsChaincode)
	stub := shim.NewMockStub("tuxedoPops", bst)
	checkInit(t, stub, []string{"Hello World"})

	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43","Outputs":null}`)
	mint(t, stub, "af5eef44907ccdcc33051d035f32f42de0d093fac2fd9d15923448f6af46bc43")
	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
	mint(t, stub, "1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254")
	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"afab4e267a433fe306d1da4608629ce9a280bde98f7004ff883383d65b9f5948","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}","{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"afab4e267a433fe306d1da4608629ce9a280bde98f7004ff883383d65b9f5948\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
	possess(t, stub, "afab4e267a433fe306d1da4608629ce9a280bde98f7004ff883383d65b9f5948", 1)
	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"92c7dff498fbe29d4b8d959a0f519a26ce43844f8871736191e5b62f8f507ea0","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Data\",\"Type\":\"Test Asset\",\"PrevCounter\":\"1adb7c0c1b464fb45860355bf8e711312c608d01202197e58116a424f74af254\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}","{\"Owners\":[\"0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7\"],\"Threshold\":1,\"Data\":\"Test possess\",\"Type\":\"Test Asset\",\"PrevCounter\":\"afab4e267a433fe306d1da4608629ce9a280bde98f7004ff883383d65b9f5948\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
	unitize(t, stub, "92c7dff498fbe29d4b8d959a0f519a26ce43844f8871736191e5b62f8f507ea0")
	checkQuery(t, stub, "74ded2036e988fc56e3cff77a40c58239591e921", `{"Address":"74ded2036e988fc56e3cff77a40c58239591e921","Counter":"92c7dff498fbe29d4b8d959a0f519a26ce43844f8871736191e5b62f8f507ea0","Outputs":["{\"Owners\":[\"0278b76afbefb1e1185bc63ed1a17dd88634e0587491f03e9a8d2d25d9ab289ee7\"],\"Threshold\":1,\"Data\":\"Test possess\",\"Type\":\"Test Asset\",\"PrevCounter\":\"afab4e267a433fe306d1da4608629ce9a280bde98f7004ff883383d65b9f5948\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)
	// checkQuery(t, stub, "10734390011641497f489cb475743b8e50d429bb", `{"Address":"10734390011641497f489cb475743b8e50d429bb","Counter":"f4a49f8c3db70c46bcc4fd2567e17aee0bae43df12ef942cd0e9e0f370d391e5","Outputs":["{\"Owners\":null,\"Threshold\":0,\"Data\":\"Test Unitize\",\"Type\":\"Test Asset\",\"PrevCounter\":\"83b298acdf5d7231597ffb776c8f027877ca89cbafa7675a3f177619b0a9ad74\",\"Creator\":\"03cc7d40833fdf46e05a7f86a6c9cf8a697a129fbae0676ad6bad71f163ea22b26\",\"Amount\":10}"]}`)

	for i := 0; i < 2; i++ {
		checkCounterSeedChange(t, stub)
	}

	registerRecipe(t, stub)
	function := "recipe"
	bytes, err := stub.MockQuery(function, []string{"test recipe"})
	if err != nil {
		fmt.Printf("Query (%s) failed. ERR: %v", function, err.Error())
		t.FailNow()
	}
	if bytes == nil {
		fmt.Printf("Query (%s) failed to get value\n", function)
		t.FailNow()
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(bytes, &jsonMap); err != nil {
		fmt.Printf("error unmarshalling json string %s", bytes)
		t.FailNow()
	}
	fmt.Printf("JSON: %s\n", jsonMap)

	checkCombine(t, stub)

	//new testing suite in progress below:
	testUsers, err := generateUsers(stub)
	if err != nil {
		t.Errorf("error generating users: (%v)\n", err.Error())
		t.FailNow()
	}

	testPopcodes, err := generatePopcodes(stub)
	if err != nil {
		t.Errorf("error generating popcodes: (%v)\n", err.Error())
		t.FailNow()
	}

	var tests = []test{
		{t, stub, testUsers, testPopcodes},
	}

	fmt.Printf("\n\n\nSTARTING NEW TESTING SUITE\n\n\n")
	altMint1(t, stub, tests[0].users.user1, tests[0].popcodes.popcode1, "data", "Water", 100)
	prevOwners := make([]*keyInfo, 1)
	newOwners := make([]*keyInfo, 1)
	newOwners[0] = tests[0].users.user1
	altPossess(t, stub, tests[0].popcodes.popcode1, prevOwners, newOwners, 0, "data")
	prevOwners[0] = newOwners[0]
	newOwners[0] = tests[0].users.user2
	fmt.Printf("\n\n\nprevOwners: (%v)\nnewOwners: (%v)\n\n\n\n", prevOwners[0], newOwners[0])
	altPossess(t, stub, tests[0].popcodes.popcode1, prevOwners, newOwners, 0, "data")
	altUnitize(t, stub, tests[0].popcodes.popcode1, tests[0].popcodes.popcode1, newOwners, "data", []int32{100}, 0)
}

// func altUnitize(t *testing.T, stub *shim.MockStub, sourcePopcode *keyInfo, destPopcode *keyInfo, owners []*keyInfo, data string, amounts []int, output int) {
