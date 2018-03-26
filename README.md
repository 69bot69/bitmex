

## Usage



# BitMEX API 自動売買プログラム ボット仕様解説
====

弊プログラムは、Bitcoinを代表とした暗号通貨（仮想通貨）を取引する最大手取引所BitMEX APIを用いた自動売買プログラムです。  

## Description
2時間足SMA9/16のゴールデン・デッドクロスで取引開始する自動売買プログラムです。トレンド転換時のみ発注を行い、売買方向反転時は建玉を解消し、新規注文を発注します。

## Requirement
- Linux
- Golang
- BitMEX API KEY/SECRET KEY

## Usage
1. ./models/\_config-sample.goを参考にKEYなど各種設定を記述ください。  
2. また、./evaluations/math2.goにて、移動平均線の数値を変更できます。初期設定はSMA9/16
3. ./order/post-order.goにて、活用証拠金の割合数値を変更できます。初期設定は500%程度を活用します。

```
$ go run main.go

又は

$ go build
$ ./bitmex.bot &
```


## Contribution
配布者ビットコインウォレット（BitPay）:
1Mm3R7BsVvT39sKUiYWK9QxZ28f1HPUjLR

配布者ビットコインキャッシュウォレット（BitPay）:
qqq2ld7yj5j9xe77xqxeclykauve70zlu53urvjctz


## Author
質問は、DMやリプライお受けしています。
初対面前提で適切な言葉遣いでお声がけいただければ幸いです。
[76c4ca47c252b58a174bc5862fe1a523](https://twitter.com/numbTrade)
