package dataDecrypt

const DES_MODE_KEY = 78
const DES_STR_KEY = "be1336de3bf53b68cb2f725869f3b6bc"

func Decrypt(inValue *[]byte, outValue *[]byte) {
	var keyByte = []byte(DES_STR_KEY)
	var temInData []byte = *inValue
	sg_crypt(inValue, outValue, len(temInData), &keyByte, uint32(len(keyByte)), DES_MODE_KEY)
}

func Encrypt(inValue *[]byte, outValue *[]byte) {
	var keyByte = []byte(DES_STR_KEY)
	var temInData []byte = *inValue
	sg_crypt(inValue, outValue, len(temInData), &keyByte, uint32(len(keyByte)), -1*DES_MODE_KEY)
}

func sg_crypt(inValue *[]byte, outValue *[]byte, ValueLen int, Key *[]byte, KeyLen uint32, ckey int) {
	var k int
	k = 0
	var cKeyInt int
	cKeyInt = ckey
	var temInData, tempOutData, KeyData []byte
	temInData = *inValue
	tempOutData = *outValue
	KeyData = *Key

	for v := 0; v < ValueLen; v++ {
		if k == int(KeyLen) {
			k = 0
		}
		var tempdata1 int
		tempdata1 = cKeyInt * (int(KeyData[k]) + k)
		tempOutData[v] = byte(int(temInData[v]) + tempdata1)
		cKeyInt = -1 * cKeyInt
	}
}
