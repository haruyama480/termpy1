# termpy1

とこぷよをCLIからできるようにしました。
暇つぶしに作ったので洗練する予定はありません。

前提
- [Go](https://go.dev/doc/install)
- 「⬤」を全角で表示してくれるターミナルエミュレータ

インストール方法
```
go install github.com/haruyama480/termpy1/cmd/toko@latest
```

実行方法
```
toko
```

使い方
- `d`,`f` 横移動
- `j`,`k` 左回転、右回転
- `space` 即落下
- `z` 一手戻る
- `q` 終了
