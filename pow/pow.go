package pow

import (
	"encoding/binary"
	"fmt"
	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/utils"
	"math/big"
)

type ProofOfWork uint64

func NewProofOfWork(data []byte) (p ProofOfWork) {
	var counter uint64 = 0
	target := constant.ShaHashSize - constant.PowDifficulty
	bigTarget := big.NewInt(1)
	bigTarget.Lsh(bigTarget, uint(target))
	var hashedData utils.SHA256
	for {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, counter)
		dataToHash := append(data, b...)
		hashedData = utils.SHA256Hash(dataToHash)
		hashedInt := big.NewInt(0)
		hashedInt.SetBytes(hashedData[:])

		if bigTarget.Cmp(hashedInt) == 1 {
			break
		}
		counter++
	}
	fmt.Printf("SHA256Hash is : %v and counter is %d \n", hashedData, counter)
	return ProofOfWork(counter)

}

func (pow ProofOfWork) Validator(data []byte) bool {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(pow))
	dataToHash := append(data, b...)
	hashedData := utils.SHA256Hash(dataToHash)
	hashedInt := big.NewInt(0)
	hashedInt.SetBytes(hashedData[:])
	target := constant.ShaHashSize - constant.PowDifficulty
	bigTarget := big.NewInt(1)
	bigTarget.Lsh(bigTarget, uint(target))
	return bigTarget.Cmp(hashedInt) == 1

}
