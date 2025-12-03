package idgen

import (
	"fmt"
	"hash/crc32"
	"net"
	"strconv"

	"github.com/yitter/idgenerator-go/idgen"
)

func NextId() string {
	rawId := idgen.NextId()
	// 在前面补0，补到17位
	IdString := fmt.Sprintf("1%017d", rawId)
	check := LuhnCheckDigit(IdString)
	return IdString + strconv.Itoa(check)
}

func init() {
	var options = idgen.NewIdGeneratorOptions(getWorkerIdFromMAC())
	idgen.SetIdGenerator(options)
}

func getWorkerIdFromMAC() uint16 {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range interfaces {
		// 跳过虚拟网卡和禁用的网卡
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		mac := iface.HardwareAddr.String()
		if mac != "" {
			// 使用 CRC32 哈希 MAC 地址
			hash := crc32.ChecksumIEEE([]byte(mac))
			workerId := uint16(hash % 64) // 默认最大 63
			return workerId
		}
	}

	return 1
}

// LuhnCheckDigit 计算Luhn校验位 (Mod 10)
// 输入一个数字字符串，返回其Luhn校验位
func LuhnCheckDigit(number string) int {
	sum := 0
	double := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	// 计算校验位使其和为10的倍数
	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

// VerifyLuhn 验证Luhn算法
// 输入一个包含校验位的数字字符串，返回是否通过验证
func VerifyLuhn(number string) bool {
	sum := 0
	double := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	return sum%10 == 0
}
