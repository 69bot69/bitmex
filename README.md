# BitMEX API自動売買プログラム ボットが完成
BitFlyerLightningAPI自動売買プログラムボット開発を行っておりましたが、この度、BitMEXという暗号通貨取引所に主戦場を変更し、BitMEX API自動売買プログラム ボットを完成させました。

このプログラムはGo言語で書かれており、移動平均線9/16のゴールデン・デッドクロスで取引開始するプログラムで、Windows/Mac/Linuxで使用可能なワンソースアプリ化が可能です。
[HowTo原文](http://program.okitama.org/2018/03/2018-03-28_bitmex-api-auto-trade-howto/)

__発注などの基本機能以外の特筆機能__  
- トレンド検出するとメール送信  
- 平均足設定  
- 取引枚数調整  
- レバレッジ調整  
- 証拠金は1000円程度でも可能  
- デモ稼働機能  

<!-- toc -->
<!--more-->

## 稼働前提
下記環境設定まではご自身のOSに従い整備してからお読みください。

- Windows OR Mac OR Linux
- Install Go
- $ go env で環境情報がでるようにPATHを通す
- Goのソースを動かすPCディレクトリを整備する
- githubからソースをダウンロードするためのgitをインストール
- BitMEX API KEY/SECRET KEY（[MEX紹介リンク](https://www.bitmex.com/register/hwAUJG)）

### Go言語インストール参考
-  [はじめての Go 言語 (on Windows)](https://qiita.com/spiegel-im-spiegel/items/dca0df389df1470bdbfa)
- [GO言語をMacで使ってみる インストール](https://qiita.com/Noah0x00/items/63e024f9b5a27276401b)
-  [まっさらなLinux CentOS7に必要なものを初期インストールし設定をする](http://program.okitama.org/2018/01/2017-12-26_cent-os-install-nginx/)

## 参考資料
__テスト稼働損益__  ※ レバレッジ調整や裁量取引を多分に含んでいます。  
![ビットコイン 自動売買プログラム ボット稼働 損益](http://program.okitama.org/img/bitmex_pl.png)


__過去チャートでの試算__  ※ 去年までの上昇相場を含みます。  
![ビットコイン 自動売買プログラム ボット稼働 バックテスト](http://program.okitama.org/img/bitmex_trial.png)

## gitでソースをダウンロード（clone）し、goで動かす
[Source at GitLab](https://gitlab.com/k-terashima/bitmex.pub)

```
// ソースをそのまま動かす
$ git clone git@gitlab.com:k-terashima/bitmex.pub.git bitmex.lib
$ cd ./bitmex.lib
$ go get
$ go run main.go
```

又は

```
// ワンソースビルド
$ git clone git@gitlab.com:k-terashima/bitmex.pub.git bitmex.lib
$ cd ./bitmex.lib
$ go get
$ go build
$ ./bitmex.lib &
```

### 諸設定が空なのでエラーが出ます
{{< hl-text yellow >}}初期設定が必要です{{< /hl-text >}}。各ファイルをテキストファイルで開いて下記初期設定欄にKEYなどを記述ください。まずは、1さえ記述すれば動きます。 記述後、当ファイル名先頭の「アンダーバー（\_）」を消してプログラムを有効にしてください。


1. 初期設定: ./models/\_config-sample.goを参考にKEYなど各種設定を記述ください。  
2. 適宜設定変更: ./evaluations/math2.goにて、移動平均線の数値を変更できます。初期設定はSMA9/16がクロスすると取引開始する設定になっています。  
3. 適宜設定変更: ./order/post-order.goにて、活用証拠金の割合数値を変更できます。初期設定は500%程度を活用する設定になっています。  


### どうなれば動いてるってこと？
__成功例__  
下記画像のようにINFOという文字が表示され、作業ファイルが報告されます。また、過去500時間程度のトレンドシグナルが表示されます。
![成功例](http://program.okitama.org/img/bitmex_success.png)


__失敗例__
下記画像のようにERRORという文字が表示され、エラーファイルが報告されます。テキストからウォレット（証拠金）情報が取得できなかったことが分かり、問題はログイン（KEY/SECRETKEY）だと分かります。
![失敗例](http://program.okitama.org/img/bitmex_error.png)

## FAQ
__Q. ウェブ上から稼働状況確認はできますか？__  
A. はい、できます。  
初期設定では、稼働後[http://localhost:8080](http://localhost:8080) で閲覧可能です。サーバーに設置した時は、ポート参照すれば適宜ドメインから閲覧いただけます。

__Q. サーバーに設置して稼働することは可能ですか？__  
A. 可能です、ぼくもサーバーからビルドしたボットを稼働しています。

__Q. 建玉の決済はいつどこでどういう取引で行いますか？__  
A. トレンドの転換点であるSMA9/16クロスシグナルで取引を開始します。例えばそれが買いシグナルだった場合、売りシグナルが検出されたら建玉を成り売りしてから新規注文を発注します。  
例外的に、取得建玉の価額*0.1以上の変化が起こった時、決済を行います。

なお、__瞬間的__ な価格変動があった場合、決済時に検知価格で決済されるとは限りません。

__Q. 儲かりますか？リスクは？__  
A. 外的な要因が排除できないこともあり、想定外のことが発生します。  
暗号通貨業界はまだまだ未発達です。取引所・税制など各所が未整備なこともあり参加者は恐る恐る、そして学習と実行・反省を繰り返しながら取引をしています。  

それを踏まえ、弊ボットを盲目的に信じず、バックテストをご自身で行ったり、裁量の範囲を検討したり。もちろん、断念したり、ご判断ください。

__Q. いきなり実戦投入ですか？__  
A. [BitMEX TestNet](https://testnet.bitmex.com) というBitMEXが整備しているデモサイトの活用をおすすめします。クローンサイトとなっており、同機能でテスト稼働が行なえます。  
デモ稼働する場合は、models/configファイル内にデモ稼働設定があります。デモ利用の際は、BitMEX本番API KEYでは動きません。BitMEX TestNetから登録し、デモ専用のAPI KEYを取得してください。
```
// true OR false 設定
Demo = true

```

__Q. 改変してもいいですか？__  
A. もちろんです、なにか面白い知見や結果を得ましたら共有してください┏○ﾍﾟｺ

__Q. 改変後販売してもいいですか？__  
A. 禁止しません。  
自社製品として販売の際は、金融商品取引業との関連を要チェックし検討ください。  


## 最後に
投資は自己責任。  
ご自身の裁量内でお楽しみください。

## 投げ銭
配布者ビットコインウォレット（BitPay）:
1JuuQma9xhsH63A99uyqArr681P4ii2DsF

配布者ビットコインキャッシュウォレット（BitPay）:
qqq2ld7yj5j9xe77xqxeclykauve70zlu53urvjctz


## 制作者
質問は、DMやリプライお受けしています。  
初対面前提で適切な言葉遣いでお声がけいただければ幸いです。  
Twitter: [@76c4ca47c252b58a174bc5862fe1a523](https://twitter.com/numbTrade)
