package oui

import (
	"strings"
)

// OUIDatabase provides MAC address to manufacturer lookup.
type OUIDatabase struct {
	// embedded map of OUI prefix to manufacturer name
	data map[string]string
}

// NewOUIDatabase creates a new OUI database with embedded data.
func NewOUIDatabase() *OUIDatabase {
	return &OUIDatabase{
		data: ouiData,
	}
}

// LookupManufacturer returns the manufacturer name for a given MAC address.
// Returns empty string if not found.
func (db *OUIDatabase) LookupManufacturer(mac string) string {
	// Normalize MAC and extract OUI prefix (first 3 octets)
	mac = strings.ToUpper(mac)
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ".", "")

	if len(mac) < 6 {
		return ""
	}

	// Try full 6-character OUI first
	oui := mac[:6]
	if manufacturer, ok := db.data[oui]; ok {
		return manufacturer
	}

	return ""
}

// Common OUI prefixes - a curated list of common manufacturers.
var ouiData = map[string]string{
	// Apple
	"000A27": "Apple", "000A95": "Apple", "000D93": "Apple", "001124": "Apple",
	"001451": "Apple", "0016CB": "Apple", "0017F2": "Apple", "0019E3": "Apple",
	"001B63": "Apple", "001CB3": "Apple", "001D4F": "Apple", "001E52": "Apple",
	"001EC2": "Apple", "001F5B": "Apple", "001FF3": "Apple", "002241": "Apple",
	"002312": "Apple", "002332": "Apple", "002436": "Apple", "00254B": "Apple",
	"002608": "Apple", "0026B0": "Apple", "0026BB": "Apple", "003065": "Apple",
	"003EE1": "Apple", "0050E4": "Apple", "006171": "Apple", "00A040": "Apple",
	"00C610": "Apple", "00CDFE": "Apple", "00F4B9": "Apple", "00F76F": "Apple",
	"041552": "Apple", "04489A": "Apple", "045453": "Apple", "046C59": "Apple",
	"04D3CF": "Apple", "04DB56": "Apple", "04E536": "Apple", "04F13E": "Apple",
	"04F7E4": "Apple", "086698": "Apple", "086D41": "Apple", "087045": "Apple",
	"0C4DE9": "Apple", "0C74C2": "Apple", "0C771A": "Apple", "0CBC9F": "Apple",
	"10417F": "Apple", "1094BB": "Apple", "10DDB1": "Apple", "14109F": "Apple",
	"148FC6": "Apple", "14BD61": "Apple", "185936": "Apple", "189EFC": "Apple",
	"18AF61": "Apple", "18AF8F": "Apple", "18E7F4": "Apple", "18EE69": "Apple",
	"1C1AC0": "Apple", "1C36BB": "Apple", "1C9148": "Apple", "1C9E46": "Apple",
	"20768F": "Apple", "209BCD": "Apple", "20A2E4": "Apple", "20AB37": "Apple",
	"20C9D0": "Apple", "24240E": "Apple", "24A074": "Apple", "24A2E1": "Apple",
	"24AB81": "Apple", "24E314": "Apple", "28CFDA": "Apple", "28E02C": "Apple",
	"28E14C": "Apple", "28E7CF": "Apple", "2C200B": "Apple", "2C3361": "Apple",
	"2C61F6": "Apple", "2CBE08": "Apple", "3010E4": "Apple", "3090AB": "Apple",
	"30F7C5": "Apple", "341298": "Apple", "3451C9": "Apple", "34C059": "Apple",
	"38484C": "Apple", "38C986": "Apple", "3C0754": "Apple", "3C15C2": "Apple",
	"3C2EF9": "Apple", "3CD0F8": "Apple", "403004": "Apple", "40331A": "Apple",
	"40A6D9": "Apple", "40B395": "Apple", "40CBC0": "Apple", "40D32D": "Apple",
	"442A60": "Apple", "44D884": "Apple", "48437C": "Apple", "484BAA": "Apple",
	"48746E": "Apple", "48A195": "Apple", "48D705": "Apple", "48E9F1": "Apple",
	"4C32FF": "Apple", "4C57CA": "Apple", "4C7C5F": "Apple", "4C8D79": "Apple",
	"500115": "Apple", "5430F9": "Apple", "54724F": "Apple", "54AE27": "Apple",
	"54E43A": "Apple", "54EA3A": "Apple", "587F57": "Apple", "5855CA": "Apple",
	"58B035": "Apple", "5C5948": "Apple", "5C8D4E": "Apple", "5C969D": "Apple",
	"5C97F3": "Apple", "60C547": "Apple", "60D9C7": "Apple", "60F445": "Apple",
	"60F81D": "Apple", "60FACD": "Apple", "60FEC5": "Apple", "64200C": "Apple",
	"6476BA": "Apple", "649ABE": "Apple", "64A3CB": "Apple", "64B0A6": "Apple",
	"64E682": "Apple", "68644B": "Apple", "68967B": "Apple", "68A86D": "Apple",
	"68D93C": "Apple", "68DBCA": "Apple", "68FEF7": "Apple", "6C19C0": "Apple",
	"6C4008": "Apple", "6C709F": "Apple", "6C94F8": "Apple", "6CC26B": "Apple",
	"700902": "Apple", "7014A6": "Apple", "70480F": "Apple", "70565C": "Apple",
	"70A2B3": "Apple", "70CD60": "Apple", "70DEE2": "Apple", "70ECE4": "Apple",
	"70F087": "Apple", "7831C1": "Apple", "786C1C": "Apple", "789F70": "Apple",
	"78A3E4": "Apple", "78CA39": "Apple", "78D75F": "Apple", "78FD94": "Apple",
	"7C04D0": "Apple", "7C11BE": "Apple", "7C5049": "Apple", "7C6D62": "Apple",
	"7CC3A1": "Apple", "7CC537": "Apple", "7CD1C3": "Apple", "7CF05F": "Apple",
	"7CFADF": "Apple", "803C69": "Apple", "80929F": "Apple", "8463D6": "Apple",
	"848506": "Apple", "84788B": "Apple", "848E0C": "Apple", "84A134": "Apple",
	"84B153": "Apple", "84FCAC": "Apple", "84FCFE": "Apple", "8866A5": "Apple",
	"88C663": "Apple", "88CB87": "Apple", "88E87F": "Apple", "8C2DAA": "Apple",
	"8C5877": "Apple", "8C7B9D": "Apple", "8C7C92": "Apple", "8C8590": "Apple",
	"8C8EF2": "Apple", "8CFABA": "Apple", "90272B": "Apple", "903C92": "Apple",
	"9060F1": "Apple", "908D6C": "Apple", "90B0ED": "Apple", "90B21F": "Apple",
	"90B931": "Apple", "90FD61": "Apple", "949426": "Apple", "94E96A": "Apple",
	"94F6A3": "Apple", "9801A7": "Apple", "9803D8": "Apple", "98B8E3": "Apple",
	"98D6BB": "Apple", "98E0D9": "Apple", "98F0AB": "Apple", "98FE94": "Apple",
	"9C04EB": "Apple", "9C207B": "Apple", "9C293F": "Apple", "9C35EB": "Apple",
	"9CF387": "Apple", "9CFC01": "Apple", "A03BE3": "Apple", "A0999B": "Apple",
	"A0EDCD": "Apple", "A43135": "Apple", "A45E60": "Apple", "A46706": "Apple",
	"A4B197": "Apple", "A4C361": "Apple", "A4D18C": "Apple", "A4D1D2": "Apple",
	"A82066": "Apple", "A85B78": "Apple", "A860B6": "Apple", "A886DD": "Apple",
	"A88808": "Apple", "A8968A": "Apple", "A8BBCF": "Apple", "A8FAD8": "Apple",
	"AC293A": "Apple", "AC3C0B": "Apple", "AC61EA": "Apple", "AC7F3E": "Apple",
	"AC87A3": "Apple", "ACBC32": "Apple", "ACFDEC": "Apple", "B03495": "Apple",
	"B065BD": "Apple", "B0702D": "Apple", "B09FBA": "Apple", "B0CA68": "Apple",
	"B418D1": "Apple", "B48B19": "Apple", "B4F0AB": "Apple", "B4F1DA": "Apple",
	"B8098A": "Apple", "B817C2": "Apple", "B841A4": "Apple", "B844D9": "Apple",
	"B88D12": "Apple", "B8C75D": "Apple", "B8E856": "Apple", "B8F6B1": "Apple",
	"B8FF61": "Apple", "BC3BAF": "Apple", "BC4CC4": "Apple", "BC5436": "Apple",
	"BC6778": "Apple", "BC9FEF": "Apple", "BCA920": "Apple", "BCEC5D": "Apple",
	"BCF5AC": "Apple", "C01ADA": "Apple", "C06394": "Apple", "C0847A": "Apple",
	"C0CECD": "Apple", "C0D012": "Apple", "C0F2FB": "Apple", "C42C03": "Apple",
	"C4B301": "Apple", "C81EE7": "Apple", "C82A14": "Apple", "C869CD": "Apple",
	"C86F1D": "Apple", "C8B5B7": "Apple", "C8BCC8": "Apple", "C8E0EB": "Apple",
	"C8F650": "Apple", "CC088D": "Apple", "CC20E8": "Apple", "CC25EF": "Apple",
	"CC29F5": "Apple", "CC4463": "Apple", "CC785F": "Apple", "CCC760": "Apple",
	"D0034B": "Apple", "D023DB": "Apple", "D02598": "Apple", "D03311": "Apple",
	"D0A637": "Apple", "D0C5F3": "Apple", "D0D2B0": "Apple", "D0E140": "Apple",
	"D4619D": "Apple", "D49A20": "Apple", "D4DCCD": "Apple", "D4F46F": "Apple",
	"D83062": "Apple", "D89695": "Apple", "D89E3F": "Apple", "D8A25E": "Apple",
	"D8BB2C": "Apple", "D8CF9C": "Apple", "D8D1CB": "Apple", "DC0C5C": "Apple",
	"DC2B2A": "Apple", "DC2B61": "Apple", "DC3714": "Apple", "DC415F": "Apple",
	"DC56E7": "Apple", "DC86D8": "Apple", "DC9B9C": "Apple", "DCA4CA": "Apple",
	"DCA904": "Apple", "E05F45": "Apple", "E0B52D": "Apple", "E0C767": "Apple",
	"E0C97A": "Apple", "E0F5C6": "Apple", "E0F847": "Apple", "E425E7": "Apple",
	"E49A79": "Apple", "E4C63D": "Apple", "E4CE8F": "Apple", "E4E0A6": "Apple",
	"E80688": "Apple", "E88D28": "Apple", "E8B2AC": "Apple", "EC852F": "Apple",
	"F02475": "Apple", "F04F7C": "Apple", "F0766F": "Apple", "F07960": "Apple",
	"F099BF": "Apple", "F0B479": "Apple", "F0C1F1": "Apple", "F0CBD1": "Apple",
	"F0D1A9": "Apple", "F0DB30": "Apple", "F0DBE2": "Apple", "F0DBF8": "Apple",
	"F0DCE2": "Apple", "F0F61C": "Apple", "F40F24": "Apple", "F41BA1": "Apple",
	"F437B7": "Apple", "F45C89": "Apple", "F4F15A": "Apple", "F4F951": "Apple",
	"F80377": "Apple", "F81EDF": "Apple", "F82793": "Apple", "F86214": "Apple",
	"FC253F": "Apple", "FCD848": "Apple", "FCE998": "Apple", "FCFC48": "Apple",

	// Samsung
	"00000F": "Samsung", "0000F0": "Samsung", "001247": "Samsung", "001377": "Samsung",
	"001599": "Samsung", "001632": "Samsung", "00166B": "Samsung", "00166C": "Samsung",
	"001785": "Samsung", "0017C9": "Samsung", "0017D5": "Samsung", "001A8A": "Samsung",
	"001B98": "Samsung", "001D25": "Samsung", "001DF6": "Samsung", "001E7D": "Samsung",
	"001EE1": "Samsung", "001EE2": "Samsung", "001FCC": "Samsung", "001FCD": "Samsung",
	"00214C": "Samsung", "002339": "Samsung", "0024E9": "Samsung", "002454": "Samsung",
	"0024D6": "Samsung", "002566": "Samsung", "00265D": "Samsung", "00265F": "Samsung",
	"002690": "Samsung", "0026BC": "Samsung", "0034DA": "Samsung", "0012FB": "Samsung",
	"6C2F2C": "Samsung", "8C71F8": "Samsung", "CC07AB": "Samsung", "F0D5BF": "Samsung",
	"5C3C27": "Samsung", "C45006": "Samsung",

	// Google
	"001A11": "Google", "3C5AB4": "Google", "54608E": "Google", "94EB2C": "Google",
	"F4F5D8": "Google", "F4F5E8": "Google", "3C2EFF": "Google", "A4C639": "Google",

	// Microsoft
	"001DD8": "Microsoft", "002481": "Microsoft", "0025AE": "Microsoft",
	"0050F2": "Microsoft", "28186D": "Microsoft", "485073": "Microsoft",
	"5882A8": "Microsoft", "60458E": "Microsoft", "7CED8D": "Microsoft",
	"B4AE2B": "Microsoft", "C83DD4": "Microsoft", "D48FAA": "Microsoft",
	"DC536C": "Microsoft",

	// Amazon
	"00FC8B": "Amazon", "0C47C9": "Amazon", "34D270": "Amazon", "40B4CD": "Amazon",
	"44650D": "Amazon", "50F5DA": "Amazon", "68378E": "Amazon", "68548E": "Amazon",
	"6C567E": "Amazon", "747548": "Amazon", "78E103": "Amazon", "84D6D0": "Amazon",
	"A002DC": "Amazon", "AC63BE": "Amazon", "B47C9C": "Amazon", "FC65DE": "Amazon",
	"FCA667": "Amazon",

	// Intel
	"001111": "Intel", "001302": "Intel", "0013CE": "Intel", "001517": "Intel",
	"0016EA": "Intel", "0016EB": "Intel", "001820": "Intel", "001E64": "Intel",
	"001E65": "Intel", "001F3B": "Intel", "001F3C": "Intel", "002170": "Intel",
	"0021D8": "Intel", "002263": "Intel", "002314": "Intel", "002564": "Intel",
	"00270E": "Intel", "086266": "Intel", "3413E8": "Intel", "3497F6": "Intel",
	"48F17F": "Intel", "5C5181": "Intel", "5C879C": "Intel", "647002": "Intel",
	"68A3C4": "Intel", "68D025": "Intel", "787B8A": "Intel", "78929C": "Intel",
	"78FF57": "Intel", "8086F2": "Intel", "8C700A": "Intel", "8C70AB": "Intel",
	"A0369F": "Intel", "A08CFA": "Intel", "A4C494": "Intel", "B4B675": "Intel",
	"B8763F": "Intel", "C0B6F9": "Intel", "D04E99": "Intel", "E03676": "Intel",
	"E4A471": "Intel", "F81654": "Intel",

	// Dell
	"0006C9": "Dell", "000874": "Dell", "000B3E": "Dell", "000BDB": "Dell",
	"000C76": "Dell", "000D56": "Dell", "000F1F": "Dell", "001143": "Dell",
	"0014E9": "Dell", "001372": "Dell", "0012C9": "Dell", "0021D0": "Dell",
	"0C29EF": "Dell", "1458D0": "Dell", "149167": "Dell", "180373": "Dell",
	"204747": "Dell", "2468A3": "Dell", "28C7CE": "Dell", "34E6AD": "Dell",
	"448A5B": "Dell", "509A4C": "Dell", "549F13": "Dell", "5C260A": "Dell",
	"646E69": "Dell", "7845C4": "Dell", "844765": "Dell", "8C16E2": "Dell",
	"90B11C": "Dell", "984BE1": "Dell", "A41F72": "Dell", "A4BADB": "Dell",
	"B08350": "Dell", "B4E10F": "Dell", "B8AC6F": "Dell", "BC305B": "Dell",
	"C80AA9": "Dell", "D4AE52": "Dell", "D4BE26": "Dell", "F04DA2": "Dell",
	"F48C50": "Dell", "F8B156": "Dell", "F8BC12": "Dell", "F8DB88": "Dell",

	// Lenovo
	"006059": "Lenovo", "00098C": "Lenovo", "000A5E": "Lenovo", "00D0C9": "Lenovo",
	"001E4F": "Lenovo", "001E68": "Lenovo", "002181": "Lenovo", "002616": "Lenovo",
	"0026F2": "Lenovo", "282CBF": "Lenovo", "5CF3FC": "Lenovo", "70F1A1": "Lenovo",
	"70720D": "Lenovo", "8C64C2": "Lenovo", "985D82": "Lenovo", "9C216A": "Lenovo",

	// HP
	"0001E6": "HP", "0001E7": "HP", "000802": "HP", "0008C7": "HP",
	"000A57": "HP", "000BCD": "HP", "000D9D": "HP", "000E7F": "HP",
	"000F20": "HP", "000F61": "HP", "001083": "HP", "0010E3": "HP",
	"001185": "HP", "00110A": "HP", "001279": "HP", "001321": "HP",
	"0014C2": "HP", "001560": "HP", "001635": "HP", "0017A4": "HP",
	"00188B": "HP", "0019BB": "HP", "001A4B": "HP", "001B78": "HP",
	"001CC4": "HP", "001E0B": "HP", "001F29": "HP", "00215A": "HP",
	"00215C": "HP", "002264": "HP", "002376": "HP", "002398": "HP",
	"00248C": "HP", "002562": "HP", "0026F1": "HP", "00306E": "HP",
	"1C98EC": "HP", "28924A": "HP", "2C27D7": "HP", "308D99": "HP",
	"38EAA7": "HP", "3CA82A": "HP", "6CC217": "HP", "70106F": "HP",
	"803F5D": "HP", "8851FB": "HP", "94659C": "HP", "986B3D": "HP",
	"98E7F4": "HP", "9CB654": "HP", "A02BB8": "HP", "A0D3C1": "HP",
	"A8E018": "HP", "B0AA77": "HP", "B4B52F": "HP", "B8AF67": "HP",
	"C03D9F": "HP", "D4C947": "HP", "D8DF7A": "HP", "E84DD0": "HP",
	"F4034F": "HP",

	// Cisco
	"0000B0": "Cisco", "0000BC": "Cisco", "00018A": "Cisco", "00018B": "Cisco",
	"000192": "Cisco", "000193": "Cisco", "000194": "Cisco", "0001C7": "Cisco",
	"0001C9": "Cisco", "000195": "Cisco", "000196": "Cisco", "000197": "Cisco",

	// TP-Link
	"001D0F": "TP-Link", "002127": "TP-Link", "00279E": "TP-Link", "1060E2": "TP-Link",
	"1C3BF3": "TP-Link", "3C4E47": "TP-Link", "3C52A1": "TP-Link", "50FA84": "TP-Link",
	"54C80F": "TP-Link", "5C628B": "TP-Link", "5C899A": "TP-Link", "600194": "TP-Link",
	"6466B3": "TP-Link", "68FF7B": "TP-Link", "741865": "TP-Link", "78A106": "TP-Link",
	"8C210A": "TP-Link", "94D9B3": "TP-Link", "98DED0": "TP-Link", "A8574E": "TP-Link",
	"AC84C6": "TP-Link", "B0A7B9": "TP-Link", "C04A00": "TP-Link", "C0E42D": "TP-Link",
	"CC3459": "TP-Link", "D46E0E": "TP-Link", "D80D17": "TP-Link", "E006E6": "TP-Link",
	"E894F6": "TP-Link", "EC086B": "TP-Link", "EC888F": "TP-Link", "F42E7F": "TP-Link",
	"F4F26D": "TP-Link", "F8D111": "TP-Link",

	// Netgear
	"000FB5": "Netgear", "00146C": "Netgear", "00184D": "Netgear", "001B2F": "Netgear",
	"001E2A": "Netgear", "001F33": "Netgear", "00223F": "Netgear", "002438": "Netgear",
	"002636": "Netgear", "204E7F": "Netgear", "28C68E": "Netgear", "2CB05D": "Netgear",
	"30469A": "Netgear", "3894ED": "Netgear", "3C3786": "Netgear", "445561": "Netgear",
	"44A56E": "Netgear", "4C09B4": "Netgear", "504A6E": "Netgear", "6038E0": "Netgear",
	"6CB0CE": "Netgear", "744401": "Netgear", "84002D": "Netgear", "9CD36D": "Netgear",
	"A00460": "Netgear", "A021B7": "Netgear", "A042E1": "Netgear", "A4F4C2": "Netgear",
	"B03956": "Netgear", "B07FB9": "Netgear", "B0B980": "Netgear", "C03F0E": "Netgear",
	"C0FFD4": "Netgear", "CC40D0": "Netgear", "E0469A": "Netgear", "E091F5": "Netgear",
	"E4F4C6": "Netgear", "F87394": "Netgear",

	// Asus
	"001731": "Asus", "001A92": "Asus", "001BFC": "Asus", "001D60": "Asus",
	"001E8C": "Asus", "001FC6": "Asus", "002354": "Asus", "002618": "Asus",
	"049226": "Asus", "0C9D92": "Asus", "107B44": "Asus", "10BF48": "Asus",
	"10C37B": "Asus", "14DDA9": "Asus", "1831BF": "Asus", "1C872C": "Asus",
	"2C4D54": "Asus", "2C56DC": "Asus", "305A3A": "Asus", "3085A9": "Asus",
	"38D547": "Asus", "40167E": "Asus", "485B39": "Asus", "50465D": "Asus",
	"54A050": "Asus", "60A44C": "Asus", "6045CB": "Asus", "708BCD": "Asus",
	"74D02B": "Asus", "7824AF": "Asus", "90E6BA": "Asus", "A036BC": "Asus",
	"AC220B": "Asus", "AC9E17": "Asus", "B06EBF": "Asus", "BC5FF4": "Asus",
	"BCEE7B": "Asus", "C86000": "Asus", "D017C2": "Asus", "D45D64": "Asus",
	"D850E6": "Asus", "E03F49": "Asus", "E8CC18": "Asus", "F07959": "Asus",
	"F42853": "Asus", "F46D04": "Asus", "FC345B": "Asus",

	// Raspberry Pi
	"B827EB": "Raspberry Pi", "DC442E": "Raspberry Pi", "DCA632": "Raspberry Pi",
	"E45F01": "Raspberry Pi",

	// Xiaomi
	"100B41": "Xiaomi", "14F65A": "Xiaomi", "286C07": "Xiaomi", "341532": "Xiaomi",
	"58445F": "Xiaomi", "5C99BF": "Xiaomi", "6402CB": "Xiaomi", "680AE2": "Xiaomi",
	"6800B0": "Xiaomi", "78D6DC": "Xiaomi", "840DB6": "Xiaomi", "88C397": "Xiaomi",
	"8C790A": "Xiaomi", "9C99A0": "Xiaomi", "9C9933": "Xiaomi", "9CD2D2": "Xiaomi",
	"A88E24": "Xiaomi", "AC150E": "Xiaomi", "B023BD": "Xiaomi", "B0E235": "Xiaomi",
	"C40B04": "Xiaomi", "C40B4A": "Xiaomi", "D0E782": "Xiaomi", "D4970B": "Xiaomi",
	"E85A09": "Xiaomi", "EC3586": "Xiaomi", "F0B429": "Xiaomi", "FC64BA": "Xiaomi",

	// Huawei
	"001E10": "Huawei", "002568": "Huawei", "00259E": "Huawei", "0025F0": "Huawei",
	"002729": "Huawei", "0034FE": "Huawei", "04021F": "Huawei", "045F7E": "Huawei",
	"04B0E7": "Huawei", "0811FC": "Huawei", "087A8C": "Huawei", "08192D": "Huawei",
	"08637C": "Huawei", "0C37DC": "Huawei", "0CCAD7": "Huawei", "100047": "Huawei",
	"10B1F8": "Huawei", "10C61F": "Huawei", "1400E1": "Huawei", "14B968": "Huawei",
	"20D4BE": "Huawei", "2469A5": "Huawei", "286ED4": "Huawei", "28A6DB": "Huawei",
	"2C55D3": "Huawei", "2CABEB": "Huawei", "306023": "Huawei", "3400A3": "Huawei",
	"34CDBE": "Huawei", "38BC01": "Huawei", "3C6200": "Huawei", "3CFBFE": "Huawei",
	"406C8F": "Huawei", "4C1FCC": "Huawei", "4CB16C": "Huawei", "50A72B": "Huawei",
	"5425EA": "Huawei", "54B121": "Huawei", "581F28": "Huawei", "582AF7": "Huawei",
	"5CB43E": "Huawei", "5CB559": "Huawei", "5CCF7F": "Huawei", "602E20": "Huawei",
	"648788": "Huawei", "687F74": "Huawei", "688F84": "Huawei", "6C5AB5": "Huawei",
	"700F6A": "Huawei", "74882A": "Huawei", "78F5FD": "Huawei", "80380D": "Huawei",
	"80B686": "Huawei", "843A5B": "Huawei", "84A8E4": "Huawei", "84BE52": "Huawei",
	"84DBFC": "Huawei", "881528": "Huawei", "88867E": "Huawei", "8C0500": "Huawei",
	"8C25F0": "Huawei", "8CC8CD": "Huawei", "900E83": "Huawei", "90171F": "Huawei",
	"9094E4": "Huawei", "942E17": "Huawei", "943FC2": "Huawei", "94772B": "Huawei",
	"949494": "Huawei", "9839DC": "Huawei", "98E551": "Huawei", "98CDAC": "Huawei",
	"9CC5D2": "Huawei", "9CD917": "Huawei", "A07898": "Huawei", "A0CC2B": "Huawei",
	"A4DCBE": "Huawei", "AC4E91": "Huawei", "AC853D": "Huawei", "B07195": "Huawei",
	"B0E5ED": "Huawei", "B49410": "Huawei", "B8BC1B": "Huawei", "BC2551": "Huawei",
	"BC7670": "Huawei", "C49300": "Huawei", "C4070B": "Huawei", "C808E9": "Huawei",
	"CC08FB": "Huawei", "CC8818": "Huawei", "D46AA8": "Huawei", "D4A148": "Huawei",
	"D8C06A": "Huawei", "DC729B": "Huawei", "E0247F": "Huawei", "E0A3AC": "Huawei",
	"E48326": "Huawei", "E47E66": "Huawei", "E8088B": "Huawei", "E8CD2D": "Huawei",
	"EC388F": "Huawei", "ECE555": "Huawei", "F4028B": "Huawei", "F44E05": "Huawei",
	"F4552C": "Huawei", "F4C714": "Huawei", "F80113": "Huawei", "F83DFF": "Huawei",
	"F8E811": "Huawei", "FC48EF": "Huawei", "FCC897": "Huawei",

	// Ubiquiti
	"04180F": "Ubiquiti", "0418D6": "Ubiquiti", "18E829": "Ubiquiti",
	"245A4C": "Ubiquiti", "24A43C": "Ubiquiti", "44D9E7": "Ubiquiti",
	"68728B": "Ubiquiti", "68D79A": "Ubiquiti", "74ACB9": "Ubiquiti",
	"788A20": "Ubiquiti", "802AA8": "Ubiquiti", "80266C": "Ubiquiti",
	"B4FBE4": "Ubiquiti", "D021F9": "Ubiquiti", "DC9FDB": "Ubiquiti",
	"E063DA": "Ubiquiti", "F09FC2": "Ubiquiti", "FCECDA": "Ubiquiti",

	// Nintendo
	"00191D": "Nintendo", "001AE9": "Nintendo", "001BEA": "Nintendo",
	"001CBE": "Nintendo", "001DBC": "Nintendo", "001E35": "Nintendo",
	"001F32": "Nintendo", "001FC5": "Nintendo", "002147": "Nintendo",
	"00224C": "Nintendo", "0023CC": "Nintendo", "002403": "Nintendo",
	"0024F3": "Nintendo", "0025A0": "Nintendo", "002709": "Nintendo",
	"34AF2C": "Nintendo", "40D28A": "Nintendo", "40F407": "Nintendo",
	"582F40": "Nintendo", "58BDA3": "Nintendo", "8C56C5": "Nintendo",
	"8CCDE8": "Nintendo", "98415C": "Nintendo", "98B6E9": "Nintendo",
	"9CE635": "Nintendo", "A438CC": "Nintendo", "B87826": "Nintendo",
	"B88AEC": "Nintendo", "CC9E00": "Nintendo", "CCFB65": "Nintendo",
	"D8601F": "Nintendo", "D86BF7": "Nintendo", "DC68EB": "Nintendo",
	"E00C7F": "Nintendo", "E0E751": "Nintendo", "E84ECE": "Nintendo",
	"ECC40D": "Nintendo",

	// Sony
	"000AD9": "Sony", "000E07": "Sony", "000FDE": "Sony", "001315": "Sony",
	"00138F": "Sony", "0013A9": "Sony", "001410": "Sony", "0014A4": "Sony",
	"001A80": "Sony", "001D0D": "Sony", "001DBF": "Sony", "001E4C": "Sony",
	"001EAA": "Sony", "001FE1": "Sony", "001FE2": "Sony", "002345": "Sony",
	"0024BE": "Sony", "0024C6": "Sony", "00D9D1": "Sony", "28A02B": "Sony",
	"28C2DD": "Sony", "2C61F3": "Sony", "30F9ED": "Sony", "48F072": "Sony",
	"509EA7": "Sony", "545CF8": "Sony", "589CFC": "Sony", "5C7A5C": "Sony",
	"68E7C2": "Sony", "70D9F6": "Sony", "785C12": "Sony", "7CBB8A": "Sony",
	"84EE86": "Sony", "8C0EE3": "Sony", "986FAD": "Sony", "9CB2B2": "Sony",
	"9CDC71": "Sony", "AC8DF1": "Sony", "B0C8AD": "Sony", "B83B85": "Sony",
	"C4B3D6": "Sony", "C8D762": "Sony", "CC982B": "Sony", "D4F2D8": "Sony",
	"E0DB55": "Sony", "E4A7C5": "Sony", "F0B46E": "Sony", "F0BF97": "Sony",
	"F8424A": "Sony", "FC0FE6": "Sony",

	// LG
	"0013E0": "LG", "0014F6": "LG", "001C62": "LG", "001E75": "LG",
	"001F6B": "LG", "002038": "LG", "00261A": "LG", "10689A": "LG",
	"10F96F": "LG", "14C913": "LG", "28949C": "LG", "30766F": "LG",
	"38B1DB": "LG", "3C2DB7": "LG", "48DBED": "LG", "50B9C0": "LG",
	"5CCD5B": "LG", "6474F6": "LG", "6C5C14": "LG", "708B78": "LG",
	"78F882": "LG", "7C2476": "LG", "800A80": "LG", "80E120": "LG",
	"88C9D0": "LG", "9024A6": "LG", "982CBC": "LG", "9C02D1": "LG",
	"A0E9DB": "LG", "B4B5FE": "LG", "BC4499": "LG", "C49A02": "LG",
	"C4DAF7": "LG", "C8BD61": "LG", "CC6DA0": "LG", "D0D3FC": "LG",
	"DCD255": "LG", "E405F4": "LG", "E83935": "LG", "EC866B": "LG",
	"F072EA": "LG", "F8F1B6": "LG", "FCF8B7": "LG",

	// D-Link
	"001195": "D-Link", "0013D4": "D-Link", "001760": "D-Link", "001B11": "D-Link",
	"001CF0": "D-Link", "00265A": "D-Link", "082E5F": "D-Link", "0840F3": "D-Link",
	"1062EB": "D-Link", "14D64D": "D-Link", "1CBDB9": "D-Link", "1CAFF7": "D-Link",
	"340804": "D-Link", "340A98": "D-Link", "3CFB96": "D-Link", "5C5BC4": "D-Link",
	"5CD998": "D-Link", "7898CC": "D-Link", "78542E": "D-Link", "84C9B2": "D-Link",
	"9CD643": "D-Link", "ACF1DF": "D-Link", "B8A386": "D-Link", "BC0F9A": "D-Link",
	"C0A0BB": "D-Link", "C8BE19": "D-Link", "C8D3A3": "D-Link", "CC988B": "D-Link",
	"CC74C4": "D-Link", "E0ADA0": "D-Link", "F0B4D2": "D-Link", "F8E903": "D-Link",

	// Linksys
	"001217": "Linksys", "001310": "Linksys", "001625": "Linksys", "00195B": "Linksys",
	"001A70": "Linksys", "001C10": "Linksys", "001D7E": "Linksys", "001E58": "Linksys",
	"002265": "Linksys", "00226B": "Linksys", "002369": "Linksys", "0024B2": "Linksys",
	"0025D3": "Linksys", "58C79C": "Linksys", "60A10A": "Linksys", "64A651": "Linksys",
	"6C90E4": "Linksys", "80D4A5": "Linksys", "9C5D12": "Linksys", "A4E9A3": "Linksys",
	"C0C1C0": "Linksys",

	// VMware
	"000569": "VMware", "000C29": "VMware", "001C14": "VMware", "005056": "VMware",
}
