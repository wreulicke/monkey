
## Monkey

[Go言語で作るインタプリタ](https://www.amazon.co.jp/Go%E8%A8%80%E8%AA%9E%E3%81%A7%E3%81%A4%E3%81%8F%E3%82%8B%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%97%E3%83%AA%E3%82%BF-Thorsten-Ball/dp/4873118220/ref=sr_1_1?adgrpid=52270124614&gclid=EAIaIQobChMInZfZycm35wIVR6qWCh0LRg_BEAAYASAAEgJCu_D_BwE&hvadid=338518266894&hvdev=c&hvlocphy=1009692&hvnetw=g&hvpos=1t1&hvqmt=e&hvrand=8174231056717738738&hvtargid=kwd-456677309977&hydadcr=27267_11561158&jp-ad-ap=0&keywords=go%E8%A8%80%E8%AA%9E%E3%81%A7%E3%81%A4%E3%81%8F%E3%82%8B%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%97%E3%83%AA%E3%82%BF&qid=1580808234&sr=8-1)をベースに実装した `monkey` 言語です。

## Usage

```
$ go run .
>> let x = 2; x
2
>> let x = 2; puts(x)
2
null
>> let [x] = [2]
[2]
>> let [x] = [2]; x
2
>> let {x} = {"x": "2"}; x
2
>> 2 | fn(x) { x }
2
>> [1, 2] | fn([x, y]) { x + y }
3
>> {"x": 1, "y": 2} | fn({x, y}) { x + y }
3
```


## 本の実装から拡張された実装

本に載っていない実装として以下を実装しています。

* runeを使ったlexer
* (tokenのみ) 小数点数
  * ParserとInterpreterはサポートしていないのでエラーで落ちます
* エスケープ形式の文字列
* パイプラインオペレータ
* 配列やハッシュ形式のDestructuring
  * let文におけるDestructuring
  * 関数リテラルにおける引数のDestructuring