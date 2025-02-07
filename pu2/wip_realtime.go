package pu2

// type RState int

// const (
// 	RStateNone RState = iota
// 	RStatePlayTsumo
// 	RStateFallTsumo
// 	RStateFallBody
// 	RStateFallOjama
// 	RStateIdle
// 	RStateDead
// )

// type RField struct {
// 	Field

// 	State RState

// 	// T: tsumo
// 	TP Pair

// 	TY      int
// 	TYDelay int // 自由落下は8フレーム。半マスの都合があるので0-16で表現する。0は設置時に利用

// 	TX       int
// 	TXRight  int // 何フレーム連続で入力されたか
// 	TXLeft   int
// 	TXTarget int // 最終的な位置
// 	TXDelay  int // 終了まで何フレームか。横移動は2フレーム

// 	TDir       Dir
// 	TDirTarget Dir // 最終的な向き
// 	TDirDetail int // 0-7

// 	// used tsumo, body, ojama
// 	FallFrameCount int

// 	IdleFrameCount int
// }

// // 1フレームで可能な操作
// type ROp struct {
// 	Seek   int
// 	Rotate int
// 	Drop   int
// }

// func (r *RField) TTsumo(p Pair) {
// 	r.TP = p

// 	r.TY = 12
// 	r.TYDelay = 8 // 12段目の0.5マス上から開始される

// 	r.TX = 2
// 	r.TXRight = 0
// 	r.TXLeft = 0
// 	r.TXTarget = 2
// 	r.TXDelay = 0

// 	r.TDir = DirUp
// 	r.TDirTarget = DirUp
// 	r.TDirDetail = 0
// }

// // TDirDetail の距離
// func dist(a, b int) int {
// 	return (a - b + 8) % 8
// }

// // 設置されたらtrue
// // FIXME
// // - 回転時の移動がサポートされない
// // - 縦回転がサポートされない
// // - 横移動がサポートされない
// // キャンセルの再現がつらすぎるのであきらめようかな...
// func (r *RField) TNext(op ROp) bool {
// 	if op.Rotate > 0 {
// 		r.TDirTarget = r.TDirTarget.Inc()
// 	} else if op.Rotate < 0 {
// 		r.TDirTarget = r.TDirTarget.Dec()
// 	}
// 	// TODO: 移動するか判定
// 	// TDirTarget に横ぷよが存在したらスライド。スライドだめなら縦回転のフラグを建てる
// 	// TDirTarget に下にぷよが存在したら1マスあげる

// 	if r.TDir != r.TDirTarget {
// 		td := int(r.TDirTarget) * 2
// 		if dist(r.TDirDetail, td) > dist(r.TDirDetail+1, td) { // 右回転のが近い
// 			r.TDirDetail = (r.TDirDetail + 1) % 8
// 		} else {
// 			r.TDirDetail = (r.TDirDetail + 7) % 8
// 		}
// 		r.TDir = Dir(r.TDirDetail / 2)
// 	}

// 	// TODO: 横移動
// 	// - 連続入力でなければy移動はなし
// 	if op.Seek > 0 {
// 	}

// 	drop := false
// 	if op.Drop > 0 {
// 		if r.TYDelay <= 8 {
// 			drop = true
// 			r.TYDelay += 8
// 		} else {
// 			r.TYDelay -= 8
// 		}
// 	} else {
// 		if r.TYDelay == 1 {
// 			drop = true
// 			r.TYDelay = 16
// 		} else {
// 			r.TYDelay--
// 		}
// 	}
// 	if drop {
// 		nexty := r.TY - 1
// 		if nexty < 0 || r.Field[r.TXTarget][r.TY] != 0 {
// 			return true
// 		}
// 	}

// 	return false
// }

// func (f *Field) RPossibleOps(x int, dir Dir, newx int, newdir int) []ROp {
// 	return nil
// 	// どのような操作によってツモをその位置に移動させられるかを判定する
// 	// 不可能なら空
// }
