// Score implementation
// ref. https://puyo-camp.jp/posts/166815

package pu2

func BonusChain(chain int) int {
	if chain <= 1 {
		return 0
	}
	switch chain {
	case 2:
		return 8
	case 3:
		return 16
	case 4:
		return 32
	default:
		return 32 + (chain-4)*32
	}
}

func BonusConnection(connect int) int {
	if connect <= 4 {
		return 0
	}
	if connect >= 11 {
		return 10
	}
	return connect - 3
}

func BonusColor(Color int) int {
	return (Color - 1) * 3
}
