# BitMEX API 自動売買プログラム ボット仕様解説
====

弊プログラムは、Bitcoinを代表とした暗号通貨（仮想通貨）を取引する最大手取引所BitMEX APIを用いた自動売買プログラムです。
Linux/Windows/Mac上で動き、Go言語で動いています。

## 簡易説明
2時間足SMA9/16のゴールデン・デッドクロスで取引開始する自動売買プログラムです。トレンド転換時のみ発注を行い、売買方向反転時は建玉を解消し、新規注文を発注します。  

### 機能
- トレンド検出するとメール送信
- 平均足設定変更
- 取引枚数調整
- レバレッジ調整

## 必要なもの
- OS: Linux/Windows/Mac
- Golang
- qct/bitmex-go
- BitMEX API KEY/SECRET KEY

## Usage  

1. 初期設定: ./models/\_config-sample.goを参考にKEYなど各種設定を記述ください。  
2. 適宜設定変更: ./evaluations/math2.goにて、移動平均線の数値を変更できます。初期設定はSMA9/16
3. 適宜設定変更: ./order/post-order.goにて、活用証拠金の割合数値を変更できます。初期設定は500%程度を活用します。

```
$ git clone git@gitlab.com:k-terashima/bitmex.pub.git bitmex.lib
$ cd ./bitmex.lib
$ go get
$ go run main.go
```

又は

```
$ git clone git@gitlab.com:k-terashima/bitmex.pub.git bitmex.lib
$ cd ./bitmex.lib
$ go get
$ go build
$ ./bitmex.lib &
```


## 投げ銭
配布者ビットコインウォレット（BitPay）:
1Mm3R7BsVvT39sKUiYWK9QxZ28f1HPUjLR

配布者ビットコインキャッシュウォレット（BitPay）:
qqq2ld7yj5j9xe77xqxeclykauve70zlu53urvjctz


## 制作者
質問は、DMやリプライお受けしています。
初対面前提で適切な言葉遣いでお声がけいただければ幸いです。
[76c4ca47c252b58a174bc5862fe1a523](https://twitter.com/numbTrade)
